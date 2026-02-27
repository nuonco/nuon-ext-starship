package resolve

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type orgResponse struct {
	Name string `json:"name"`
}

type installResponse struct {
	Name string `json:"name"`
}

var httpClient = &http.Client{Timeout: 3 * time.Second}

func OrgName(ctx context.Context, apiURL, apiToken, orgID string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL+"/v1/orgs/current", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+apiToken)
	req.Header.Set("X-Nuon-Org-ID", orgID)

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API returned %d: %s", resp.StatusCode, string(body))
	}

	var org orgResponse
	if err := json.NewDecoder(resp.Body).Decode(&org); err != nil {
		return "", err
	}
	return org.Name, nil
}

func InstallName(ctx context.Context, apiURL, apiToken, orgID, installID string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL+"/v1/installs/"+installID, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+apiToken)
	req.Header.Set("X-Nuon-Org-ID", orgID)

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API returned %d: %s", resp.StatusCode, string(body))
	}

	var install installResponse
	if err := json.NewDecoder(resp.Body).Decode(&install); err != nil {
		return "", err
	}
	return install.Name, nil
}
