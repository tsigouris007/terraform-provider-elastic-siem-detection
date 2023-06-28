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

func generateTestExceptionContainer() string {
	ruleContent := transferobjects.ExceptionContainer{}
	ruleContent.ID = "myTestID"
	ruleContent.ListID = "12345678-abcd-efgh-ijkl-1234567890ab"
	ruleContent.Description = "Test Container Description"
	ruleContent.Name = "Test Container Name"
	ruleContent.NamespaceType = "single"
	ruleContent.Type = "detection"

	str, err := json.Marshal(ruleContent)
	if err != nil {
		fmt.Println(err)
		return "{}"
	}
	objStr := string(str)

	return objStr
}

func TestAccExceptionContainerResource(t *testing.T) {

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
				Config: testAccExceptionContainerResourceConfig(generateTestExceptionContainer(), "test"),
				Check: resource.ComposeAggregateTestCheckFunc(
					fakeserver.TestAccCheckRestapiObjectExists("elastic-siem-detection_exception_container.test", "id", client),
					resource.TestCheckResourceAttr("elastic-siem-detection_exception_container.test", "exception_container_content", generateTestExceptionContainer()),
				),
				ExpectNonEmptyPlan: true, // stubbed
			},
			// ImportState testing
			{
				ResourceName:      "elastic-siem-detection_exception_container.test",
				ImportState:       true,
				ImportStateVerify: true,
				// This is not normally necessary, but is here because this
				// example code does not have an actual upstream service.
				// Once the Read method is able to refresh information from
				// the upstream service, this can be removed.
				ImportStateVerifyIgnore: []string{"rule_content"},
			},
			// Update and Read testing
			// {
			// 	Config: testAccExceptionContainerResourceConfig(generateTestExceptionContainer(), "test"),
			// 	Check: resource.ComposeAggregateTestCheckFunc(
			// 		resource.TestCheckResourceAttr("elastic-siem-detection_exception_container.test", "exception_container_content", generateTestExceptionContainer()),
			// 	),
			// 	ExpectNonEmptyPlan: true, // stubbed
			// },
			// Delete testing automatically occurs in TestCase
		},
	})

	svr.Shutdown()
}

func testAccExceptionContainerResourceConfig(ruleContent string, name string) string {
	content := strconv.Quote(string(ruleContent))
	return fmt.Sprintf(`%s
resource "elastic-siem-detection_exception_container" "%s" {
  exception_container_content = %s
}
`, providerConfig, name, content)
}
