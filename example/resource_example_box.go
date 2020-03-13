package example

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceExampleBoxCreate,
		Read:   resourceExampleBoxRead,
		Update: resourceExampleBoxUpdate,
		Delete: resourceExampleBoxDelete,

		Schema: map[string]*schema.Schema{
			"bundle": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					if !strings.HasPrefix(v, "b-") {
						errs = append(errs, fmt.Errorf(`%q must start with "b-", got: %s`, key, v))
					}
					return
				},
			},
			"logo": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "nothing.png",
			},
			"size": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"width": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  100,
						},
						"height": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  50,
						},
					},
				},
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceExampleBoxCreate(d *schema.ResourceData, m interface{}) error {
	// Generate "unique" id
	_id, err := rand.Int(rand.Reader, big.NewInt(1000))
	if err != nil {
		return err
	}
	id := fmt.Sprint(_id)

	c := m.(*Client)

	logo := d.Get("logo")
	size := d.Get("size").(*schema.Set).List()

	log.Printf("[DEBUG] Logo: %+v\n", logo)
	log.Printf("[DEBUG] Size: %+v\n", size)

	if err := c.CreateBox(&Box{
		Id:     id,
		Bundle: d.Get("bundle").(string),
	}); err != nil {
		return err
	}

	d.SetId(id)

	return resourceExampleBoxRead(d, m)
}

func resourceExampleBoxRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)

	t, err := c.ReadBox(d.Id())

	// If the resource does not exist, inform Terraform. We want to immediately
	// return here to prevent further processing.
	if err != nil {
		d.SetId("")
		return nil
	}

	d.Set("bundle", t.Bundle)
	return nil
}

func resourceExampleBoxUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)

	if err := c.UpdateBox(&Box{
		Id:     d.Id(),
		Bundle: d.Get("bundle").(string),
	}); err != nil {
		return err
	}

	return resourceExampleBoxRead(d, m)
}

func resourceExampleBoxDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)

	err := c.DeleteBox(d.Id())
	if err != nil {
		return err
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")
	return nil
}
