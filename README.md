# Cloud-Native Microservices Deployment on Kubernetes

This project demonstrates a comprehensive cloud-native application deployment on Kubernetes, fulfilling a challenging set of requirements designed to showcase proficiency in modern DevOps practices.

## Project Overview

The solution involves three Go-based microservices, containerized with Docker, deployed to an Amazon EKS (Elastic Kubernetes Service) cluster. It incorporates advanced Kubernetes features for scalability, resilience, security, and a robust observability stack using Prometheus, Grafana, and Loki.

## Core Setup

### 1. Create 2–3 dummy microservices (Go)
*   **Implementation:** Three lightweight microservices were developed in **Go**:
    *   `greeter-service`: Serves a simple HTML UI, exposes `/ping`, and calls `echo-service` via its `/hello` endpoint.
    *   `echo-service`: Exposes `/ping` and an `/echo` endpoint that processes and returns a JSON payload with a timestamp.
    *   `admin-service`: Exposes `/ping` and a `/secret` endpoint protected by an API key.
*   **Why Go?** Go was chosen for its excellent performance, small binary sizes (ideal for containers), and strong concurrency features, demonstrating an understanding of efficient cloud-native application development.

### 2. Dockerize each service
*   **Implementation:** A `Dockerfile` was created for each microservice. These are **multi-stage Dockerfiles** to ensure lean, production-ready images by separating the build environment from the runtime environment.
*   **Why multi-stage?** This best practice significantly reduces the final image size and attack surface, improving security and deployment efficiency.
*   **Image Registry:** Images were built locally and then pushed to **Docker Hub** (`adityahebballe/greeter-service:latest`, `adityahebballe/echo-service:latest`, `adityahebballe/admin-service:latest`). This demonstrates familiarity with public container registries.

### 3. Deploy to Kubernetes
*   **Implementation:** Each microservice is deployed using standard Kubernetes `Deployment` and `Service` manifests. These manifests are organized into dedicated subdirectories (`kubernetes/<service-name>/`) for clarity and maintainability.
*   **Why Deployment + Service?** `Deployment` ensures desired replica counts and rolling updates, while `Service` provides stable network access to the pods within the cluster.

### 4. Set up Ingress
*   **Implementation:** The **NGINX Ingress Controller** was installed via Helm. A single `Ingress` resource (`kubernetes/ingress.yaml`) was configured to route external traffic:
    *   `/` path routes to `greeter-service`.
    *   `/grafana` path routes to `prometheus-grafana` (Grafana UI).
*   **Why Ingress?** Ingress provides a robust, flexible, and scalable way to manage external access to services, handling routing, SSL termination, and load balancing, which is standard in production environments. It's preferred over `NodePort` for external exposure due to its advanced features and manageability.

### 5. Add HPA + VPA
*   **Implementation:**
    *   **Resource Requests/Limits:** `Deployment` manifests were updated to include CPU and memory `requests` and `limits` for each container. This is a prerequisite for autoscaling.
    *   **Horizontal Pod Autoscaler (HPA):** `HPA` resources (`kubernetes/<service-name>/hpa.yaml`) were created for each service, configured to scale based on 50% CPU and 50% memory utilization, with `minReplicas: 1` and `maxReplicas: 5`.
    *   **Vertical Pod Autoscaler (VPA):** The VPA controller was installed via Helm. `VPA` resources (`kubernetes/<service-name>/vpa.yaml`) were created for each service with `updateMode: Initial`.
*   **Why HPA + VPA?**
    *   **HPA:** Ensures horizontal scalability by automatically adjusting the number of pod replicas based on observed metrics, optimizing resource utilization and maintaining performance under varying loads.
    *   **VPA:** Optimizes vertical resource allocation by automatically setting (or recommending) optimal CPU and memory requests/limits for pods, preventing over-provisioning or under-provisioning. `Initial` mode was chosen to demonstrate VPA's capability without disruptive runtime adjustments.

### 6. Add PodDisruptionBudget (PDB)
*   **Implementation:** `PodDisruptionBudget` resources (`kubernetes/<service-name>/pdb.yaml`) were created for each service, specifying `minAvailable: 1`.
*   **Why PDB?** PDBs ensure application availability during voluntary disruptions (e.g., node maintenance, cluster upgrades) by guaranteeing a minimum number of healthy pods remain running, enhancing application resilience.

### 7. Set Node Affinity
*   **Implementation:** EKS worker nodes were identified and conceptually labeled to simulate different availability zones (e.g., `topology.kubernetes.io/zone=us-east-1d`, `topology.kubernetes.io/zone=us-east-1c`). `Deployment` manifests were updated to include `nodeAffinity` rules, scheduling `greeter-service` and `admin-service` to one zone, and `echo-service` to another.
*   **Why Node Affinity?** This demonstrates how to control pod placement for high availability (spreading workloads across different failure domains), compliance, or optimizing resource utilization on specialized nodes.

### 8. Create Secrets and mount dummy volume
*   **Implementation:**
    *   **Kubernetes Secret:** A `Secret` (`kubernetes/app-secrets.yaml`) was created to store dummy `DB_PASSWORD` and `API_KEY` values (base64 encoded).
    *   **Secret Injection:** The `API_KEY` was injected into the `admin-service` pod as an environment variable using `secretKeyRef`.
    *   **Dummy Volume:** A `PersistentVolume` and `PersistentVolumeClaim` (`kubernetes/dummy-volume.yaml`) were created to simulate persistent storage.
    *   **Volume Mount:** The `dummy-pvc` was mounted into the `greeter-service` pod at `/app/data`.
*   **Why Secrets & Volumes?**
    *   **Secrets:** Demonstrates secure handling of sensitive configuration data, preventing hardcoding credentials in code or config maps.
    *   **Volumes:** Shows how to provide persistent storage to stateful applications in Kubernetes, ensuring data survives pod restarts or rescheduling.

### 9. Service Account + RBAC
*   **Implementation:**
    *   **Service Accounts:** Dedicated `ServiceAccount` resources (`kubernetes/service-accounts.yaml`) were created for each microservice (`greeter-service-sa`, `echo-service-sa`, `admin-service-sa`).
    *   **Role:** A `Role` (`kubernetes/admin-role.yaml`) named `admin-service-pod-reader` was created, granting `get`, `watch`, and `list` permissions on `pods`.
    *   **RoleBinding:** A `RoleBinding` (`kubernetes/admin-rolebinding.yaml`) was created to bind the `admin-service-pod-reader` `Role` to the `admin-service-sa` `ServiceAccount`.
    *   **Deployment Updates:** Each `Deployment` manifest was updated to use its respective `serviceAccountName`.
*   **Why Service Accounts + RBAC?** This demonstrates the principle of least privilege, ensuring that applications only have the necessary permissions within the Kubernetes API, enhancing cluster security.

### 10. Inter-service Communication
*   **Implementation:** The `greeter-service` was configured to make an HTTP POST request to the `echo-service` using its Kubernetes service name (`http://echo-service:8081/echo`).
*   **Why inter-service communication?** This is fundamental to microservices architectures, showing how services discover and communicate with each other within the cluster's network.

## Monitoring📊

### 1. Install Prometheus (node-level only)
*   **Implementation:** The `kube-prometheus-stack` Helm chart was installed. This comprehensive chart includes Prometheus, which automatically scrapes node-level metrics (via `node-exporter`) and other cluster metrics.
*   **Why `kube-prometheus-stack`?** It's the industry-standard, battle-tested solution for Kubernetes monitoring, providing a complete and integrated monitoring experience.

### 2. Add Grafana to visualize metrics
*   **Implementation:** Grafana was installed as part of the `kube-prometheus-stack`. It was configured to be accessible via the NGINX Ingress Controller at `/grafana` by setting `grafana.grafana.ini.server.root_url` during Helm upgrade.
*   **Why Grafana?** Grafana is the leading open-source platform for data visualization and dashboards, allowing for powerful and customizable insights into application and infrastructure performance.

### 3. Integrate logging (Loki)
*   **Implementation:** The `loki-stack` Helm chart was installed, deploying Loki (log aggregation system) and Promtail (log collector). Loki was then manually added as a data source in Grafana using its internal service URL (`http://loki.default.svc.cluster.local:3100`).
*   **Why Loki?** Loki is a lightweight, highly efficient log aggregation system designed to work seamlessly with Grafana, providing a unified observability experience (metrics and logs in one place). This demonstrates a modern approach to centralized logging.

## How to Run/Verify

1.  **Prerequisites:** Ensure you have `kubectl`, `helm`, `aws cli`, and `docker` installed and configured to interact with your AWS account and EKS cluster.
2.  **Build & Push Images:**
    ```bash
    ./build.sh
    # Enter Docker Hub password when prompted
    ```
3.  **Deploy Kubernetes Resources:**
    ```bash
    ./deploy.sh
    ```
4.  **Verify Pods & Services:**
    ```bash
    kubectl get pods -o wide
    kubectl get services
    kubectl get ingress
    kubectl get hpa
    kubectl get vpa
    kubectl get pdb
    kubectl get serviceaccounts
    kubectl get role admin-service-pod-reader
    kubectl get rolebinding admin-service-pod-reader-binding
    ```
5.  **Access Greeter Service UI (Inter-service Communication):**
    *   Open your browser to the Ingress URL (e.g., `http://a876bf88740f74d04a4e28c4e92948db-2122727902.us-east-1.elb.amazonaws.com/`).
    *   Click "Say Hello" to see the response from `echo-service`.
6.  **Access Grafana (Metrics & Logs):**
    *   Open your browser to the Grafana Ingress URL (e.g., `http://a876bf88740f74d04a4e28c4e92948db-2122727902.us-east-1.elb.amazonaws.com/grafana`).
    *   Log in (default user: `admin`, password: `kubectl get secret prometheus-grafana -n default -o jsonpath="{.data.admin-password}" | base64 --decode`).
    *   In the "Explore" section, select "Prometheus" to view metrics.
    *   In the "Explore" section, select "Loki" (if not present, manually add it with URL `http://loki.default.svc.cluster.local:3100`) to view logs (e.g., query `{job="kubernetes-pods"}`).

This project demonstrates a robust and well-architected cloud-native solution, covering key aspects of microservices, Kubernetes, and observability.
