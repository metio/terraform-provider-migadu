package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/metio/terraform-provider-migadu/internal/client"
	"golang.org/x/net/idna"
	"strings"
)

var (
	_ resource.Resource                = &aliasResource{}
	_ resource.ResourceWithConfigure   = &aliasResource{}
	_ resource.ResourceWithImportState = &aliasResource{}
)

func NewAliasResource() resource.Resource {
	return &aliasResource{}
}

type aliasResource struct {
	migaduClient *client.MigaduClient
}

type aliasResourceModel struct {
	ID               types.String `tfsdk:"id"`
	LocalPart        types.String `tfsdk:"local_part"`
	DomainName       types.String `tfsdk:"domain_name"`
	Address          types.String `tfsdk:"address"`
	Destinations     types.List   `tfsdk:"destinations"`
	DestinationsIDN  types.List   `tfsdk:"destinations_idn"`
	IsInternal       types.Bool   `tfsdk:"is_internal"`
	Expirable        types.Bool   `tfsdk:"expirable"`
	ExpiresOn        types.String `tfsdk:"expires_on"`
	RemoveUponExpiry types.Bool   `tfsdk:"remove_upon_expiry"`
}

func (r *aliasResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_alias"
}

func (r *aliasResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manage a single alias.",
		MarkdownDescription: "Manage a single alias.",
		Attributes: map[string]schema.Attribute{
			"domain_name": schema.StringAttribute{
				Description:         "The domain name of the alias to manage.",
				MarkdownDescription: "The domain name of the alias to manage.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"local_part": schema.StringAttribute{
				Description:         "The local part of the alias to manage.",
				MarkdownDescription: "The local part of the alias to manage.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				Description:         "Contains the full email address 'local_part@domain_name'.",
				MarkdownDescription: "Contains the full email address `local_part@domain_name`.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"address": schema.StringAttribute{
				Description:         "Contains the full email address 'local_part@domain_name' as returned by the Migadu API. This might be different from the 'id' attribute in case you are using international domain names. The Migadu API always returns the punycode version of a domain.",
				MarkdownDescription: "Contains the full email address `local_part@domain_name` as returned by the Migadu API. This might be different from the `id` attribute in case you are using international domain names. The Migadu API always returns the punycode version of a domain.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"destinations": schema.ListAttribute{
				Description:         "List of email addresses that act as destinations of the alias.",
				MarkdownDescription: "List of email addresses that act as destinations of the alias.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				Validators: []validator.List{
					listvalidator.ConflictsWith(path.MatchRoot("destinations_idn")),
					listvalidator.SizeAtLeast(1),
				},
			},
			"destinations_idn": schema.ListAttribute{
				Description:         "List of email addresses that act as destinations of the alias. Use this attribute instead of 'destinations' in case you are dealing with international domain names (IDN).",
				MarkdownDescription: "List of email addresses that act as destinations of the alias. Use this attribute instead of `destinations` in case you are dealing with international domain names (IDN).",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				Validators: []validator.List{
					listvalidator.ConflictsWith(path.MatchRoot("destinations")),
					listvalidator.SizeAtLeast(1),
				},
			},
			"is_internal": schema.BoolAttribute{
				Description:         "Internal aliases can only receive emails from Migadu email servers.",
				MarkdownDescription: "Internal aliases can only receive emails from Migadu email servers.",
				Optional:            true,
				Computed:            true,
			},
			"expirable": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"expires_on": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"remove_upon_expiry": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
		},
	}
}
func (r *aliasResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	migaduClient, ok := req.ProviderData.(*client.MigaduClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.MigaduClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.migaduClient = migaduClient
}

func (r *aliasResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan aliasResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var destinations []string
	if !plan.Destinations.IsUnknown() {
		resp.Diagnostics.Append(plan.Destinations.ElementsAs(ctx, &destinations, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	if !plan.DestinationsIDN.IsUnknown() {
		resp.Diagnostics.Append(plan.DestinationsIDN.ElementsAs(ctx, &destinations, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	alias := &client.Alias{
		LocalPart:        plan.LocalPart.ValueString(),
		Destinations:     destinations,
		IsInternal:       plan.IsInternal.ValueBool(),
		Expirable:        plan.Expirable.ValueBool(),
		ExpiresOn:        plan.ExpiresOn.ValueString(),
		RemoveUponExpiry: plan.RemoveUponExpiry.ValueBool(),
	}

	createdAlias, err := r.migaduClient.CreateAlias(ctx, plan.DomainName.ValueString(), alias)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating alias",
			"Could not create alias, unexpected error: "+err.Error(),
		)
		return
	}

	receivedDestinations, diags := types.ListValueFrom(ctx, types.StringType, createdAlias.Destinations)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	receivedDestinationsIDN, diags := types.ListValueFrom(ctx, types.StringType, convertToUnicode(createdAlias.Destinations, &resp.Diagnostics))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.Destinations = receivedDestinations
	plan.DestinationsIDN = receivedDestinationsIDN
	plan.ID = types.StringValue(fmt.Sprintf("%s@%s", plan.LocalPart.ValueString(), plan.DomainName.ValueString()))
	plan.Address = types.StringValue(createdAlias.Address)
	plan.IsInternal = types.BoolValue(createdAlias.IsInternal)
	plan.Expirable = types.BoolValue(createdAlias.Expirable)
	plan.ExpiresOn = types.StringValue(createdAlias.ExpiresOn)
	plan.RemoveUponExpiry = types.BoolValue(createdAlias.RemoveUponExpiry)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *aliasResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state aliasResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	alias, err := r.migaduClient.GetAlias(ctx, state.DomainName.ValueString(), state.LocalPart.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading alias",
			fmt.Sprintf("Could not read alias %s@%s: %v", state.LocalPart.ValueString(), state.DomainName.ValueString(), err),
		)
		return
	}

	destinations, diags := types.ListValueFrom(ctx, types.StringType, alias.Destinations)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	destinationsIDN, diags := types.ListValueFrom(ctx, types.StringType, convertToUnicode(alias.Destinations, &resp.Diagnostics))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.Destinations = destinations
	state.DestinationsIDN = destinationsIDN
	state.ID = types.StringValue(fmt.Sprintf("%s@%s", state.LocalPart.ValueString(), state.DomainName.ValueString()))
	state.Address = types.StringValue(alias.Address)
	state.IsInternal = types.BoolValue(alias.IsInternal)
	state.Expirable = types.BoolValue(alias.Expirable)
	state.ExpiresOn = types.StringValue(alias.ExpiresOn)
	state.RemoveUponExpiry = types.BoolValue(alias.RemoveUponExpiry)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func convertToUnicode(destinations []string, diag *diag.Diagnostics) []string {
	var asciiDestinations []string
	for _, dest := range destinations {
		parts := strings.Split(dest, "@")

		ascii, punyErr := idna.ToUnicode(parts[1])
		if punyErr == nil {
			asciiDestinations = append(asciiDestinations, fmt.Sprintf("%s@%s", parts[0], ascii))
		} else {
			diag.AddError("Could not convert to unicode", dest)
		}
	}
	return asciiDestinations
}

func (r *aliasResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan aliasResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var destinations []string
	if !plan.Destinations.IsUnknown() {
		resp.Diagnostics.Append(plan.Destinations.ElementsAs(ctx, &destinations, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	if !plan.DestinationsIDN.IsUnknown() {
		resp.Diagnostics.Append(plan.DestinationsIDN.ElementsAs(ctx, &destinations, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	alias := &client.Alias{
		LocalPart:        plan.LocalPart.ValueString(),
		Destinations:     destinations,
		IsInternal:       plan.IsInternal.ValueBool(),
		Expirable:        plan.Expirable.ValueBool(),
		ExpiresOn:        plan.ExpiresOn.ValueString(),
		RemoveUponExpiry: plan.RemoveUponExpiry.ValueBool(),
	}

	updatedAlias, err := r.migaduClient.UpdateAlias(ctx, plan.DomainName.ValueString(), plan.LocalPart.ValueString(), alias)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating alias",
			"Could not update alias "+plan.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	receivedDestinations, diags := types.ListValueFrom(ctx, types.StringType, updatedAlias.Destinations)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	receivedDestinationsIDN, diags := types.ListValueFrom(ctx, types.StringType, convertToUnicode(updatedAlias.Destinations, &resp.Diagnostics))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.Destinations = receivedDestinations
	plan.DestinationsIDN = receivedDestinationsIDN
	plan.ID = types.StringValue(fmt.Sprintf("%s@%s", plan.LocalPart.ValueString(), plan.DomainName.ValueString()))
	plan.Address = types.StringValue(updatedAlias.Address)
	plan.IsInternal = types.BoolValue(updatedAlias.IsInternal)
	plan.Expirable = types.BoolValue(updatedAlias.Expirable)
	plan.ExpiresOn = types.StringValue(updatedAlias.ExpiresOn)
	plan.RemoveUponExpiry = types.BoolValue(updatedAlias.RemoveUponExpiry)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *aliasResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state aliasResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.migaduClient.DeleteAlias(ctx, state.DomainName.ValueString(), state.LocalPart.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting alias",
			"Could not delete alias "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}
}

func (r *aliasResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, "@")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Error importing alias",
			fmt.Sprintf("Expected import identifier with format: 'local_part@domain_name' Got: '%q'", req.ID),
		)
		return
	}

	localPart := idParts[0]
	domainName := idParts[1]
	tflog.Trace(ctx, "parsed import ID", map[string]interface{}{
		"local_part":  localPart,
		"domain_name": domainName,
	})

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	resource.ImportStatePassthroughID(ctx, path.Root("address"), req, resp)

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("local_part"), localPart)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("domain_name"), domainName)...)
}
