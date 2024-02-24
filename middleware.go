package domain

import (
	"context"
	"log/slog"
	"net/http"
)

type Middleware struct {
	log *slog.Logger
	*Service
}

const (
	AuthDataKey = "AuthData"
)

func NewMiddleware(log *slog.Logger, s *Service) *Middleware {
	return &Middleware{log: log, Service: s}
}

func (m *Middleware) AuthHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authData, code, err := m.Check(r)
		if err != nil {
			m.log.Info(err.Error())
			w.WriteHeader(code)
			return
		}
		if code == http.StatusUnauthorized {
			m.log.Info("AuthHandler: User Unauthorized")
			w.WriteHeader(code)
			return
		}
		if code == http.StatusNotFound {
			m.log.Error("Auth Server Not Found")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if authData.ID == 0 {
			m.log.Error("AuthHandler: User id is nil")
			w.WriteHeader(http.StatusInternalServerError)
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
