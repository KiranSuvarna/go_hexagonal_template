package usecase

import (
	"github.com/google/uuid"
	"github.com/hex_microservice_template/internal/core/domain"
	"github.com/hex_microservice_template/internal/core/ports/inbound"
	"github.com/hex_microservice_template/internal/core/ports/outbound"
	"github.com/pkg/errors"
	"gopkg.in/dealancer/validate.v2"
	"time"
)

var (
	ErrRedirectNotFound = errors.New("REDIRECT_NOT_FOUND")
	ErrRedirectInvalid  = errors.New("REDIRECT_INVALID")
)

type redirectService struct {
	redirectRepo outbound.RedirectRepository
}

func NewRedirectService(redirectRepo outbound.RedirectRepository) inbound.RedirectService {
	return &redirectService{
		redirectRepo,
	}
}

func (r *redirectService) Find(code string) (*domain.Redirect, error) {
	return r.redirectRepo.Find(code)
}

func (r *redirectService) Store(redirect *domain.Redirect) error {
	if err := validate.Validate(redirect); err != nil {
		return errors.Wrap(ErrRedirectInvalid, "usecase.Redirect.Store")
	}
	uuid, _ := uuid.NewUUID()
	redirect.Code = uuid.String()
	redirect.CreatedAt = time.Now().UTC().Unix()
	return r.redirectRepo.Store(redirect)
}
