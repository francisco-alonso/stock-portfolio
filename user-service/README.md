# Project Setup Guide for GCP Microservices System

This README provides a detailed, step-by-step guide for setting up a microservices system using Google Cloud Platform (GCP), with services deployed on Compute Engine VMs, PostgreSQL on Cloud SQL, and Docker for service management.

---

## **Prerequisites**

- **GCP Project**: Ensure you have a GCP project created.
- **Billing**: Enabled billing for the project.
- **gcloud CLI**: Installed and authenticated with your GCP account.
- **Docker**: Installed on your local machine for image building.

---

## **1. Set Environment Variables**

```bash
export PROJECT_ID="your-project-id"
export REGION="your-region"  # e.g., us-central1
export ZONE="your-zone"      # e.g., us-central1-a
export VPC_NETWORK="default" # Assuming default VPC
export DB_INSTANCE_NAME="trading-db-instance"
export DB_NAME="trading_db"
export DB_USER="trading_user"
export DB_PASSWORD="your-db-password"
export DB_PRIVATE_IP="<PRIVATE_IP_OF_CLOUD_SQL>"
```

---

## **2. Create Compute Engine VMs**

```bash
gcloud compute instances create vm1 \
    --project=$PROJECT_ID \
    --zone=$ZONE \
    --machine-type=e2-medium \
    --image-family=debian-11 \
    --image-project=debian-cloud \
    --network=$VPC_NETWORK \
    --subnet=default \
    --tags=http-server,https-server \
    --boot-disk-size=20GB

gcloud compute instances create vm2 \
    --project=$PROJECT_ID \
    --zone=$ZONE \
    --machine-type=e2-medium \
    --image-family=debian-11 \
    --image-project=debian-cloud \
    --network=$VPC_NETWORK \
    --subnet=default \
    --tags=http-server,https-server \
    --boot-disk-size=20GB
```

---

## **3. Install Docker on Each VM**

SSH into each VM and run the following commands:

```bash
sudo apt update
sudo apt install -y ca-certificates curl gnupg

# Add Docker’s official GPG key
sudo install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/debian/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg

# Set up the Docker repository
echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/debian $(lsb_release -cs) stable" | \
sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

# Install Docker Engine
sudo apt update
sudo apt install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

# Start and enable Docker
sudo systemctl start docker
sudo systemctl enable docker

# Verify Docker installation
docker --version
```

---

## **4. Set Up Cloud SQL (PostgreSQL)**

### Create Cloud SQL Instance
```bash
gcloud sql instances create $DB_INSTANCE_NAME \
    --database-version=POSTGRES_14 \
    --cpu=2 \
    --memory=4GB \
    --region=$REGION \
    --network=$VPC_NETWORK \
    --no-assign-ip \
    --enable-ipv4
```

### Set Root Password
```bash
gcloud sql users set-password postgres --instance=$DB_INSTANCE_NAME --password=$DB_PASSWORD
```

### Create Database and User
```bash
gcloud sql databases create $DB_NAME --instance=$DB_INSTANCE_NAME

gcloud sql users create $DB_USER --instance=$DB_INSTANCE_NAME --password=$DB_PASSWORD
```

---

## **5. Configure Access to Cloud SQL**

- Ensure that the Cloud SQL instance is associated with the correct VPC.
- Verify that the VMs can connect using the private IP.
- Test the connection from VM:

```bash
psql "host=$DB_PRIVATE_IP user=$DB_USER password=$DB_PASSWORD dbname=$DB_NAME sslmode=disable"
```

---

## **6. Create Initial Database Tables**

Inside the `psql` session:

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert test data
INSERT INTO users (username, email) VALUES ('alice', 'alice@example.com');

-- Verify the insertion
SELECT * FROM users;
```

---

## **7. Build and Run User Service in Docker**

### .env File Example
```env
DB_HOST=$DB_PRIVATE_IP
DB_USER=$DB_USER
DB_PASSWORD=$DB_PASSWORD
DB_NAME=$DB_NAME
```

### Run Docker Container
```bash
docker build -t user-service .

docker run -d \
    --name user-service \
    --env-file .env \
    -p 8080:8080 \
    user-service
```

### Test the Create-User Endpoint
```bash
curl -X POST http://localhost:8080/create-user -d '{"username":"bob","email":"bob@example.com"}' -H "Content-Type: application/json"
```

---

## **8. Debugging Common Issues**
- Ensure firewall rules allow traffic for required ports (like 8080).
- Verify Docker logs using `docker logs user-service`.
- Use `gcloud sql instances describe $DB_INSTANCE_NAME` to confirm correct IP configurations.

---

This guide ensures every step of the setup is well-documented, from VM creation to service deployment with secure environment variable management. For additional steps or issues, update the guide accordingly.

