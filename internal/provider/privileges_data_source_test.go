package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"terraform-provider-elastic-siem/internal/fakeserver"
	"terraform-provider-elastic-siem/internal/helpers"
	"testing"
)

func TestAccPrivilegesDataSource(t *testing.T) {
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

	client.SendRequest("POST", "/api/detection_engine/privileges", `
    {
  "username": "elastic",
  "has_all_requested": true,
  "cluster": {
    "monitor_ml": true,
    "manage_ccr": true,
    "manage_index_templates": true,
    "monitor_watcher": true,
    "monitor_transform": true,
    "read_ilm": true,
    "manage_api_key": true,
    "manage_security": true,
    "manage_own_api_key": true,
    "manage_saml": true,
    "all": true,
    "manage_ilm": true,
    "manage_ingest_pipelines": true,
    "read_ccr": true,
    "manage_rollup": true,
    "monitor": true,
    "manage_watcher": true,
    "manage": true,
    "manage_transform": true,
    "manage_token": true,
    "manage_ml": true,
    "manage_pipeline": true,
    "monitor_rollup": true,
    "transport_client": true,
    "create_snapshot": true
  },
  "index": {
    ".alerts-security.alerts-default": {
      "all": true,
      "create": true,
      "create_doc": true,
      "create_index": true,
      "delete": true,
      "delete_index": true,
      "index": true,
      "maintenance": true,
      "manage": true,
      "manage_follow_index": true,
      "manage_ilm": true,
      "manage_leader_index": true,
      "monitor": true,
      "read": true,
      "read_cross_cluster": true,
      "view_index_metadata": true,
      "write": true
    }
  },
  "application": {},
  "is_authenticated": true,
  "has_encryption_key": true
}
  `)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			svr.StartInBackground()
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccPrivilegesDataSourceConfig("test"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.elastic-siem_privileges.test", "id", helpers.Sha256String("elastic")),
				),
			},
		},
	})
}

func testAccPrivilegesDataSourceConfig(name string) string {
	return fmt.Sprintf(`%s
data "elastic-siem_privileges" "%s" {
}
`, providerConfig, name)
}
