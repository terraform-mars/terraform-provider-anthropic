package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/terraform-mars/terraform-provider-anthropic/internal/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &InviteResource{}
var _ resource.ResourceWithImportState = &InviteResource{}

func NewInviteResource() resource.Resource {
	return &InviteResource{}
}

// InviteResource defines the resource implementation.
type InviteResource struct {
	client *client.Client
}

// InviteResourceModel describes the resource data model.
type InviteResourceModel struct {
	ID        types.String `tfsdk:"id"`
	Email     types.String `tfsdk:"email"`
	Role      types.String `tfsdk:"role"`
	Status    types.String `tfsdk:"status"`
	CreatedAt types.String `tfsdk:"created_at"`
	ExpiresAt types.String `tfsdk:"expires_at"`
	InviterID types.String `tfsdk:"inviter_id"`
}

func (r *InviteResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_invite"
}

func (r *InviteResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an invitation to join the Anthropic organization. Invites allow you to add new users to your organization.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the invite.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"email": schema.StringAttribute{
				Description: "The email address to send the invitation to.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"role": schema.StringAttribute{
				Description: "The role to assign to the invited user. Valid values: user, admin, developer.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("user", "admin", "developer"),
				},
			},
			"status": schema.StringAttribute{
				Description: "The status of the invite (pending, accepted, expired, deleted).",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "The timestamp when the invite was created.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"expires_at": schema.StringAttribute{
				Description: "The timestamp when the invite expires.",
				Computed:    true,
			},
			"inviter_id": schema.StringAttribute{
				Description: "The ID of the user who created the invite.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *InviteResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *InviteResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data InviteResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	invite, err := r.client.CreateInvite(ctx, &client.CreateInviteRequest{
		Email: data.Email.ValueString(),
		Role:  data.Role.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create invite: %s", err))
		return
	}

	data.ID = types.StringValue(invite.ID)
	data.Email = types.StringValue(invite.Email)
	data.Role = types.StringValue(invite.Role)
	data.Status = types.StringValue(invite.Status)
	data.CreatedAt = types.StringValue(invite.CreatedAt)
	data.ExpiresAt = types.StringValue(invite.ExpiresAt)

	if invite.InviterID != "" {
		data.InviterID = types.StringValue(invite.InviterID)
	} else {
		data.InviterID = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *InviteResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data InviteResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	invite, err := r.client.GetInvite(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read invite: %s", err))
		return
	}

	data.Email = types.StringValue(invite.Email)
	data.Role = types.StringValue(invite.Role)
	data.Status = types.StringValue(invite.Status)
	data.ExpiresAt = types.StringValue(invite.ExpiresAt)

	if invite.InviterID != "" {
		data.InviterID = types.StringValue(invite.InviterID)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *InviteResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Invites cannot be updated, only created or deleted
	resp.Diagnostics.AddError(
		"Update Not Supported",
		"Invites cannot be updated. Delete and recreate the invite with the new settings.",
	)
}

func (r *InviteResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data InviteResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteInvite(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete invite: %s", err))
		return
	}
}

func (r *InviteResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
