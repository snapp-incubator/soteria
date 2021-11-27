#!/usr/bin/env bash

set -e

echo "Rolling out soteria..."

current_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
path_to_k8s="$current_dir/soteria"

helm upgrade --install soteria "$path_to_k8s" \
	-f "$path_to_k8s/values.yaml" \
	-f "$path_to_k8s/values.ode.yaml"
