package provider

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"terraform-provider-elastic-siem-detection/internal/fakeserver"
	"terraform-provider-elastic-siem-detection/internal/provider/transferobjects"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func generateTestRule() string {
	ruleContent := transferobjects.DetectionRule{
		RuleID: "7CE764F6-36A7-4E72-AB8B-166170CD1C93",
		ID:     "testID",
	}
	str, err := json.Marshal(ruleContent)
	if err != nil {
		fmt.Println(err)
		return "{}"
	}
	return string(str)
}

func TestAccDetectionRuleResource(t *testing.T) {

	debug := true
	apiServerObjects := make(map[string]map[string]interface{})

	svr := fakeserver.NewFakeServer(test_post, apiServerObjects, true, debug, "")
	test_url := fmt.Sprintf(`http://%s:%d`, test_host, test_post)
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
					fakeserver.TestAccCheckRestapiObjectExists("elastic-siem_detection_rule.test", "id", client),
					resource.TestCheckResourceAttr("elastic-siem_detection_rule.test", "rule_content", generateTestRule()),
				),
			},
			// ImportState testing
			{
				ResourceName:      "elastic-siem_detection_rule.test",
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
					resource.TestCheckResourceAttr("elastic-siem_detection_rule.test", "rule_content", generateTestRule()),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})

	svr.Shutdown()
}

func testAccDetectionRuleResourceConfig(ruleContent string, name string) string {
	content := strconv.Quote(string(ruleContent))
	return fmt.Sprintf(`%s
resource "elastic-siem_detection_rule" "%s" {
  rule_content = %s
}
`, providerConfig, name, content)
}
