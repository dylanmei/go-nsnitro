package nsnitro

type Server struct {
	Name  string `json:"name"`
	IP    string `json:"ipaddress"`
	State string `json:"state,omitempty"`
}

func NewServer(name, ip string) Server {
	return Server{Name: name, IP: ip}
}

func (c *Client) GetServer(name string) (Server, error) {
	servers := []Server{}
	err := c.fetch(nsrequest{
		Type: "server", Name: name,
	}, &servers)
	if err != nil {
		return Server{}, err
	}
	return servers[0], nil
}

func (c *Client) GetServers() ([]Server, error) {
	servers := []Server{}
	err := c.fetch(nsrequest{Type: "server"}, &servers)
	return servers, err
}

func (c *Client) AddServer(server Server) error {
	return c.do("POST",
		nsrequest{
			Type: "server",
		},
		&nsresource{
			Server: &server,
		})
}

func (c *Client) RemoveServer(name string) error {
	return c.do("DELETE",
		nsrequest{
			Type: "server",
			Name: name,
		}, nil)
}
