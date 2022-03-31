# Build the manager binary
FROM --platform=$BUILDPLATFORM golang:1.18-alpine as builder
RUN apk update
RUN apk add make bash
WORKDIR /workspace
COPY go.mod go.sum .
RUN go mod download

# Copy the go source
COPY Makefile Makefile
COPY cmd cmd
COPY pkg pkg

ARG TARGETOS TARGETARCH TARGETVARIANT
RUN if [[ -n "${TARGETVARIANT}" ]]; then export GOARM=${TARGETVARIANT}; fi
# Build
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH \
    GO_BUILD_OUTPUT_PATH=/build/bin/kubeserial \
    make kubeserial

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /build/bin/kubeserial .
USER 65532:65532

ENTRYPOINT ["/kubeserial"]
