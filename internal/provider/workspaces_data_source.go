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
var _ datasource.DataSource = &WorkspacesDataSource{}

func NewWorkspacesDataSource() datasource.DataSource {
	return &WorkspacesDataSource{}
}

// WorkspacesDataSource defines the data source implementation.
type WorkspacesDataSource struct {
	client *client.Client
}

// WorkspacesDataSourceModel describes the data source data model.
type WorkspacesDataSourceModel struct {
	Workspaces []WorkspaceModel `tfsdk:"workspaces"`
}

// WorkspaceModel describes a single workspace in the list.
type WorkspaceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	DisplayName types.String `tfsdk:"display_name"`
	CreatedAt   types.String `tfsdk:"created_at"`
	ArchivedAt  types.String `tfsdk:"archived_at"`
}

func (d *WorkspacesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workspaces"
}

func (d *WorkspacesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a list of all workspaces in the Anthropic organization.",
		Attributes: map[string]schema.Attribute{
			"workspaces": schema.ListNestedAttribute{
				Description: "List of workspaces.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The unique identifier of the workspace.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "The name of the workspace.",
							Computed:    true,
						},
						"display_name": schema.StringAttribute{
							Description: "The display name of the workspace.",
							Computed:    true,
						},
						"created_at": schema.StringAttribute{
							Description: "The timestamp when the workspace was created.",
							Computed:    true,
						},
						"archived_at": schema.StringAttribute{
							Description: "The timestamp when the workspace was archived, if applicable.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *WorkspacesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *WorkspacesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data WorkspacesDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Fetch all workspaces with pagination
	var allWorkspaces []client.Workspace
	var afterID string

	for {
		workspaces, err := d.client.ListWorkspaces(ctx, 100, "", afterID)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list workspaces: %s", err))
			return
		}

		allWorkspaces = append(allWorkspaces, workspaces.Data...)

		if !workspaces.HasMore || workspaces.LastID == nil {
			break
		}
		afterID = *workspaces.LastID
	}

	// Convert to model
	data.Workspaces = make([]WorkspaceModel, len(allWorkspaces))
	for i, ws := range allWorkspaces {
		data.Workspaces[i] = WorkspaceModel{
			ID:          types.StringValue(ws.ID),
			Name:        types.StringValue(ws.Name),
			DisplayName: types.StringValue(ws.DisplayName),
			CreatedAt:   types.StringValue(ws.CreatedAt),
		}
		if ws.ArchivedAt != "" {
			data.Workspaces[i].ArchivedAt = types.StringValue(ws.ArchivedAt)
		} else {
			data.Workspaces[i].ArchivedAt = types.StringNull()
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
