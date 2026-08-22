package notify

type SlackClient struct {
	webhookURL string
}

func NewSlackClient(webhookURL string) *SlackClient {
	return &SlackClient{webhookURL: webhookURL}
}

func (c *SlackClient) Notify(message string) error {
	return nil
}
