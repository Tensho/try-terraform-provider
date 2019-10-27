package example

import (
	"fmt"
	"testing"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("example_box", &resource.Sweeper{
		Name:	"example_box",
		F:		func (region string) error {
			fmt.Println("Sweep, sweep, sweep")
			return nil
		},
	})
}

func TestAccExampleBox_basic(t *testing.T) {
	var b Box

	bundle := fmt.Sprintf("b-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckExampleBoxDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccExampleBoxConfig(bundle),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExampleBoxExists("example_box.agent", &b),
					testAccCheckExampleBoxValues(&b, bundle),
					resource.TestCheckResourceAttr("example_box.agent", "bundle", bundle),
				),
			},
		},
	})
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("EXAMPLE_API_KEY"); v != "" {
		t.Fatal("EXAMPLE_API_KEY must be set for acceptance tests")
	}
}

func testAccExampleBoxConfig(bundle string) string {
	return fmt.Sprintf(`
resource example_box agent {
  bundle = "%s"
}`, bundle)
}

func testAccCheckExampleBoxExists(name string, b *Box) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		c := testAccProvider.Meta().(*Client)
		resp, err := c.ReadBox(rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("Box (%s) not found", rs.Primary.ID)
		}

		*b = *resp

		return nil
	}
}

func testAccCheckExampleBoxValues(b *Box, bundle string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		t := *b
		if t.Bundle != bundle {
			return fmt.Errorf("bad bundle, expected %q, got: %#v", bundle, *b)
		}
		return nil
	}
}

func testAccCheckExampleBoxDestroy(s *terraform.State) error {
  c := testAccProvider.Meta().(*Client)

  for _, rs := range s.RootModule().Resources {
    if rs.Type != "example_box" {
      continue
    }

    _, err := c.ReadBox(rs.Primary.ID)
    if err.Error() != fmt.Sprintf("Box (%s) not found", rs.Primary.ID) {
    	return fmt.Errorf("Box (%s) still exists", rs.Primary.ID)
    }
  }

  return nil
}
