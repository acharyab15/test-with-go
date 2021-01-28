package stripe

type Customer struct {
	ID string
}

type Client struct{}

func (c *Client) Customer(tok string) (*Customer, error) {
	return nil, nil
}
