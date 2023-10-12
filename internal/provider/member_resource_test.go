package provider

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccMemberResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccMemberResourceConfig(2, "test@fake.com", "", true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("bitwarden_member.test", "type", "2"),
					resource.TestCheckResourceAttr("bitwarden_member.test", "access_all", "true"),
					resource.TestCheckResourceAttr("bitwarden_member.test", "external_id", ""),
					resource.TestCheckResourceAttr("bitwarden_member.test", "email", "test@fake.com"),
					resource.TestCheckResourceAttr("bitwarden_member.test", "name", ""),
					resource.TestCheckResourceAttrSet("bitwarden_member.test", "id"),
					resource.TestCheckResourceAttrSet("bitwarden_member.test", "status"),
					resource.TestCheckResourceAttrSet("bitwarden_member.test", "last_updated"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "bitwarden_member.test",
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
				Config: testAccMemberResourceConfig(3, "test@fake.com", "external-two", false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("bitwarden_member.test", "type", "3"),
					resource.TestCheckResourceAttr("bitwarden_member.test", "access_all", "false"),
					resource.TestCheckResourceAttr("bitwarden_member.test", "external_id", "external-two"),
					resource.TestCheckResourceAttr("bitwarden_member.test", "email", "test@fake.com"),
					resource.TestCheckResourceAttr("bitwarden_member.test", "name", ""),
					resource.TestCheckResourceAttrSet("bitwarden_member.test", "id"),
					resource.TestCheckResourceAttrSet("bitwarden_member.test", "status"),
					resource.TestCheckResourceAttrSet("bitwarden_member.test", "last_updated"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccMemberResourceConfig(mtype int64, email, externalId string, accessAll bool) string {
	builder := strings.Builder{}
	builder.WriteString("resource \"bitwarden_member\" \"test\" {\n")
	builder.WriteString(fmt.Sprintf("type = %d\n", mtype))
	builder.WriteString(fmt.Sprintf("email = %[1]q\n", email))
	builder.WriteString(fmt.Sprintf("access_all = %v\n", accessAll))
	if len(externalId) != 0 {
		builder.WriteString(fmt.Sprintf("external_id = %[1]q\n", externalId))
	}
	builder.WriteString("}")

	return builder.String()
}
