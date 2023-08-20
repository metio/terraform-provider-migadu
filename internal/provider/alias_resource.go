/*
 * SPDX-FileCopyrightText: The terraform-provider-migadu Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/metio/terraform-provider-migadu/internal/provider/custom_types"
	"github.com/metio/terraform-provider-migadu/migadu/client"
	"github.com/metio/terraform-provider-migadu/migadu/model"
	"net/http"
	"strings"
)

var (
	_ resource.Resource                = (*AliasResource)(nil)
	_ resource.ResourceWithConfigure   = (*AliasResource)(nil)
	_ resource.ResourceWithImportState = (*AliasResource)(nil)
)

func NewAliasResource() resource.Resource {
	return &AliasResource{}
}

type AliasResource struct {
	MigaduClient *client.MigaduClient
}

type AliasResourceModel struct {
	ID               custom_types.EmailAddressValue    `tfsdk:"id"`
	LocalPart        types.String                      `tfsdk:"local_part"`
	DomainName       custom_types.DomainNameValue      `tfsdk:"domain_name"`
	Address          custom_types.EmailAddressValue    `tfsdk:"address"`
	Destinations     custom_types.EmailAddressSetValue `tfsdk:"destinations"`
	IsInternal       types.Bool                        `tfsdk:"is_internal"`
	Expirable        types.Bool                        `tfsdk:"expirable"`
	ExpiresOn        types.String                      `tfsdk:"expires_on"`
	RemoveUponExpiry types.Bool                        `tfsdk:"remove_upon_expiry"`
}

func (r *AliasResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_alias"
}

func (r *AliasResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description:         "Provides an email alias.",
		MarkdownDescription: "Provides an email alias.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Contains the value 'local_part@domain_name'.",
				MarkdownDescription: "Contains the value `local_part@domain_name`.",
				Required:            false,
				Optional:            false,
				Computed:            true,
				CustomType:          custom_types.EmailAddressType{},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"local_part": schema.StringAttribute{
				Description:         "The local part of the alias.",
				MarkdownDescription: "The local part of the alias.",
				Required:            true,
				Optional:            false,
				Computed:            false,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"domain_name": schema.StringAttribute{
				Description:         "The domain name of the alias.",
				MarkdownDescription: "The domain name of the alias.",
				Required:            true,
				Optional:            false,
				Computed:            false,
				CustomType:          custom_types.DomainNameType{},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"address": schema.StringAttribute{
				Description:         "The email address 'local_part@domain_name' as returned by the Migadu API. This might be different from the 'id' attribute in case you are using international domain names. The Migadu API always returns the punycode version of a domain.",
				MarkdownDescription: "The email address `local_part@domain_name` as returned by the Migadu API. This might be different from the `id` attribute in case you are using international domain names. The Migadu API always returns the punycode version of a domain.",
				Required:            false,
				Optional:            false,
				Computed:            true,
				CustomType:          custom_types.EmailAddressType{},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"destinations": schema.SetAttribute{
				Description:         "Set of email addresses that act as destinations of the alias.",
				MarkdownDescription: "Set of email addresses that act as destinations of the alias.",
				Required:            true,
				Optional:            false,
				Computed:            false,
				CustomType: custom_types.EmailAddressSetType{
					SetType: types.SetType{
						ElemType: custom_types.EmailAddressType{},
					},
				},
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
			},
			"is_internal": schema.BoolAttribute{
				Description:         "Internal aliases can only receive emails from Migadu email servers.",
				MarkdownDescription: "Internal aliases can only receive emails from Migadu email servers.",
				Required:            false,
				Optional:            true,
				Computed:            true,
			},
			"expirable": schema.BoolAttribute{
				Description:         "Whether this alias expires at some time.",
				MarkdownDescription: "Whether this alias expires at some time.",
				Required:            false,
				Optional:            true,
				Computed:            true,
			},
			"expires_on": schema.StringAttribute{
				Description:         "The expiration date of this alias.",
				MarkdownDescription: "The expiration date of this alias.",
				Required:            false,
				Optional:            true,
				Computed:            true,
			},
			"remove_upon_expiry": schema.BoolAttribute{
				Description:         "Whether to remove this alias upon expiry.",
				MarkdownDescription: "Whether to remove this alias upon expiry.",
				Required:            false,
				Optional:            true,
				Computed:            true,
			},
		},
	}
}
func (r *AliasResource) Configure(_ context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}

	if migaduClient, ok := request.ProviderData.(*client.MigaduClient); ok {
		r.MigaduClient = migaduClient
	} else {
		response.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.MigaduClient, got: %T. Please report this issue to the provider developers.", request.ProviderData),
		)
	}
}

func (r *AliasResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan AliasResourceModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	var destinations []string
	response.Diagnostics.Append(plan.Destinations.ElementsAs(ctx, &destinations, false)...)
	if response.Diagnostics.HasError() {
		return
	}

	alias := &model.Alias{
		LocalPart:        plan.LocalPart.ValueString(),
		Destinations:     destinations,
		IsInternal:       plan.IsInternal.ValueBool(),
		Expirable:        plan.Expirable.ValueBool(),
		ExpiresOn:        plan.ExpiresOn.ValueString(),
		RemoveUponExpiry: plan.RemoveUponExpiry.ValueBool(),
	}

	createdAlias, err := r.MigaduClient.CreateAlias(ctx, plan.DomainName.ValueString(), alias)
	if err != nil {
		response.Diagnostics.Append(AliasCreateError(err))
		return
	}

	plan.ID = custom_types.NewEmailAddressValue(CreateAliasID(plan.LocalPart, plan.DomainName))
	plan.Address = custom_types.NewEmailAddressValue(createdAlias.Address)
	plan.IsInternal = types.BoolValue(createdAlias.IsInternal)
	plan.Expirable = types.BoolValue(createdAlias.Expirable)
	plan.ExpiresOn = types.StringValue(createdAlias.ExpiresOn)
	plan.RemoveUponExpiry = types.BoolValue(createdAlias.RemoveUponExpiry)

	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
}

func (r *AliasResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state AliasResourceModel
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	alias, err := r.MigaduClient.GetAlias(ctx, state.DomainName.ValueString(), state.LocalPart.ValueString())
	if err != nil {
		var requestError *client.RequestError
		if errors.As(err, &requestError) {
			if requestError.StatusCode == http.StatusNotFound {
				response.State.RemoveResource(ctx)
				return
			}
		}
		response.Diagnostics.Append(AliasReadError(err))
		return
	}

	receivedDestinations, diags := custom_types.NewEmailAddressSetValueFrom(ctx, alias.Destinations)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if equal, _ := state.Destinations.SetSemanticEquals(ctx, receivedDestinations); !equal {
		state.Destinations = receivedDestinations
	}

	state.ID = custom_types.NewEmailAddressValue(CreateAliasID(state.LocalPart, state.DomainName))
	state.Address = custom_types.NewEmailAddressValue(alias.Address)
	state.IsInternal = types.BoolValue(alias.IsInternal)
	state.Expirable = types.BoolValue(alias.Expirable)
	state.ExpiresOn = types.StringValue(alias.ExpiresOn)
	state.RemoveUponExpiry = types.BoolValue(alias.RemoveUponExpiry)

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (r *AliasResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan AliasResourceModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	var destinations []string
	response.Diagnostics.Append(plan.Destinations.ElementsAs(ctx, &destinations, false)...)
	if response.Diagnostics.HasError() {
		return
	}

	alias := &model.Alias{
		LocalPart:        plan.LocalPart.ValueString(),
		Destinations:     destinations,
		IsInternal:       plan.IsInternal.ValueBool(),
		Expirable:        plan.Expirable.ValueBool(),
		ExpiresOn:        plan.ExpiresOn.ValueString(),
		RemoveUponExpiry: plan.RemoveUponExpiry.ValueBool(),
	}

	updatedAlias, err := r.MigaduClient.UpdateAlias(ctx, plan.DomainName.ValueString(), plan.LocalPart.ValueString(), alias)
	if err != nil {
		response.Diagnostics.Append(AliasUpdateError(err))
		return
	}

	plan.ID = custom_types.NewEmailAddressValue(CreateAliasID(plan.LocalPart, plan.DomainName))
	plan.Address = custom_types.NewEmailAddressValue(updatedAlias.Address)
	plan.IsInternal = types.BoolValue(updatedAlias.IsInternal)
	plan.Expirable = types.BoolValue(updatedAlias.Expirable)
	plan.ExpiresOn = types.StringValue(updatedAlias.ExpiresOn)
	plan.RemoveUponExpiry = types.BoolValue(updatedAlias.RemoveUponExpiry)

	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
}

func (r *AliasResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state AliasResourceModel
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	_, err := r.MigaduClient.DeleteAlias(ctx, state.DomainName.ValueString(), state.LocalPart.ValueString())
	if err != nil {
		response.Diagnostics.Append(AliasDeleteError(err))
		return
	}
}

func (r *AliasResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	idParts := strings.Split(request.ID, "@")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		response.Diagnostics.Append(AliasImportError(request.ID))
		return
	}

	localPart := idParts[0]
	domainName := idParts[1]
	tflog.Trace(ctx, "parsed import ID", map[string]interface{}{
		"local_part":  localPart,
		"domain_name": domainName,
	})

	resource.ImportStatePassthroughID(ctx, path.Root("id"), request, response)

	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("local_part"), localPart)...)
	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("domain_name"), domainName)...)
}
