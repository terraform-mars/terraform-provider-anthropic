package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/terraform-mars/terraform-provider-anthropic/internal/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &APIKeysDataSource{}

func NewAPIKeysDataSource() datasource.DataSource {
	return &APIKeysDataSource{}
}

// APIKeysDataSource defines the data source implementation.
type APIKeysDataSource struct {
	client *client.Client
}

// APIKeysDataSourceModel describes the data source data model.
type APIKeysDataSourceModel struct {
	WorkspaceID types.String    `tfsdk:"workspace_id"`
	Status      types.String    `tfsdk:"status"`
	APIKeys     []APIKeyModel   `tfsdk:"api_keys"`
}

// APIKeyModel describes a single API key in the list.
type APIKeyModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	WorkspaceID types.String `tfsdk:"workspace_id"`
	Status      types.String `tfsdk:"status"`
	Hint        types.String `tfsdk:"hint"`
	CreatedAt   types.String `tfsdk:"created_at"`
}

func (d *APIKeysDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_api_keys"
}

func (d *APIKeysDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a list of API keys in the Anthropic organization, optionally filtered by workspace or status.",
		Attributes: map[string]schema.Attribute{
			"workspace_id": schema.StringAttribute{
				Description: "Filter API keys by workspace ID.",
				Optional:    true,
			},
			"status": schema.StringAttribute{
				Description: "Filter API keys by status (active, inactive, archived).",
				Optional:    true,
			},
			"api_keys": schema.ListNestedAttribute{
				Description: "List of API keys.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The unique identifier of the API key.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "The name of the API key.",
							Computed:    true,
						},
						"workspace_id": schema.StringAttribute{
							Description: "The ID of the workspace this API key belongs to.",
							Computed:    true,
						},
						"status": schema.StringAttribute{
							Description: "The status of the API key.",
							Computed:    true,
						},
						"hint": schema.StringAttribute{
							Description: "The last 4 characters of the API key.",
							Computed:    true,
						},
						"created_at": schema.StringAttribute{
							Description: "The timestamp when the API key was created.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *APIKeysDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = c
}

func (d *APIKeysDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data APIKeysDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get filter values
	var workspaceID, status string
	if !data.WorkspaceID.IsNull() {
		workspaceID = data.WorkspaceID.ValueString()
	}
	if !data.Status.IsNull() {
		status = data.Status.ValueString()
	}

	// Fetch all API keys with pagination
	var allAPIKeys []client.APIKey
	var afterID string

	for {
		apiKeys, err := d.client.ListAPIKeys(ctx, 100, "", afterID, status, workspaceID)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list API keys: %s", err))
			return
		}

		allAPIKeys = append(allAPIKeys, apiKeys.Data...)

		if !apiKeys.HasMore || apiKeys.LastID == nil {
			break
		}
		afterID = *apiKeys.LastID
	}

	// Convert to model
	data.APIKeys = make([]APIKeyModel, len(allAPIKeys))
	for i, key := range allAPIKeys {
		data.APIKeys[i] = APIKeyModel{
			ID:        types.StringValue(key.ID),
			Name:      types.StringValue(key.Name),
			Status:    types.StringValue(key.Status),
			Hint:      types.StringValue(key.Hint),
			CreatedAt: types.StringValue(key.CreatedAt),
		}
		if key.WorkspaceID != "" {
			data.APIKeys[i].WorkspaceID = types.StringValue(key.WorkspaceID)
		} else {
			data.APIKeys[i].WorkspaceID = types.StringNull()
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
