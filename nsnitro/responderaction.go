package nsnitro

type ResponderAction struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Target string `json:"target,omitempty"`
}

func (c *Client) GetResponderAction(name string) (ResponderAction, error) {
	actions := []ResponderAction{}
	err := c.fetch(nsrequest{
		Type: "responderaction", Name: name,
	}, &actions)

	if err != nil {
		return ResponderAction{}, err
	}

	return actions[0], nil
}

func (c *Client) GetResponderActions() ([]ResponderAction, error) {
	actions := []ResponderAction{}
	err := c.fetch(nsrequest{Type: "responderaction"}, &actions)
	return actions, err
}
