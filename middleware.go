package domain

import (
	"context"
	"github.com/Linkify-Company/common_utils/logger"
	"github.com/Linkify-Company/common_utils/response"
	"net/http"
)

type Middleware struct {
	log logger.Logger
	*Service
}

const (
	AuthDataKey = "AuthData"
)

func NewMiddleware(log logger.Logger, s *Service) *Middleware {
	return &Middleware{log: log, Service: s}
}

func (m *Middleware) AuthHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authData, err := m.Check(r)
		if err != nil {
			response.Error(w, err.JoinLoc("AuthHandler"), m.log)
			return
		}

		setCtx(r, authData)
		next.ServeHTTP(w, r)
	})
}

func GetAuthData(req *http.Request) (*AuthData, bool) {
	v, ok := req.Context().Value(AuthDataKey).(*AuthData)
	return v, ok
}

func setCtx(r *http.Request, data *AuthData) {
	ctx := r.Context()
	ctx = context.WithValue(ctx, AuthDataKey, data)
	*r = *r.WithContext(ctx)
}
