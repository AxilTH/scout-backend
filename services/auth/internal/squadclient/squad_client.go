// internal/squad_client.go
package squadclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// SquadClient provides methods to interact with Squad Service API
type SquadClient struct {
	baseURL    string
	httpClient *http.Client
	timeout    time.Duration
}

// NewSquadClient creates a new Squad Service client
func NewSquadClient(baseURL string, timeout time.Duration) *SquadClient {
	return &SquadClient{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: timeout},
		timeout:    timeout,
	}
}

// CreateMembership creates a new squad membership
func (c *SquadClient) CreateMembership(ctx context.Context, userID, squadID, roleID int64) error {
	// Формируем запрос согласно API Squad Service
	// На основе анализа Squad Service handler, ожидаем:
	// POST /api/v1/memberships
	// {
	//   "user_id": 123,
	//   "squad_id": 456,
	//   "role_id": 789
	// }

	membershipReq := struct {
		UserID int64 `json:"user_id"`
		SquadID int64 `json:"squad_id"`
		RoleID int64 `json:"role_id"`
	}{
		UserID: userID,
		SquadID: squadID,
		RoleID: roleID,
	}

	jsonData, err := json.Marshal(membershipReq)
	if err != nil {
		return fmt.Errorf("failed to marshal membership request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		fmt.Sprintf("%s/api/v1/memberships", c.baseURL), bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	switch resp.StatusCode {
	case http.StatusCreated: // 201
		return nil
	case http.StatusBadRequest: // 400
		// Попытка прочитать тело ошибки для лучшего диагностирования
		var errorResp struct {
			Error string `json:"error"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err == nil {
			return fmt.Errorf("bad request: %s", errorResp.Error)
		}
		return fmt.Errorf("bad request")
	case http.StatusUnauthorized: // 401
		return fmt.Errorf("unauthorized")
	case http.StatusForbidden: // 403
		return fmt.Errorf("forbidden")
	case http.StatusConflict: // 409
		return fmt.Errorf("conflict: membership may already exist")
	default:
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
}

// GetUserSquads returns squad IDs for a user
func (c *SquadClient) GetUserSquads(ctx context.Context, userID int64) ([]int64, error) {
	// GET /api/v1/users/{userID}/squads
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		fmt.Sprintf("%s/api/v1/users/%d/squads", c.baseURL, userID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var squadIDs []int64
	if err := json.NewDecoder(resp.Body).Decode(&squadIDs); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return squadIDs, nil
}

// ValidateMembership checks if user belongs to squad
func (c *SquadClient) ValidateMembership(ctx context.Context, userID, squadID int64) (bool, error) {
	squads, err := c.GetUserSquads(ctx, userID)
	if err != nil {
		return false, err
	}
	for _, id := range squads {
		if id == squadID {
			return true, nil
		}
	}
	return false, nil
}