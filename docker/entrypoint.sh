#!/bin/sh
set -eu

mkdir -p /data

if [ "$(id -u)" = "0" ]; then
  chown -R agent-runtime:agent-runtime /data
  exec su-exec agent-runtime /usr/local/bin/agent-runtime "$@"
fi

exec /usr/local/bin/agent-runtime "$@"
