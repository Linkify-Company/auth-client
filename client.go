package domain

import (
	"encoding/json"
	"fmt"
	"github.com/Linkify-Company/common_utils/errify"
	"net/http"
	"time"
)

type Service struct {
	host    string
	port    int
	timeout time.Duration
}

func NewClient(
	host string,
	port int,
	timeout time.Duration,
) *Service {
	return &Service{
		host:    host,
		port:    port,
		timeout: timeout,
	}
}

func (s *Service) Check(r *http.Request) (*AuthData, errify.IError) {
	client := &http.Client{Timeout: s.timeout}
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprint(s.host, s.port, "/srv-auth/api/v1/auth/check"),
		nil,
	)
	if err != nil {
		return nil, errify.NewInternalServerError(err.Error(), "Check/NewRequest")
	}
	req.Header.Set("Content-Type", "application/json")

	for _, cookie := range r.Cookies() {
		req.AddCookie(cookie)
	}
	req.Header = r.Header
	req.TLS = r.TLS

	resp, err := client.Do(req)
	if err != nil {
		return nil, errify.NewInternalServerError(err.Error(), "Check/Do")
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusUnauthorized {
			return nil, errify.NewUnauthorizedError("unauthorized", "Unauthorized", "Check/Do (StatusCode)")
		}
		if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusMethodNotAllowed {
			return nil, errify.NewInternalServerError("Auth service not found", "Check/Do (StatusCode)")
		}
		return nil, errify.NewInternalServerError(resp.Status, "Check/Do (StatusCode)")
	}

	auth := struct {
		Value AuthData `json:"value"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(&auth)
	if err != nil {
		return nil, errify.NewInternalServerError(err.Error(), "Check/NewDecoder")
	}
	if auth.Value.ID <= 0 {
		return nil, errify.NewInternalServerError("user is empty", "Check")
	}

	return &auth.Value, nil
}
