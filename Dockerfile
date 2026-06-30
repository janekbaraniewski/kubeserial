# syntax=docker/dockerfile:latest

# Build the manager binary
FROM --platform=$BUILDPLATFORM golang:1.26-alpine AS builder
RUN apk add --no-cache make bash
WORKDIR /workspace
COPY go.mod go.sum ./
RUN go mod download

# Copy the go source
COPY Makefile.build Makefile
COPY cmd/manager cmd/manager
COPY pkg pkg

ARG TARGETOS TARGETARCH TARGETVARIANT
RUN if [[ -n "${TARGETVARIANT}" ]]; then export GOARM=${TARGETVARIANT}; fi
# Build
RUN --mount=type=cache,target=/root/.cache/go-build CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    KUBESERIAL_BUILD_OUTPUT_PATH=/build/bin/kubeserial \
    make kubeserial

FROM gcr.io/distroless/static:nonroot
LABEL org.opencontainers.image.source=https://github.com/janekbaraniewski/kubeserial
WORKDIR /
COPY --from=builder /build/bin/kubeserial .

USER 65532:65532
ENTRYPOINT ["/kubeserial"]
