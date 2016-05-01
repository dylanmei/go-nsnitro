package nsnitro

type LBVServerServiceGroupBinding struct {
	LBVServerName    string `json:"name"`
	ServiceGroupName string `json:"servicegroupname"`
}

func (c *Client) GetLBVServerServiceGroupBindings(lbvserverName string) ([]LBVServerServiceGroupBinding, error) {
	bindings := []LBVServerServiceGroupBinding{}
	err := c.fetch(nsrequest{
		Type: "lbvserver_servicegroup_binding", Name: lbvserverName,
	}, &bindings)

	return bindings, err
}

func (c *Client) BindLBVServerToServiceGroup(lbvserverName, serviceGroupName string) error {
	binding := LBVServerServiceGroupBinding{LBVServerName: lbvserverName, ServiceGroupName: serviceGroupName}
	return c.do("PUT",
		nsrequest{
			Name: lbvserverName,
			Type: "lbvserver_servicegroup_binding",
		},
		&nsresource{
			LBVServerServiceGroupBinding: &binding,
		})
}

func (c *Client) UnbindLBVServerFromServiceGroup(lbvserverName, serviceGroupName string) error {
	return c.do("DELETE",
		nsrequest{
			Name: serviceGroupName,
			Type: "lbvserver_servicegroup_binding",
			Query: map[string]string{
				"args": "servicegroupname:" + serviceGroupName,
			},
		}, nil)
}
