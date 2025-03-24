# Recording for Secret Controller Operator
![Asciinema Recording](https://asciinema.org/a/e588722f-c2f9-425f-bd57-8defdf5a0c1e.png)
# Kubernetes Secret Generator Controller

This project implements a Kubernetes custom controller that generates secrets based on a custom resource definition (CRD). The controller allows you to define `CustomSecret` objects, which can specify different types of secrets like basic authentication credentials and JWT tokens. The controller automatically generates and manages these secrets within the Kubernetes cluster and can rotate them based on a specified period.

## Features
- **Custom Secrets**: Automatically generates secrets such as basic-auth credentials or JWT tokens.
- **Secret Rotation**: Periodically rotates secrets based on a defined rotation period.
- **RBAC Support**: Includes RBAC configuration for managing permissions.
- **Status Updates**: Updates the `CustomSecret` status with the name of the generated secret and the last updated timestamp.
- **Flexible Secret Types**: Support for multiple secret types, including basic-auth and JWT.

## Prerequisites
- Kubernetes cluster (minikube, EKS, GKE, or any other Kubernetes provider)
- Kubectl CLI configured to interact with the cluster
- Docker (for building images)
- Go (to build the controller)
- Controller-Manager (optional, if deploying to an existing cluster)

## Installation and Deployment

### 1. Build and Push the Controller Image

First, you need to build and push the controller's Docker image to a container registry.

```bash
# Build the image
make docker-build

# Tag the image for your registry
docker tag <your_image> <your_registry>/<your_image>:<tag>

# Push the image to the registry
docker push <your_registry>/<your_image>:<tag>
```

Replace `<your_image>`, `<your_registry>`, and `<tag>` with your image name, registry, and desired tag.

### 2. Apply RBAC Configurations

The controller uses RBAC (Role-Based Access Control) for permissions. Apply the RBAC resources to your cluster:

```bash
kubectl apply -f config/rbac/role.yaml
kubectl apply -f config/rbac/role_binding.yaml
```

### 3. Update the Deployment Configuration

The `config/default/deployment.yaml` file contains the deployment configuration for the controller. Update the `image` field in the YAML to use the image you just pushed to your registry:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: customsecret-controller
  namespace: default
spec:
  replicas: 1
  template:
    metadata:
      labels:
        control-plane: controller-manager
    spec:
      containers:
        - name: customsecret-controller
          image: <your_registry>/<your_image>:<tag>  # Replace with your image URL
          imagePullPolicy: Always
          ports:
            - containerPort: 9443
```

### 4. Deploy the Controller

Apply the deployment to your Kubernetes cluster:

```bash
kubectl apply -f config/default/deployment.yaml
```

This will deploy the controller to your Kubernetes cluster.

### 5. Create a CustomSecret Resource

To use the controller, create a `CustomSecret` resource. Below is an example YAML manifest for a `CustomSecret` object:

```yaml
apiVersion: app.mydomain.com/v1
kind: CustomSecret
metadata:
  name: my-secret
  namespace: default
spec:
  SecretType: "basic-auth"
  Username: "admin"
  PasswordLength: 16
  RotationPeriod: "24h"
```

Apply the `CustomSecret` object to the cluster:

```bash
kubectl apply -f customsecret.yaml
```

### 6. Monitor the Controller Logs

You can monitor the controller's logs to ensure it is functioning correctly:

```bash
kubectl logs -f <pod_name> -n default
```

### 7. Verify Secret Generation

After applying the `CustomSecret`, the controller will generate the secret and update the status of the `CustomSecret` resource. You can check the status of the `CustomSecret`:

```bash
kubectl get customsecrets -n default
```

### 8. Debugging

If the controller is not behaving as expected, check the logs of the controller pod:

```bash
kubectl logs -f <pod_name> -n default
```

You can also use `kubectl describe` to get more detailed information about the `CustomSecret` resource:

```bash
kubectl describe customsecret my-secret -n default
```

## How it Works

1. **CustomSecret Object**: You define a `CustomSecret` object that specifies the type of secret to generate (e.g., `basic-auth` or `jwt`), along with other details like the secret's username and password length.
   
2. **Secret Generation**: The controller listens for changes to the `CustomSecret` resource and generates the requested secret. It currently supports:
   - Basic Authentication secrets (`username` and `password`)
   - JWT secrets (a dummy JWT token for now)
   
3. **Secret Rotation**: If a rotation period is specified (e.g., `"24h"`), the controller will requeue the reconciliation loop to regenerate and update the secret after the specified duration.

4. **Status Update**: After generating the secret, the controller updates the `CustomSecret` status with the name of the generated secret and the last update timestamp.

## Contributing

1. Fork the repository.
2. Create a new branch.
3. Implement your changes.
4. Run tests and check your code.
5. Create a pull request to the `main` branch.

## License

This project is licensed under the Apache License, Version 2.0 – see the [LICENSE](LICENSE) file for details.
