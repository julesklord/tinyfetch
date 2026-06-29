#!/usr/bin/env bash
# Kubernetes current context and namespace plugin for tinyfetch
set -euo pipefail

# Check if kubectl exists
if ! command -v kubectl >/dev/null 2>&1; then
  exit 0
fi

# Check if kubeconfig file exists in default locations
kubeconfig_path="${KUBECONFIG:-$HOME/.kube/config}"
if [ ! -f "$kubeconfig_path" ] && [ ! -d "$HOME/.kube" ]; then
  exit 0
fi

# Get current context
context=$(kubectl config current-context 2>/dev/null || echo "")

# Exit if no context active
if [ -z "$context" ]; then
  exit 0
fi

# Get current namespace
namespace=$(kubectl config view --minify --output 'jsonpath={..namespace}' 2>/dev/null || echo "default")
[ -z "$namespace" ] && namespace="default"

# Get API Server
server=$(kubectl config view -o jsonpath='{.clusters[?(@.name=="'"$context"'")].cluster.server}' 2>/dev/null || echo "n/a")

# ANSI colors
ESC=$(printf '\033')
CYAN="${ESC}[01;36m"
RESTORE="${ESC}[0m"

echo "K8s: ${CYAN}󱏚 ${RESTORE} $context"
echo "Context: $context"
echo "Namespace: $namespace"
echo "API Server: $server"
echo "Config Path: $kubeconfig_path"
