# Distributed Task Scheduler

A cloud-native microservices backend for processing tasks like image resizing. Built with Go, Docker, and Kubernetes, this system is designed for scalability, reliability, and ease of deployment. It features a REST API for task submission, a task queue using NATS, worker pods for task processing, MinIO for result storage, and Prometheus/Grafana for monitoring.

## Features

- **Scalable Microservices**: Modular services including API, workers, NATS, MinIO, and monitoring components orchestrated with Kubernetes.
- **Auto-Scaling**: Horizontal Pod Autoscaler (HPA) dynamically scales worker pods (2 to 5 replicas) based on CPU utilization thresholds.
- **Task Processing**: Efficient image resizing tasks with results securely stored in MinIO.
- **Monitoring**: Prometheus and Grafana track task throughput, system metrics, and overall health.
- **Cloud-Ready**: Engineered for deployment on managed Kubernetes platforms like Google Kubernetes Engine (GKE).

## Tech Stack

- **Backend**: Go
- **Containerization**: Docker
- **Orchestration**: Kubernetes (Minikube)
- **Queue**: NATS
- **Storage**: MinIO
- **Monitoring**: Prometheus, Grafana

## Current Status

- **Phase 4 Complete**: Successfully deployed on Kubernetes with HPA. The system processed 500+ tasks, storing 637 processed objects in MinIO, and demonstrated stability on a 16GB RAM setup.

## Deployment Instructions

### Prerequisites

- **Minikube** installed locally.
- Minimum system requirements: 6 CPUs, 12GB RAM.

### Setup

1. **Clone the Repository**:
   ```bash
   git clone https://github.com/Nexusrex18/distributed-task-scheduler.git
   cd distributed-task-scheduler
