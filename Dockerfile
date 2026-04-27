ARG BUILDPLATFORM=linux/amd64
ARG TARGETPLATFORM=linux/amd64

FROM --platform=$BUILDPLATFORM golang:1.24-alpine AS build

WORKDIR /src
COPY go.mod ./
COPY cmd ./cmd
COPY internal ./internal

ARG TARGETOS=linux
ARG TARGETARCH=amd64
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -trimpath -ldflags="-s -w" -o /out/agent-runtime ./cmd/agent-runtime

FROM --platform=$TARGETPLATFORM alpine:3.21

RUN addgroup -S agent-runtime && adduser -S -G agent-runtime agent-runtime \
    && mkdir -p /data /etc/agent-runtime \
    && chown -R agent-runtime:agent-runtime /data /etc/agent-runtime

COPY --from=build /out/agent-runtime /usr/local/bin/agent-runtime
COPY configs/container.json /etc/agent-runtime/config.json

USER agent-runtime
EXPOSE 8080
ENV AGENT_RUNTIME_CONFIG=/etc/agent-runtime/config.json

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD wget -qO- http://127.0.0.1:8080/api/health >/dev/null || exit 1

ENTRYPOINT ["/usr/local/bin/agent-runtime"]
