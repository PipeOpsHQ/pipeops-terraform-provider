# Getting Started with PipeOps Terraform Provider

## 🚀 Quick Start (5 minutes)

### Step 1: Build the Provider

```bash
# Navigate to the provider directory
cd pipeops-terraform-provider

# Build the provider
make build

# Install locally for Terraform
make install
```

### Step 2: Configure Your PipeOps API Token

```bash
# Export your API token
export PIPEOPS_API_TOKEN="your-pipeops-api-token-here"
```

### Step 3: Create Your First Terraform Configuration

Create a file named `main.tf`:

```hcl
terraform {
  required_providers {
    pipeops = {
      source = "PipeOpsHQ/pipeops"
    }
  }
}

provider "pipeops" {
  # token is read from PIPEOPS_API_TOKEN environment variable
}

variable "workspace_id" {
  description = "Your PipeOps workspace ID"
  type        = string
}

# Create a project
resource "pipeops_project" "my_app" {
  name         = "my-first-app"
  description  = "My first application managed by Terraform"
  workspace_id = var.workspace_id
  repo_url     = "https://github.com/yourusername/your-repo"
  repo_branch  = "main"
}

# Create production environment
resource "pipeops_environment" "production" {
  name        = "production"
  project_id  = pipeops_project.my_app.id
  description = "Production environment"
}

# Create staging environment
resource "pipeops_environment" "staging" {
  name        = "staging"
  project_id  = pipeops_project.my_app.id
  description = "Staging environment"
}

# Output project ID
output "project_id" {
  value       = pipeops_project.my_app.id
  description = "The ID of the created project"
}
```

Create `terraform.tfvars`:

```hcl
workspace_id = "your-workspace-id-here"
```

### Step 4: Deploy!

```bash
# Initialize Terraform
terraform init

# Preview changes
terraform plan

# Apply configuration
terraform apply

# When you're done, clean up
terraform destroy
```

## 📚 Learn by Example

### Example 1: Basic Project Setup

```hcl
resource "pipeops_project" "web_app" {
  name         = "production-web-app"
  workspace_id = var.workspace_id
  description  = "Main production application"
  repo_url     = "https://github.com/company/web-app"
  repo_branch  = "main"
}
```

### Example 2: Multi-Environment Setup

```hcl
locals {
  environments = ["development", "staging", "production"]
}

resource "pipeops_environment" "envs" {
  for_each = toset(local.environments)
  
  name        = each.key
  project_id  = pipeops_project.web_app.id
  description = "${title(each.key)} environment"
}
```

### Example 3: Server Provisioning

```hcl
resource "pipeops_server" "app_servers" {
  name        = "app-server-${var.environment}"
  project_id  = pipeops_project.web_app.id
  server_type = "web"
  region      = var.region
  size        = var.server_size
}
```

### Example 4: Using Data Sources

```hcl
# Query an existing project
data "pipeops_project" "existing" {
  id = "project-id-here"
}

# Use its data
resource "pipeops_environment" "new_env" {
  name       = "new-environment"
  project_id = data.pipeops_project.existing.id
}
```

## 🔧 Common Workflows

### Workflow 1: New Application Deployment

```bash
# 1. Write your Terraform config
cat > main.tf << 'EOF'
resource "pipeops_project" "app" {
  name         = "my-app"
  workspace_id = var.workspace_id
  repo_url     = "https://github.com/user/app"
}
EOF

# 2. Initialize and apply
terraform init
terraform apply -auto-approve

# 3. Get the project ID
terraform output project_id
```

### Workflow 2: Import Existing Resources

```bash
# Import an existing project
terraform import pipeops_project.existing project-uuid-here

# Import an existing environment
terraform import pipeops_environment.prod env-uuid-here

# Verify the state
terraform show
```

### Workflow 3: Managing Multiple Environments

```bash
# Use workspaces for different stages
terraform workspace new development
terraform workspace new staging
terraform workspace new production

# Switch between them
terraform workspace select production
terraform apply
```

## 🛠️ Troubleshooting

### Issue: Provider not found

**Solution:**
```bash
# Rebuild and reinstall
make clean
make install

# Verify installation
ls ~/.terraform.d/plugins/registry.terraform.io/PipeOpsHQ/pipeops/
```

### Issue: Authentication failed

**Solution:**
```bash
# Verify your token
echo $PIPEOPS_API_TOKEN

# If empty, export it
export PIPEOPS_API_TOKEN="your-token"

# Or set it in the provider block
provider "pipeops" {
  token = "your-token"  # Not recommended for production
}
```

### Issue: Resource not found

**Solution:**
```bash
# Refresh Terraform state
terraform refresh

# Re-import if needed
terraform import pipeops_project.app project-id
```

## 📖 Next Steps

1. **Explore Examples**
   ```bash
   cd examples/basic
   cat README.md
   ```

2. **Read Full Documentation**
   - [README.md](README.md) - Complete feature list
   - [CONTRIBUTING.md](CONTRIBUTING.md) - Development guide
   - [PROJECT_SUMMARY.md](PROJECT_SUMMARY.md) - Architecture details

3. **Join the Community**
   - Report issues on GitHub
   - Contribute new resources
   - Share your use cases

## 💡 Pro Tips

1. **Use Variables**
   ```hcl
   variable "environment" {
     type = string
   }
   
   resource "pipeops_project" "app" {
     name = "app-${var.environment}"
   }
   ```

2. **Organize Code**
   ```
   terraform/
   ├── main.tf
   ├── variables.tf
   ├── outputs.tf
   └── terraform.tfvars
   ```

3. **Use Modules**
   ```hcl
   module "application" {
     source = "./modules/application"
     name   = "my-app"
   }
   ```

4. **State Management**
   ```hcl
   terraform {
     backend "s3" {
       bucket = "terraform-state"
       key    = "pipeops/terraform.tfstate"
     }
   }
   ```

## 🎯 Common Use Cases

### Use Case 1: CI/CD Integration

```yaml
# .github/workflows/deploy.yml
- name: Terraform Apply
  run: |
    cd terraform
    terraform init
    terraform apply -auto-approve
  env:
    PIPEOPS_API_TOKEN: ${{ secrets.PIPEOPS_TOKEN }}
```

### Use Case 2: Multi-Region Deployment

```hcl
locals {
  regions = ["us-east-1", "eu-west-1", "ap-south-1"]
}

resource "pipeops_server" "regional" {
  for_each = toset(local.regions)
  
  name   = "server-${each.key}"
  region = each.key
}
```

### Use Case 3: Disaster Recovery

```hcl
resource "pipeops_project" "primary" {
  name = "primary-app"
}

resource "pipeops_project" "dr" {
  name = "dr-app"
  # Mirror configuration
}
```

## ✅ Checklist for Production

- [ ] Store state remotely (S3, Terraform Cloud)
- [ ] Use workspaces or separate state files per environment
- [ ] Never commit API tokens to version control
- [ ] Use CI/CD for automated deployments
- [ ] Enable state locking
- [ ] Document your infrastructure
- [ ] Set up monitoring and alerts

## 🆘 Getting Help

- 📖 [Official Documentation](README.md)
- 🐛 [Report Issues](https://github.com/PipeOpsHQ/terraform-provider-pipeops/issues)
- 💬 [PipeOps Community](https://pipeops.io/community)
- 📧 Email: support@pipeops.io

---

**Ready to automate your DevOps workflow?** Start building! 🚀
