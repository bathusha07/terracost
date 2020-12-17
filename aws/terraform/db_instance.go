package terraform

import (
	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"
	"strings"

	"github.com/cycloidio/cost-estimation/aws/region"
	"github.com/cycloidio/cost-estimation/price"
	"github.com/cycloidio/cost-estimation/product"
	"github.com/cycloidio/cost-estimation/query"
	"github.com/cycloidio/cost-estimation/util"
)

type DBInstance struct {
	providerKey string

	region       region.Code
	instanceType string

	databaseEngine   string
	databaseEdition  string
	licenseModel     string
	deploymentOption string

	storageType      string
	allocatedStorage decimal.Decimal
	storageIOPS      decimal.Decimal
}

type dbInstanceValues struct {
	InstanceClass string `mapstructure:"instance_class"`
	AvailabilityZone string `mapstructure:"availability_zone"`
	Engine string `mapstructure:"engine"`
	LicenseModel string `mapstructure:"license_model"`
	MultiAZ bool `mapstructure:"multi_az"`
	AllocatedStorage float64 `mapstructure:"allocated_storage"`
	StorageType string `mapstructure:"storage_type"`
	IOPS float64 `mapstructure:"iops"`
}

type dbType struct {
	engine, edition string
}

var dbTypeMap = map[string]dbType{
	"aurora":        {"Aurora MySQL", ""},
	"aurora-mysql":  {"Aurora MySQL", ""},
	"mariadb":       {"MariaDB", ""},
	"mysql":         {"MySQL", ""},
	"postgres":      {"PostgreSQL", ""},
	"oracle-se":     {"Oracle", "Standard"},
	"oracle-se1":    {"Oracle", "Standard One"},
	"oracle-se2":    {"Oracle", "Standard Two"},
	"oracle-ee":     {"Oracle", "Enterprise"},
	"sqlserver-se":  {"SQL Server", "Standard"},
	"sqlserver-ee":  {"SQL Server", "Enterprise"},
	"sqlserver-ex":  {"SQL Server", "Express"},
	"sqlserver-web": {"SQL Server", "Web"},
}
var licenseModelMap = map[string]string{
	"license-included":       "License included",
	"bring-your-own-license": "Bring your own license",
}

func decodeDBInstanceValues(tfVals map[string]interface{}) (dbInstanceValues, error) {
	var v dbInstanceValues
	if err := mapstructure.Decode(tfVals, &v); err != nil {
		return v, err
	}
	return v, nil
}

// NewInstance creates a new Instance from Terraform values.
func (p *Provider) newDBInstance(vals dbInstanceValues) *DBInstance {
	dbType := dbTypeMap[vals.Engine]
	licenseModel := licenseModelMap[vals.LicenseModel]

	deploymentOption := "Single-AZ"
	if vals.MultiAZ {
		deploymentOption = "Multi-AZ"
	}

	inst := &DBInstance{
		providerKey:      p.key,
		region:           p.region,
		instanceType:     vals.InstanceClass,
		databaseEngine:   dbType.engine,
		databaseEdition:  dbType.edition,
		licenseModel:     licenseModel,
		deploymentOption: deploymentOption,
		allocatedStorage: decimal.NewFromFloat(vals.AllocatedStorage),
		storageType:      vals.StorageType,
		storageIOPS:      decimal.NewFromFloat(vals.IOPS),
	}

	if reg := region.NewFromZone(vals.AvailabilityZone); reg.Valid() {
		inst.region = reg
	}

	return inst
}

// Components returns the price component queries that make up this Instance.
func (inst *DBInstance) Components() []query.Component {
	components := []query.Component{inst.databaseInstanceComponent(), inst.storageComponent()}

	if strings.HasPrefix(inst.storageType, "io") {
		components = append(components, inst.iopsComponent())
	}

	return components
}

func (inst *DBInstance) databaseInstanceComponent() query.Component {
	instClass := inst.instanceType
	attrFilters := []*product.AttributeFilter{
		{Key: "instanceType", Value: util.StringPtr(inst.instanceType)},
		{Key: "deploymentOption", Value: util.StringPtr(inst.deploymentOption)},
		{Key: "databaseEngine", Value: util.StringPtr(inst.databaseEngine)},
	}

	if inst.databaseEdition != "" {
		f := &product.AttributeFilter{Key: "databaseEdition", Value: util.StringPtr(inst.databaseEdition)}
		attrFilters = append(attrFilters, f)
	}

	if inst.licenseModel != "" {
		f := &product.AttributeFilter{Key: "licenseModel", Value: util.StringPtr(inst.licenseModel)}
		attrFilters = append(attrFilters, f)
	}

	return query.Component{
		Name:           "Database instance",
		Details:        []string{inst.deploymentOption, instClass},
		HourlyQuantity: decimal.NewFromInt(1),
		ProductFilter: &product.Filter{
			Provider:         util.StringPtr(inst.providerKey),
			Service:          util.StringPtr("AmazonRDS"),
			Family:           util.StringPtr("Database Instance"),
			Location:         util.StringPtr(inst.region.String()),
			AttributeFilters: attrFilters,
		},
		PriceFilter: &price.Filter{
			Unit: util.StringPtr("Hrs"),
			AttributeFilters: []*price.AttributeFilter{
				{Key: "purchaseOption", Value: util.StringPtr("on_demand")},
			},
		},
	}
}

func (inst *DBInstance) storageComponent() query.Component {
	var volumeType string
	switch inst.storageType {
	case "standard":
		volumeType = "Magnetic"
	case "io1", "io2":
		volumeType = "Provisioned IOPS"
	default:
		volumeType = "General Purpose"
	}

	return query.Component{
		Name:            "Database storage",
		Details:         []string{volumeType},
		MonthlyQuantity: inst.allocatedStorage,
		Unit:            "GB",
		ProductFilter: &product.Filter{
			Provider: util.StringPtr(inst.providerKey),
			Service:  util.StringPtr("AmazonRDS"),
			Family:   util.StringPtr("Database Storage"),
			Location: util.StringPtr(inst.region.String()),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "deploymentOption", Value: util.StringPtr(inst.deploymentOption)},
				{Key: "volumeType", Value: util.StringPtr(volumeType)},
			},
		},
	}
}

func (inst *DBInstance) iopsComponent() query.Component {
	return query.Component{
		Name:            "Database IOPS",
		MonthlyQuantity: inst.storageIOPS,
		Unit:            "IOPS",
		ProductFilter: &product.Filter{
			Provider: util.StringPtr(inst.providerKey),
			Service:  util.StringPtr("AmazonRDS"),
			Family:   util.StringPtr("Provisioned IOPS"),
			Location: util.StringPtr(inst.region.String()),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "deploymentOption", Value: util.StringPtr(inst.deploymentOption)},
			},
		},
	}
}
