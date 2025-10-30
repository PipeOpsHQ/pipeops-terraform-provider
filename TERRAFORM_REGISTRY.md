# Publishing to Terraform Registry

This guide covers how to publish the PipeOps Terraform Provider to the Terraform Registry and how to use it from GitHub directly.

## Option 1: Publishing to Terraform Registry (Recommended)

### Prerequisites

1. **GitHub Repository Requirements**
   - Repository must be public
   - Repository name format: `terraform-provider-{NAME}`
   - Example: `terraform-provider-pipeops`

2. **GPG Key for Signing**
   ```bash
   # Generate a GPG key
   gpg --full-generate-key
   
   # Export the public key
   gpg --armor --export your-email@example.com
   
   # Export the private key (for GitHub secrets)
   gpg --armor --export-secret-keys your-email@example.com
   ```

3. **GitHub Secrets**
   Add these secrets to your GitHub repository:
   - `GPG_PRIVATE_KEY` - Your GPG private key
   - `GPG_PASSPHRASE` - Your GPG key passphrase

### Steps to Publish

#### 1. Prepare Your Repository

```bash
# Ensure your repository is clean
git status

# Ensure all tests pass
make test

# Build and verify
make build
```

#### 2. Create and Push a Version Tag

```bash
# Create a git tag (semantic versioning)
git tag -a v1.0.0 -m "Release v1.0.0"

# Push the tag
git push origin v1.0.0
```

This will trigger the GitHub Actions workflow that:
- Builds binaries for all platforms
- Signs them with your GPG key
- Creates a GitHub release
- Prepares files for Terraform Registry

#### 3. Register on Terraform Registry

1. Go to https://registry.terraform.io/
2. Sign in with GitHub
3. Click "Publish" → "Provider"
4. Select your repository: `PipeOpsHQ/terraform-provider-pipeops`
5. Terraform Registry will automatically:
   - Detect your releases
   - Verify GPG signatures
   - Generate documentation
   - Publish the provider

#### 4. Verify Publication

Once published, users can use it:

```hcl
terraform {
  required_providers {
    pipeops = {
      source  = "PipeOpsHQ/pipeops"
      version = "~> 1.0"
    }
  }
}

provider "pipeops" {
  token = var.pipeops_token
}
```

### Version Management

Follow semantic versioning:
- `v1.0.0` - Major release (breaking changes)
- `v1.1.0` - Minor release (new features, backward compatible)
- `v1.1.1` - Patch release (bug fixes)

```bash
# For a new feature
git tag -a v1.1.0 -m "Add webhook resource"
git push origin v1.1.0

# For a bug fix
git tag -a v1.1.1 -m "Fix environment creation bug"
git push origin v1.1.1
```

## Option 2: Using Provider from GitHub (Alternative)

If you want to use the provider directly from GitHub without publishing to the registry:

### A. Using GitHub Releases

#### 1. Create Local Provider Directory

```bash
mkdir -p ~/.terraform.d/plugins/github.com/PipeOpsHQ/pipeops/1.0.0/darwin_arm64
```

#### 2. Download and Install Binary

```bash
# Download from GitHub releases
VERSION="v1.0.0"
OS_ARCH="darwin_arm64"  # or linux_amd64, windows_amd64, etc.

curl -L "https://github.com/PipeOpsHQ/terraform-provider-pipeops/releases/download/${VERSION}/terraform-provider-pipeops_${VERSION}_${OS_ARCH}.zip" \
  -o provider.zip

# Extract
unzip provider.zip

# Move to plugins directory
mv terraform-provider-pipeops_v1.0.0 \
  ~/.terraform.d/plugins/github.com/PipeOpsHQ/pipeops/1.0.0/${OS_ARCH}/
```

#### 3. Configure Terraform to Use GitHub Source

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

### B. Using Local Development Version

For development and testing:

#### 1. Build and Install Locally

```bash
# Build the provider
make build

# Install to local plugins directory
make install
```

#### 2. Use Development Override

Create `~/.terraformrc`:

```hcl
provider_installation {
  dev_overrides {
    "registry.terraform.io/PipeOpsHQ/pipeops" = "/path/to/your/provider/binary"
  }
  direct {}
}
```

#### 3. Use in Terraform

```hcl
terraform {
  required_providers {
    pipeops = {
      source = "registry.terraform.io/PipeOpsHQ/pipeops"
    }
  }
}
```

### C. Using Git Module Source (Advanced)

You can also reference the provider as a module:

```hcl
# This requires building the provider first
module "pipeops_setup" {
  source = "github.com/PipeOpsHQ/terraform-provider-pipeops//examples/basic"
}
```

## CI/CD Integration

### GitHub Actions Example

```yaml
name: Terraform Deploy

on:
  push:
    branches: [main]

jobs:
  terraform:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Setup Terraform
        uses: hashicorp/setup-terraform@v3
        with:
          terraform_version: 1.6.0
      
      - name: Terraform Init
        run: terraform init
        env:
          PIPEOPS_API_TOKEN: ${{ secrets.PIPEOPS_API_TOKEN }}
      
      - name: Terraform Plan
        run: terraform plan
        env:
          PIPEOPS_API_TOKEN: ${{ secrets.PIPEOPS_API_TOKEN }}
      
      - name: Terraform Apply
        if: github.ref == 'refs/heads/main'
        run: terraform apply -auto-approve
        env:
          PIPEOPS_API_TOKEN: ${{ secrets.PIPEOPS_API_TOKEN }}
```

### GitLab CI Example

```yaml
stages:
  - plan
  - apply

variables:
  TF_VERSION: "1.6.0"

before_script:
  - apk add --update curl unzip
  - curl -o terraform.zip https://releases.hashicorp.com/terraform/${TF_VERSION}/terraform_${TF_VERSION}_linux_amd64.zip
  - unzip terraform.zip
  - mv terraform /usr/local/bin/

terraform_plan:
  stage: plan
  script:
    - terraform init
    - terraform plan -out=plan.tfplan
  artifacts:
    paths:
      - plan.tfplan

terraform_apply:
  stage: apply
  script:
    - terraform init
    - terraform apply plan.tfplan
  only:
    - main
  when: manual
```

## Comparison: Registry vs GitHub

| Feature | Terraform Registry | GitHub Direct |
|---------|-------------------|---------------|
| **Ease of Use** | ✅ Simplest | ⚠️ Manual setup |
| **Auto Updates** | ✅ Version constraints | ❌ Manual |
| **Documentation** | ✅ Auto-generated | ⚠️ Manual |
| **Discovery** | ✅ Public listing | ❌ Need URL |
| **Enterprise** | ✅ Terraform Cloud integration | ✅ Works everywhere |
| **Private Use** | ❌ Must be public | ✅ Can be private |

## Best Practices

### For Registry Publication

1. **Semantic Versioning**
   - Always use semantic versioning
   - Document breaking changes in CHANGELOG.md
   - Use pre-release tags for beta versions (v1.0.0-beta.1)

2. **Documentation**
   - Keep README.md up to date
   - Add examples for each resource
   - Document all attributes

3. **Testing**
   - Run acceptance tests before tagging
   - Test on multiple platforms
   - Verify examples work

4. **Security**
   - Sign all releases with GPG
   - Use GitHub secrets for credentials
   - Never commit private keys

### For GitHub Usage

1. **Versioning**
   - Always specify exact versions
   - Pin to specific tags or commits
   - Document version in your README

2. **Distribution**
   - Provide pre-built binaries
   - Include SHA256 checksums
   - Support major platforms

3. **Updates**
   - Notify users of updates
   - Maintain CHANGELOG.md
   - Use GitHub Releases

## Troubleshooting

### Registry Publication Issues

**Issue: GPG signature verification failed**
```bash
# Verify your GPG key
gpg --list-keys

# Re-export and update GitHub secret
gpg --armor --export-secret-keys your-email@example.com
```

**Issue: Provider not appearing in registry**
- Ensure repository is public
- Verify naming convention: `terraform-provider-{name}`
- Check that releases have proper format

### GitHub Direct Usage Issues

**Issue: Provider not found**
```bash
# Verify plugin directory
ls ~/.terraform.d/plugins/

# Check terraform init output for errors
terraform init -upgrade
```

**Issue: Version conflict**
```hcl
# Use exact version
required_providers {
  pipeops = {
    source  = "github.com/PipeOpsHQ/pipeops"
    version = "= 1.0.0"
  }
}
```

## Release Checklist

- [ ] All tests passing
- [ ] Documentation updated
- [ ] CHANGELOG.md updated
- [ ] Examples tested
- [ ] Version bumped in code
- [ ] Git tag created
- [ ] GitHub release created
- [ ] Binaries uploaded
- [ ] Registry updated (if applicable)
- [ ] Announcement made

## Additional Resources

- [Terraform Registry Requirements](https://www.terraform.io/registry/providers/publishing)
- [Provider Development Guide](https://developer.hashicorp.com/terraform/plugin)
- [Semantic Versioning](https://semver.org/)
- [GPG Signing Guide](https://docs.github.com/en/authentication/managing-commit-signature-verification)

## Support

For issues with:
- **Registry publication**: Check [Terraform Registry Docs](https://www.terraform.io/registry/providers/publishing)
- **Provider functionality**: Open an issue on [GitHub](https://github.com/PipeOpsHQ/terraform-provider-pipeops/issues)
- **PipeOps API**: Contact [PipeOps Support](https://pipeops.io/support)
