# syntax=docker/dockerfile:1
FROM golang:1.27@sha256:5f1b98fa02fe396ebce2a0227a0e08309e11792915b9dd2e1f71053d7e7e6666 AS base

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
