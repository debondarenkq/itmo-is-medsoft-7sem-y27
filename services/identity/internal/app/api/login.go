package api

import (
	"context"
	"fmt"
	"net"
	"time"

	contract "github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/api/security"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/identity/internal/models"
	"github.com/google/uuid"
)

func (implementation *Implementation) limited(ctx context.Context) bool {
	implementation.mu.Lock()
	defer implementation.mu.Unlock()
	now := time.Now()
	for key, window := range implementation.attempts {
		if now.Sub(window.start) > time.Minute {
			delete(implementation.attempts, key)
		}
	}
	address, _ := ctx.Value(remoteAddressKey).(string)
	host, _, _ := net.SplitHostPort(address)
	window := implementation.attempts[host]
	if window.start.IsZero() {
		window.start = now
	}
	window.count++
	implementation.attempts[host] = window
	return window.count > 20
}

func (implementation *Implementation) auditFailure(ctx context.Context, err error, kind models.AuditEventKind, message string) error {
	address, _ := ctx.Value(remoteAddressKey).(string)
	if auditErr := implementation.Module.Audit(ctx, "anonymous:"+address+":request:"+uuid.NewString(), message, kind); auditErr != nil {
		return auditErr
	}
	return err
}

func (implementation *Implementation) Login(ctx context.Context, request contract.LoginRequestObject) (contract.LoginResponseObject, error) {
	if implementation.limited(ctx) {
		return nil, implementation.auditFailure(ctx, models.ErrUnauthorized, models.AuditEventKindLoginFailed, "Превышен лимит попыток входа")
	}
	input := request.Body
	code := ""
	if input.Totp != nil {
		code = *input.Totp
	}
	session, err := implementation.Module.Login(ctx, input.Username, input.Password, code)
	if err != nil {
		return nil, implementation.auditFailure(ctx, err, models.AuditEventKindLoginFailed, "Неуспешная попытка входа: "+input.Username)
	}
	response, err := toSession(session)
	if err != nil {
		return nil, err
	}
	return contract.Login200JSONResponse(response), nil
}

func toStaffID(id *string) (*uuid.UUID, error) {
	if id == nil {
		return nil, nil
	}
	value, err := uuid.Parse(*id)
	if err != nil {
		return nil, fmt.Errorf("map staff id: %w", err)
	}
	return &value, nil
}

func toSession(s models.Session) (contract.Session, error) {
	id, err := uuid.Parse(s.AccountID)
	if err != nil {
		return contract.Session{}, err
	}
	staffID, err := toStaffID(s.StaffID)
	if err != nil {
		return contract.Session{}, err
	}
	return contract.Session{AccountID: id, StaffID: staffID, AccessToken: s.AccessToken, RefreshToken: s.RefreshToken, ExpiresAt: s.ExpiresAt, Role: s.Role}, nil
}
