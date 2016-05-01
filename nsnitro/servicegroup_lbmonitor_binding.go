package nsnitro

type ServiceGroupLBMonitorBinding struct {
	ServiceGroupName string `json:"servicegroupname"`
	MonitorName      string `json:"monitor_name"`
	State            string `json:"state,omitempty"`
}

func (c *Client) GetServiceGroupLBMonitorBindings(serviceGroupName string) ([]ServiceGroupLBMonitorBinding, error) {
	bindings := []ServiceGroupLBMonitorBinding{}
	err := c.fetch(nsrequest{
		Type: "servicegroup_lbmonitor_binding", Name: serviceGroupName,
	}, &bindings)
	return bindings, err
}

func (c *Client) BindServiceGroupToLBMonitor(serviceGroupName, monitorName string) error {
	binding := ServiceGroupLBMonitorBinding{ServiceGroupName: serviceGroupName, MonitorName: monitorName}
	return c.do("PUT",
		nsrequest{
			Name: serviceGroupName,
			Type: "servicegroup_lbmonitor_binding",
		},
		&nsresource{
			ServiceGroupLBMonitorBinding: &binding,
		})
}

func (c *Client) UnbindServiceGroupFromLBMonitor(serviceGroupName, monitorName string) error {
	return c.do("DELETE",
		nsrequest{
			Name: serviceGroupName,
			Type: "servicegroup_lbmonitor_binding",
			Query: map[string]string{
				"args": "monitor_name:" + monitorName,
			},
		}, nil)
}
