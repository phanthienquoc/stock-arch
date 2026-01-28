# CI/CD Pipeline Flow - Trading System

## 🔄 Workflow Architecture

```
Developer
    ↓
git push v1.x.x tag
    ↓
[1] Build & Push Images (Automatic)
    ├─ Builds all 5 services
    ├─ Pushes to Docker Hub
    └─ Collects image digests
    ↓
[2] Auto Deploy to Staging (Automatic)
    ├─ Triggered by Build completion
    ├─ Pulls latest images
    ├─ Deploys to staging VPS
    └─ Notifications sent
    ↓
QA Testing / Verification
    ↓
[3] Promote to Production (Manual - workflow_dispatch)
    ├─ Developer triggers manually
    ├─ Pulls production images
    ├─ Deploys to production VPS
    └─ Production notification
```

---

## 📋 Workflow Files

### 1. **build.yml** - Build & Push Images
**Trigger**: Push tag matching `v*.*.*` (e.g., `v1.1.1`)

**Steps**:
- ✅ Extract version from tag
- ✅ Docker login to Docker Hub
- ✅ Build each service (matrix strategy)
- ✅ Push images with version tag
- ✅ Extract and collect digests
- ✅ Upload digest artifacts

**Outputs**: Digest artifacts for each service

**Example**:
```bash
git tag v1.1.1
git push origin v1.1.1
```

---

### 2. **deploy.yml** - Auto Deploy to Staging
**Trigger**: Automatic when `build.yml` completes successfully

**Steps**:
- ✅ Download digest artifacts from build workflow
- ✅ Parse digests
- ✅ SSH to staging VPS
- ✅ Pull latest Docker images
- ✅ Restart all services
- ✅ Send success notification

**Deployment target**: Staging VPS (`app/trading`)

**Automatic**: No manual intervention needed

---

### 3. **promote-to-prod.yml** - Manual Production Deployment
**Trigger**: Manual workflow dispatch (`workflow_dispatch`)

**Steps**:
- ✅ User provides version tag input
- ✅ Checkout the specific release tag
- ✅ Verify Docker images exist
- ✅ SSH to production VPS
- ✅ Pull production Docker images
- ✅ Restart services in production
- ✅ Send deployment summary

**Deployment target**: Production VPS (`app/trading`)

**Activation**:
1. Go to GitHub Actions
2. Click "Promote to Production" workflow
3. Click "Run workflow"
4. Enter version tag (e.g., `v1.1.1`)
5. Click "Run workflow"

---

## 🚀 Complete Example Flow

### Step 1: Developer Pushes Release Tag
```bash
# Make code changes
git add .
git commit -m "feat: add new feature"

# Create and push version tag
git tag v1.1.1
git push origin v1.1.1
```

### Step 2: Build Workflow Runs (Automatic)
```
✅ build.yml triggered
  ├─ Extract version: 1.1.1
  ├─ Docker login
  ├─ Build order-core (Go)
  ├─ Build order-state (Python)
  ├─ Build risk-engine (Python)
  ├─ Build tp-sl-engine (Python)
  ├─ Build notify-service (Python)
  ├─ Collect digests
  └─ Upload artifacts
```

### Step 3: Staging Deployment (Automatic)
```
✅ deploy.yml triggered (workflow_run)
  ├─ Download digests
  ├─ SSH to staging VPS
  ├─ Pull latest images
  ├─ docker compose pull
  ├─ docker compose up -d
  └─ Notification: "Staging deployment successful!"
```

### Step 4: QA Testing
```
🧪 Test in staging environment
   - Health checks
   - Integration tests
   - Performance validation
```

### Step 5: Production Promotion (Manual)
```
Developer triggers via GitHub Actions UI:

✅ promote-to-prod.yml triggered (workflow_dispatch)
  ├─ Checkout v1.1.1 tag
  ├─ Verify images
  ├─ SSH to production VPS
  ├─ Pull production images
  ├─ docker compose -f docker-compose.prod.yml up -d
  └─ Notification: "Production deployment successful!"
```

---

## 📊 Service Configuration

Each service has:
- **Dockerfile** (optimized for size)
- **requirements.txt** (Python only)
- **go.mod/go.sum** (Go only)
- **Docker Hub image** (e.g., `alexanderloveiris/order-core`)

### Services:
1. **order-core** - Go (fast!)
2. **order-state** - Python
3. **risk-engine** - Python
4. **tp-sl-engine** - Python
5. **notify-service** - Python

---

## 🔐 Required Secrets

Set these in GitHub repository Settings → Secrets:

```
DOCKERHUB_USERNAME    - Docker Hub username
DOCKERHUB_TOKEN       - Docker Hub personal access token
VPS_HOST              - Staging/Prod VPS IP or hostname
VPS_USER              - VPS SSH username
VPS_SSH_KEY           - VPS SSH private key
```

---

## 📝 Environment Files

### Staging VPS (`app/trading`)
```
.env                          (environment variables)
docker-compose.yml      (development setup)
docker-compose.prod.yml       (production setup)
```

### Digests Update Flow
```
build.yml
  ↓ creates digests.txt
deploy.yml
  ↓ reads digests.txt
  ↓ updates .env
VPS
  ↓ docker compose pull
  ↓ docker compose up -d
```

---

## ✅ Workflow Status Checks

**Build Workflow** (`build.yml`):
- [x] Checkout code
- [x] Extract version from tag
- [x] Docker login
- [x] Build all 5 services (parallel)
- [x] Collect digests
- [x] Upload artifacts

**Deploy Workflow** (`deploy.yml`):
- [x] Wait for build completion
- [x] Download digests
- [x] SSH to VPS
- [x] Pull images
- [x] Restart services
- [x] Verify deployment

**Promote Workflow** (`promote-to-prod.yml`):
- [x] Manual trigger via GitHub UI
- [x] Version input validation
- [x] Image verification
- [x] Production deployment
- [x] Deployment summary

---

## 🚨 Troubleshooting

### Build fails
```
Check: .github/workflows/build.yml logs
Verify: Docker Hub credentials
Verify: Dockerfiles are correct
```

### Auto-deploy to staging fails
```
Check: deploy.yml logs in GitHub Actions
Verify: VPS_HOST, VPS_USER, VPS_SSH_KEY secrets
Verify: VPS has Docker & docker-compose installed
Verify: app/trading directory exists and is writable
```

### Promote to prod fails
```
Check: promote-to-prod.yml logs
Verify: Version tag exists (e.g., v1.1.1)
Verify: Images were built for that version
Verify: Production environment approval (if needed)
```

---

## 🎯 Key Points

✅ **Automatic**: Build → Staging deployment  
✅ **Manual**: Staging → Production (controlled release)  
✅ **Versioned**: Every tag creates a release  
✅ **Auditable**: Every deployment logged in GitHub  
✅ **Rollback**: Push new tag with previous version  
✅ **Isolated**: Staging and prod are separate  

---

## 📞 Support

For workflow issues, check:
1. GitHub Actions tab → Workflow runs
2. VPS logs: `docker logs <container>`
3. SSH key permissions: `chmod 600 ~/.ssh/key`
4. Network access: `ssh -i key user@host` test

---

**Last Updated**: January 28, 2026  
**Status**: ✅ Production Ready
