# syntax=docker/dockerfile:latest

# Build the manager binary
FROM --platform=$BUILDPLATFORM golang:1.18-alpine as builder
RUN apk update
RUN apk add make bash
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
RUN --mount=type=cache,target=/root/.cache/go-build GOOS=$TARGETOS GOARCH=$TARGETARCH \
    KUBESERIAL_BUILD_OUTPUT_PATH=/build/bin/kubeserial \
    make kubeserial

FROM alpine
LABEL org.opencontainers.image.source https://github.com/janekbaraniewski/kubeserial
WORKDIR /
ENV USER_NAME=kubeserial \
    USER_UID=1001
COPY --from=builder /build/bin/kubeserial .
COPY entrypoint entrypoint


ENTRYPOINT ["/entrypoint", "/kubeserial"]
USER ${USER_UID}
