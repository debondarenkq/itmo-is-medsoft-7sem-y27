package modules

import (
	"context"
	"crypto/ed25519"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type StaffDirectory interface {
	Get(context.Context, string) (models.Staff, error)
	All(context.Context) ([]models.Staff, error)
	UpdateRole(context.Context, string, string) (models.Staff, error)
}

type SecurityStorage interface {
	ServiceAccountStorage
	AccountByID(context.Context, string) (models.Account, error)
	Account(context.Context, string) (models.Account, error)
	RefreshAccount(context.Context, string) (models.Account, error)
	SetupSecret(context.Context, string) (bool, string, error)
	Setup(context.Context, models.Account, string, string, int64) error
	Issue(context.Context, models.Account, string, string, int64) error
	CreateAccount(context.Context, models.Account, string) (models.Account, error)
	ChangeAccount(context.Context, string, *bool, string, string) (models.Account, error)
	Logout(context.Context, string, string) error
	AppendAudit(context.Context, string, string, models.AuditEventKind) error
	SearchAccounts(context.Context, models.AccountFilter) ([]models.Account, int, error)
	SearchAuditEvents(context.Context, models.AuditFilter) ([]models.AuditEvent, int, error)
	UpdateAccount(context.Context, string, models.AccountChanges, string) (models.Account, error)
}

type Security struct {
	Directory StaffDirectory
	Store     SecurityStorage
	key       ed25519.PrivateKey
	dummyHash []byte
}

type ServiceAccountStorage interface {
	ServiceAccountByKey(context.Context, string) (models.Account, error)
	ServiceAccounts(context.Context) ([]models.Account, error)
	SaveServiceAccount(context.Context, models.Account, bool, string) (models.Account, error)
}

func NewSecurity(store SecurityStorage, key ed25519.PrivateKey) *Security {
	dummy, _ := bcrypt.GenerateFromPassword([]byte("invalid-login-password"), 12)
	return &Security{Store: store, key: key, dummyHash: dummy}
}
