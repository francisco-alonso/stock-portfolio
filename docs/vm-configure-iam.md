# Configuring GCP VM Instance to Pull and Run Images from gcr.io

This guide explains how to set up a Google Cloud Platform (GCP) VM instance to pull container images from `gcr.io` and run them using Docker. 

---

## Prerequisites
1. A GCP project with the required container images hosted in `gcr.io`.
2. A VM instance with Compute Engine enabled.
3. `gcloud` CLI installed and authenticated.

---

## Steps

### 1. Assign IAM Permissions to the Service Account
The VM instance's service account needs the `Storage Object Viewer` role to access container images from `gcr.io`.

#### Command:
```bash
gcloud projects add-iam-policy-binding [PROJECT-ID] \
  --member="serviceAccount:[SERVICE-ACCOUNT-NAME]@[PROJECT-ID].iam.gserviceaccount.com" \
  --role="roles/storage.objectViewer"
```
### 2. Verify or Update the Service Account Attached to the VM

Ensure that your VM instance is using the correct service account.

Steps:
1. Navigate to Compute Engine > VM instances in the GCP Console.   
2. Locate your VM instance and verify the attached service account.
3. If the service account needs to be updated: 
    Stop the VM instance.
    Edit the VM settings and update the service account.
    Restart the VM.

### 3. Authenticate Docker
Configure Docker on the VM to use `gcr.io` as an authenticated registry.

```bash
gcloud auth configure-docker
```

### 4. Pull and Run the Container Image

```bash
docker pull gcr.io/[PROJECT-ID]/[IMAGE-NAME]:[TAG]
docker run [OPTIONS] gcr.io/[PROJECT-ID]/[IMAGE-NAME]:[TAG]
```

### Example

```bash
# Assign IAM permissions
gcloud projects add-iam-policy-binding my-project \
  --member="serviceAccount:my-vm-service-account@my-project.iam.gserviceaccount.com" \
  --role="roles/storage.objectViewer"

# Authenticate Docker
gcloud auth configure-docker

# Pull and run the image
docker pull gcr.io/my-project/my-container:latest
docker run -d -p 8080:80 gcr.io/my-project/my-container:latest
```