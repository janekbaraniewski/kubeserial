# End-to-End (e2e) Testing for kubeserial

Status: **VERIFIED**. The full suite (E1-E7, 8 specs, 0 skipped) has been run
green against a live kind cluster (kind v0.18.0, `kindest/node` v1.27.1) via
`make test-e2e`, including device simulation on a real node and the webhook
admission round-trip. See [Verified run](#0-verified-run) for the passing output
and [Open items](#6-open-items-and-caveats) for remaining caveats.

## 0. Verified run

`make test-e2e` performs the entire flow from scratch (create cluster -> build 3
images -> kind load -> install cert-manager + self-signed Issuer -> install both
Helm charts -> run specs -> teardown) and passed:

```
==> Running e2e specs
Ran 8 of 8 Specs in 55.419 seconds
SUCCESS! -- 8 Passed | 0 Failed | 0 Pending | 0 Skipped
==> e2e suite passed
==> Deleting kind cluster kubeserial-e2e
```

Specs that ran (all passed, none skipped):

- E1 SerialDevice CRD registered
- E2 manager Deployment has a ready replica
- E3 device with no backing /dev node stays not-Available
- E4 device appears -> Available + Free flip True, Status.NodeName set
- E5 device disappears -> Available flips back to False, NodeName cleared
- E6 webhook injects socat bridge when device is Free
- E7 webhook does not inject when device not Free, and not when no annotation (2 specs)

---

## 1. What we are testing and why e2e (not just unit/envtest)

kubeserial has three runtime components:

| Component | Binary | Runs as | What it does |
|-----------|--------|---------|--------------|
| **manager** | `cmd/manager` | Deployment (controller-runtime) | Reconciles `SerialDevice`, `Manager`, `ManagerScheduleRequest`, `KubeSerial` CRs. Creates the device-monitor DaemonSet, gateway, etc. |
| **device-monitor** | `cmd/device-monitor` | DaemonSet (privileged, hostPath `/dev`) | Every 1s, `os.Stat("/dev/<SerialDevice.metadata.name>")`. If present, flips the `SerialDevice` `Available`/`Free` conditions to `True` and sets `Status.NodeName`. If it disappears, flips them back. |
| **webhook** | `cmd/webhook` | Deployment + `MutatingWebhookConfiguration` | On pod create with annotation `app.kubeserial.com/inject-device: <name>`, if that device's `Free` condition is `True`, rewrites `containers[0]` command to wrap it with a `socat` bridge to `tcp:<device>-gateway:3333`. |

The existing test pyramid:

- **Unit tests** (`pkg/controllers/*_test.go`, `pkg/monitor`, `pkg/webhooks`) use
  the controller-runtime **fake client** and an **afero in-memory `FileSystem`**
  (`utils.InMemoryFS`). These already cover reconciler logic and the monitor's
  state-transition logic with a faked `/dev`.
- **Integration tests** (`pkg/controllers/integration_tests/`, currently
  commented out) target **envtest** (a real apiserver + etcd, no kubelet).

What unit/envtest **cannot** cover, and what therefore justifies e2e:

1. **The device-monitor against a real `/dev`.** envtest has no kubelet, so it
   cannot schedule a DaemonSet, and the monitor's `os.Stat("/dev/...")` runs
   against the real OS filesystem (`utils.OsFS`), not afero. The
   appear/disappear -> condition-flip behavior is only exercised end-to-end when
   a real pod runs on a real node with a real (or simulated) device node.
2. **RBAC, ServiceAccount, leader election, image wiring** -- only real on a cluster.
3. **The mutating webhook in the admission path** -- including TLS cert wiring,
   `MutatingWebhookConfiguration` registration, and the apiserver actually
   calling the webhook on pod create. (The handler logic itself is unit-testable;
   the *registration + admission round-trip* is not.)

---

## 2. The hard problem: simulating a serial/USB device on a node

The monitor's only contract with hardware is **one line**
(`pkg/monitor/monitor.go`):

```go
if _, err := m.fs.Stat("/dev/" + name); os.IsNotExist(err) { return false }
```

Note: `name` is the **`SerialDevice.metadata.name`**, not `idVendor`/`idProduct`.
The monitor does **not** open the device, does **not** issue USB/tty ioctls, and
does **not** query udev or netlink. The `udev-monitor` sidecar container in the
DaemonSet runs `systemd-udevd` so that *real* hardware gets a stable
`/dev/<name>` node via udev rules, but the Go monitor itself only cares that the
path exists.

This single fact drives the entire strategy.

### Options evaluated

#### Option A -- umockdev (record/replay udev+sysfs+ioctl) -- REJECTED for the monitor path

[umockdev](https://github.com/martinpitt/umockdev) is the gold standard for
mocking udev/sysfs/ioctl device trees. **However**, its mock devices are visible
**only** to processes launched under `umockdev-run` (or wrapped with
`LD_PRELOAD=libumockdev-preload.so`). The preload library intercepts
`open()`/`stat()`/`ioctl()`/netlink and re-routes them into a sandbox; processes
**outside** that wrapper see nothing. Confirmed from the upstream README:

> "the preload library intercepts access to /sys, /dev/, /proc/, the kernel's
> netlink socket and ioctl() and re-routes them into the sandbox ... stat()
> delivers a block/char device with appropriate major/minor."

Consequences for us:

- To use umockdev for the monitor, we would have to launch the **device-monitor
  binary itself** under `umockdev-run`. That conflicts with the shipped
  `Dockerfile.monitor` (systemd + `systemd-udevd` as PID 1) and would require a
  test-only image/entrypoint. It also means we would be testing the binary in a
  non-production launch configuration.
- umockdev's value is in faking the **udev/sysfs/ioctl** layer. But the monitor
  consumes **none** of that -- only `stat()`. So umockdev would add a large,
  process-isolated, LD_PRELOAD-shaped dependency to fake a layer the code under
  test never reads.

**Verdict:** umockdev is the right tool *if and when* we add tests for the
`udev-monitor`/`systemd-udevd` sidecar and real udev-rule-driven naming
(`98-devices.rules`), or if the Go monitor ever starts issuing ioctls. It is the
**wrong** tool for the current `stat()`-only contract. We keep it documented as
the chosen approach for a *future* "udev rules produce the right `/dev/<name>`"
test, runnable as a standalone `go test` under `umockdev-run` on a Linux runner.

Concrete future recipe (for lane U1):

```sh
# inside a Linux container/runner with umockdev installed
umockdev-run \
  --device tty-usb-serial.umockdev \
  -- go test ./test/udev/...
# where tty-usb-serial.umockdev was produced on real hardware via:
#   umockdev-record /dev/ttyUSB0 > tty-usb-serial.umockdev
# and the test asserts udev applies 98-devices.rules to create /dev/<symlink>.
```

#### Option B -- Real device node / PTY on the node's hostPath `/dev` -- CHOSEN

The DaemonSet mounts the **host's `/dev`** (`hostPath: /dev`) into both
containers, privileged. On kind, "the host" is the kind **node container**. So
if we create a node at `/dev/<name>` *inside the kind node container*, the
hostPath mount surfaces it into the monitor pod, and `os.Stat` succeeds -- exactly
as a real plugged-in device would (post-udev-rename).

Ways to create that node, in order of preference:

1. **A privileged "device-simulator" Job/Pod** scheduled on the node that runs
   `socat PTY,link=/dev/<name> PTY` (creates a real char device pair), or simply
   `mkfifo /dev/<name>` / `mknod`, then sleeps. Deleting the pod (or having it
   `rm` the node) simulates unplug. This keeps the whole harness *inside*
   Kubernetes, so it works identically on kind, k3d, or a real cluster, and needs
   no `docker exec`. **This is the primary mechanism.**
2. **`docker exec <kind-node> sh -c 'socat ... &'`** from the test process. Simpler
   to reason about, but kind-specific and bypasses Kubernetes. Used as a fallback
   helper for local debugging.

Why a PTY (`socat PTY,...`) over a plain file:

- A PTY produces a genuine **character device** with a tty major/minor, which is
  the closest cheap analogue to `/dev/ttyUSB0`. If the monitor (or a future
  webhook/socat consumer) ever opens it, reads/writes succeed.
- A plain `touch`/`mkfifo` also satisfies `os.Stat` (sufficient for *today's*
  monitor), and is the trivial fallback when `socat` is unavailable.

**Trade-off / honest caveat:** this does **not** exercise udev or the
`udev-monitor` sidecar at all -- we create the final `/dev/<name>` directly,
skipping the discovery/rename pipeline. That pipeline is what Option A would
cover. The test matrix below makes this split explicit.

#### Option C -- tty0tty / null-modem emulator -- REJECTED (overkill)

`tty0tty` creates connected `/dev/tnt*` null-modem pairs. Useful if we needed two
ends of a serial link (e.g. to assert the `socat` bridge actually passes bytes).
For the monitor's presence check it adds a kernel-module/build dependency for no
benefit. Revisit only for a "data actually flows through the gateway" test.

#### Option D -- Kernel gadget / `g_serial`, `dummy_hcd`, `usbip`, QEMU USB passthrough -- REJECTED for CI

All require loading kernel modules (`modprobe g_serial`, `dummy_hcd`) or a real
USB host controller. GitHub Actions `ubuntu-latest` runners do **not** reliably
allow loading arbitrary out-of-tree/host modules, and kind nodes share the host
kernel. `usbip` needs a real exporting host. QEMU USB passthrough needs a VM and
real hardware. These are appropriate for a **self-hosted bare-metal runner**,
not for hosted CI. Documented as the path for a future hardware-in-the-loop lane.

### Decision

| Layer under test | Mechanism | Lane |
|---|---|---|
| Monitor presence check (`stat /dev/<name>`) -> condition flip | **Option B**: PTY/file on node `/dev` via privileged simulator pod | kind, hosted CI |
| udev rule -> correct `/dev/<name>` naming (sidecar) | **Option A**: umockdev under `umockdev-run`, standalone `go test` | Linux CI, future |
| Real serial bytes through gateway/socat | Option C (tty0tty) | self-hosted, future |
| Real USB enumeration | Option D | self-hosted bare metal, future |

---

## 3. Cluster substrate and CI architecture

### envtest vs kind vs k3d vs minikube

- **envtest** -- apiserver + etcd only, **no kubelet**. Cannot run a DaemonSet or
  mount hostPath `/dev`. Great for the reconciler integration tests (and that is
  where they belong), useless for the device-monitor e2e. **Not used for e2e.**
- **kind** -- Kubernetes-in-Docker. Nodes are containers; their `/dev` is what the
  DaemonSet's hostPath mounts, so we can inject simulated devices into a node by
  creating nodes inside the node container. First-class GitHub Actions support
  via [`helm/kind-action`](https://github.com/helm/kind-action) (latest v1.14.0,
  kind v0.31.0 / k8s v1.35.0 at time of writing). **Chosen.**
- **k3d** -- k3s-in-Docker, also container nodes; would work similarly. Rejected
  only to avoid a second tool; kind is the kubebuilder/controller-runtime default
  and what most operator e2e suites use. Easy to swap (`kind create` -> `k3d
  cluster create`) since the harness shells out.
- **minikube** -- the repo already references a `MINIKUBE_PROFILE` for local dev.
  Heavier in CI (VM driver) and device injection is driver-dependent. Kept as a
  local-only option; CI uses kind.

### Surfacing a simulated device into a kind node

The kind node is a Docker container. Its `/dev` is a container `/dev`. The
DaemonSet's `hostPath: /dev` therefore points at the **node container's** `/dev`.
We create `/dev/<name>` there via the privileged simulator pod (section 2 Option
B), which kind/k3d/real-cluster all support identically. No host-machine `/dev`
pollution; everything is inside the ephemeral node container.

**Mount-propagation question (resolved):** the design doc originally flagged a
risk that a device created in the simulator pod's hostPath `/dev` mount might not
propagate to the separate monitor pod. Verified on kind v0.18 / node v1.27.1:
**no extra handling is needed.** A `socat PTY,link=/dev/<name>` created in the
simulator pod is immediately visible inside the monitor pod via the shared
hostPath, and the monitor's `os.Stat` (which follows the symlink to the PTY)
succeeds. On unplug (pod deletion) the PTY closes; `os.Stat` on the now-dangling
symlink returns `IsNotExist`, so the monitor correctly reports the device gone
even though the symlink itself lingers. No `/host-dev` indirection was required.

### Prerequisite: cert-manager (webhook TLS)

The kubeserial chart's webhook ships a `cert-manager.io/v1` `Certificate` and a
`MutatingWebhookConfiguration` annotated with `cert-manager.io/inject-ca-from`.
So the e2e harness installs **cert-manager** and a **self-signed `Issuer`** in
the install namespace before the chart, and passes
`--set certManagerIssuer.name=kubeserial-selfsigned`. cert-manager's CA injector
then populates the webhook's `caBundle`, which is what makes the apiserver trust
and call the webhook. `hack/e2e.sh` does this automatically.

### Image build note (Makefile quirk)

`make docker-local` works for the manager image but the `device-monitor` and
`injector-webhook` docker targets hardcode a `--cache-to/--cache-from=type=
registry` build cache that requires registry auth and fails offline (the manager
target uses `?=` so it can be overridden; the other two use `=` so it cannot).
`hack/e2e.sh` therefore invokes `docker buildx build ... --load` directly for all
three images. Fixing the Makefile to use `?=` consistently is a separate cleanup
left to the maintainer (it is product/build tooling, not the e2e harness).

### CI flow (`.github/workflows/e2e.yml`)

```
checkout -> setup-go -> build 3 images (make *-docker-local, --load)
        -> helm/kind-action (create cluster)
        -> kind load docker-image (all 3 images into node)
        -> helm install kubeserial-crds + kubeserial chart
        -> go test ./test/e2e (Ginkgo) -- installs simulator, asserts conditions
        -> (always) kind export logs + delete cluster
```

Initially gated on `workflow_dispatch` (manual), so it does not block every PR
while it stabilizes.

---

## 4. Test framework: Ginkgo/Gomega vs plain `go test`

**Chosen: Ginkgo v2 + Gomega.** Rationale:

- It is the controller-runtime / kubebuilder convention; the repo already vendors
  `github.com/onsi/ginkgo/v2` and `gomega`, and the (commented-out) integration
  suite was already written for it.
- `Eventually(...).Should(...)` is the natural fit for asserting on
  asynchronously-reconciled status conditions (the monitor's 1s poll, controller
  requeues) without hand-rolled retry loops.
- `BeforeSuite`/`AfterSuite` cleanly own cluster-scoped setup/teardown and log
  export on failure.

Plain `go test` would work (Ginkgo suites are just `go test`), but we would
re-implement polling and structured reporting. The suite is still invoked through
`go test ./test/e2e/...` so no extra CLI is mandatory in CI.

---

## 5. Test matrix

| ID | Behavior | Needs real device? | Lane | Status |
|----|----------|--------------------|------|--------|
| E1 | CRDs install; `SerialDevice` list succeeds | no | kind smoke | **passing** |
| E2 | Manager Deployment has a ready replica | no | kind smoke | **passing** |
| E3 | `SerialDevice` with **no** device node stays **not-Available** | no | kind smoke | **passing** |
| E4 | Device appears on node -> monitor flips `Available`+`Free` to True, sets `NodeName` | **yes (Option B)** | kind device | **passing** |
| E5 | Device disappears -> conditions flip back, `NodeName` cleared | **yes (Option B)** | kind device | **passing** |
| E6 | Webhook injects `socat` wrapper when device `Free` | no real hw (sets `Free=True` on status + webhook deployed) | kind | **passing** |
| E7 | Webhook does **not** inject when device not `Free`, and not when no annotation | no | kind | **passing** (2 specs) |
| U1 | udev rule renames raw device to `/dev/<name>` | n/a (umockdev) | Linux CI | documented only (Option A), no code |

The **webhook path (E6/E7) needs no real hardware** -- it reads the `Free`
condition off the CR, which the suite sets directly via the status subresource.

---

## 6. Open items and caveats

What is fully working (verified, see [section 0](#0-verified-run)): E1-E7 on
kind via `make test-e2e`, including device simulation and webhook admission.

Remaining items, in priority order:

1. **U1 (udev lane) is documented but not implemented.** The current device
   simulation (Option B) creates the final `/dev/<name>` directly and so does
   **not** exercise udev, the `udev-monitor` sidecar, or `98-devices.rules`. The
   umockdev-under-`umockdev-run` approach (Option A) is the way to cover that and
   is described in section 2; it has no code yet.
2. **kind version in CI vs local.** Locally this was verified with the
   environment's kind v0.18.0, which caps Kubernetes at v1.27.1 (the node image
   is pinned in `hack/e2e.sh`). CI uses `helm/kind-action` (newer kind); the
   workflow passes `E2E_NODE_IMAGE=""` so kind picks its own newer default node
   image. Confirm a CI run is green on the newer Kubernetes.
3. **Makefile build-cache quirk (product/build tooling, not the harness).** The
   `device-monitor` and `injector-webhook` `*-docker` targets hardcode a
   registry build cache and cannot be overridden, so `hack/e2e.sh` calls
   `docker buildx build --load` directly. A maintainer cleanup to use `?=`
   consistently would let `make docker-local` work offline.
4. **Single-node scheduling.** kind's single control-plane node was schedulable
   without untainting in the verified environment. On stock kind the
   control-plane may carry `node-role.kubernetes.io/control-plane:NoSchedule`;
   if the chart pods stay Pending, untaint with
   `kubectl taint nodes --all node-role.kubernetes.io/control-plane-`.

---

## 7. Running it

```sh
# Local, requires Docker + kind + helm + kubectl on PATH:
make test-e2e

# Skip cluster create/teardown (reuse an existing kind cluster):
E2E_SKIP_CLUSTER_SETUP=true E2E_KUBECONTEXT=kind-kubeserial-e2e \
  go test ./test/e2e/... -v
```

Environment knobs honored by the suite (see `test/e2e/helpers.go`):

| Env var | Default | Meaning |
|---------|---------|---------|
| `E2E_SKIP_CLUSTER_SETUP` | `false` | If `true`, the suite assumes a cluster + chart are already installed and only runs specs. |
| `E2E_KIND_CLUSTER` | `kubeserial-e2e` | kind cluster name. |
| `E2E_NAMESPACE` | `kubeserial` | Namespace the chart installs into. |
| `E2E_IMAGE_TAG` | `local` | Image tag built by `make *-docker-local` and loaded with `kind load`. |
| `E2E_SKIP_DEVICE_SIM` | `false` | If `true`, device specs (E4/E5) `Skip()`. Default `false` (device sim is verified on kind); set `true` on substrates without hostPath `/dev`. |
| `E2E_NODE_IMAGE` | pinned v1.27.1 | `kindest/node` image. Empty string -> let kind choose its default (use with newer kind in CI). |
| `E2E_CERT_MANAGER_VERSION` | `v1.14.5` | cert-manager release installed before the chart (webhook certs depend on it). |
