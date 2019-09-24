#!/bin/bash
#
# Perfoms basic liveness check for the cilium-agent

set -euo pipefail

cilium status --brief --timeout 2s 2>&1 /dev/null
cilium service list 2>&1 /dev/null
