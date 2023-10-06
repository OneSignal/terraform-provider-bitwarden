package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-bitwarden/internal/bitwarden"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &bitwardenProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &bitwardenProvider{
			version: version,
		}
	}
}

// bitwardenProvider is the provider implementation.
type bitwardenProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// bitwardenProviderModel maps provider schema data to a Go type.
type bitwardenProviderModel struct {
	ClientID          types.String `tfsdk:"client_id"`
	ClientSecret      types.String `tfsdk:"client_secret"`
	APIUrl            types.String `tfsdk:"api_url"`
	AuthenticationUrl types.String `tfsdk:"authentication_url"`
}

// Metadata returns the provider type name.
func (p *bitwardenProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "bitwarden"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *bitwardenProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Bitwarden API Provider, focussing on Group and Group Members for organisations",
		Attributes: map[string]schema.Attribute{
			"client_id": schema.StringAttribute{
				MarkdownDescription: "The client_id of your organisation, can also be configured as `BITWARDEN_CLIENT_ID`. See [docs](https://bitwarden.com/help/public-api/#authentication) for more information",
				Optional:            true,
			},
			"client_secret": schema.StringAttribute{
				MarkdownDescription: "The client_secret of your organisation, can also be configured as `BITWARDEN_CLIENT_SECRET`. See [docs](https://bitwarden.com/help/public-api/#authentication) for more information",
				Optional:            true,
				Sensitive:           true,
			},
			"api_url": schema.StringAttribute{
				MarkdownDescription: "The Bitwarden API URL, defaults to `https://api.bitwarden.com`, can also be configured as `BITWARDEN_API_URL`. See [docs](https://bitwarden.com/help/public-api/#endpoints) for more information",
				Optional:            true,
			},
			"authentication_url": schema.StringAttribute{
				MarkdownDescription: "The Bitwarden Authentication URL, defaults to `https://identity.bitwarden.com/connect/token`, can also be configured as `BITWARDEN_AUTH_URL`. See [docs](https://bitwarden.com/help/public-api/#authentication-endpoints) for more information",
				Optional:            true,
			},
		},
	}
}

// Configure prepares a Bitwarden API client for data sources and resources.
func (p *bitwardenProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config bitwardenProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.ClientID.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("client_id"),
			"Unknown Bitwarden client_id",
			"The provider cannot create the Bitwarden API client as there is an unknown configuration value for the Bitwarden Client ID. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the BITWARDEN_CLIENT_ID environment variable.",
		)
	}

	if config.ClientSecret.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("client_secret"),
			"Unknown Bitwarden client_secret",
			"The provider cannot create the Bitwarden API client as there is an unknown configuration value for the Bitwarden client_secret. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the BITWARDEN_CLIENT_SECRET environment variable.",
		)
	}

	if config.APIUrl.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_url"),
			"Unknown Bitwarden API URL",
			"The provider cannot create the Bitwarden API client as there is an unknown configuration value for the Bitwarden API URL. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the BITWARDEN_API_URL environment variable.",
		)
	}

	if config.AuthenticationUrl.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("authentication_url"),
			"Unknown Bitwarden Authentication URL",
			"The provider cannot create the Bitwarden API client as there is an unknown configuration value for the Bitwarden Authentication URL. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the BITWARDEN_AUTHENTICATION_URL environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.
	clientId := os.Getenv("BITWARDEN_CLIENT_ID")
	clientSecret := os.Getenv("BITWARDEN_CLIENT_SECRET")
	apiUrl := os.Getenv("BITWARDEN_API_URL")
	authUrl := os.Getenv("BITWARDEN_AUTHENTICATION_URL")

	if !config.ClientID.IsNull() {
		clientId = config.ClientID.ValueString()
	}

	if !config.ClientSecret.IsNull() {
		clientSecret = config.ClientSecret.ValueString()
	}

	if !config.APIUrl.IsNull() {
		apiUrl = config.APIUrl.ValueString()
	} else if len(apiUrl) == 0 {
		apiUrl = "https://api.bitwarden.com/public"
	}

	if !config.AuthenticationUrl.IsNull() {
		authUrl = config.AuthenticationUrl.ValueString()
	} else if len(authUrl) == 0 {
		authUrl = "https://identity.bitwarden.com/connect/token"
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.
	if clientId == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("client_id"),
			"Missing Bitwarden client_id",
			"The provider cannot create the Bitwarden API client as there is a missing or empty value for the Bitwarden client_id. "+
				"Set the client_id value in the configuration or use the BITWARDEN_CLIENT_ID environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if clientSecret == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("client_secret"),
			"Missing Bitwarden client_secret",
			"The provider cannot create the Bitwarden API client as there is a missing or empty value for the Bitwarden client_secret. "+
				"Set the client_secret value in the configuration or use the BITWARDEN_CLIENT_SECRET environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if apiUrl == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_url"),
			"Missing Bitwarden API URL",
			"The provider cannot create the Bitwarden API client as there is a missing or empty value for the Bitwarden API URL. "+
				"Set the api_url value in the configuration or use the BITWARDEN_API_URL environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if authUrl == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("authentication_url"),
			"Missing Bitwarden Authentication URL",
			"The provider cannot create the Bitwarden API client as there is a missing or empty value for the Bitwarden Authentication URL. "+
				"Set the authentication_url value in the configuration or use the BITWARDEN_AUTHENTICATION_URL environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "bitwarden_client_id", clientId)
	ctx = tflog.SetField(ctx, "bitwarden_api_url", apiUrl)
	ctx = tflog.SetField(ctx, "bitwarden_authentication_url", authUrl)
	ctx = tflog.SetField(ctx, "bitwarden_client_secret", clientSecret)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "bitwarden_client_secret")

	tflog.Debug(ctx, "Creating Bitwarden API client")

	// Create a new Bitwarden client using the configuration values
	client, err := bitwarden.NewClient(ctx, clientId, clientSecret, apiUrl, authUrl)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Bitwarden API Client",
			"An unexpected error occurred when creating the Bitwarden API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Bitwarden Client Error: "+err.Error(),
		)
		return
	}

	// Make the Bitwarden client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client
}

// DataSources defines the data sources implemented in the provider.
func (p *bitwardenProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}

// Resources defines the resources implemented in the provider.
func (p *bitwardenProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewGroupResource,
	}
}
