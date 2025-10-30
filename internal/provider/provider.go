package provider

import (
	"context"
	"os"
	"time"

	"github.com/PipeOpsHQ/pipeops-go-sdk/pipeops"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ provider.Provider = &PipeOpsProvider{}

type PipeOpsProvider struct {
	version string
}

type PipeOpsProviderModel struct {
	Token      types.String `tfsdk:"token"`
	BaseURL    types.String `tfsdk:"base_url"`
	Timeout    types.Int64  `tfsdk:"timeout"`
	MaxRetries types.Int64  `tfsdk:"max_retries"`
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &PipeOpsProvider{
			version: version,
		}
	}
}

func (p *PipeOpsProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "pipeops"
	resp.Version = p.version
}

func (p *PipeOpsProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Terraform provider for managing PipeOps infrastructure and deployments.",
		Attributes: map[string]schema.Attribute{
			"token": schema.StringAttribute{
				Description: "PipeOps API token for authentication. Can also be set via PIPEOPS_API_TOKEN environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
			"base_url": schema.StringAttribute{
				Description: "PipeOps API base URL. Defaults to https://api.pipeops.io. Can also be set via PIPEOPS_BASE_URL environment variable.",
				Optional:    true,
			},
			"timeout": schema.Int64Attribute{
				Description: "Timeout in seconds for API requests. Defaults to 30.",
				Optional:    true,
			},
			"max_retries": schema.Int64Attribute{
				Description: "Maximum number of retries for failed API requests. Defaults to 3.",
				Optional:    true,
			},
		},
	}
}

func (p *PipeOpsProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config PipeOpsProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get token from config or environment
	token := os.Getenv("PIPEOPS_API_TOKEN")
	if !config.Token.IsNull() {
		token = config.Token.ValueString()
	}

	if token == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Missing PipeOps API Token",
			"The provider cannot create the PipeOps API client as there is a missing or empty value for the PipeOps API token. "+
				"Set the token value in the configuration or use the PIPEOPS_API_TOKEN environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
		return
	}

	// Get base URL from config or environment
	baseURL := os.Getenv("PIPEOPS_BASE_URL")
	if !config.BaseURL.IsNull() {
		baseURL = config.BaseURL.ValueString()
	}
	if baseURL == "" {
		baseURL = "https://api.pipeops.io"
	}

	// Get timeout from config or use default
	timeout := int64(30)
	if !config.Timeout.IsNull() {
		timeout = config.Timeout.ValueInt64()
	}

	// Get max retries from config or use default
	maxRetries := int64(3)
	if !config.MaxRetries.IsNull() {
		maxRetries = config.MaxRetries.ValueInt64()
	}

	// Create PipeOps client
	client, err := pipeops.NewClient(
		baseURL,
		pipeops.WithTimeout(time.Duration(timeout)*time.Second),
		pipeops.WithMaxRetries(int(maxRetries)),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create PipeOps API Client",
			"An unexpected error occurred when creating the PipeOps API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"PipeOps Client Error: "+err.Error(),
		)
		return
	}

	// Set the authentication token
	client.SetToken(token)

	// Make the client available to resources and data sources
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *PipeOpsProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewProjectResource,
		NewEnvironmentResource,
		NewServerResource,
	}
}

func (p *PipeOpsProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewProjectDataSource,
	}
}
