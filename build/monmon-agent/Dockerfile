# Build
FROM golang:1.22-bookworm AS build

RUN apt-get update \
  && DEBIAN_FRONTEND=noninteractive \
    apt-get install --no-install-recommends --assume-yes \
      protobuf-compiler \
      libprotobuf-dev

RUN mkdir -p /opt/monmon
WORKDIR /opt/monmon

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .
RUN make build

# Release
FROM alpine:3

COPY --from=build /opt/monmon/bin/monmon-agent /bin/monmon-agent
COPY --from=build /opt/monmon/config/monmon.yaml /etc/monmon/agent.yaml

CMD ["/bin/monmon-agent", "start", "-config", "/etc/monmon/agent.yaml"]
