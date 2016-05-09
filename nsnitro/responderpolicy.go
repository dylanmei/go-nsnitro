package nsnitro

type ResponderPolicy struct {
	Name   string `json:"name"`
	Action string `json:"action"`
	Rule   string `json:"rule"`
}

func (c *Client) GetResponderPolicy(name string) (ResponderPolicy, error) {
	policies := []ResponderPolicy{}
	err := c.fetch(nsrequest{
		Type: "responderpolicy", Name: name,
	}, &policies)

	if err != nil {
		return ResponderPolicy{}, err
	}

	return policies[0], nil
}

func (c *Client) GetResponderPolicies() ([]ResponderPolicy, error) {
	policies := []ResponderPolicy{}
	err := c.fetch(nsrequest{Type: "responderpolicy"}, &policies)
	return policies, err
}
