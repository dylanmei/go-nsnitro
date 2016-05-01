package nsnitro

type LBMonitor struct {
	Name     string `json:"monitorname"`
	Type     string `json:"type"`
	Send     string `json:"send,omitempty"`
	Recv     string `json:"recv,omitempty"`
	Port     int    `json:"destport,omitempty"`
	Interval int    `json:"interval,omitempty"`
}

func NewLBMonitor(name, monitorType string) LBMonitor {
	return LBMonitor{Name: name, Type: monitorType}
}

func (c *Client) GetLBMonitor(name string) (LBMonitor, error) {
	lbmonitors := []LBMonitor{}
	err := c.fetch(nsrequest{
		Type: "lbmonitor", Name: name,
	}, &lbmonitors)

	if err != nil {
		return LBMonitor{}, err
	}
	return lbmonitors[0], nil
}

func (c *Client) GetLBMonitors() ([]LBMonitor, error) {
	lbmonitors := []LBMonitor{}
	err := c.fetch(nsrequest{Type: "lbmonitor"}, &lbmonitors)
	return lbmonitors, err
}

func (c *Client) AddLBMonitor(lbmonitor LBMonitor) error {
	return c.do("POST",
		nsrequest{
			Type: "lbmonitor",
		},
		&nsresource{
			LBMonitor: &lbmonitor,
		})
}

func (c *Client) RemoveLBMonitor(name, monitorType string) error {
	return c.do("DELETE",
		nsrequest{
			Type:  "lbmonitor",
			Name:  name,
			Query: map[string]string{"args": "type:" + monitorType},
		}, nil)
}
