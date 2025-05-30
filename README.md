# Exercise 02 - Containerization Technology

This project demonstrates comprehensive containerization practices using Docker, including multi-stage builds, CI/CD integration with GitLab, and reverse proxy configurations. The setup containerizes the Beetroot API using Infrastructure as Code principles with Docker Compose orchestration and automated deployment pipelines.

## Table of Contents

- [Overview](#overview)
- [Prerequisites](#prerequisites)
- [Project Structure](#project-structure)
- [Quick Start](#quick-start)
- [Components](#components)
- [Multi-Stage Dockerfile Architecture](#multi-stage-dockerfile-architecture)
- [Reverse Proxy Configurations](#reverse-proxy-configurations)
- [GitLab CI/CD Integration](#gitlab-cicd-integration)
- [Configuration](#configuration)
- [Troubleshooting](#troubleshooting)
- [Daily Operations](#daily-operations)
- [Technical Details](#technical-details)

## Overview

This project automates the containerization and deployment of the Beetroot API with the following architecture:
- **Multi-Stage Docker Build**: Optimized Go compilation with minimal runtime image
- **Container Registry**: GitLab-integrated automated image building and storage
- **Reverse Proxy Solutions**: Traefik, Nginx, and Apache configurations
- **Network Isolation**: API accessible only through reverse proxies
- **CI/CD Automation**: GitLab pipelines for automated builds and deployments
- **Environment Management**: .env file configurations for different deployment scenarios

## Prerequisites

Before starting, ensure you have the following installed and configured:

- **Docker Engine** (20.10 or later)
- **Docker Compose** (v2.0 or later)
- **GitLab Account**: With container registry access
- **Git**: For repository management
- **MTU Configuration**: Docker daemon configured with MTU 1442

### Required System Configuration

Ensure your Docker daemon configuration includes:
```json
{
  "mtu": 1442
}
```

## Project Structure

```
exercise-02-containerization-technology/
├── README.md --------------------------------> # This file
├── Dockerfile -------------------------------> # Multi-stage Beetroot container build
├── .gitlab-ci.yml ---------------------------> # CI/CD pipeline configuration
├── beetroot/ --------------------------------> # Base Beetroot application setup
│   ├── docker-compose.yml -------------------> # Basic container orchestration
│   ├── .env ---------------------------------> # Environment variables
│   └── .env.example -------------------------> # Environment template
├── data/ ------------------------------------> # Application data
│   └── beetroot.json ------------------------> # API data file
├── traefik/ ---------------------------------> # Traefik reverse proxy stack
│   ├── docker-compose.yml -------------------> # Traefik orchestration
│   ├── traefik.yml --------------------------> # Traefik static configuration
│   ├── config.yml ---------------------------> # Traefik dynamic configuration
│   ├── .env ---------------------------------> # Traefik environment variables
│   └── .env.example -------------------------> # Traefik environment template
├── nginx/ -----------------------------------> # Nginx reverse proxy stack
│   ├── docker-compose.yml -------------------> # Nginx orchestration
│   ├── nginx.conf ---------------------------> # Nginx proxy configuration
│   ├── .env ---------------------------------> # Nginx environment variables
│   └── .env.example -------------------------> # Nginx environment template
└── apache/ ----------------------------------> # Apache reverse proxy stack
    ├── docker-compose.yml -------------------> # Apache orchestration
    ├── apache.conf --------------------------> # Apache proxy configuration
    ├── .env ---------------------------------> # Apache environment variables
    └── .env.example -------------------------> # Apache environment template
```

## Quick Start

### Automated Deployment

1. **Clone the repository:**
   ```bash
   git clone <repository-url>
   cd exercise-02-containerization-technology
   ```

2. **Set up environment variables:**
   ```bash
   # Copy example environment files
   cp beetroot/.env.example beetroot/.env
   cp traefik/.env.example traefik/.env
   cp nginx/.env.example nginx/.env
   cp apache/.env.example apache/.env
   
   # Edit .env files with your configurations
   nano beetroot/.env
   ```

3. **Build and deploy with Traefik (recommended):**
   ```bash
   cd traefik/
   docker compose --env-file .env up -d
   ```

4. **Access the API:**
   - Visit: `http://localhost`
   - API endpoint: `http://localhost/api`

### Alternative Proxy Deployments

**Nginx deployment:**
```bash
cd nginx/
docker compose --env-file .env up -d
```

**Apache deployment:**
```bash
cd apache/
docker compose --env-file .env up -d
```

## Components

### Container Components

| Component | Description | Image | Port |
|-----------|-------------|-------|------|
| **Beetroot API** | Go-based REST API application | Custom built (Alpine-based) | 8080 |
| **Traefik** | Modern reverse proxy and load balancer | traefik:v2.10 | 80 |
| **Nginx** | High-performance web server and reverse proxy | nginx:1.25 | 80 |
| **Apache** | Flexible web server with reverse proxy | httpd:2.4 | 80 |

### Network Architecture

- **External Access**: `Internet` → **Reverse Proxy:80** → **Beetroot API:8080**
- **Container Isolation**: API not directly accessible from external network
- **Service Discovery**: Container-to-container communication via Docker networking

## Multi-Stage Dockerfile Architecture

The Dockerfile implements a two-stage build process for optimal security and size:

### Stage 1: Builder (golang:1.21-alpine)
- **Purpose**: Compile Go source code into static binary
- **Components**: Git, Go toolchain, tzdata for timezone support
- **Process**: Clone repository → Install dependencies → Build binary

### Stage 2: Runtime (alpine:3.18)
- **Purpose**: Minimal runtime environment
- **Components**: Only tzdata (temporary), compiled binary
- **Security**: No build tools, minimal attack surface
- **Size Optimization**: ~19.7MB final image

### Build Flow

<img width="600" alt="Multi-stage build flow for creating the Beetroot API container image" src="https://github.com/user-attachments/assets/d27d44d4-4cc9-44a5-87ed-a1c28b77d290" />

*Figure 1: Multi-stage build flow for creating the Beetroot API container image.*

### Network Topology

<img width="500" alt="Topology of Beetroot API deployment using a reverse proxy" src="https://github.com/user-attachments/assets/ba0cf006-49c9-428d-849d-f7b796f28d57" />

*Figure 2: Topology of Beetroot API deployment using a reverse proxy (Traefik, Nginx, or Apache).*

## Reverse Proxy Configurations

### Traefik Configuration

**Modern Features:**
- Dynamic service discovery
- Automatic HTTPS (when configured)
- Dashboard and API monitoring
- File-based configuration with hot reload

**Key Files:**
- `traefik.yml`: Static configuration (entry points, providers)
- `config.yml`: Dynamic routing rules
- Labels-based service discovery in docker-compose

### Nginx Configuration

**Traditional Approach:**
- High-performance HTTP server
- Robust reverse proxy capabilities
- Custom `nginx.conf` configuration
- Upstream load balancing support

**Configuration Highlights:**
```nginx
location / {
    proxy_pass http://beetroot:8080;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
}
```

### Apache Configuration

**Enterprise Features:**
- Module-based architecture
- Flexible virtual host configuration
- ProxyPass/ProxyPassReverse directives
- Extensive logging capabilities

**Configuration Highlights:**
```apache
ProxyPass "/" "http://beetroot:8080/"
ProxyPassReverse "/" "http://beetroot:8080/"
ProxyPreserveHost On
```

## GitLab CI/CD Integration

### Pipeline Architecture

The CI/CD pipeline uses **Kaniko** for secure, unprivileged container building:

```yaml
build:
  stage: build
  image:
    name: gcr.io/kaniko-project/executor:v1.23.2-debug
    entrypoint: [""]
  script:
    - /kaniko/executor \
      --context "${CI_PROJECT_DIR}" \
      --dockerfile "${CI_PROJECT_DIR}/Dockerfile" \
      --destination "${CI_REGISTRY_IMAGE}:${CI_COMMIT_TAG}"
  rules:
    - changes:
      - Dockerfile
```

### Automation Features

- **Trigger Conditions**: Automatic builds on Dockerfile changes
- **Secure Building**: No privileged Docker daemon required
- **Registry Integration**: Direct push to GitLab Container Registry
- **Version Management**: Git tag-based image versioning

## Configuration

### Environment Variables (.env files)

#### Beetroot Application (`beetroot/.env`)
```bash
REGISTRY_IMAGE=registry.gitlab.com/username/project
TAG=latest
BEETROOT_JSON_PATH=/data/beetroot.json
TZ=Europe/Oslo
```

#### Traefik Configuration (`traefik/.env`)
```bash
REGISTRY_IMAGE=registry.gitlab.com/username/project
TAG=latest
TRAEFIK_API_DASHBOARD=true
TRAEFIK_LOG_LEVEL=DEBUG
```

#### Nginx Configuration (`nginx/.env`)
```bash
REGISTRY_IMAGE=registry.gitlab.com/username/project
TAG=latest
NGINX_HOST=localhost
NGINX_PORT=80
```

#### Apache Configuration (`apache/.env`)
```bash
REGISTRY_IMAGE=registry.gitlab.com/username/project
TAG=latest
APACHE_SERVER_NAME=localhost
APACHE_LOG_LEVEL=warn
```

### Data Configuration

The `data/beetroot.json` file contains the API data served by the application:
```json
{
  "status": "online",
  "message": "Beetroot API is running",
  "timestamp": "2025-01-01T00:00:00Z"
}
```

## Troubleshooting

### Common Issues

| Issue | Symptoms | Solution |
|-------|----------|----------|
| **Image tag case sensitivity** | GitLab push failures | Use lowercase tags only: `latest` not `Latest` |
| **Binary not found** | Container startup failures | Verify COPY paths in Dockerfile match build output |
| **Port conflicts** | "Port already in use" errors | Stop conflicting containers: `docker compose down` |
| **Environment variables** | Configuration not loading | Check .env file syntax and placement |
| **Network isolation** | Direct API access possible | Verify no published ports on beetroot service |
| **MTU configuration** | Network connectivity issues | Ensure Docker daemon MTU is set to 1442 |

### Diagnostic Commands

#### Container Diagnostics
```bash
# Check running containers
docker ps -a

# View container logs
docker logs <container-name>

# Inspect container configuration
docker inspect <container-name>

# Test API accessibility
curl -i http://localhost/
curl -i http://localhost:8080/  # Should fail (no direct access)
```

#### Image Diagnostics
```bash
# List built images
docker images

# Check image layers and size
docker history <image-name>

# Verify multi-stage build
docker build --target builder -t beetroot-builder .
docker run --rm beetroot-builder ls -la /app/Beetroot/cmd/beetroot/
```

#### Network Diagnostics
```bash
# Check Docker networks
docker network ls

# Inspect compose network
docker network inspect <compose-network>

# Test internal connectivity
docker exec <proxy-container> ping beetroot
```

## Daily Operations

### Development Workflow

1. **Local Development:**
   ```bash
   # Build image locally
   docker build -t beetroot-local .
   
   # Test with different proxies
   cd traefik/ && docker compose up -d
   cd ../nginx/ && docker compose up -d
   cd ../apache/ && docker compose up -d
   ```

2. **Testing Changes:**
   ```bash
   # Rebuild and restart services
   docker compose build
   docker compose up -d --force-recreate
   
   # View logs for debugging
   docker compose logs -f
   ```

3. **Production Deployment:**
   ```bash
   # Tag and push to registry
   git tag v1.0.0
   git push origin v1.0.0
   
   # Deploy with production image
   docker compose --env-file .env up -d
   ```

### Maintenance Tasks

#### Starting Services
```bash
# Start specific proxy stack
cd <proxy-directory>/
docker compose --env-file .env up -d

# Verify services are running
docker compose ps
curl -i http://localhost/
```

#### Stopping Services
```bash
# Stop and remove containers
docker compose down

# Stop and remove containers + volumes
docker compose down -v

# Remove unused images
docker image prune -f
```

#### Updating Configurations
```bash
# Reload proxy configurations (Traefik auto-reloads)
# Nginx/Apache require restart:
docker compose restart nginx
docker compose restart apache

# Update environment variables
nano .env
docker compose up -d --force-recreate
```

## Technical Details

### Security Considerations

- **Multi-stage builds**: Eliminate build tools from runtime image
- **Non-root execution**: Alpine base with minimal permissions
- **Network isolation**: API only accessible through reverse proxy
- **Version pinning**: Specific image tags prevent supply chain attacks
- **Secret management**: Environment variables for sensitive configuration

### Performance Optimization

- **Image size**: 19.7MB final image using Alpine Linux
- **Build caching**: Docker layer caching for faster rebuilds
- **Static binary**: Go compilation produces self-contained executable
- **Reverse proxy caching**: Optional caching headers for static content

### Container Orchestration

- **Service dependencies**: Proper container startup order
- **Health checks**: Built-in container health monitoring
- **Volume management**: Persistent data and configuration mounting
- **Network segmentation**: Isolated container communication

### CI/CD Best Practices

- **Kaniko builds**: Secure, unprivileged container building
- **Registry integration**: Automated image storage and versioning
- **Change detection**: Builds triggered only on relevant file changes
- **Tag management**: Git tag-based image versioning strategy

## Additional Notes

### Customization Options

- **Timezone configuration**: Modify `TZ` build argument for different regions
- **API data**: Update `data/beetroot.json` for custom API responses
- **Proxy settings**: Adjust reverse proxy configurations for specific requirements
- **Resource limits**: Add memory/CPU constraints in docker-compose files

### Extension Possibilities

- **HTTPS support**: Add SSL/TLS termination at proxy level
- **Monitoring**: Integrate Prometheus metrics and Grafana dashboards
- **Load balancing**: Scale API containers with proxy load balancing
- **Health checks**: Implement application health endpoints

### Production Considerations

This setup is designed for educational purposes. For production deployment:
- Implement proper SSL/TLS certificates
- Configure authentication and authorization
- Set up centralized logging and monitoring
- Implement backup and disaster recovery procedures
- Use container orchestration platforms (Kubernetes, Docker Swarm)

---

**Project**: Exercise 02 - Containerization Technology  
**Course**: IKT114 - IT Orchestration  
**Institution**: University of Agder  
**Authors**: Yusef Said & Eirik André Lindseth

## Version History

- **v1.0**: Initial multi-stage Docker implementation with reverse proxy configurations
- **v1.1**: Added GitLab CI/CD integration and automated container registry deployment
