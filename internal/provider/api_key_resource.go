package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/terraform-mars/terraform-provider-anthropic/internal/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &APIKeyResource{}
var _ resource.ResourceWithImportState = &APIKeyResource{}

func NewAPIKeyResource() resource.Resource {
	return &APIKeyResource{}
}

// APIKeyResource defines the resource implementation.
type APIKeyResource struct {
	client *client.Client
}

// APIKeyResourceModel describes the resource data model.
type APIKeyResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	WorkspaceID types.String `tfsdk:"workspace_id"`
	Status      types.String `tfsdk:"status"`
	Hint        types.String `tfsdk:"hint"`
	Key         types.String `tfsdk:"key"`
	CreatedAt   types.String `tfsdk:"created_at"`
}

func (r *APIKeyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_api_key"
}

func (r *APIKeyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an Anthropic API key. API keys are used to authenticate requests to the Anthropic API and can be scoped to specific workspaces.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the API key.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the API key.",
				Required:    true,
			},
			"workspace_id": schema.StringAttribute{
				Description: "The ID of the workspace this API key belongs to. If not specified, the key is organization-wide.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"status": schema.StringAttribute{
				Description: "The status of the API key (active, inactive).",
				Optional:    true,
				Computed:    true,
			},
			"hint": schema.StringAttribute{
				Description: "The last 4 characters of the API key for identification.",
				Computed:    true,
			},
			"key": schema.StringAttribute{
				Description: "The full API key value. Only available immediately after creation.",
				Computed:    true,
				Sensitive:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_at": schema.StringAttribute{
				Description: "The timestamp when the API key was created.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *APIKeyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = c
}

func (r *APIKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data APIKeyResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	createReq := &client.CreateAPIKeyRequest{
		Name: data.Name.ValueString(),
	}
	if !data.WorkspaceID.IsNull() {
		createReq.WorkspaceID = data.WorkspaceID.ValueString()
	}

	apiKey, err := r.client.CreateAPIKey(ctx, createReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create API key: %s", err))
		return
	}

	data.ID = types.StringValue(apiKey.ID)
	data.Name = types.StringValue(apiKey.Name)
	data.Status = types.StringValue(apiKey.Status)
	data.Hint = types.StringValue(apiKey.Hint)
	data.CreatedAt = types.StringValue(apiKey.CreatedAt)

	// The key is only returned on creation
	if apiKey.Key != "" {
		data.Key = types.StringValue(apiKey.Key)
	} else {
		data.Key = types.StringNull()
	}

	if apiKey.WorkspaceID != "" {
		data.WorkspaceID = types.StringValue(apiKey.WorkspaceID)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *APIKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data APIKeyResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	apiKey, err := r.client.GetAPIKey(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read API key: %s", err))
		return
	}

	data.Name = types.StringValue(apiKey.Name)
	data.Status = types.StringValue(apiKey.Status)
	data.Hint = types.StringValue(apiKey.Hint)
	data.CreatedAt = types.StringValue(apiKey.CreatedAt)

	if apiKey.WorkspaceID != "" {
		data.WorkspaceID = types.StringValue(apiKey.WorkspaceID)
	}

	// The key is not returned on read, preserve existing value
	// data.Key stays as-is from state

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *APIKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data APIKeyResourceModel
	var state APIKeyResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := &client.UpdateAPIKeyRequest{}

	// Check if name changed
	if !data.Name.Equal(state.Name) {
		updateReq.Name = data.Name.ValueString()
	}

	// Check if status changed
	if !data.Status.Equal(state.Status) && !data.Status.IsNull() {
		updateReq.Status = data.Status.ValueString()
	}

	apiKey, err := r.client.UpdateAPIKey(ctx, data.ID.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update API key: %s", err))
		return
	}

	data.Name = types.StringValue(apiKey.Name)
	data.Status = types.StringValue(apiKey.Status)
	data.Hint = types.StringValue(apiKey.Hint)

	// Preserve the key from state since it's not returned on update
	data.Key = state.Key

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *APIKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data APIKeyResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteAPIKey(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete API key: %s", err))
		return
	}
}

func (r *APIKeyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
