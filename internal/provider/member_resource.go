package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-bitwarden/internal/bitwarden"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &memberResource{}
	_ resource.ResourceWithConfigure   = &memberResource{}
	_ resource.ResourceWithImportState = &memberResource{}
)

// NewMemberResource is a helper function to simplify the provider implementation.
func NewMemberResource() resource.Resource {
	return &memberResource{}
}

// memberResource is the resource implementation.
type memberResource struct {
	client *bitwarden.Client
}

type memberResourceModel struct {
	Type       types.Int64  `tfsdk:"type"`
	AccessAll  types.Bool   `tfsdk:"access_all"`
	ExternalId types.String `tfsdk:"external_id"`
	Email      types.String `tfsdk:"email"`

	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Status      types.Int64  `tfsdk:"status"`
	LastUpdated types.String `tfsdk:"last_updated"`
}

// Metadata returns the resource type name.
func (r *memberResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_member"
}

func (r *memberResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(bitwarden.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *bitwarden.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = &client
}

// Schema defines the schema for the resource.
func (r *memberResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The Bitwarden member resources manage the members aka users within an Bitwarden organization. We leverage the public [Bitwarden API](https://bitwarden.com/help/api/]",
		Attributes: map[string]schema.Attribute{
			"type": schema.Int64Attribute{
				Required:    true,
				Description: "Defines Member type, needs to be one of:\n    Owner = 0,\n    Admin = 1,\n    User = 2,\n    Manager = 3,\n    Custom = 4.\n    See https://github.com/bitwarden/server/blob/master/src/Core/Enums/OrganizationUserType.cs",
				Validators: []validator.Int64{
					// TODO: check if we can switch to custom type later for easier autocomplete in TF
					int64validator.OneOf([]int64{0, 1, 2, 3, 4}...),
				},
			},
			"access_all": schema.BoolAttribute{
				Computed:    true,
				Optional:    true,
				Description: "Determines if this member can access all collections within the organization, or only the associated collections. If set to {true}, this option overrides any collection assignments",
				Default:     booldefault.StaticBool(false),
			},
			"external_id": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "External identifier for reference or linking this member to another system, such as a user directory",
				Default:     stringdefault.StaticString(""),
			},
			"email": schema.StringAttribute{
				Required:    true,
				Description: "The member's email address.",
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The member's unique identifier within the organization",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "The member's name, set from their user account profile",
			},
			"status": schema.Int64Attribute{
				Computed:    true,
				Description: "The member's status within the organisation, is one of the following:\n    Invited = 0,\n    Accepted = 1,\n    Confirmed = 2,\n    Revoked = -1.\n    See https://github.com/bitwarden/server/blob/master/src/Core/Enums/OrganizationUserStatusType.cs",
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *memberResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan memberResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	member := bitwarden.Member{
		Type:       bitwarden.OrganizationUserType(plan.Type.ValueInt64()),
		AccessAll:  plan.AccessAll.ValueBool(),
		ExternalId: plan.ExternalId.ValueString(),
		Email:      plan.Email.ValueString(),
	}

	// Create new member
	newMember, err := (*r.client).CreateMember(ctx, member)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating member",
			"Could not create member, unexpected error: "+err.Error(),
		)
		return
	}

	// TODO: check how we can keep some items hidden from user. Like for example ID, that field is automatically configured
	plan.Type = types.Int64Value(int64(newMember.Type))
	plan.AccessAll = types.BoolValue(newMember.AccessAll)
	plan.ExternalId = types.StringValue(newMember.ExternalId)
	plan.Email = types.StringValue(newMember.Email)

	plan.ID = types.StringValue(newMember.ID)
	plan.Name = types.StringValue(newMember.Name)
	plan.Status = types.Int64Value(newMember.Status)

	plan.LastUpdated = types.StringValue(time.Now().UTC().Format(time.RFC850))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *memberResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state memberResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed member value from BitWarden
	member, err := (*r.client).GetMember(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Bitwarden member",
			"Could not read Bitwarden member ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Overwrite member with refreshed state
	state.Type = types.Int64Value(int64(member.Type))
	state.AccessAll = types.BoolValue(member.AccessAll)
	state.ExternalId = types.StringValue(member.ExternalId)
	state.Email = types.StringValue(member.Email)

	state.ID = types.StringValue(member.ID)
	state.Name = types.StringValue(member.Name)
	state.Status = types.Int64Value(member.Status)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *memberResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan memberResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	member := bitwarden.Member{
		Type:       bitwarden.OrganizationUserType(plan.Type.ValueInt64()),
		AccessAll:  plan.AccessAll.ValueBool(),
		ExternalId: plan.ExternalId.ValueString(),
		Email:      plan.Email.ValueString(),
	}

	// Update existing member
	newMember, err := (*r.client).UpdateMember(ctx, plan.ID.ValueString(), member)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Bitwarden member",
			"Could not update member, unexpected error: "+err.Error(),
		)
		return
	}

	// Update resource state with updated items and timestamp
	plan.Type = types.Int64Value(int64(newMember.Type))
	plan.AccessAll = types.BoolValue(newMember.AccessAll)
	plan.ExternalId = types.StringValue(newMember.ExternalId)
	plan.Email = types.StringValue(newMember.Email)

	plan.ID = types.StringValue(newMember.ID)
	plan.Name = types.StringValue(newMember.Name)
	plan.Status = types.Int64Value(newMember.Status)

	plan.LastUpdated = types.StringValue(time.Now().UTC().Format(time.RFC850))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *memberResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state memberResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing member
	err := (*r.client).DeleteMember(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Bitwarden member",
			"Could not delete member, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *memberResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
