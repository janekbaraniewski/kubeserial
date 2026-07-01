# Changelog

## [0.2.4](https://github.com/janekbaraniewski/kubeserial/compare/0.2.3...0.2.4) (2026-07-01)


### Bug Fixes

* **api:** remove stray dupa field from ManagerSpec ([1d4d973](https://github.com/janekbaraniewski/kubeserial/commit/1d4d973dbc4682caa09e4d939a6a879c997adb3c))
* **ci:** pin docs mdBook toolchain to preprocessor-compatible versions ([886d786](https://github.com/janekbaraniewski/kubeserial/commit/886d786f9fc27a348cf4f064b44f1a1cfee14091))
* **ci:** unbreak lint, chart-testing, and multi-arch docker ([a7c16d7](https://github.com/janekbaraniewski/kubeserial/commit/a7c16d7084b3fa4fdf3b1bf615330fdd5d49c2ec))
* webhook decoder wiring, monitor list error, device type naming ([f6e69c9](https://github.com/janekbaraniewski/kubeserial/commit/f6e69c96d82662432baf3d52bbf08be19b2fc1a9))


### Refactors

* idiomatic cleanups in monitor and utils ([afa5bbd](https://github.com/janekbaraniewski/kubeserial/commit/afa5bbdb61404e170cb92fb8b5cc9249f6be4817))


### CI/CD

* add release-please for automated release PRs ([c761882](https://github.com/janekbaraniewski/kubeserial/commit/c7618825e93d51d6f60a1ce3be26ed5409fb60d8))
* add Renovate config with k8s lockstep grouping and automerge ([c4417db](https://github.com/janekbaraniewski/kubeserial/commit/c4417db21d644fc1c7b80f942821f02d69356b19))
* **e2e:** run e2e suite on every PR to master ([a6c6242](https://github.com/janekbaraniewski/kubeserial/commit/a6c6242086abba0f73d245e1f82ff5b06bf017d5))
* modernize release workflows (attestations, signing, OCI charts) ([ab7915a](https://github.com/janekbaraniewski/kubeserial/commit/ab7915ad05219b7fcd155764192b0af45cacf6a9))
* modernize workflows, Docker images, golangci-lint v2, envtest ([d94690c](https://github.com/janekbaraniewski/kubeserial/commit/d94690c16f4c6bd26999bc82ef96cb110a37645a))
