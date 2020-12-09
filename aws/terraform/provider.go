package terraform

import (
	"fmt"

	"github.com/cycloidio/cost-estimation/aws/region"
	"github.com/cycloidio/cost-estimation/query"
	"github.com/cycloidio/cost-estimation/terraform"
)

// Provider is an implementation of the terraform.Provider, used to extract component queries from
// terraform resources.
type Provider struct {
	key    string
	region region.Code
}

// NewProvider returns a new Provider with the provided default region and a query key.
func NewProvider(key string, regionCode region.Code) (*Provider, error) {
	if !regionCode.Valid() {
		return nil, fmt.Errorf("invalid AWS region: %q", regionCode)
	}
	return &Provider{key: key, region: regionCode}, nil
}

// ResourceComponents returns Component queries for a given terraform.Resource.
func (p *Provider) ResourceComponents(tfRes terraform.Resource) []query.Component {
	switch tfRes.Type {
	case "aws_instance":
		vals, err := decodeInstanceValues(tfRes.Values)
		if err != nil {
			return nil
		}
		return p.newInstance(vals).Components()
	case "aws_ebs_volume":
		vals, err := decodeVolumeValues(tfRes.Values)
		if err != nil {
			return nil
		}
		return p.newVolume(vals).Components()
	default:
		return nil
	}
}
