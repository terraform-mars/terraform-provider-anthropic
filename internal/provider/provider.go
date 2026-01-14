package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/terraform-mars/terraform-provider-anthropic/internal/client"
)

// Ensure AnthropicProvider satisfies various provider interfaces.
var _ provider.Provider = &AnthropicProvider{}

// AnthropicProvider defines the provider implementation.
type AnthropicProvider struct {
	version string
}

// AnthropicProviderModel describes the provider data model.
type AnthropicProviderModel struct {
	AdminKey types.String `tfsdk:"admin_key"`
	BaseURL  types.String `tfsdk:"base_url"`
}

func (p *AnthropicProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "anthropic"
	resp.Version = p.version
}

func (p *AnthropicProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The Anthropic provider allows you to manage Anthropic organization resources such as workspaces, API keys, and members using the Admin API.",
		Attributes: map[string]schema.Attribute{
			"admin_key": schema.StringAttribute{
				Description: "The Anthropic Admin API key. Can also be set via the ANTHROPIC_ADMIN_KEY environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
			"base_url": schema.StringAttribute{
				Description: "The base URL for the Anthropic API. Defaults to https://api.anthropic.com. Can also be set via the ANTHROPIC_BASE_URL environment variable.",
				Optional:    true,
			},
		},
	}
}

func (p *AnthropicProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config AnthropicProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get admin key from config or environment
	adminKey := os.Getenv("ANTHROPIC_ADMIN_KEY")
	if !config.AdminKey.IsNull() {
		adminKey = config.AdminKey.ValueString()
	}

	if adminKey == "" {
		resp.Diagnostics.AddError(
			"Missing Admin Key",
			"The admin_key must be set in the provider configuration or via the ANTHROPIC_ADMIN_KEY environment variable.",
		)
		return
	}

	// Get base URL from config or environment
	baseURL := os.Getenv("ANTHROPIC_BASE_URL")
	if !config.BaseURL.IsNull() {
		baseURL = config.BaseURL.ValueString()
	}

	// Create the client
	c := client.NewClient(adminKey)
	if baseURL != "" {
		c.WithBaseURL(baseURL)
	}

	// Make the client available to data sources and resources
	resp.DataSourceData = c
	resp.ResourceData = c
}

func (p *AnthropicProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewWorkspaceResource,
		NewAPIKeyResource,
		NewWorkspaceMemberResource,
		NewInviteResource,
	}
}

func (p *AnthropicProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewWorkspaceDataSource,
		NewWorkspacesDataSource,
		NewAPIKeyDataSource,
		NewAPIKeysDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &AnthropicProvider{
			version: version,
		}
	}
}
