## **Prerequisites**

- **GCP Project**: Ensure you have a GCP project created.
- **Billing**: Enabled billing for the project.
- **gcloud CLI**: Installed and authenticated with your GCP account.
- **Docker**: Installed on your local machine for image building.
- **PostgreSQL Client**: Installed for local database access.
- **IAM Permissions**: Ensure proper access to Cloud SQL and Secret Manager.

## **1. Set Environment Variables**

```bash
export PROJECT_ID="your-project-id"
export REGION="your-region"  # e.g., us-central1
export ZONE="your-zone"      # e.g., us-central1-a
export VPC_NETWORK="default" # Assuming default VPC
export DB_INSTANCE_NAME="trading-db-instance"
export DB_NAME="trading_db"
export DB_USER="trading_user"
export DB_PRIVATE_IP="<PRIVATE_IP_OF_CLOUD_SQL>"
export SECRET_ID="portfolio-db-password" # when running the service this will be retrieved form secret manager
export TOPIC_ID="your-topic-id" # when order is submitted
```

## **2. Create Compute Engine VMs**

We will create 2 VMs. `vm1` will run user management service and `vm2` will run trading, portfolio services.

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
    --service-account=vm1-service-account@$PROJECT_ID.iam.gserviceaccount.com \
    --scopes=https://www.googleapis.com/auth/cloud-platform

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
     --service-account=vm2-service-account@$PROJECT_ID.iam.gserviceaccount.com \
    --scopes=https://www.googleapis.com/auth/cloud-platform
```

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
gcloud sql users set-password postgres --instance=$DB_INSTANCE_NAME --password=$SECRET_ID
```

### Create Database and User
```bash
gcloud sql databases create $DB_NAME --instance=$DB_INSTANCE_NAME

gcloud sql users create $DB_USER --instance=$DB_INSTANCE_NAME --password=$SECRET_ID
```

## **5. Store Portfolio Service DB Password in Secret Manager**

```bash
gcloud secrets create $SECRET_ID --replication-policy="automatic"
gcloud secrets versions add $SECRET_ID --data-file=password.txt
```

Ensure the VM’s service account has access to Secret Manager:
```bash
gcloud projects add-iam-policy-binding $PROJECT_ID \
    --member="serviceAccount:vm2-service-account@$PROJECT_ID.iam.gserviceaccount.com" \
    --role="roles/secretmanager.secretAccessor"
gcloud projects add-iam-policy-binding $PROJECT_ID \
    --member="serviceAccount:vm1-service-account@$PROJECT_ID.iam.gserviceaccount.com" \
    --role="roles/secretmanager.secretAccessor"
```

## **6. Database Schema**

### **User Service Schema**
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE
);
```

**Go Model:**
```go
type User struct {
    ID       int64
    Username string
    Email    string
}
```

### **Portfolio Service Schema**
```sql
CREATE TABLE positions (
    user_id VARCHAR(100) NOT NULL,
    asset VARCHAR(50) NOT NULL,
    quantity INT NOT NULL,
    price FLOAT NOT NULL,
    PRIMARY KEY (user_id, asset)
);
```

**Go Model:**
```go
type Position struct {
    UserID   string  `json:"user_id"`
    Asset    string  `json:"asset"`
    Quantity int     `json:"quantity"`
    Price    float64 `json:"price"`
}
```

## **7. Running Containers Locally with Cloud SQL Access**

### **Set Up Cloud SQL Proxy**

To connect to Cloud SQL from your local machine, install and run Cloud SQL Proxy:

```bash
gcloud auth login
gcloud auth application-default login

wget https://dl.google.com/cloudsql/cloud_sql_proxy.linux.amd64 -O cloud_sql_proxy
chmod +x cloud_sql_proxy

./cloud_sql_proxy -instances=$PROJECT_ID:$REGION:$DB_INSTANCE_NAME=tcp:5432 &
```

### **Run Containers Locally**

```bash
docker run -d \
    --name user-service \
    -p 8080:8080 \
    -e DB_HOST=127.0.0.1 \
    -e DB_USER=$DB_USER \
    -e DB_PASSWORD=$SECRET_ID \
    -e DB_NAME=$DB_NAME \
    user-service


docker run -d \
    --name portfolio-service \
    -p 8082:8082 \
    -e DB_HOST=127.0.0.1 \
    -e DB_USER=$DB_USER \
    -e DB_PASSWORD=$SECRET_ID \
    -e DB_NAME=$DB_NAME \
    portfolio-service
```
