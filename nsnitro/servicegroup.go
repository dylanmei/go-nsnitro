package nsnitro

type ServiceGroup struct {
	Name  string `json:"servicegroupname"`
	Type  string `json:"servicetype"`
	State string `json:"state,omitempty"`
}

func NewServiceGroup(name, serviceType string) ServiceGroup {
	return ServiceGroup{Name: name, Type: serviceType}
}

func (c *Client) GetServiceGroup(name string) (ServiceGroup, error) {
	servicegroups := []ServiceGroup{}
	err := c.fetch(nsrequest{
		Type: "servicegroup", Name: name,
	}, &servicegroups)
	if err != nil {
		return ServiceGroup{}, err
	}
	return servicegroups[0], nil
}

func (c *Client) GetServiceGroups() ([]ServiceGroup, error) {
	servicegroups := []ServiceGroup{}
	err := c.fetch(nsrequest{Type: "servicegroup"}, &servicegroups)
	return servicegroups, err
}

func (c *Client) RemoveServiceGroup(name string) error {
	return c.do("DELETE",
		nsrequest{
			Type: "servicegroup",
			Name: name,
		}, nil)
}

func (c *Client) AddServiceGroup(serviceGroup ServiceGroup) error {
	return c.do("POST",
		nsrequest{
			Type: "servicegroup",
		},
		&nsresource{
			ServiceGroup: &serviceGroup,
		})
}
