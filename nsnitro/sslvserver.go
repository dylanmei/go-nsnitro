package nsnitro

type SSLVServer struct {
	Name           string `json:"vservername"`
	SessionReuse   string `json:"sessreuse,omitempty"`
	SessionTimeout string `json:"sesstimeout,omitempty"`
	ClientAuth     string `json:"clientauth"`
	CipherRedirect string `json:"cipherredirect"`
	SSLRedirect    string `json:"sslredirect,omitempty"`
	DH             string `json:"dh"`
	SSL2           string `json:"ssl2,omitempty"`
	SSL3           string `json:"ssl3,omitempty"`
	TLS1           string `json:"tls1,omitempty"`
	TLS11          string `json:"tls11,omitempty"`
	TLS12          string `json:"tls12,omitempty"`
}

func (c *Client) GetSSLVServer(name string) (SSLVServer, error) {
	sslvservers := []SSLVServer{}
	err := c.fetch(nsrequest{
		Type: "sslvserver", Name: name,
	}, &sslvservers)

	if err != nil {
		return SSLVServer{}, err
	}

	return sslvservers[0], nil
}

func (c *Client) GetSSLVServers() ([]SSLVServer, error) {
	sslvservers := []SSLVServer{}
	err := c.fetch(nsrequest{Type: "sslvserver"}, &sslvservers)
	return sslvservers, err
}
