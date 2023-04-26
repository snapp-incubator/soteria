#!/usr/bin/env bash

set -e

company="$(oc project -q | cut -d- -f1)"
ode="$(oc project -q | cut -d- -f3)"

echo "Rolling out soteria on $company-ode-$ode"

current_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
path_to_k8s="$current_dir/soteria"

helm upgrade --install soteria "$path_to_k8s" \
	-f "$path_to_k8s/values.yaml" \
	-f "$path_to_k8s/values.ode.$company.yaml"
