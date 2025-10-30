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
	_ resource.Resource                = &ServerResource{}
	_ resource.ResourceWithConfigure   = &ServerResource{}
	_ resource.ResourceWithImportState = &ServerResource{}
)

func NewServerResource() resource.Resource {
	return &ServerResource{}
}

type ServerResource struct {
	client *pipeops.Client
}

type ServerResourceModel struct {
	ID         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	ProjectID  types.String `tfsdk:"project_id"`
	ServerType types.String `tfsdk:"server_type"`
	Region     types.String `tfsdk:"region"`
	Size       types.String `tfsdk:"size"`
	Status     types.String `tfsdk:"status"`
	IPAddress  types.String `tfsdk:"ip_address"`
	CreatedAt  types.String `tfsdk:"created_at"`
	UpdatedAt  types.String `tfsdk:"updated_at"`
}

func (r *ServerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server"
}

func (r *ServerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a PipeOps server.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Server ID",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Server name",
				Required:    true,
			},
			"project_id": schema.StringAttribute{
				Description: "Project ID",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"server_type": schema.StringAttribute{
				Description: "Server type",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"region": schema.StringAttribute{
				Description: "Server region",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"size": schema.StringAttribute{
				Description: "Server size/tier",
				Optional:    true,
			},
			"status": schema.StringAttribute{
				Description: "Server status",
				Computed:    true,
			},
			"ip_address": schema.StringAttribute{
				Description: "Server IP address",
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

func (r *ServerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ServerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := &pipeops.CreateServerRequest{
		Name:        plan.Name.ValueString(),
		Provider:    plan.ServerType.ValueString(),
		Region:      plan.Region.ValueString(),
		WorkspaceID: plan.ProjectID.ValueString(), // Using project_id as workspace_id
	}

	if !plan.Size.IsNull() {
		createReq.InstanceType = plan.Size.ValueString()
	}

	server, _, err := r.client.Servers.Create(ctx, createReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating server",
			"Could not create server, unexpected error: "+err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(server.Data.Server.ID)
	plan.Name = types.StringValue(server.Data.Server.Name)
	plan.ProjectID = types.StringValue(server.Data.Server.WorkspaceID)
	if server.Data.Server.Provider != "" {
		plan.ServerType = types.StringValue(server.Data.Server.Provider)
	}
	if server.Data.Server.Region != "" {
		plan.Region = types.StringValue(server.Data.Server.Region)
	}
	if server.Data.Server.Status != "" {
		plan.Status = types.StringValue(server.Data.Server.Status)
	}
	if server.Data.Server.CreatedAt != nil {
		plan.CreatedAt = types.StringValue(server.Data.Server.CreatedAt.String())
	}
	if server.Data.Server.UpdatedAt != nil {
		plan.UpdatedAt = types.StringValue(server.Data.Server.UpdatedAt.String())
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *ServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ServerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	server, httpResp, err := r.client.Servers.Get(ctx, state.ID.ValueString())
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading server",
			"Could not read server ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	state.Name = types.StringValue(server.Data.Server.Name)
	state.ProjectID = types.StringValue(server.Data.Server.WorkspaceID)
	if server.Data.Server.Status != "" {
		state.Status = types.StringValue(server.Data.Server.Status)
	}
	if server.Data.Server.Provider != "" {
		state.ServerType = types.StringValue(server.Data.Server.Provider)
	}
	if server.Data.Server.Region != "" {
		state.Region = types.StringValue(server.Data.Server.Region)
	}
	if server.Data.Server.UpdatedAt != nil {
		state.UpdatedAt = types.StringValue(server.Data.Server.UpdatedAt.String())
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *ServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ServerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Server updates may not be supported in the API currently
	// Just refresh the state
	server, _, err := r.client.Servers.Get(ctx, plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading server for update",
			"Could not read server ID "+plan.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	plan.Name = types.StringValue(server.Data.Server.Name)
	if server.Data.Server.UpdatedAt != nil {
		plan.UpdatedAt = types.StringValue(server.Data.Server.UpdatedAt.String())
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *ServerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ServerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Servers.Delete(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting server",
			"Could not delete server ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}
}

func (r *ServerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
