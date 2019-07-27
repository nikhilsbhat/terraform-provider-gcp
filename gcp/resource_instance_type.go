package gcp

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform/helper/schema"
)

// DecodeJSON will be called by provider to convert json to map
func instancetype() *schema.Resource {
	return &schema.Resource{
		Create: getInstanceType,
		Delete: schema.RemoveFromState,
		Read:   resourceReadItem,

		Schema: map[string]*schema.Schema{
			"machine_type": {
				Type:        schema.TypeString,
				Description: "Name of the machine type in gcp of which the information has to be fetched",
				Required:    true,
				ForceNew:    true,
			},

			"available_cpu": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"zone": {
				Type:        schema.TypeString,
				Description: "Name of the gcp zone which has to be used while initializing the client",
				Optional:    true,
				Computed:    true,
			},

			"maximum_persistent_disks": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// GetInstanceType is the core function that fetches the details of the type of instace selected.
func getInstanceType(d *schema.ResourceData, m interface{}) error {

	client := m.(*gcloudAuth)

	inputType := new(instypeInput)
	inputType.machineType = d.Get("machine_type").(string)
	intype, err := client.instanceType(*inputType)
	if err != nil {
		return errwrap.Wrapf("Error from gcp "+getStringOfMessage(err), err)
	}

	d.Set("available_cpu", intype.AvailableCpus)
	d.Set("name", intype.Name)
	d.Set("zone", intype.Zone)
	d.Set("maximum_persistent_disks", intype.MaximumPersistentDisks)
	d.SetId(intype.Name)
	return nil
}

func resourceReadItem(d *schema.ResourceData, m interface{}) error {

	client := m.(*gcloudAuth)

	inputType := new(instypeInput)
	inputType.machineType = d.Id()
	intype, err := client.instanceType(*inputType)
	if err != nil {
		return errwrap.Wrapf(getStringOfMessage(err), err)
	}

	d.Set("available_cpu", intype.AvailableCpus)
	d.Set("name", intype.Name)
	d.Set("zone", intype.Zone)
	d.Set("maximum_persistent_disks", intype.MaximumPersistentDisks)
	d.SetId(intype.Name)
	return nil
}

func writeFunc(data *gcloudAuth) error {
	file, err := os.Create("/Users/nikhil.bhat/terraform-codes/gcp/out.txt")
	if err != nil {
		return errwrap.Wrapf(getStringOfMessage(err), err)
	}
	defer file.Close()
	dataValue := fmt.Sprintf("%v", data)
	jsonDataVal, err := json.MarshalIndent(dataValue, "", " ")
	if err != nil {
		return errwrap.Wrapf("from json marshal"+getStringOfMessage(err), err)
	}
	_, err = file.Write(jsonDataVal)
	if err != nil {
		return errwrap.Wrapf(getStringOfMessage(err), err)
	}
	return nil
}
