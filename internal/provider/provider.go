package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-elastic-siem/internal/helpers"
)

// Ensure ElasticSiemProvider satisfies various provider interfaces.
var _ provider.Provider = &ElasticSiemProvider{}

// ElasticSiemProvider defines the provider implementation.
type ElasticSiemProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// ScaffoldingProviderModel describes the provider data model.
type ScaffoldingProviderModel struct {
	Hostname types.String `tfsdk:"hostname"`
	UseTLS   types.Bool   `tfsdk:"tls"`
	Port     types.Int64  `tfsdk:"port"`
	Username types.String `tfsdk:"user"`
	Password types.String `tfsdk:"password"`
}

func (p *ElasticSiemProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "elastic-siem"
	resp.Version = p.version
}

func (p *ElasticSiemProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"hostname": schema.StringAttribute{
				MarkdownDescription: "The Kibana host name",
				Required:            true,
			},
			"tls": schema.BoolAttribute{
				MarkdownDescription: "Connect to host using TLS or unencrypted",
				Optional:            true,
			},
			"port": schema.Int64Attribute{
				MarkdownDescription: "Connect to host on a custom port",
				Optional:            true,
			},
			"user": schema.StringAttribute{
				MarkdownDescription: "The username to authenticate to Kiaba and interact with the SIEM",
				Required:            true,
				Sensitive:           true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "The password to authenticate to Kiaba and interact with the SIEM",
				Required:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *ElasticSiemProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data ScaffoldingProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	// if data.Endpoint.IsNull() { /* ... */ }

	hostname := "localhost"
	port := 443
	useTls := true
	username := "elastic"
	password := "elastic"

	if !data.Hostname.IsNull() {
		hostname = data.Hostname.ValueString()
	}

	if !data.Port.IsNull() {
		port = int(data.Port.ValueInt64())
	}

	if !data.UseTLS.IsNull() {
		useTls = data.UseTLS.ValueBool()
	}

	if !data.Username.IsNull() {
		username = data.Username.ValueString()
	}

	if !data.Password.IsNull() {
		password = data.Password.ValueString()
	}

	// Example client configuration for data sources and resources
	client := helpers.NewClient(&helpers.NewClientInput{
		Hostname: hostname,
		Port:     port,
		UseTls:   useTls,
		Username: username,
		Password: password,
	})

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *ElasticSiemProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewDetectionRuleResource,
		NewExceptionItemResource,
		NewExceptionContainerResource,
	}
}

func (p *ElasticSiemProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewPrivilegesDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &ElasticSiemProvider{
			version: version,
		}
	}
}
