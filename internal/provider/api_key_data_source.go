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
var _ datasource.DataSource = &APIKeyDataSource{}

func NewAPIKeyDataSource() datasource.DataSource {
	return &APIKeyDataSource{}
}

// APIKeyDataSource defines the data source implementation.
type APIKeyDataSource struct {
	client *client.Client
}

// APIKeyDataSourceModel describes the data source data model.
type APIKeyDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	WorkspaceID types.String `tfsdk:"workspace_id"`
	Status      types.String `tfsdk:"status"`
	Hint        types.String `tfsdk:"hint"`
	CreatedAt   types.String `tfsdk:"created_at"`
}

func (d *APIKeyDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_api_key"
}

func (d *APIKeyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves information about an existing Anthropic API key.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the API key.",
				Required:    true,
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
				Description: "The status of the API key (active, inactive, archived).",
				Computed:    true,
			},
			"hint": schema.StringAttribute{
				Description: "The last 4 characters of the API key for identification.",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "The timestamp when the API key was created.",
				Computed:    true,
			},
		},
	}
}

func (d *APIKeyDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *APIKeyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data APIKeyDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	apiKey, err := d.client.GetAPIKey(ctx, data.ID.ValueString())
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
	} else {
		data.WorkspaceID = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
