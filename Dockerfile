# syntax=docker/dockerfile:1
FROM golang:1.27@sha256:f20d83a83e18dbd2a4714b8ffd4b30b287b6243962dd75608e86abcc63ba4efc AS base

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
