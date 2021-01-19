package client

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/kanodoe/dremio_scrapper/internal/entity"
)

// NewLoggingMiddleware ...
func NewLoggingMiddleware(logger log.Logger) Middleware {
	return func(next DremioClient) DremioClient {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   DremioClient
	logger log.Logger
}

func (l *loggingMiddleware) Login(ctx context.Context) (dlr *entity.DremioLoginResponse, err error) {
	defer func(begin time.Time) {
		l.logger.Log("method", "Login", time.Since(begin), "Dremio Login Response", dlr, "err", err)
	}(time.Now())
	return l.next.Login(ctx)
}
