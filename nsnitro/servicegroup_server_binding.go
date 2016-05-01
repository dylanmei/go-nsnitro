package nsnitro

import "fmt"

type ServiceGroupServerBinding struct {
	ServiceGroupName string `json:"servicegroupname"`
	ServerName       string `json:"servername"`
	Port             int    `json:"port"`
	State            string `json:"state,omitempty"`
}

func (c *Client) GetServiceGroupServerBindings(serviceGroupName string) ([]ServiceGroupServerBinding, error) {
	bindings := []ServiceGroupServerBinding{}
	err := c.fetch(nsrequest{
		Type: "servicegroup_servicegroupmember_binding", Name: serviceGroupName,
	}, &bindings)
	return bindings, err
}

func (c *Client) BindServiceGroupToServer(serviceGroupName, serverName string, port int) error {
	binding := ServiceGroupServerBinding{
		ServiceGroupName: serviceGroupName,
		ServerName:       serverName,
		Port:             port,
	}
	return c.do("PUT",
		nsrequest{
			Name: serviceGroupName,
			Type: "servicegroup_servicegroupmember_binding",
			Query: map[string]string{
				"action": "bind",
			},
		},
		&nsresource{
			ServiceGroupServerBinding: &binding,
		})
}

func (c *Client) UnbindServiceGroupFromServer(serviceGroupName, serverName string, port int) error {
	return c.do("DELETE",
		nsrequest{
			Name: serviceGroupName,
			Type: "servicegroup_servicegroupmember_binding",
			Query: map[string]string{
				"args": fmt.Sprintf("servername:%s,port:%d", serverName, port),
			},
		}, nil)
}
