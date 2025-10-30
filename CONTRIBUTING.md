# Contributing to Terraform Provider for PipeOps

Thank you for your interest in contributing to the PipeOps Terraform Provider!

## Development Setup

### Prerequisites

- [Go](https://golang.org/doc/install) >= 1.21
- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Git](https://git-scm.com/downloads)
- PipeOps account and API token

### Getting Started

1. Clone the repository:
   ```bash
   git clone https://github.com/PipeOpsHQ/terraform-provider-pipeops.git
   cd terraform-provider-pipeops
   ```

2. Install dependencies:
   ```bash
   make deps
   ```

3. Build the provider:
   ```bash
   make build
   ```

4. Install locally for testing:
   ```bash
   make install
   ```

## Development Workflow

### Running Tests

```bash
# Run unit tests
make test

# Run acceptance tests (requires API credentials)
export PIPEOPS_API_TOKEN="your-token"
make testacc
```

### Code Format and Linting

```bash
# Format code
make fmt

# Run linter
make lint
```

### Building

```bash
# Build for current platform
make build

# Build for all platforms
make build-all
```

## Project Structure

```
.
├── internal/
│   ├── provider/          # Provider configuration
│   ├── resources/         # Terraform resources
│   └── datasources/       # Terraform data sources
├── examples/              # Usage examples
├── docs/                  # Documentation
└── main.go               # Provider entrypoint
```

## Adding New Resources

1. Create the resource file in `internal/resources/`:
   ```go
   package resources

   import (
       "github.com/PipeOpsHQ/pipeops-go-sdk/pipeops"
       "github.com/hashicorp/terraform-plugin-framework/resource"
   )

   func NewMyResource() resource.Resource {
       return &MyResource{}
   }

   type MyResource struct {
       client *pipeops.Client
   }
   ```

2. Register the resource in `internal/provider/provider.go`:
   ```go
   func (p *PipeOpsProvider) Resources(ctx context.Context) []func() resource.Resource {
       return []func() resource.Resource{
           NewProjectResource,
           NewMyResource, // Add your new resource
       }
   }
   ```

3. Add tests and documentation

## Adding New Data Sources

Similar to resources, but place in `internal/datasources/` and register in the `DataSources()` function.

## Testing

### Unit Tests

Write unit tests for individual functions and methods.

### Acceptance Tests

Acceptance tests interact with the real PipeOps API:

```go
func TestAccProject_basic(t *testing.T) {
    resource.Test(t, resource.TestCase{
        ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
        Steps: []resource.TestStep{
            {
                Config: testAccProjectConfig_basic(),
                Check: resource.ComposeTestCheckFunc(
                    resource.TestCheckResourceAttr("pipeops_project.test", "name", "test-project"),
                ),
            },
        },
    })
}
```

## Documentation

- Update README.md for user-facing changes
- Add examples in the `examples/` directory
- Use `make docs` to generate provider documentation

## Pull Request Process

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/my-feature`
3. Make your changes
4. Run tests: `make test`
5. Format code: `make fmt`
6. Commit with clear messages
7. Push to your fork
8. Open a Pull Request

### PR Guidelines

- Include tests for new features
- Update documentation
- Follow existing code style
- Keep changes focused and atomic
- Write clear commit messages

## Code Style

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` for formatting
- Follow Terraform provider conventions
- Add comments for exported functions

## Common Tasks

### Debug Provider Locally

```bash
# Build with debug flags
go build -gcflags="all=-N -l" -o terraform-provider-pipeops

# Run with dlv
dlv exec --headless --listen=:2345 --api-version=2 ./terraform-provider-pipeops -- -debug
```

### Test with Local Changes

```bash
# Build and install
make install

# Use in your Terraform config
cd examples/basic
terraform init
terraform plan
```

## Getting Help

- Check [existing issues](https://github.com/PipeOpsHQ/terraform-provider-pipeops/issues)
- Review [PipeOps API documentation](https://api.pipeops.io/docs)
- Ask questions in pull requests

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
