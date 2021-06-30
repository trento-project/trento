package cloud

type CloudInstance struct {
	Provider string      `mapstructure:"provider,omitempty"`
	Metadata interface{} `mapstructure:"metadata,omitempty"`
}

// Implement this method using dmidecode
func IdentifyCloudProvider() string {
	return "azure"
}

func NewCloudInstance() (*CloudInstance, error) {
	var err error
	var cloudMetadata interface{}

	provider := IdentifyCloudProvider()
	cInst := &CloudInstance{
		Provider: provider,
	}

	switch provider {
	case "azure":
		cloudMetadata, err = NewAzureMetadata()
		if err != nil {
			return nil, err
		}
	}

	cInst.Metadata = cloudMetadata

	return cInst, nil

}
