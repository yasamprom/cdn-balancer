package balancerservice

import "context"

type usecases interface {
	GetLink(ctx context.Context, uri string) (string, error)
}
