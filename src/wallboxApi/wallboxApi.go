package wallboxApi

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"
)

type BoolInt bool

func (b *BoolInt) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		*b = false
		return nil
	}

	var boolean bool
	if err := json.Unmarshal(data, &boolean); err == nil {
		*b = BoolInt(boolean)
		return nil
	}

	var integer int
	if err := json.Unmarshal(data, &integer); err == nil {
		*b = integer != 0
		return nil
	}

	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		switch str {
		case "1", "true", "True", "TRUE":
			*b = true
			return nil
		case "0", "false", "False", "FALSE":
			*b = false
			return nil
		}
	}

	return fmt.Errorf("invalid BoolInt value: %s", string(data))
}

type Status struct {
	Car int       `json:"car"`
	Ust int       `json:"trx"`
	Amp int       `json:"amp"`
	Nrg []float64 `json:"nrg"`
	Wh  float64   `json:"wh"`
	Dws float64   `json:"dws"`
	Psm int       `json:"psm"`
	Frc int       `json:"frc"`
	Tma []float64 `json:"tma"`
	Alw BoolInt   `json:"alw"`
}

func buildURL(host, port, filter string) string {
	u := url.URL{
		Scheme: "http",
		Host:   net.JoinHostPort(host, port),
		Path:   "/api/status",
	}
	if filter != "" {
		q := url.Values{}
		q.Set("filter", filter)
		u.RawQuery = q.Encode()
	}
	return u.String()
}

func FetchStatus(host, port, filter string) (*Status, error) {
	url := buildURL(host, port, filter)
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wallbox API returned status %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var status Status
	if err := json.Unmarshal(body, &status); err != nil {
		return nil, fmt.Errorf("failed to parse wallbox API response: %w", err)
	}

	return &status, nil
}
