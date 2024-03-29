package domain

import (
	"context"
	"github.com/Linkify-Company/common_utils/errify"
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
		authData, err := m.Check(r, m.log)
		if err != nil {
			response.Error(w, err.JoinLoc("AuthHandler"), m.log)
			return
		}

		setCtx(r, authData)
		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) AuthFuncWithRoles(roles ...Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := GetAuthData(r)
			if !ok {
				response.Error(w,
					errify.NewInternalServerError("could not get data from context",
						"AuthFuncWithRoles/GetAuthData"), m.log)
				return
			}
			for _, role := range roles {
				if user.Role == role {
					next.ServeHTTP(w, r)
					return
				}
			}
			response.Error(w,
				errify.NewUnauthorizedError("not enough authority",
					"User is not authorized", "AuthFuncWithRoles/GetAuthData"), m.log)
		})
	}
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
