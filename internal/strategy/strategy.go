package strategy

import "load-balancer/internal/backend"

type Strategy interface {
	Next() *backend.Backend
}
