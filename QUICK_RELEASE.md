# Quick Release Guide

## TL;DR - Release in 3 Commands

```bash
# 1. Ensure everything is ready
make test && make build

# 2. Create and push tag
git tag -a v1.0.0 -m "Release v1.0.0" && git push origin v1.0.0

# 3. Done! 🎉
# GitHub Actions will automatically:
# - Build binaries for all platforms
# - Create GitHub Release
# - Sign with GPG
# - Publish to Terraform Registry
```

## Two Installation Methods After Release

### Method 1: Terraform Registry (Recommended)

Users install with:

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

Then just:
```bash
terraform init  # That's it!
```

### Method 2: Direct from GitHub

Users can also install from GitHub Releases:

```bash
VERSION="v1.0.0"
OS_ARCH="darwin_arm64"  # or linux_amd64, windows_amd64, etc.

# Download
curl -L "https://github.com/PipeOpsHQ/terraform-provider-pipeops/releases/download/${VERSION}/terraform-provider-pipeops_${VERSION}_${OS_ARCH}.zip" -o provider.zip

# Install
unzip provider.zip
mkdir -p ~/.terraform.d/plugins/github.com/PipeOpsHQ/pipeops/1.0.0/${OS_ARCH}
mv terraform-provider-pipeops_* ~/.terraform.d/plugins/github.com/PipeOpsHQ/pipeops/1.0.0/${OS_ARCH}/
```

Then use in Terraform:

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

## First-Time Setup (One Time Only)

### 1. Generate GPG Key

```bash
gpg --full-generate-key
# Select: (1) RSA and RSA
# Key size: 4096
# Valid for: 0 (doesn't expire)
# Enter your name and email
```

### 2. Export Keys

```bash
# Get your key ID
gpg --list-secret-keys --keyid-format=long

# Export private key for GitHub
gpg --armor --export-secret-keys YOUR_EMAIL > private_key.asc

# Export public key for Terraform Registry
gpg --armor --export YOUR_EMAIL > public_key.asc
```

### 3. Add to GitHub Secrets

Go to: `Settings → Secrets and variables → Actions`

Add these secrets:
- `GPG_PRIVATE_KEY` = Content of `private_key.asc`
- `GPG_PASSPHRASE` = Your GPG passphrase

### 4. Register on Terraform Registry (Optional)

1. Go to https://registry.terraform.io
2. Sign in with GitHub
3. Click "Publish" → "Provider"
4. Select: `PipeOpsHQ/terraform-provider-pipeops`
5. Add your GPG public key at: https://registry.terraform.io/settings/gpg-keys

Done! Now every time you push a tag, the provider is automatically published.

## Version Numbers

Follow semantic versioning:

```bash
# New features (backward compatible)
git tag -a v1.1.0 -m "Release v1.1.0"

# Bug fixes
git tag -a v1.0.1 -m "Release v1.0.1"

# Breaking changes
git tag -a v2.0.0 -m "Release v2.0.0"
```

## Pre-Release Checklist

```bash
# 1. Tests pass
make test

# 2. Build works
make build

# 3. Update CHANGELOG.md
echo "## [1.0.0] - $(date +%Y-%m-%d)" >> CHANGELOG.md
echo "### Added" >> CHANGELOG.md
echo "- Feature description" >> CHANGELOG.md

# 4. Commit changes
git add CHANGELOG.md
git commit -m "chore: prepare v1.0.0 release"
git push
```

## Release Commands

```bash
# Latest stable release
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0

# Pre-release (beta)
git tag -a v1.0.0-beta.1 -m "Beta release v1.0.0-beta.1"
git push origin v1.0.0-beta.1

# Release candidate
git tag -a v1.0.0-rc.1 -m "Release candidate v1.0.0-rc.1"
git push origin v1.0.0-rc.1
```

## Verify Release

```bash
# Check GitHub Release
open https://github.com/PipeOpsHQ/terraform-provider-pipeops/releases

# Check workflow status
open https://github.com/PipeOpsHQ/terraform-provider-pipeops/actions

# Test installation (after ~5 minutes)
mkdir test-install && cd test-install
cat > main.tf <<EOF
terraform {
  required_providers {
    pipeops = {
      source  = "PipeOpsHQ/pipeops"
      version = "~> 1.0"
    }
  }
}
EOF
terraform init
```

## Hotfix Release

```bash
# Create hotfix branch
git checkout -b hotfix/v1.0.1 v1.0.0

# Fix the bug
# ... make changes ...
git commit -m "fix: critical bug"

# Merge to main
git checkout main
git merge hotfix/v1.0.1
git push

# Tag and release
git tag -a v1.0.1 -m "Hotfix v1.0.1"
git push origin v1.0.1
```

## Rollback

```bash
# Delete release and tag
gh release delete v1.0.0 --yes
git tag -d v1.0.0
git push origin :refs/tags/v1.0.0
```

## Common Issues

### Release workflow failed?

Check: https://github.com/PipeOpsHQ/terraform-provider-pipeops/actions

Common fixes:
- Verify GPG secrets are set
- Check tag format (must be v1.0.0)
- Ensure tests pass

### Not appearing on Terraform Registry?

- Wait 10-15 minutes
- Verify GPG signature on release
- Check tag follows semver (v1.0.0)
- Ensure repository is public

### Binary not working?

```bash
# Make executable
chmod +x terraform-provider-pipeops_v1.0.0

# Test
./terraform-provider-pipeops_v1.0.0 -version
```

## Complete Example

```bash
# 1. Prepare
git checkout main
git pull
make test
make build

# 2. Update changelog
cat >> CHANGELOG.md <<EOF

## [1.1.0] - $(date +%Y-%m-%d)

### Added
- New webhook resource
- Enhanced error messages

### Fixed
- Environment creation bug
EOF

git add CHANGELOG.md
git commit -m "chore: prepare v1.1.0 release"
git push

# 3. Release
git tag -a v1.1.0 -m "Release v1.1.0"
git push origin v1.1.0

# 4. Verify (after 5 mins)
open https://github.com/PipeOpsHQ/terraform-provider-pipeops/releases/latest

# 5. Announce
echo "🎉 PipeOps Terraform Provider v1.1.0 is out!"
```

## Links

- **Releases**: https://github.com/PipeOpsHQ/terraform-provider-pipeops/releases
- **Actions**: https://github.com/PipeOpsHQ/terraform-provider-pipeops/actions
- **Registry**: https://registry.terraform.io/providers/PipeOpsHQ/pipeops
- **Full Guide**: [RELEASE_PROCESS.md](RELEASE_PROCESS.md)

---

**That's it!** Tag → Push → Automated Release 🚀
