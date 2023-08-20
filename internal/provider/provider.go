/*
 * SPDX-FileCopyrightText: The terraform-provider-migadu Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/metio/terraform-provider-migadu/migadu/client"
	"os"
	"strconv"
	"time"
)

var (
	_ provider.Provider = (*MigaduProvider)(nil)
)

type MigaduProvider struct{}

type MigaduProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	Token    types.String `tfsdk:"token"`
	Username types.String `tfsdk:"username"`
	Timeout  types.Int64  `tfsdk:"timeout"`
}

func New() provider.Provider {
	return &MigaduProvider{}
}

func (p *MigaduProvider) Metadata(_ context.Context, _ provider.MetadataRequest, response *provider.MetadataResponse) {
	response.TypeName = "migadu"
}

func (p *MigaduProvider) Schema(_ context.Context, _ provider.SchemaRequest, response *provider.SchemaResponse) {
	response.Schema = schema.Schema{
		Description:         "Provider for the Migadu API. Requires Terraform 1.0 or later.",
		MarkdownDescription: "Provider for the [Migadu](https://www.migadu.com/api/) API. Requires Terraform 1.0 or later.",
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				Description:         "The API endpoint to use. Can be specified with the 'MIGADU_ENDPOINT' environment variable. Defaults to 'https://api.migadu.com/v1/'. Take a look at https://www.migadu.com/api/#api-requests for more information.",
				MarkdownDescription: "The API endpoint to use. Can be specified with the `MIGADU_ENDPOINT` environment variable. Defaults to `https://api.migadu.com/v1/`. Take a look at https://www.migadu.com/api/#api-requests for more information.",
				Optional:            true,
			},
			"token": schema.StringAttribute{
				Description:         "The API key to use. Can be specified with the 'MIGADU_TOKEN' environment variable. Take a look at https://www.migadu.com/api/#api-keys for more information.",
				MarkdownDescription: "The API key to use. Can be specified with the `MIGADU_TOKEN` environment variable. Take a look at https://www.migadu.com/api/#api-keys for more information.",
				Optional:            true,
				Sensitive:           true,
			},
			"username": schema.StringAttribute{
				Description:         "The username to use. Can be specified with the 'MIGADU_USERNAME' environment variable. Take a look at https://www.migadu.com/api/#api-requests for more information.",
				MarkdownDescription: "The username to use. Can be specified with the `MIGADU_USERNAME` environment variable. Take a look at https://www.migadu.com/api/#api-requests for more information.",
				Optional:            true,
				Sensitive:           true,
			},
			"timeout": schema.Int64Attribute{
				Description:         "The timeout to apply for HTTP requests in seconds. Can be specified with the 'MIGADU_TIMEOUT' environment variable. Defaults to '10'.",
				MarkdownDescription: "The timeout to apply for HTTP requests in seconds. Can be specified with the `MIGADU_TIMEOUT` environment variable. Defaults to `10`.",
				Optional:            true,
			},
		},
	}
}

func (p *MigaduProvider) Configure(ctx context.Context, request provider.ConfigureRequest, response *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Migadu client")

	var config MigaduProviderModel
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}

	if config.Endpoint.IsUnknown() {
		response.Diagnostics.AddAttributeError(
			path.Root("endpoint"),
			"Unknown Migadu API Endpoint",
			"The provider cannot create the Migadu API client as there is an unknown configuration value for the Migadu API endpoint. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the MIGADU_ENDPOINT environment variable.",
		)
	}

	if config.Username.IsUnknown() {
		response.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Unknown Migadu API Username",
			"The provider cannot create the Migadu API client as there is an unknown configuration value for the Migadu API username. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the MIGADU_USERNAME environment variable.",
		)
	}

	if config.Token.IsUnknown() {
		response.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Unknown Migadu API Token",
			"The provider cannot create the Migadu API client as there is an unknown configuration value for the Migadu API token. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the MIGADU_TOKEN environment variable.",
		)
	}

	if config.Timeout.IsUnknown() {
		response.Diagnostics.AddAttributeError(
			path.Root("timeout"),
			"Unknown Migadu API Timeout",
			"The provider cannot create the Migadu API client as there is an unknown configuration value for the Migadu API timeout. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the MIGADU_TIMEOUT environment variable.",
		)
	}

	if response.Diagnostics.HasError() {
		return
	}

	endpoint := os.Getenv("MIGADU_ENDPOINT")
	username := os.Getenv("MIGADU_USERNAME")
	token := os.Getenv("MIGADU_TOKEN")
	timeout := os.Getenv("MIGADU_TIMEOUT")

	if !config.Endpoint.IsNull() {
		endpoint = config.Endpoint.ValueString()
	}

	if !config.Username.IsNull() {
		username = config.Username.ValueString()
	}

	if !config.Token.IsNull() {
		token = config.Token.ValueString()
	}

	if !config.Timeout.IsNull() {
		timeout = strconv.FormatInt(config.Timeout.ValueInt64(), 10)
	}

	if endpoint == "" {
		endpoint = "https://api.migadu.com/v1/"
	}

	if timeout == "" {
		timeout = "10"
	}

	if username == "" {
		response.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing Migadu API Username",
			"The provider cannot create the Migadu API client as there is a missing or empty value for the Migadu API username. "+
				"Set the username value in the configuration or use the MIGADU_USERNAME environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if token == "" {
		response.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Missing Migadu API Token",
			"The provider cannot create the Migadu API client as there is a missing or empty value for the Migadu API token. "+
				"Set the password value in the configuration or use the MIGADU_TOKEN environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	duration, err := time.ParseDuration(fmt.Sprintf("%ss", timeout))
	if err != nil {
		response.Diagnostics.AddAttributeError(
			path.Root("timeout"),
			"Invalid Migadu API Timeout",
			"The supplied timeout value cannot be parsed into a duration: "+err.Error(),
		)
	}

	if response.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "migadu_endpoint", endpoint)
	ctx = tflog.SetField(ctx, "migadu_username", username)
	ctx = tflog.SetField(ctx, "migadu_token", token)
	ctx = tflog.SetField(ctx, "migadu_timeout", timeout)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "migadu_username")
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "migadu_token")

	tflog.Debug(ctx, "Creating Migadu client")

	c, err := client.New(&endpoint, &username, &token, duration)
	if err != nil {
		response.Diagnostics.AddError(
			"Unable to Create Migadu API Client",
			"An unexpected error occurred when creating the Migadu API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Migadu Client Error: "+err.Error(),
		)
		return
	}

	response.DataSourceData = c
	response.ResourceData = c

	tflog.Info(ctx, "Configured Migadu client")
}

func (p *MigaduProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewAliasDataSource,
		NewAliasesDataSource,
		NewIdentitiesDataSource,
		NewIdentityDataSource,
		NewMailboxDataSource,
		NewMailboxesDataSource,
		NewRewriteRuleDataSource,
		NewRewriteRulesDataSource,
	}
}

func (p *MigaduProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewAliasResource,
		NewIdentityResource,
		NewMailboxResource,
		NewRewriteRuleResource,
	}
}
