package carrier

type Label struct {
	TrackingID string
	CarrierRef string
}

// Carrier is implemented by any shipping provider integration.
type Carrier interface {
	Ship(weightGrams int, destZIP string) (Label, error)
	Cancel(trackingID string) error
	RateQuote(weightGrams int, destZIP string) (cents int64, err error)
}

type FedExClient struct {
	apiKey string
}

func NewFedExClient(apiKey string) *FedExClient {
	return &FedExClient{apiKey: apiKey}
}

func (c *FedExClient) Ship(weightGrams int, destZIP string) (Label, error) {
	return Label{TrackingID: "FX-000", CarrierRef: "fedex"}, nil
}

func (c *FedExClient) Cancel(trackingID string) error {
	return nil
}

func (c *FedExClient) RateQuote(weightGrams int, destZIP string) (cents int64, err error) {
	return 999, nil
}
