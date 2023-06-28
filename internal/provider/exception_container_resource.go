package provider

import (
	"context"
	"fmt"
	"strings"
	"terraform-provider-elastic-siem/internal/helpers"
	"terraform-provider-elastic-siem/internal/provider/transferobjects"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &ExceptionContainerResource{}
var _ resource.ResourceWithImportState = &ExceptionContainerResource{}

func NewExceptionContainerResource() resource.Resource {
	return &ExceptionContainerResource{}
}

// ExceptionContainerResource defines the resource implementation.
type ExceptionContainerResource struct {
	client *helpers.Client
}

// ExceptionContainerResourceModel describes the resource data model.
type ExceptionContainerResourceModel struct {
	RuleContent types.String `tfsdk:"exception_container_content"`
	Id          types.String `tfsdk:"id"`
}

func (r *ExceptionContainerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_exception_container"
}

func (r *ExceptionContainerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Exception container resource",

		Attributes: map[string]schema.Attribute{
			"exception_container_content": schema.StringAttribute{
				MarkdownDescription: "The content of the exception container (JSON encoded string)",
				Required:            true,
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Exception container identifier (in UUID format)",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *ExceptionContainerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*helpers.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"[Configure][ExceptionContainer] Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *ExceptionContainerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *ExceptionContainerResourceModel
	var body *transferobjects.ExceptionContainer
	var itemsToRemove []string = []string{}

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Process the content
	err := helpers.ObjectFromJSON(data.RuleContent.ValueString(), &body)
	if err != nil {
		resp.Diagnostics.AddError("[Create][ExceptionContainer] Parser Error", fmt.Sprintf("Unable to parse file, got error: %s", err))
		return
	}

	// Create via API
	var response transferobjects.ExceptionContainerResponse
	if err := r.client.Post("/exception_lists", body, &response, itemsToRemove); err != nil {
		resp.Diagnostics.AddError("[Create][ExceptionContainer] Client Error", fmt.Sprintf("Error during request, got error: \n%s", err))
		return
	}

	// Save id into the Terraform state
	data.Id = types.StringValue(response.ID)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ExceptionContainerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *ExceptionContainerResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get via API
	var response transferobjects.ExceptionContainerResponse
	path := fmt.Sprintf("/exception_lists?id=%s", data.Id.ValueString())

	if err := r.client.Get(path, &response); err != nil {
		if strings.Contains(err.Error(), "404") {
			resp.Diagnostics.AddWarning("[Read][ExceptionContainer] Client Error", fmt.Sprintf("Resource not found. Will try to recreate if needed. Got error: %s", err))
			data.Id = types.StringNull()
			resp.State.RemoveResource(ctx)
			return
		} else {
			resp.Diagnostics.AddError("[Read][ExceptionContainer] Client Error", fmt.Sprintf("Error during request, got error: %s", err))
			return
		}
	}

	// Remove immutable or deprecated objects
	var itemsToRemove []string
	itemsToRemove = append(itemsToRemove, "id")
	itemsToRemove = append(itemsToRemove, "created_by")
	itemsToRemove = append(itemsToRemove, "created_at")
	itemsToRemove = append(itemsToRemove, "updated_by")
	itemsToRemove = append(itemsToRemove, "updated_at")

	// Update the state in case of diffs
	jsonStr, err := helpers.JSONfromObject(response.ExceptionContainer, itemsToRemove)
	if err != nil {
		resp.Diagnostics.AddError("[Read][ExceptionItem] Marshal Error", fmt.Sprintf("Error while marshalling the updated state Exception Item Content, got error: %s", err))
		return
	}

	data.RuleContent = types.StringValue(jsonStr)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ExceptionContainerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *ExceptionContainerResourceModel
	var body *transferobjects.ExceptionContainer
	var itemsToRemove []string = []string{}

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Process the content
	err := helpers.ObjectFromJSON(data.RuleContent.ValueString(), &body)
	if err != nil {
		resp.Diagnostics.AddError("[Update][ExceptionContainer] Parser Error", fmt.Sprintf("Unable to parse file, got error: %s", err))
		return
	}

	body.ID = data.Id.ValueString()

	// Update via API
	var response transferobjects.ExceptionContainerResponse
	if err := r.client.Put("/exception_lists", body, &response, itemsToRemove); err != nil {
		resp.Diagnostics.AddError("[Update][ExceptionContainer] Client Error", fmt.Sprintf("Error during request, got error: \n%s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ExceptionContainerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *ExceptionContainerResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get via API
	apiPath := fmt.Sprintf("/exception_lists?id=%s", data.Id.ValueString())
	if err := r.client.Delete(apiPath); err != nil {
		if strings.Contains(err.Error(), "404") {
			resp.Diagnostics.AddWarning("[Delete][ExceptionContainer] Client Error", fmt.Sprintf("Resource not found. Will destroy if needed. Got error: %s", err))
			data.Id = types.StringNull()
			resp.State.RemoveResource(ctx)
			return
		} else {
			resp.Diagnostics.AddError("[Delete][ExceptionContainer] Client Error", fmt.Sprintf("Error during request, got error: %s", err))
			return
		}
	}
}

func (r *ExceptionContainerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
