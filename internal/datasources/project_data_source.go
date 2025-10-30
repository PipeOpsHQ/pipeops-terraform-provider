package datasources

import (
	"context"
	"fmt"

	"github.com/PipeOpsHQ/pipeops-go-sdk/pipeops"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ProjectDataSource{}
	_ datasource.DataSourceWithConfigure = &ProjectDataSource{}
)

func NewProjectDataSource() datasource.DataSource {
	return &ProjectDataSource{}
}

type ProjectDataSource struct {
	client *pipeops.Client
}

type ProjectDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	WorkspaceID types.String `tfsdk:"workspace_id"`
	RepoURL     types.String `tfsdk:"repo_url"`
	RepoBranch  types.String `tfsdk:"repo_branch"`
	CreatedAt   types.String `tfsdk:"created_at"`
	UpdatedAt   types.String `tfsdk:"updated_at"`
}

func (d *ProjectDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

func (d *ProjectDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches information about a PipeOps project.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Project ID",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "Project name",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "Project description",
				Computed:    true,
			},
			"workspace_id": schema.StringAttribute{
				Description: "Workspace ID",
				Computed:    true,
			},
			"repo_url": schema.StringAttribute{
				Description: "Repository URL",
				Computed:    true,
			},
			"repo_branch": schema.StringAttribute{
				Description: "Repository branch",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "Creation timestamp",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "Last update timestamp",
				Computed:    true,
			},
		},
	}
}

func (d *ProjectDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*pipeops.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *pipeops.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *ProjectDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ProjectDataSourceModel

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	project, _, err := d.client.Projects.Get(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading project",
			"Could not read project ID "+data.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	data.Name = types.StringValue(project.Data.Project.Name)
	data.WorkspaceID = types.StringValue(project.Data.Project.WorkspaceID)

	if project.Data.Project.Description != "" {
		data.Description = types.StringValue(project.Data.Project.Description)
	}
	if project.Data.Project.Repository != "" {
		data.RepoURL = types.StringValue(project.Data.Project.Repository)
	}
	if project.Data.Project.Branch != "" {
		data.RepoBranch = types.StringValue(project.Data.Project.Branch)
	}
	if project.Data.Project.CreatedAt != nil {
		data.CreatedAt = types.StringValue(project.Data.Project.CreatedAt.String())
	}
	if project.Data.Project.UpdatedAt != nil {
		data.UpdatedAt = types.StringValue(project.Data.Project.UpdatedAt.String())
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
