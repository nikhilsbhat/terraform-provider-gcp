package gcp

import (
	"fmt"
	"net/http"

	"context"

	"github.com/hashicorp/terraform/helper/schema"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/compute/v1"
)

var (
	auth = new(gcloudAuth)
)

type gcpSVCred struct {
	Type                string `json:"type,omitempty"`
	ProjectID           string `json:"project_id,omitempty"`
	PrivateKeyID        string `json:"private_key_id,omitempty"`
	PrivateKey          string `json:"private_key,omitempty"`
	ClientEmail         string `json:"client_email,omitempty"`
	ClientID            string `json:"client_id,omitempty"`
	AuthURI             string `json:"auth_uri,omitempty"`
	TokenURI            string `json:"token_uri,omitempty"`
	AuthProviderCertURL string `json:"auth_provider_x509_cert_url,omitempty"`
	ClientCertURL       string `json:"client_x509_cert_url,omitempty"`
}

type gcloudAuth struct {
	GCPSVCauth *gcpSVCred
	ProjectID  string
	Scopes     []string
	JSONPath   string
	Zone       string
	RawJSON    []byte
	Client     *http.Client
}

func getGCPClient(d *schema.ResourceData) (interface{}, error) {

	credspath := d.Get("credentials").(string)
	if credspath != "" {
		client, err := getCustomClient(credspath)
		if err != nil {
			return auth, err
		}
		client.Zone = d.Get("zone").(string)
		client.ProjectID = client.GCPSVCauth.ProjectID
		return client, nil
	}
	client, err := getDefaultClient()
	if err != nil {
		return nil, fmt.Errorf("Unable to fetch the default client")
	}
	zone := d.Get("zone").(string)
	if zone == "" {
		return nil, fmt.Errorf("Zone is not set and hence you cannot initialize client")
	}
	project := d.Get("project").(string)
	if project == "" {
		return nil, fmt.Errorf("Project ID is not set and hence you cannot initialize client")
	}
	client.Zone = zone
	client.ProjectID = project
	return client, nil
}

func getDefaultClient() (*gcloudAuth, error) {

	ctx := context.Background()
	client, err := google.DefaultClient(ctx, compute.CloudPlatformScope)
	if err != nil {
		return auth, err
	}
	auth.Client = client
	return auth, nil
}

func getCustomClient(path string) (*gcloudAuth, error) {

	jsonCont, err := readFile(path)
	if err != nil {
		return nil, err
	}

	var jsonAuth gcpSVCred

	if decodneuerr := jsonDecode(jsonCont, &jsonAuth); decodneuerr != nil {
		return nil, decodneuerr
	}

	auth.GCPSVCauth = &jsonAuth
	auth.JSONPath = path
	auth.Scopes = []string{compute.CloudPlatformScope}
	client := auth.getClient()
	if client == nil {
		return auth, fmt.Errorf("Unbale to initialize custom client")
	}
	auth.Client = client
	return auth, nil
}

func (auth *gcloudAuth) getClient() *http.Client {

	conf := &jwt.Config{
		Email:      auth.GCPSVCauth.ClientEmail,
		PrivateKey: []byte(auth.GCPSVCauth.PrivateKey),
		Scopes:     auth.Scopes,
		TokenURL:   auth.GCPSVCauth.TokenURI,
		Subject:    auth.GCPSVCauth.ClientEmail,
	}

	client := conf.Client(oauth2.NoContext)
	return client
}
