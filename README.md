# Blind GPU Scheduler

A Kubernetes scheduler extender that enforces zero-trust placement of sensitive GPU workloads using node attestation evidence instead of mutable labels.

- Trusts cryptographic attestation (via SPIRE) rather than node labels.
- Pods declare a required attestation hash using an annotation.
- The extender filters candidate nodes to those with an exact attestation hash match.

## Status

MVP extender with mock SPIRE client (hash map) and Kubernetes manifests for a local demo. SPIRE integration wiring and e2e tests forthcoming.

## Repository Structure

- `cmd/extender/` — main for the scheduler extender HTTP server
- `pkg/extender/` — extender logic and types
- `pkg/spire/` — SPIRE client interface and mock implementation
- `pkg/logging/` — structured logging setup
- `deploy/` — Kubernetes manifests (extender Deployment/Service, scheduler config, RBAC)
- `examples/` — example Pod with required attestation hash
- `docs/` — architecture and threat model

## How It Works

1. Users add an annotation to the Pod:
   - `attestation-hash.my-company.com/required-hash: "<expected-hash>"`
2. The default scheduler calls this extender with candidate nodes.
3. The extender queries SPIRE (mocked initially) to obtain each node's attestation hash.
4. Only nodes with a matching hash are returned; scheduling is constrained to that set.

## Quickstart (Local Demo)

Prereqs:
- Docker
- Go 1.21+
- kubectl
- kind or minikube

### Build & Run Locally

```bash
make build
./bin/extender
```

Environment variables:
- `BGS_LISTEN_ADDR` (default `:8000`)
- `BGS_REQUIRED_ANNOTATION` (default `attestation-hash.my-company.com/required-hash`)
- `BGS_NODE_HASH_MAP` — JSON map of nodeName->attestationHash, e.g. `{"kind-worker":"abc123"}`

### Container Image

```bash
make docker-build IMG=ghcr.io/your-org/blind-gpu-scheduler:dev
```

### Deploy Extender to Cluster

```bash
# Create namespace
kubectl create ns blind-gpu-scheduler || true

# Deploy RBAC, Service, Deployment
kubectl apply -f deploy/rbac.yaml
kubectl apply -f deploy/extender-deployment.yaml
kubectl apply -f deploy/extender-service.yaml
```

Configure kube-scheduler with the extender (varies by environment). See `deploy/scheduler-config.yaml` and `docs/ARCHITECTURE.md` for guidance.

### Example Pod

```bash
kubectl apply -f examples/secure-gpu-pod.yaml
```

If no node's attestation hash matches the required hash, the Pod remains Pending.

## Roadmap

- Replace mock SPIRE client with real SPIRE server queries
- Add freshness validation and anti-replay (nonce, timestamp)
- Observability (Prometheus metrics)
- e2e tests with kind

## Security Notes

- This extender is designed to be non-ignorable by the scheduler; if no nodes match, the Pod is unscheduled, not silently downgraded.
- Treat the SPIRE Server as the trust anchor; harden access and audit logs.
