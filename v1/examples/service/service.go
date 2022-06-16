package service

// Dummy interface that simulates getting profile data
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Service
type Service interface {
	UpdateProfile(string) error
}

type service struct{}

func (s service) UpdateProfile(id string) error {
	return nil
}

func NewService() Service {
	return service{}
}
