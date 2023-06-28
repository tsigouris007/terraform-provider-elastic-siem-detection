package provider

import (
	"context"
	"fmt"
	"terraform-provider-elastic-siem/internal/helpers"
	"terraform-provider-elastic-siem/internal/provider/transferobjects"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &PrivilegesDataSource{}

func NewPrivilegesDataSource() datasource.DataSource {
	return &PrivilegesDataSource{}
}

// PrivilegesDataSource defines the data source implementation.
type PrivilegesDataSource struct {
	client *helpers.Client
}

// PrivilegesDataSourceModel describes the data source data model.
type PrivilegesDataSourceModel struct {
	Id types.String `tfsdk:"id"`
}

func (d *PrivilegesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_privileges"
}

func (d *PrivilegesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Privileges data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Privileges identifier",
				Computed:            true,
			},
		},
	}
}

func (d *PrivilegesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*helpers.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"[Configure][Privileges] Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *PrivilegesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data PrivilegesDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get the privileges through the API
	var response transferobjects.PrivilegesResponse
	if err := d.client.Get("/detection_engine/privileges", &response); err != nil {
		resp.Diagnostics.AddError("[Read][Privileges] Client Error", fmt.Sprintf("Error during request, got error: %s", err))
		return
	}

	// Save id into the Terraform state.
	data.Id = types.StringValue(helpers.Sha256String(response.Username))

	// Write logs using the tflog package
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
