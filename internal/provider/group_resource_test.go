package provider

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccGroupResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccGroupResourceConfig("one", "", true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("bitwarden_group.test", "name", "one"),
					resource.TestCheckResourceAttr("bitwarden_group.test", "access_all", "true"),
					resource.TestCheckResourceAttrSet("bitwarden_group.test", "id"),
					resource.TestCheckResourceAttrSet("bitwarden_group.test", "last_updated"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "bitwarden_group.test",
				ImportState:       true,
				ImportStateVerify: true,
				// This is not normally necessary, but is here because this
				// example code does not have an actual upstream service.
				// Once the Read method is able to refresh information from
				// the upstream service, this can be removed.
				ImportStateVerifyIgnore: []string{"last_updated"},
			},
			//Update and Read testing
			{
				Config: testAccGroupResourceConfig("two", "external-two", false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("bitwarden_group.test", "name", "two"),
					resource.TestCheckResourceAttr("bitwarden_group.test", "external_id", "external-two"),
					resource.TestCheckResourceAttr("bitwarden_group.test", "access_all", "false"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccGroupResourceConfig(name, externalId string, accessAll bool) string {
	builder := strings.Builder{}
	builder.WriteString("resource \"bitwarden_group\" \"test\" {\n")
	builder.WriteString(fmt.Sprintf("name = %[1]q\n", name))
	builder.WriteString(fmt.Sprintf("access_all = %v\n", accessAll))
	if len(externalId) != 0 {
		builder.WriteString(fmt.Sprintf("external_id = %[1]q\n", externalId))
	}
	builder.WriteString("}")

	return builder.String()
}
