package inbound

import (
	"github.com/hex_microservice_template/pkg/core/domain"
)

type RedirectService interface {
	Find(code string) (*domain.Redirect, error)
	Store(redirect *domain.Redirect) error
}
