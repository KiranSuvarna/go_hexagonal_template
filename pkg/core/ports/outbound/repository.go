package outbound

import (
	"github.com/hex_microservice_template/pkg/core/domain"
)

type RedirectRepository interface {
	Find(code string) (*domain.Redirect, error)
	Store(redirect *domain.Redirect) error
}
