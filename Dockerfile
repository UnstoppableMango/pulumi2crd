# syntax=docker/dockerfile:1
FROM golang:1.27@sha256:e432b43af23a9328d56a7c499be0476810aa344acbcf65fc7c455d4ff5a40602 AS base

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
