package homeassistantApi

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"wallbox-monitor/variables"
)

func CallHomeAssisatantApiSet(apiCall string, data *strings.Reader) error {
	client := &http.Client{}
	req, err := http.NewRequest("POST", variables.UrlHomeAssistant+apiCall, data)
	if err != nil {
		slog.Error("Failed to create request", "error", err)
		return err
	}

	req.Header.Set("Authorization", "Bearer "+variables.TokenHomeAssistant)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		slog.Error("API request failed", "error", err)
		return err
	}
	defer resp.Body.Close()
	bodyText, _ := io.ReadAll(resp.Body)
	_ = bodyText

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("API request failed: %s", resp.Status)
		slog.Error("API request failed", "status", resp.Status, "url", req.URL, "method", req.Method)
		return err
	}

	if variables.HomeAssistantLogApiCallResult200 {
		slog.Info("API request successful", "status", resp.Status, "url", req.URL, "method", req.Method)
	}

	return nil
}
