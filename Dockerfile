# syntax=docker/dockerfile:1
FROM golang:1.24@sha256:d2d2bc1c84f7e60d7d2438a3836ae7d0c847f4888464e7ec9ba3a1339a1ee804 AS base

ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH

FROM base AS download
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

FROM download AS build
WORKDIR /src
COPY pkg ./pkg
COPY main.go ./

RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /out/pulumi2crd

FROM --platform=$BUILDPLATFORM scratch
COPY --from=build /out/pulumi2crd /usr/bin/
ENTRYPOINT [ "/usr/bin/pulumi2crd" ]
