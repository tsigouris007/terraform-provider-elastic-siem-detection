package provider

import (
	"context"
	"fmt"
	"strings"
	"terraform-provider-elastic-siem-detection/internal/helpers"
	"terraform-provider-elastic-siem-detection/internal/provider/transferobjects"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &ExceptionItemResource{}
var _ resource.ResourceWithImportState = &ExceptionItemResource{}

func NewExceptionItemResource() resource.Resource {
	return &ExceptionItemResource{}
}

// ExceptionItemResource defines the resource implementation.
type ExceptionItemResource struct {
	client *helpers.Client
}

// ExceptionItemResourceModel describes the resource data model.
type ExceptionItemResourceModel struct {
	RuleContent types.String `tfsdk:"exception_item_content"`
	Id          types.String `tfsdk:"id"`
}

func (r *ExceptionItemResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_exception_item"
}

func (r *ExceptionItemResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Exception item resource",

		Attributes: map[string]schema.Attribute{
			"exception_item_content": schema.StringAttribute{
				MarkdownDescription: "The content of the exception item (JSON encoded string)",
				Required:            true,
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Exception item identifier (in UUID format)",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *ExceptionItemResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*helpers.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"[Configure][ExceptionItem] Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *ExceptionItemResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *ExceptionItemResourceModel
	var body transferobjects.ExceptionItem
	var itemsToRemove []string

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Process the content
	err := helpers.ObjectFromJSON(data.RuleContent.ValueString(), &body)
	if err != nil {
		resp.Diagnostics.AddError("[Create][ExceptionItem] Parser Error", fmt.Sprintf("Unable to parse file, got error: %s", err))
		return
	}

	// Completely disabling comments on creation due to API inconsistencies
	if len(body.Comments) > 0 {
		resp.Diagnostics.AddWarning("[Create][ExceptionItem] Comments unsupported", fmt.Sprintf("Due to Elastic API inconsistencies comments are not supported."))
		itemsToRemove = append(itemsToRemove, "comments")
	}

	// Create via API
	var response transferobjects.ExceptionItemResponse
	if err := r.client.Post("/exception_lists/items", body, &response, itemsToRemove); err != nil {
		resp.Diagnostics.AddError("[Create][ExceptionItem] Client Error", fmt.Sprintf("Error during request, got error: \n%s", err))
		return
	}

	// Save id into the Terraform state
	data.Id = types.StringValue(response.ID)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ExceptionItemResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *ExceptionItemResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get via API
	var response transferobjects.ExceptionItemResponse
	path := fmt.Sprintf("/exception_lists/items?id=%s", data.Id.ValueString())

	if err := r.client.Get(path, &response); err != nil {
		if strings.Contains(err.Error(), "404") {
			resp.Diagnostics.AddWarning("[Read][ExceptionItem] Client Error", fmt.Sprintf("Resource not found. Will try to recreate if needed. Got error: %s", err))
			data.Id = types.StringNull()
			resp.State.RemoveResource(ctx)
			return
		} else {
			resp.Diagnostics.AddError("[Read][ExceptionItem] Client Error", fmt.Sprintf("Error during request, got error: %s", err))
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

	// Ignore comments
	if len(response.Comments) >= 0 {
		itemsToRemove = append(itemsToRemove, "threshold")
	}

	// Update the state in case of diffs
	jsonStr, err := helpers.JSONfromObject(response.ExceptionItemBase, itemsToRemove)
	if err != nil {
		resp.Diagnostics.AddError("[Read][ExceptionItem] Marshal Error", fmt.Sprintf("Error while marshalling the updated state Exception Item Content, got error: %s", err))
		return
	}

	data.RuleContent = types.StringValue(jsonStr)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ExceptionItemResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *ExceptionItemResourceModel
	var body *transferobjects.ExceptionItem
	var itemsToRemove []string

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Process the content
	err := helpers.ObjectFromJSON(data.RuleContent.ValueString(), &body)
	if err != nil {
		resp.Diagnostics.AddError("[Update][ExceptionItem] Parser Error", fmt.Sprintf("Unable to parse file, got error: %s", err))
		return
	}

	body.ID = data.Id.ValueString()
	body.ListID = "" // This should always be empty in the PUT request

	// Completely disabling comments on update due to API inconsistencies
	if len(body.Comments) > 0 {
		resp.Diagnostics.AddWarning("[Update][ExceptionItem] Comments unsupported", fmt.Sprintf("Due to Elastic API inconsistencies comments are not supported."))
		itemsToRemove = append(itemsToRemove, "comments")
	}

	// Update via API
	var response transferobjects.ExceptionItemResponse
	if err := r.client.Put("/exception_lists/items", body, &response, itemsToRemove); err != nil {
		resp.Diagnostics.AddError("[Update][ExceptionItem] Client Error", fmt.Sprintf("Error during request, got error: \n%s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ExceptionItemResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *ExceptionItemResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get via API
	path := fmt.Sprintf("/exception_lists/items?id=%s", data.Id.ValueString())
	if err := r.client.Delete(path); err != nil {
		if strings.Contains(err.Error(), "404") {
			resp.Diagnostics.AddWarning("[Delete][ExceptionItem] Client Error", fmt.Sprintf("Resource not found. Will destroy if needed. Got error: %s", err))
			data.Id = types.StringNull()
			resp.State.RemoveResource(ctx)
			return
		} else {
			resp.Diagnostics.AddError("[Delete][ExceptionItem] Client Error", fmt.Sprintf("Error during request, got error: %s", err))
			return
		}
	}
}

func (r *ExceptionItemResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
