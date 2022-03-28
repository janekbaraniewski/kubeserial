FROM --platform=$BUILDPLATFORM golang:1.13-alpine AS build
WORKDIR /go/src/github.com/janekbaraniewski/kubeserial
COPY go.mod go.sum .
RUN go mod download
COPY . .
ARG TARGETOS TARGETARCH TARGETVARIANT
RUN if [[ -n "${TARGETVARIANT}" ]]; then export GOARM=${TARGETVARIANT}; fi
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -o /build/bin/kubeserial cmd/manager/main.go

FROM scratch
ENV OPERATOR=/usr/local/bin/kubeserial \
    USER_UID=1001 \
    USER_NAME=kubeserial
COPY --from=build /build/bin/kubeserial ${OPERATOR}
COPY build/bin/entrypoint /usr/local/bin/entrypoint

ENTRYPOINT [ "/usr/local/bin/entrypoint" ]

USER ${USER_UID}
