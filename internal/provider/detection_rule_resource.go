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
var _ resource.Resource = &DetectionRuleResource{}
var _ resource.ResourceWithImportState = &DetectionRuleResource{}

func NewDetectionRuleResource() resource.Resource {
	return &DetectionRuleResource{}
}

// DetectionRuleResource defines the resource implementation.
type DetectionRuleResource struct {
	client *helpers.Client
}

// DetectionRuleResourceModel describes the resource data model.
type DetectionRuleResourceModel struct {
	RuleContent types.String `tfsdk:"rule_content"`
	Id          types.String `tfsdk:"id"`
}

func (r *DetectionRuleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_detection_rule"
}

func (r *DetectionRuleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Detection rule resource",

		Attributes: map[string]schema.Attribute{
			"rule_content": schema.StringAttribute{
				MarkdownDescription: "The content of the rule (JSON encoded string)",
				Required:            true,
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Rule identifier (in UUID format)",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *DetectionRuleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*helpers.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"[Configure][DetectionRule] Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *DetectionRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *DetectionRuleResourceModel
	var body *transferobjects.DetectionRule
	var itemsToRemove []string

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Process the content
	err := helpers.ObjectFromJSON(data.RuleContent.ValueString(), &body)
	if err != nil {
		resp.Diagnostics.AddError("[Create][DetectionRule] Parser Error", fmt.Sprintf("Unable to parse file, got error: %s", err))
		return
	}

	// If the threshold is empty the request has a "threshold":{} empty field that should be removed for the API request to work
	if len(body.Threshold.Field) == 0 {
		itemsToRemove = append(itemsToRemove, "threshold")
	}

	if len(body.ExceptionsList) == 0 {
		itemsToRemove = append(itemsToRemove, "exceptions_list")
	}

	// Create via API
	var response transferobjects.DetectionRuleResponse
	if err := r.client.Post("/detection_engine/rules", body, &response, itemsToRemove); err != nil {
		resp.Diagnostics.AddError("[Create][DetectionRule] Client Error", fmt.Sprintf("Error during request, got error: \n%s", err))
		return
	}

	// Save id into the Terraform state
	data.Id = types.StringValue(response.ID)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DetectionRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *DetectionRuleResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get via API
	var response transferobjects.DetectionRuleResponse
	path := fmt.Sprintf("/detection_engine/rules?id=%s", data.Id.ValueString())
	if err := r.client.Get(path, &response); err != nil {
		if strings.Contains(err.Error(), "404") {
			resp.Diagnostics.AddWarning("[Read][DetectionRule] Client Error", fmt.Sprintf("Resource not found. Will try to recreate if needed. Got error: %s", err))
			data.Id = types.StringNull()
			resp.State.RemoveResource(ctx)
			return
		} else {
			resp.Diagnostics.AddError("[Read][DetectionRule] Client Error", fmt.Sprintf("Error during request, got error: %s", err))
			return
		}
	}

	// Remove immutable or deprecated objects
	var itemsToRemove []string
	itemsToRemove = append(itemsToRemove, "created_by")
	itemsToRemove = append(itemsToRemove, "created_at")
	itemsToRemove = append(itemsToRemove, "updated_by")
	itemsToRemove = append(itemsToRemove, "updated_at")
	itemsToRemove = append(itemsToRemove, "version")
	itemsToRemove = append(itemsToRemove, "to")
	itemsToRemove = append(itemsToRemove, "id")
	itemsToRemove = append(itemsToRemove, "immutable")
	itemsToRemove = append(itemsToRemove, "throttle")

	// If threhold: {} remove it
	if len(response.Threshold.Field) == 0 {
		itemsToRemove = append(itemsToRemove, "threshold")
	}

	// If exception_list: [] remove it
	if len(response.DetectionRule.ExceptionsList) == 0 {
		itemsToRemove = append(itemsToRemove, "exceptions_list")
	}

	// Update the current state in case of diffs
	jsonStr, err := helpers.JSONfromObject(response.DetectionRule, itemsToRemove)
	if err != nil {
		resp.Diagnostics.AddError("[Read][DetectionRule] JSONfromObject Error", fmt.Sprintf("Error while JSONfromObject the updated state Rule Content, got error: %s", err))
		return
	}

	data.RuleContent = types.StringValue(jsonStr)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DetectionRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *DetectionRuleResourceModel
	// var stateData *DetectionRuleResourceModel
	var body *transferobjects.DetectionRule
	var itemsToRemove []string

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Process the content
	err := helpers.ObjectFromJSON(data.RuleContent.ValueString(), &body)
	if err != nil {
		resp.Diagnostics.AddError("[Update][DetectionRule] Parser Error", fmt.Sprintf("Unable to parse state file, got error: %s", err))
		return
	}

	// If threhold: {} remove it
	if len(body.Threshold.Field) == 0 {
		itemsToRemove = append(itemsToRemove, "threshold")
	}

	// If exception_list: [] remove it
	if len(body.ExceptionsList) == 0 {
		itemsToRemove = append(itemsToRemove, "exceptions_list")
	}

	if !helpers.CheckIfKeyExists(body, "rule_id") {
		body.ID = data.Id.ValueString()
	}

	// Update via API
	var response transferobjects.DetectionRuleResponse
	if err := r.client.Put("/detection_engine/rules", body, &response, itemsToRemove); err != nil {
		resp.Diagnostics.AddError("[Update][DetectionRule] Client Error", fmt.Sprintf("Error during request, got error: \n%s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DetectionRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *DetectionRuleResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get via API
	path := fmt.Sprintf("/detection_engine/rules?id=%s", data.Id.ValueString())
	if err := r.client.Delete(path); err != nil {
		if strings.Contains(err.Error(), "404") {
			resp.Diagnostics.AddWarning("[Delete][DetectionRule] Client Error", fmt.Sprintf("Resource not found. Will destroy if needed. Got error: %s", err))
			data.Id = types.StringNull()
			resp.State.RemoveResource(ctx)
			return
		} else {
			resp.Diagnostics.AddError("[Delete][DetectionRule] Client Error", fmt.Sprintf("Error during request, got error: %s", err))
			return
		}
	}
}

func (r *DetectionRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
