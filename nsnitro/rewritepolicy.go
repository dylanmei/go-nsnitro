package nsnitro

type RewritePolicy struct {
	Name   string `json:"name"`
	Action string `json:"action"`
	Rule   string `json:"rule"`
}

func (c *Client) GetRewritePolicy(name string) (RewritePolicy, error) {
	policies := []RewritePolicy{}
	err := c.fetch(nsrequest{
		Type: "rewritepolicy", Name: name,
	}, &policies)

	if err != nil {
		return RewritePolicy{}, err
	}

	return policies[0], nil
}

func (c *Client) GetRewritePolicies() ([]RewritePolicy, error) {
	policies := []RewritePolicy{}
	err := c.fetch(nsrequest{Type: "rewritepolicy"}, &policies)
	return policies, err
}
