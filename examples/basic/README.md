# Basic PipeOps Terraform Example

This example demonstrates basic usage of the PipeOps Terraform provider.

## What This Example Creates

- A PipeOps project
- Three environments (production, staging, development)
- A server resource
- Demonstrates data source usage

## Prerequisites

- Terraform >= 1.0
- PipeOps API token
- PipeOps workspace ID

## Usage

1. Set your API token:
   ```bash
   export PIPEOPS_API_TOKEN="your-token-here"
   ```

2. Copy and configure variables:
   ```bash
   cp terraform.tfvars.example terraform.tfvars
   # Edit terraform.tfvars and add your workspace_id
   ```

3. Initialize Terraform:
   ```bash
   terraform init
   ```

4. Preview changes:
   ```bash
   terraform plan
   ```

5. Apply the configuration:
   ```bash
   terraform apply
   ```

6. View outputs:
   ```bash
   terraform output
   ```

## Cleanup

To destroy all resources:
```bash
terraform destroy
```

## What You'll Learn

- How to configure the PipeOps provider
- How to create projects and environments
- How to manage servers
- How to use data sources
- How to organize outputs
