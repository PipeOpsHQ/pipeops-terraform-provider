# Terraform Provider for PipeOps

A Terraform provider for managing PipeOps infrastructure and DevOps automations.

## Features

- **Project Management**: Create and manage PipeOps projects
- **Environment Management**: Configure development, staging, and production environments
- **Server Management**: Provision and manage infrastructure servers
- **Data Sources**: Query existing PipeOps resources

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.21 (for development)
- PipeOps API token

## Installation

### Option 1: From Terraform Registry (Recommended)

Once published, you can use it directly:

```hcl
terraform {
  required_providers {
    pipeops = {
      source  = "PipeOpsHQ/pipeops"
      version = "~> 1.0"
    }
  }
}
```

### Option 2: From GitHub Releases

Download the binary for your platform from [GitHub Releases](https://github.com/PipeOpsHQ/terraform-provider-pipeops/releases):

```bash
# Example for macOS ARM64
VERSION="v1.0.0"
OS_ARCH="darwin_arm64"

# Download
curl -L "https://github.com/PipeOpsHQ/terraform-provider-pipeops/releases/download/${VERSION}/terraform-provider-pipeops_${VERSION}_${OS_ARCH}.zip" -o provider.zip

# Extract and install
unzip provider.zip
mkdir -p ~/.terraform.d/plugins/github.com/PipeOpsHQ/pipeops/1.0.0/${OS_ARCH}
mv terraform-provider-pipeops_* ~/.terraform.d/plugins/github.com/PipeOpsHQ/pipeops/1.0.0/${OS_ARCH}/
```

Then use in your Terraform:

```hcl
terraform {
  required_providers {
    pipeops = {
      source  = "github.com/PipeOpsHQ/pipeops"
      version = "~> 1.0"
    }
  }
}
```

### Option 3: Build from Source (Development)

```bash
# Clone the repository
git clone https://github.com/PipeOpsHQ/terraform-provider-pipeops.git
cd terraform-provider-pipeops

# Build the provider
make build

# Install to local Terraform plugins directory
make install
```

See [TERRAFORM_REGISTRY.md](TERRAFORM_REGISTRY.md) for detailed publishing and installation instructions.

## Usage

### Provider Configuration

```hcl
provider "pipeops" {
  token    = "your-api-token"  # or set PIPEOPS_API_TOKEN
  base_url = "https://api.pipeops.io"  # optional
}
```

### Environment Variables

- `PIPEOPS_API_TOKEN`: Your PipeOps API token
- `PIPEOPS_BASE_URL`: PipeOps API base URL (default: https://api.pipeops.io)

### Example: Create a Project

```hcl
resource "pipeops_project" "example" {
  name         = "my-application"
  description  = "My awesome application"
  workspace_id = "workspace-123"
  repo_url     = "https://github.com/user/repo"
  repo_branch  = "main"
}
```

### Example: Create an Environment

```hcl
resource "pipeops_environment" "production" {
  name        = "production"
  project_id  = pipeops_project.example.id
  description = "Production environment"
  type        = "production"
}
```

### Example: Create a Server

```hcl
resource "pipeops_server" "app_server" {
  name        = "app-server-01"
  project_id  = pipeops_project.example.id
  server_type = "web"
  region      = "us-east-1"
  size        = "small"
}
```

### Example: Data Source

```hcl
data "pipeops_project" "existing" {
  id = "project-123"
}

output "project_name" {
  value = data.pipeops_project.existing.name
}
```

## Complete Example

See the [examples](./examples) directory for complete working examples.

```hcl
terraform {
  required_providers {
    pipeops = {
      source = "PipeOpsHQ/pipeops"
    }
  }
}

provider "pipeops" {
  token = var.pipeops_token
}

resource "pipeops_project" "app" {
  name         = "production-app"
  workspace_id = var.workspace_id
  description  = "Production application"
}

resource "pipeops_environment" "prod" {
  name       = "production"
  project_id = pipeops_project.app.id
  type       = "production"
}

resource "pipeops_environment" "staging" {
  name       = "staging"
  project_id = pipeops_project.app.id
  type       = "staging"
}

resource "pipeops_server" "web" {
  name       = "web-server"
  project_id = pipeops_project.app.id
  region     = "us-east-1"
  size       = "medium"
}
```

## Development

### Building

```bash
make build
```

### Testing

```bash
make test
```

### Local Testing

```bash
# Build and install locally
make install

# Use in your Terraform configuration
cd examples/basic
terraform init
terraform plan
terraform apply
```

## Resources

### `pipeops_project`

Manages a PipeOps project.

**Arguments:**
- `name` (Required): Project name
- `workspace_id` (Required): Workspace ID
- `description` (Optional): Project description
- `repo_url` (Optional): Repository URL
- `repo_branch` (Optional): Repository branch

**Attributes:**
- `id`: Project ID
- `created_at`: Creation timestamp
- `updated_at`: Last update timestamp

### `pipeops_environment`

Manages a PipeOps environment.

**Arguments:**
- `name` (Required): Environment name
- `project_id` (Required): Project ID
- `description` (Optional): Environment description
- `type` (Optional): Environment type

**Attributes:**
- `id`: Environment ID
- `created_at`: Creation timestamp
- `updated_at`: Last update timestamp

### `pipeops_server`

Manages a PipeOps server.

**Arguments:**
- `name` (Required): Server name
- `project_id` (Required): Project ID
- `server_type` (Optional): Server type
- `region` (Optional): Server region
- `size` (Optional): Server size

**Attributes:**
- `id`: Server ID
- `status`: Server status
- `ip_address`: Server IP address
- `created_at`: Creation timestamp
- `updated_at`: Last update timestamp

## Data Sources

### `pipeops_project`

Fetches information about a PipeOps project.

**Arguments:**
- `id` (Required): Project ID

**Attributes:**
- All project attributes

## Import

Resources can be imported using their ID:

```bash
terraform import pipeops_project.example project-123
terraform import pipeops_environment.prod env-456
terraform import pipeops_server.web server-789
```

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for details.

## License

This provider is distributed under the terms specified in the LICENSE file.

## Support

- [Documentation](https://docs.pipeops.io)
- [PipeOps API Documentation](https://api.pipeops.io/docs)
- [Issues](https://github.com/PipeOpsHQ/terraform-provider-pipeops/issues)
