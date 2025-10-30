package resources

import (
	"context"
	"fmt"

	"github.com/PipeOpsHQ/pipeops-go-sdk/pipeops"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &ProjectResource{}
	_ resource.ResourceWithConfigure   = &ProjectResource{}
	_ resource.ResourceWithImportState = &ProjectResource{}
)

func NewProjectResource() resource.Resource {
	return &ProjectResource{}
}

type ProjectResource struct {
	client *pipeops.Client
}

type ProjectResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	WorkspaceID types.String `tfsdk:"workspace_id"`
	RepoURL     types.String `tfsdk:"repo_url"`
	RepoBranch  types.String `tfsdk:"repo_branch"`
	CreatedAt   types.String `tfsdk:"created_at"`
	UpdatedAt   types.String `tfsdk:"updated_at"`
}

func (r *ProjectResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

func (r *ProjectResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a PipeOps project.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Project ID",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Project name",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "Project description",
				Optional:    true,
			},
			"workspace_id": schema.StringAttribute{
				Description: "Workspace ID",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"repo_url": schema.StringAttribute{
				Description: "Repository URL",
				Optional:    true,
			},
			"repo_branch": schema.StringAttribute{
				Description: "Repository branch",
				Optional:    true,
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

func (r *ProjectResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*pipeops.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *pipeops.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *ProjectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ProjectResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := &pipeops.CreateProjectRequest{
		Name:          plan.Name.ValueString(),
		ServerID:      "", // Will need server_id in real usage
		EnvironmentID: "", // Will need environment_id in real usage
		Repository:    plan.RepoURL.ValueString(),
		Branch:        plan.RepoBranch.ValueString(),
	}

	if !plan.Description.IsNull() {
		createReq.Description = plan.Description.ValueString()
	}

	project, _, err := r.client.Projects.Create(ctx, createReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating project",
			"Could not create project, unexpected error: "+err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(project.Data.Project.ID)
	plan.Name = types.StringValue(project.Data.Project.Name)
	plan.WorkspaceID = types.StringValue(project.Data.Project.WorkspaceID)
	if project.Data.Project.Description != "" {
		plan.Description = types.StringValue(project.Data.Project.Description)
	}
	if project.Data.Project.Repository != "" {
		plan.RepoURL = types.StringValue(project.Data.Project.Repository)
	}
	if project.Data.Project.Branch != "" {
		plan.RepoBranch = types.StringValue(project.Data.Project.Branch)
	}
	if project.Data.Project.CreatedAt != nil {
		plan.CreatedAt = types.StringValue(project.Data.Project.CreatedAt.String())
	}
	if project.Data.Project.UpdatedAt != nil {
		plan.UpdatedAt = types.StringValue(project.Data.Project.UpdatedAt.String())
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *ProjectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ProjectResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	project, httpResp, err := r.client.Projects.Get(ctx, state.ID.ValueString())
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading project",
			"Could not read project ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	state.Name = types.StringValue(project.Data.Project.Name)
	state.WorkspaceID = types.StringValue(project.Data.Project.WorkspaceID)
	if project.Data.Project.Description != "" {
		state.Description = types.StringValue(project.Data.Project.Description)
	}
	if project.Data.Project.Repository != "" {
		state.RepoURL = types.StringValue(project.Data.Project.Repository)
	}
	if project.Data.Project.Branch != "" {
		state.RepoBranch = types.StringValue(project.Data.Project.Branch)
	}
	if project.Data.Project.CreatedAt != nil {
		state.CreatedAt = types.StringValue(project.Data.Project.CreatedAt.String())
	}
	if project.Data.Project.UpdatedAt != nil {
		state.UpdatedAt = types.StringValue(project.Data.Project.UpdatedAt.String())
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *ProjectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ProjectResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := &pipeops.UpdateProjectRequest{
		Name: plan.Name.ValueString(),
	}

	if !plan.Description.IsNull() {
		updateReq.Description = plan.Description.ValueString()
	}

	project, _, err := r.client.Projects.Update(ctx, plan.ID.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating project",
			"Could not update project ID "+plan.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	plan.Name = types.StringValue(project.Data.Project.Name)
	if project.Data.Project.Description != "" {
		plan.Description = types.StringValue(project.Data.Project.Description)
	}
	if project.Data.Project.UpdatedAt != nil {
		plan.UpdatedAt = types.StringValue(project.Data.Project.UpdatedAt.String())
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *ProjectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ProjectResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Projects.Delete(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting project",
			"Could not delete project ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}
}

func (r *ProjectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
