package modules

import (
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/laboratory/internal/pkg/adapters"
	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/laboratory/internal/pkg/storage"
)

type Module struct {
	Store       *storage.Repository
	Directories *adapters.Directories
}
