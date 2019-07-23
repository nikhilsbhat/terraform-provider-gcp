package gcp

import (
	"golang.org/x/net/context"

	//"golang.org/x/oauth2/google"
	"strconv"

	"google.golang.org/api/compute/v1"
)

type instype struct {
	CreationTimestamp            string `json:"creationTimestamp,omitempty"`
	Description                  string `json:"description,omitempty"`
	AvailableCpus                string `json:"guestCpus,omitempty"`
	ID                           string `json:"id,omitempty"`
	ImageSpaceGb                 string `json:"imageSpaceGb,omitempty"`
	Kind                         string `json:"kind,omitempty"`
	MaximumPersistentDisks       string `json:"maximumPersistentDisks,omitempty"`
	MaximumPersistentDisksSizeGb string `json:"maximumPersistentDisksSizeGb,omitempty"`
	MemoryMb                     string `json:"memoryMb,omitempty"`
	Name                         string `json:"name,omitempty"`
	SelfLink                     string `json:"selfLink,omitempty"`
	Zone                         string `json:"zone,omitempty"`
	MachineTypeRaw               *compute.MachineType
}

type instypeInput struct {
	machineType string
	returnraw   bool
	zone        string
}

var (
	// Instype holds the info of instance type selected
	Instype = new(instype)
)

func (auth *gcloudAuth) instanceType(input instypeInput) (instype, error) {
	ctx := context.Background()

	auth.Scopes = []string{compute.CloudPlatformScope}
	client := auth.getClient()

	computeService, err := compute.New(client)
	if err != nil {
		return instype{}, err
	}

	resp, err := computeService.MachineTypes.Get(auth.GCPSVCauth.ProjectID, auth.Zone, input.machineType).Context(ctx).Do()
	if err != nil {
		return instype{}, err
	}

	if input.returnraw == true {
		Instype.MachineTypeRaw = resp
		return *Instype, nil
	}

	Instype.CreationTimestamp = resp.CreationTimestamp
	Instype.Description = resp.Description
	Instype.AvailableCpus = strconv.Itoa(int(resp.GuestCpus))
	Instype.ID = strconv.Itoa(int(resp.Id))
	Instype.ImageSpaceGb = strconv.Itoa(int(resp.ImageSpaceGb))
	Instype.Kind = resp.Kind
	Instype.MaximumPersistentDisks = strconv.Itoa(int(resp.MaximumPersistentDisks))
	Instype.MaximumPersistentDisksSizeGb = strconv.Itoa(int(resp.MaximumPersistentDisksSizeGb))
	Instype.MemoryMb = strconv.Itoa(int(resp.MemoryMb))
	Instype.Name = resp.Name
	Instype.SelfLink = resp.SelfLink
	Instype.Zone = resp.Zone

	return *Instype, nil
}
