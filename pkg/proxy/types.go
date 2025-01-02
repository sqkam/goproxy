package proxy

import "context"

type Service interface {
	Run(ctx context.Context)
}
