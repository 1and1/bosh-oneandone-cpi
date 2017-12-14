package action


// Environment used to create an instance
type Environment map[string]interface{}

// StemcellCloudProperties holds the CPI specific stemcell properties
// defined in stemcell's manifest
type StemcellCloudProperties struct {
	Name           string `json:"name,omitempty"`
	ImageID        string `json:"image-id,omitempty"`
	ImageSourceURL string `json:"image-source-url,omitempty"`
}

// NetworkCloudProperties holds the CPI specific network properties
// defined in cloud config
type NetworkCloudProperties struct {
	VcnName    string `json:"vcn,omitempty"`
	SubnetName string `json:"subnet_name,omitempty"`
}

// VMCloudProperties holds the CPI specific properties
// defined in cloud-config for creating a instance
type VMCloudProperties struct {
	Name       string  `json:"name,omitempty"`
	Datacenter string  `json:"datacenter,omitempty"`
	Cores      int     `json:"cores,omitempty"`
	DiskSize   int     `json:"diskSize,omitempty"`
	Ram        float32 `json:"ram,omitempty"`
	SSHKey     string  `json:"rsa_key,omitempty"`
}

func (o VMCloudProperties) Validate() error {
	//err := o.Agent.Validate()
	//if err != nil {
	//	return bosherr.WrapError(err, "Validating Agent configuration")
	//}
	//
	//err = o.Softlayer.Validate()
	//if err != nil {
	//	return bosherr.WrapError(err, "Validating SoftLayer configuration")
	//}

	return nil
}
