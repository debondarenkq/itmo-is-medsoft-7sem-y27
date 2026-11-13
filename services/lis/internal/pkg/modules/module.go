package modules

import (
	"net/http"
	"sync"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/lis/internal/pkg/storage"
)

type Module struct {
	Gateway string
	Store   *storage.Repository
	HTTP    *http.Client
	mu      sync.Mutex
}
