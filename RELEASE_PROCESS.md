# Release Process Guide

This document outlines the complete process for releasing new versions of the PipeOps Terraform Provider.

## Prerequisites

Before releasing, ensure you have:

1. **GPG Key Setup**
   ```bash
   # Generate GPG key if you don't have one
   gpg --full-generate-key
   
   # List your keys
   gpg --list-secret-keys --keyid-format=long
   
   # Export public key (for Terraform Registry)
   gpg --armor --export YOUR_KEY_ID > public_key.asc
   ```

2. **GitHub Secrets Configured**
   - `GPG_PRIVATE_KEY` - Your private GPG key
   - `GPG_PASSPHRASE` - GPG key passphrase
   - `GITHUB_TOKEN` - Automatically provided by GitHub

3. **Repository Access**
   - Write access to the repository
   - Ability to create tags and releases

## Release Types

Follow [Semantic Versioning](https://semver.org/):

- **Major (v2.0.0)** - Breaking changes
- **Minor (v1.1.0)** - New features, backward compatible
- **Patch (v1.0.1)** - Bug fixes, backward compatible

## Step-by-Step Release Process

### 1. Prepare the Release

#### Update Version Documentation

Create/update `CHANGELOG.md`:

```markdown
## [1.0.0] - 2025-01-30

### Added
- Initial release
- Project resource with full CRUD operations
- Environment resource
- Server resource
- Project data source

### Changed
- N/A

### Fixed
- N/A
```

#### Run Tests

```bash
# Run all tests
make test

# Run acceptance tests if available
make testacc

# Build and verify
make build

# Test examples
cd examples/basic
terraform init
terraform validate
cd ../..
```

#### Update Documentation

```bash
# Ensure README is up to date
# Verify all examples work
# Check that GETTING_STARTED.md reflects current state
```

### 2. Create the Release

#### Method A: Using Git Tags (Recommended)

```bash
# Ensure you're on main/master branch
git checkout main
git pull origin main

# Create and push a tag
VERSION="v1.0.0"
git tag -a $VERSION -m "Release $VERSION"
git push origin $VERSION
```

This will automatically:
1. Trigger the GitHub Actions workflow
2. Build binaries for all platforms
3. Sign releases with GPG
4. Create a GitHub Release
5. Publish to Terraform Registry (if registered)

#### Method B: Using GitHub CLI

```bash
# Install GitHub CLI if needed
brew install gh

# Create release
gh release create v1.0.0 \
  --title "v1.0.0" \
  --notes "See CHANGELOG.md for details"
```

### 3. Verify the Release

#### Check GitHub Release

1. Go to https://github.com/PipeOpsHQ/terraform-provider-pipeops/releases
2. Verify the release exists
3. Check that all platform binaries are present:
   - darwin_amd64
   - darwin_arm64
   - linux_amd64
   - linux_arm64
   - windows_amd64
4. Verify checksums file is present
5. Verify GPG signature is present

#### Test the Release

```bash
# Download and test a binary
VERSION="v1.0.0"
OS_ARCH="darwin_arm64"

curl -L "https://github.com/PipeOpsHQ/terraform-provider-pipeops/releases/download/${VERSION}/terraform-provider-pipeops_${VERSION}_${OS_ARCH}.zip" -o provider.zip

unzip provider.zip
./terraform-provider-pipeops_v1.0.0 -version
```

#### Verify Terraform Registry (if published)

1. Go to https://registry.terraform.io/providers/PipeOpsHQ/pipeops
2. Check that the new version appears
3. Verify documentation is generated
4. Test installation:

```bash
# Create a test directory
mkdir test-release
cd test-release

# Create main.tf
cat > main.tf <<EOF
terraform {
  required_providers {
    pipeops = {
      source  = "PipeOpsHQ/pipeops"
      version = "~> ${VERSION}"
    }
  }
}
EOF

# Test init
terraform init
```

### 4. Announce the Release

#### Update Documentation Sites

- [ ] Update main README if needed
- [ ] Update examples if API changed
- [ ] Update GETTING_STARTED.md if workflow changed

#### Notify Users

```bash
# Post to GitHub Discussions (if enabled)
# Tweet about the release
# Update company blog/changelog
# Notify in Slack/Discord channels
```

## Publishing to Terraform Registry

### First Time Setup

1. **Register on Terraform Registry**
   - Go to https://registry.terraform.io
   - Sign in with GitHub
   - Click "Publish" → "Provider"
   - Select repository: `PipeOpsHQ/terraform-provider-pipeops`

2. **Add GPG Public Key**
   ```bash
   # Get your GPG key ID
   gpg --list-keys
   
   # Export public key
   gpg --armor --export YOUR_EMAIL > public_key.asc
   ```
   - Go to https://registry.terraform.io/settings/gpg-keys
   - Add your public key

3. **Verify Repository Settings**
   - Repository must be public
   - Repository name: `terraform-provider-{name}`
   - Must have proper tags (v1.0.0 format)
   - Must have signed releases

### Subsequent Releases

Once registered, Terraform Registry automatically:
1. Detects new tagged releases
2. Downloads and verifies binaries
3. Generates documentation from schema
4. Makes the version available

**Note:** It may take 5-10 minutes for a new release to appear on the registry.

## Rollback Process

If you need to rollback a release:

### Delete a GitHub Release

```bash
# Using GitHub CLI
gh release delete v1.0.0 --yes

# Or manually through GitHub UI
# Go to Releases → Click release → Delete
```

### Remove a Tag

```bash
# Delete local tag
git tag -d v1.0.0

# Delete remote tag
git push origin :refs/tags/v1.0.0
```

### Revert Code Changes

```bash
# Create a revert tag
git revert <commit-hash>
git tag -a v1.0.1 -m "Revert changes from v1.0.0"
git push origin v1.0.1
```

## Hotfix Process

For critical bug fixes:

```bash
# Create hotfix branch from tag
git checkout -b hotfix/v1.0.1 v1.0.0

# Make fixes
git commit -m "fix: critical bug in resource creation"

# Merge to main
git checkout main
git merge hotfix/v1.0.1

# Tag and release
git tag -a v1.0.1 -m "Hotfix: Critical bug fixes"
git push origin v1.0.1

# Clean up
git branch -d hotfix/v1.0.1
```

## Release Checklist

Before tagging:
- [ ] All tests pass
- [ ] Examples work
- [ ] CHANGELOG.md updated
- [ ] Documentation updated
- [ ] Version bumped (if applicable)
- [ ] Branch is up to date

After tagging:
- [ ] GitHub Release created successfully
- [ ] All binaries built
- [ ] GPG signatures present
- [ ] Checksums verified
- [ ] Terraform Registry updated (if published)
- [ ] Documentation generated
- [ ] Installation tested
- [ ] Announcement prepared

## Troubleshooting

### Release Workflow Failed

Check GitHub Actions logs:
1. Go to Actions tab
2. Find the failed workflow
3. Review error messages
4. Common issues:
   - GPG key not configured
   - Invalid version format
   - Build failures

### GPG Signing Failed

```bash
# Verify GPG key in GitHub secrets
# Re-export and update secret
gpg --armor --export-secret-keys YOUR_EMAIL

# Test signing locally
echo "test" | gpg --sign --armor
```

### Registry Not Updating

- Wait 10-15 minutes
- Check release format (must be proper semantic version)
- Verify GPG signature on release
- Check Terraform Registry status page
- Contact Terraform Registry support if needed

### Binary Not Working

```bash
# Check binary is executable
chmod +x terraform-provider-pipeops_v1.0.0

# Verify version
./terraform-provider-pipeops_v1.0.0 -version

# Check for missing dependencies
ldd terraform-provider-pipeops_v1.0.0  # Linux
otool -L terraform-provider-pipeops_v1.0.0  # macOS
```

## Continuous Deployment

For automated releases on every merge to main:

```yaml
# .github/workflows/auto-release.yml
name: Auto Release

on:
  push:
    branches:
      - main
    paths-ignore:
      - '**.md'
      - 'docs/**'

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0
      
      - name: Determine version
        id: version
        run: |
          # Your versioning logic here
          echo "version=v1.0.${GITHUB_RUN_NUMBER}" >> $GITHUB_OUTPUT
      
      - name: Create tag
        run: |
          git tag ${{ steps.version.outputs.version }}
          git push origin ${{ steps.version.outputs.version }}
```

## Best Practices

1. **Version Naming**
   - Always use `v` prefix (v1.0.0, not 1.0.0)
   - Follow semantic versioning strictly
   - Use pre-release tags for beta versions (v1.0.0-beta.1)

2. **Release Notes**
   - Clear description of changes
   - Migration guide for breaking changes
   - Examples of new features
   - Credits to contributors

3. **Testing**
   - Test on multiple platforms before release
   - Verify examples work
   - Run acceptance tests
   - Test upgrade path from previous version

4. **Documentation**
   - Update before releasing
   - Include migration guides
   - Keep changelog up to date
   - Update version in examples

5. **Communication**
   - Announce breaking changes in advance
   - Provide deprecation warnings
   - Give users time to migrate
   - Be available for questions

## Support

For release issues:
- Check [GitHub Actions logs](https://github.com/PipeOpsHQ/terraform-provider-pipeops/actions)
- Review [Terraform Registry docs](https://www.terraform.io/registry/providers/publishing)
- Open an issue for help
- Contact maintainers

---

**Last Updated:** 2025-01-30
**Maintained by:** PipeOps Team
