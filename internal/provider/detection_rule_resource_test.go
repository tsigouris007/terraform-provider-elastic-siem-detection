package provider

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"terraform-provider-elastic-siem-detection/internal/fakeserver"
	"terraform-provider-elastic-siem-detection/internal/provider/transferobjects"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func generateTestRule() string {
	ruleContent := transferobjects.DetectionRule{}
	ruleContent.ID = "myTestID"
	ruleContent.RuleID = "12345678-abcd-edfg-hijk-1234567890ab"
	ruleContent.Description = "Test Rule Description"
	ruleContent.Name = "Test Rule Name"
	ruleContent.RiskScore = 21
	ruleContent.Severity = "low"
	ruleContent.Type = "query"

	str, err := json.Marshal(ruleContent)
	if err != nil {
		fmt.Println(err)
		return "{}"
	}
	objStr := string(str)
	// Remove empty objects
	objStr = strings.Replace(objStr, "\"exceptions_list\":null,", "", 1)
	objStr = strings.Replace(objStr, ",\"threshold\":{}", "", 1)

	return objStr
}

func generateInvalidTestRule() string {
	str, _ := json.Marshal("{invalid_json}")
	return string(str)
}

func TestAccDetectionRuleResource(t *testing.T) {

	debug := true
	apiServerObjects := make(map[string]map[string]interface{})

	svr := fakeserver.NewFakeServer(test_port, apiServerObjects, true, debug, "")
	test_url := fmt.Sprintf(`http://%s:%d`, test_host, test_port)
	os.Setenv("REST_API_URI", test_url)

	opt := &fakeserver.ApiClientOpt{
		Uri:                 test_url,
		Insecure:            false,
		Username:            "",
		Password:            "",
		Headers:             make(map[string]string),
		Timeout:             2,
		IdAttribute:         "id",
		CopyKeys:            make([]string, 0),
		WriteReturnsObject:  false,
		CreateReturnsObject: false,
		Debug:               debug,
	}
	client, err := fakeserver.NewAPIClient(opt)
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			svr.StartInBackground()
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccDetectionRuleResourceConfig(generateTestRule(), "test"),
				Check: resource.ComposeAggregateTestCheckFunc(
					fakeserver.TestAccCheckRestapiObjectExists("elastic-siem-detection_detection_rule.test", "id", client),
					resource.TestCheckResourceAttr("elastic-siem-detection_detection_rule.test", "rule_content", generateTestRule()),
				),
				ExpectNonEmptyPlan: true, // stubbed
			},
			// Invalid rule_content
			{
				Config:      testAccDetectionRuleResourceConfig(generateInvalidTestRule(), "test"),
				ExpectError: regexp.MustCompile(`Parser Error`),
			},
			// ImportState testing
			{
				ResourceName:      "elastic-siem-detection_detection_rule.test",
				ImportState:       true,
				ImportStateVerify: true,
				// This is not normally necessary, but is here because this
				// example code does not have an actual upstream service.
				// Once the Read method is able to refresh information from
				// the upstream service, this can be removed.
				ImportStateVerifyIgnore: []string{"rule_content"},
			},
			// Update and Read testing
			{
				Config: testAccDetectionRuleResourceConfig(generateTestRule(), "test"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("elastic-siem-detection_detection_rule.test", "rule_content", generateTestRule()),
				),
				ExpectNonEmptyPlan: true, // stubbed
			},
			// Delete testing automatically occurs in TestCase
		},
	})

	svr.Shutdown()
}

func testAccDetectionRuleResourceConfig(ruleContent string, name string) string {
	content := strconv.Quote(string(ruleContent))
	return fmt.Sprintf(`%s
resource "elastic-siem-detection_detection_rule" "%s" {
  rule_content = %s
}
`, providerConfig, name, content)
}
