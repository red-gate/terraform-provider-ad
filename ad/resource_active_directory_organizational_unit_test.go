package ad

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/go-ldap/ldap/v3"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAdOU_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckShort(t)
			testAccPreCheck(t)
			testAccResourceAdOUPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAdOUDestroy("ad_organizational_unit.test"),
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccResourceAdOUConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAdOUExists("ad_organizational_unit.test"),
					resource.TestCheckResourceAttr(
						"ad_organizational_unit.test", "ou_name", "terraform"),
				),
			},
		},
	})
}

func testAccResourceAdOUPreCheck(t *testing.T) {
	if v := os.Getenv("AD_OU_DOMAIN"); v == "" {
		t.Fatal("AD_OU_DOMAIN must be set for acceptance tests")
	}
}

func testAccCheckAdOUDestroy(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No AD OU ID is set")
		}
		client := testAccProvider.Meta().(*ldap.Conn)
		domain := rs.Primary.Attributes["domain"]
		var dnOfOU string
		domainArr := strings.Split(domain, ".")
		dnOfOU = "dc=" + domainArr[0]
		for index, item := range domainArr {
			if index == 0 {
				continue
			}
			dnOfOU += ",dc=" + item
		}
		searchRequest := ldap.NewSearchRequest(
			dnOfOU, //"cn=code1,cn=OU,dc=terraform,dc=local", // The base dn to search
			ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
			"(&(objectClass=organizationalUNit)(cn="+rs.Primary.Attributes["ou_name"]+"))", // The filter to apply
			[]string{"dn", "cn"}, // A list attributes to retrieve
			nil,
		)
		sr, err := client.Search(searchRequest)
		if err != nil {
			return err
		}
		if len(sr.Entries) == 0 {
			return nil
		}

		return fmt.Errorf("OU AD still exists")
	}

}

func testAccCheckAdOUExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No AD OU ID is set")
		}
		client := testAccProvider.Meta().(*ldap.Conn)
		domain := rs.Primary.Attributes["domain"]
		var dnOfOU string
		domainArr := strings.Split(domain, ".")
		dnOfOU = "dc=" + domainArr[0]
		for index, item := range domainArr {
			if index == 0 {
				continue
			}
			dnOfOU += ",dc=" + item
		}
		searchRequest := ldap.NewSearchRequest(
			dnOfOU, //"cn=code1,cn=OUs,dc=terraform,dc=local", // The base dn to search
			ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
			"(&(objectClass=organizationalUNit)(cn="+rs.Primary.Attributes["ou_name"]+"))", // The filter to apply
			[]string{"dn", "cn"}, // A list attributes to retrieve
			nil,
		)
		sr, err := client.Search(searchRequest)
		if err != nil {
			return err
		}
		if len(sr.Entries) > 0 {
			return nil
		}
		return nil
	}
}

func testAccResourceAdOUConfig() string {
	return fmt.Sprintf(`
provider "ad" {
  domain   = "%s"
  ip       = "%s"
  user     = "%s"
  password = "%s"
}
resource "ad_organizational_unit" "test" {
  domain = "%s"
  ou_name = "terraform"
}`,
		os.Getenv("AD_DOMAIN"),
		os.Getenv("AD_IP"),
		os.Getenv("AD_USER"),
		os.Getenv("AD_PASSWORD"),
		os.Getenv("AD_OU_DOMAIN"))
}
