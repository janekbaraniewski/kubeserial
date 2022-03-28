FROM --platform=$BUILDPLATFORM golang:1.13-alpine AS build
WORKDIR /go/src/github.com/janekbaraniewski/kubeserial
COPY go.mod go.sum .
RUN go mod download
COPY . .
ARG TARGETOS TARGETARCH TARGETVARIANT
RUN if [[ -n "${TARGETVARIANT}" ]]; then export GOARM=${TARGETVARIANT}; fi
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -o /build/bin/kubeserial cmd/manager/main.go
CMD ["/build/bin/kubeserial"]
