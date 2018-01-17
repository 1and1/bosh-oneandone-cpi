package action

// DiskCloudProperties holds the CPI specific disk properties
type DiskCloudProperties struct {
	Datacenter string `json:"datacenter,omitempty"`
}

// Environment used to create an instance
type Environment map[string]interface{}

// StemcellCloudProperties holds the CPI specific stemcell properties
// defined in stemcell's manifest
type StemcellCloudProperties struct {
	Name           string `json:"name,omitempty"`
	ImageID        string `json:"image-id,omitempty"`
	ImageSourceURL string `json:"image-source-url,omitempty"`
	OSType         string `json:"os-type,omitempty"`
	Architecture   int    `json:"architecture,omitempty"`
}

// NetworkCloudProperties holds the CPI specific network properties
// defined in cloud config
type NetworkCloudProperties struct {
	OpenPorts []Rule `json:"open-ports,omitempty"`
}

type Rule struct {
	PortFrom *int   `json:"port-from,omitempty"`
	PortTo   *int   `json:"port-to,omitempty"`
	Source   string `json:"source,omitempty"`
}

// VMCloudProperties holds the CPI specific properties
// defined in cloud-config for creating a instance
type VMCloudProperties struct {
	Name           string  `json:"name,omitempty"`
	Datacenter     string  `json:"datacenter,omitempty"`
	InstanceFlavor string  `json:"flavor,omitempty"`
	Cores          int     `json:"cores,omitempty"`
	DiskSize       int     `json:"diskSize,omitempty"`
	Ram            float32 `json:"ram,omitempty"`
	SSHKey         string  `json:"rsa_key,omitempty"`
	PublicIP       string  `json:"public_ip,omitempty"`
}
