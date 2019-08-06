package gcp

import (
	"fmt"
	"os"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"credentials": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("GOOGLE_APPLICATION_CREDENTIALS", ""),
				ValidateFunc: validateCredentials,
			},

			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("GOOGLE_PROJECT", ""),
			},

			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("GOOGLE_REGION", ""),
			},

			"zone": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("GOOGLE_ZONE", ""),
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"gcp_instance_type": instancetype(),
		},

		ConfigureFunc: getGCPClient,
	}
}

func validateCredentials(v interface{}, k string) (warnings []string, errors []error) {
	if v == nil || v.(string) == "" {
		return
	}
	creds := v.(string)
	// if this is a path and we can stat it, assume it's ok
	if _, err := os.Stat(creds); err != nil {
		errors = append(errors,
			fmt.Errorf("Unable to locate the credentials: %s", creds))
	}
	return
}
