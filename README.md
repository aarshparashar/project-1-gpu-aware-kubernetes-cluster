# GPU-Aware Kubernetes Cluster Simulator

This project is a modular platform that simulates a GPU-aware Kubernetes cluster. It is designed to demonstrate DevOps/SRE principles, specifically around infrastructure, scheduling, and observability, without requiring physical GPU hardware.

## Architecture

The system is composed of several key components:

1.  **Control Plane**: A REST API service that accepts workload submissions, handles namespaces/tenants, and tracks cluster state (nodes, pending/running jobs).
2.  **GPU-Aware Scheduler**: Determines optimal placement for jobs based on simulated GPU availability, memory, priority, and queue times.
3.  **Worker Node Simulator**: Simulates worker nodes with specific GPU capabilities (e.g., GPU type, memory) and utilization metrics.
4.  **Observability Stack**: Integrated Prometheus metrics, structured logs, and Grafana dashboards for deep visibility into queue depth, scheduling latency, and node utilization.

## Scheduling Model

The scheduling model mimics real-world constraints by considering:
*   **Resource Requirements:** Memory and compute profile.
*   **GPU Labels:** Specific simulated capabilities per node.
*   **Priority & Fair-Share:** Multi-tenant quotas and admission control to prevent noisy-neighbor issues.

GPU awareness is simulated via node metadata and labels rather than actual CUDA/NVIDIA driver interactions, keeping the platform lightweight and easy to run locally.

## Getting Started (Local Execution)

*(Instructions for Docker Compose/Kind will be added here)*

## Observability

*(Instructions for accessing Grafana dashboards and Prometheus metrics will be added here)*
