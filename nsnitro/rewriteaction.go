package nsnitro

type RewriteAction struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	Target     string `json:"target,omitempty"`
	Expression string `json:"stringBuilderExpr,omitempty"`
}

func (c *Client) GetRewriteAction(name string) (RewriteAction, error) {
	actions := []RewriteAction{}
	err := c.fetch(nsrequest{
		Type: "rewriteaction", Name: name,
	}, &actions)

	if err != nil {
		return RewriteAction{}, err
	}

	return actions[0], nil
}

func (c *Client) GetRewriteActions() ([]RewriteAction, error) {
	actions := []RewriteAction{}
	err := c.fetch(nsrequest{Type: "rewriteaction"}, &actions)
	return actions, err
}
