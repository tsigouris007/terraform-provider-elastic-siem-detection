package provider

import (
	"fmt"
	"os"
	"strings"
	"terraform-provider-elastic-siem-detection/internal/fakeserver"
	"terraform-provider-elastic-siem-detection/internal/provider/transferobjects"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func generateTestExceptionContainer() transferobjects.ExceptionContainer {
	ruleContent := transferobjects.ExceptionContainer{
		Name:          "test container",
		NamespaceType: "single",
		Tags:          []string{"asdf", "fdsa"},
		Type:          "detection",
		ListID:        "7CE764F6-36A7-4E72-AB8B-166170CD1C93",
		Description:   "test description",
		ID:            "generatedTestID", // needs to be this string
	}
	return ruleContent
}

func TestAccExceptionContainerResource(t *testing.T) {

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
				Config: testAccExceptionContainerResourceConfig(generateTestExceptionContainer(), "test"),
				Check: resource.ComposeAggregateTestCheckFunc(
					fakeserver.TestAccCheckRestapiObjectExists("elastic-siem_exception_container.test", "id", client),
					resource.TestCheckResourceAttr("elastic-siem_exception_container.test", "namespace_type", generateTestExceptionContainer().NamespaceType),
				),
			},
			// ImportState testing
			//{
			//	ResourceName:      "elastic-siem_exception_container.test",
			//	ImportState:       true,
			//	ImportStateVerify: true,
			//	// This is not normally necessary, but is here because this
			//	// example code does not have an actual upstream service.
			//	// Once the Read method is able to refresh information from
			//	// the upstream service, this can be removed.
			//	ImportStateVerifyIgnore: []string{"rule_content"},
			//},
			// Update and Read testing
			{
				Config: testAccExceptionContainerResourceConfig(generateTestExceptionContainer(), "test"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("elastic-siem_exception_container.test", "namespace_type", generateTestExceptionContainer().NamespaceType),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})

	svr.Shutdown()
}

func testAccExceptionContainerResourceConfig(ruleContent transferobjects.ExceptionContainer, name string) string {
	return fmt.Sprintf(`%s
resource "elastic-siem_exception_container" "%s" {
  description = "%s"
  name = "%s"
  list_id = "%s"
  type = "%s"
  namespace_type = "%s"
  tags = ["%s"]
}
`, providerConfig, name,
		ruleContent.Description,
		ruleContent.Name,
		ruleContent.ListID,
		ruleContent.Type,
		ruleContent.NamespaceType,
		strings.Join(ruleContent.Tags, "\", \""),
	)
}
