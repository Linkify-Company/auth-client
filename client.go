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

type Parser struct {
	Value   interface{} `json:"value"`
	Message string      `json:"message"`
}

func (s *Service) Check(r *http.Request) (*AuthData, int, error) {
	client := &http.Client{Timeout: s.timeout}
	req, err := http.NewRequest(
		http.MethodGet,
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
	var parser = Parser{
		Value: &user,
	}
	err = json.NewDecoder(resp.Body).Decode(&parser)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if user.ID <= 0 {
		return nil, http.StatusInternalServerError, fmt.Errorf("user is empty")
	}

	return &user, http.StatusOK, err
}
