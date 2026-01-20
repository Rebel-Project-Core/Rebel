FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG VERSION
ARG COMMIT
ARG DATE

RUN CGO_ENABLED=0 go build -ldflags "-s -w -X 'credo/version.Version=${VERSION}' -X 'credo/version.Commit=${COMMIT}' -X 'credo/version.BuildDate=${DATE}'" -o /usr/bin/credo .

FROM ubuntu:24.04 AS final

ENV \
	DEBIAN_FRONTEND=noninteractive \
	LANG="C.UTF-8"

SHELL ["/bin/bash", "-o", "pipefail", "-c"]

RUN set -x \
	&& apt-get update -yq --no-install-recommends \
	&& apt-get install -yq --no-install-recommends \
	build-essential \
	gfortran \
	ca-certificates \
	&& rm -rf /var/lib/apt/lists/* \
	&& apt-get clean

COPY --from=builder /usr/bin/credo /usr/bin/credo

RUN set -x \
	&& chmod +x /usr/bin/credo \
	&& mkdir -p /workdir

WORKDIR /workdir
