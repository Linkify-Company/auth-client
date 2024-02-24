package domain

import (
	"encoding/json"
	"fmt"
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

func (s *Service) Check(r *http.Request) (*AuthData, int, error) {
	client := &http.Client{Timeout: s.timeout}
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprint(s.host, s.port, "/srv-auth/api/v1/auth/check"),
		nil,
	)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	req.Header.Set("Content-Type", "application/json")

	for _, cookie := range r.Cookies() {
		req.AddCookie(cookie)
	}
	req.TLS = r.TLS

	resp, err := client.Do(req)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode, nil
	}

	var user AuthData
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return &user, http.StatusOK, err
}
