terraform {
  required_providers {
    pipeops = {
      source = "PipeOpsHQ/pipeops"
    }
  }
}

provider "pipeops" {
  # token is read from PIPEOPS_API_TOKEN environment variable
  # or can be set explicitly:
  # token = "your-api-token"
}

# Variables
variable "workspace_id" {
  description = "PipeOps workspace ID"
  type        = string
}

# Create a project
resource "pipeops_project" "example" {
  name         = "my-terraform-project"
  description  = "Project managed by Terraform"
  workspace_id = var.workspace_id
  repo_url     = "https://github.com/example/repo"
  repo_branch  = "main"
}

# Create production environment
resource "pipeops_environment" "production" {
  name        = "production"
  project_id  = pipeops_project.example.id
  description = "Production environment"
  type        = "production"
}

# Create staging environment
resource "pipeops_environment" "staging" {
  name        = "staging"
  project_id  = pipeops_project.example.id
  description = "Staging environment"
  type        = "staging"
}

# Create development environment
resource "pipeops_environment" "development" {
  name        = "development"
  project_id  = pipeops_project.example.id
  description = "Development environment"
  type        = "development"
}

# Create a server
resource "pipeops_server" "app_server" {
  name        = "app-server-01"
  project_id  = pipeops_project.example.id
  server_type = "web"
  region      = "us-east-1"
  size        = "small"
}

# Data source example - fetch project details
data "pipeops_project" "example" {
  id = pipeops_project.example.id
}

# Outputs
output "project_id" {
  description = "Created project ID"
  value       = pipeops_project.example.id
}

output "project_name" {
  description = "Project name"
  value       = data.pipeops_project.example.name
}

output "production_env_id" {
  description = "Production environment ID"
  value       = pipeops_environment.production.id
}

output "staging_env_id" {
  description = "Staging environment ID"
  value       = pipeops_environment.staging.id
}

output "server_id" {
  description = "Server ID"
  value       = pipeops_server.app_server.id
}

output "server_status" {
  description = "Server status"
  value       = pipeops_server.app_server.status
}

output "server_ip" {
  description = "Server IP address"
  value       = pipeops_server.app_server.ip_address
}
