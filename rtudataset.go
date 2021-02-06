package rtuapi

type DataSet interface {
	ToJson() []byte
}

type RmDeviceConfiguration struct {
	Host     string
	Username string
	Password string
}

func (rmdc *RmDeviceConfiguration) ToJson() []byte {
	return ToJSON(rmdc)
}
