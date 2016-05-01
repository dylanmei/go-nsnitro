package nsnitro

type ServiceGroupServerBinding struct {
	ServiceGroupName string `json:"servicegroupname"`
	ServerName       string `json:"servername"`
	State            string `json:"state"`
}

func (c *Client) GetServiceGroupServerBindings(serviceGroupName string) ([]ServiceGroupServerBinding, error) {
	bindings := []ServiceGroupServerBinding{}
	err := c.fetch(nsrequest{
		Type: "servicegroup_servicegroupmember_binding", Name: serviceGroupName,
	}, &bindings)
	return bindings, err
}
