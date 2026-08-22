package pricing

// RateProvider is implemented by anything that can look up a spot exchange
// rate between two currency codes.
type RateProvider interface {
	FetchRate(base, quote string) (float64, error)
}

type Service struct {
	rates RateProvider
}

func NewService(rates RateProvider) *Service {
	return &Service{rates: rates}
}

func (s *Service) Convert(amount float64, base, quote string) (float64, error) {
	rate, err := s.rates.FetchRate(base, quote)
	if err != nil {
		return 0, err
	}
	return amount * rate, nil
}
