package redash

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// RedashQuery represents a Redash query
type RedashQuery struct {
	ID                int                    `json:"id"`
	Name              string                 `json:"name"`
	Description       string                 `json:"description"`
	Query             string                 `json:"query"`
	DataSourceID      int                    `json:"data_source_id"`
	LatestQueryDataID int                    `json:"latest_query_data_id"`
	IsArchived        bool                   `json:"is_archived"`
	CreatedAt         string                 `json:"created_at"`
	UpdatedAt         string                 `json:"updated_at"`
	Runtime           float64                `json:"runtime"`
	Options           map[string]interface{} `json:"options"`
	Visualizations    []RedashVisualization  `json:"visualizations"`
}

// CreateQueryRequest represents a request to create a new query
type CreateQueryRequest struct {
	Name         string                 `json:"name"`
	DataSourceId int                    `json:"data_source_id"`
	Query        string                 `json:"query"`
	Description  *string                `json:"description,omitempty"`
	Options      map[string]interface{} `json:"options,omitempty"`
	Schedule     map[string]interface{} `json:"schedule,omitempty"`
	Tags         []string               `json:"tags,omitempty"`
}

// UpdateQueryRequest represents a request to update an existing query
type UpdateQueryRequest struct {
	Name         *string                `json:"name,omitempty"`
	DataSourceId *int                   `json:"data_source_id,omitempty"`
	Query        *string                `json:"query,omitempty"`
	Description  *string                `json:"description,omitempty"`
	Options      map[string]interface{} `json:"options,omitempty"`
	Schedule     map[string]interface{} `json:"schedule,omitempty"`
	Tags         []string               `json:"tags,omitempty"`
	IsArchived   *bool                  `json:"is_archived,omitempty"`
	IsDraft      *bool                  `json:"is_draft,omitempty"`
}

// RedashVisualization represents a Redash visualization
type RedashVisualization struct {
	ID          int                    `json:"id"`
	Type        string                 `json:"type"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Options     map[string]interface{} `json:"options"`
	QueryID     int                    `json:"query_id"`
}

// RedashQueryResult represents the result of a query execution
type RedashQueryResult struct {
	ID           int    `json:"id"`
	QueryID      int    `json:"query_id"`
	DataSourceID int    `json:"data_source_id"`
	QueryHash    string `json:"query_hash"`
	Query        string `json:"query"`
	Data         struct {
		Columns []struct {
			Name         string `json:"name"`
			Type         string `json:"type"`
			FriendlyName string `json:"friendly_name"`
		} `json:"columns"`
		Rows []map[string]interface{} `json:"rows"`
	} `json:"data"`
	Runtime     float64 `json:"runtime"`
	RetrievedAt string  `json:"retrieved_at"`
}

// RedashDashboard represents a Redash dashboard
type RedashDashboard struct {
	ID                      int      `json:"id"`
	Name                    string   `json:"name"`
	Slug                    string   `json:"slug"`
	Tags                    []string `json:"tags"`
	IsArchived              bool     `json:"is_archived"`
	IsDraft                 bool     `json:"is_draft"`
	CreatedAt               string   `json:"created_at"`
	UpdatedAt               string   `json:"updated_at"`
	Version                 int      `json:"version"`
	DashboardFiltersEnabled bool     `json:"dashboard_filters_enabled"`
	Widgets                 []struct {
		ID            int `json:"id"`
		Visualization *struct {
			ID          int                    `json:"id"`
			Type        string                 `json:"type"`
			Name        string                 `json:"name"`
			Description string                 `json:"description"`
			Options     map[string]interface{} `json:"options"`
			QueryID     int                    `json:"query_id"`
		} `json:"visualization,omitempty"`
		Text        *string                `json:"text,omitempty"`
		Width       int                    `json:"width"`
		Options     map[string]interface{} `json:"options"`
		DashboardID int                    `json:"dashboard_id"`
	} `json:"widgets"`
}

// RedashClient handles communication with the Redash API
type RedashClient struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
}

// NewRedashClient creates a new Redash client
func NewRedashClient() (*RedashClient, error) {
	baseURL := os.Getenv("REDASH_URL")
	if baseURL == "" {
		return nil, fmt.Errorf("REDASH_URL environment variable is required")
	}

	apiKey := os.Getenv("REDASH_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("REDASH_API_KEY environment variable is required")
	}

	// Ensure baseURL ends with a slash
	if baseURL[len(baseURL)-1] != '/' {
		baseURL += "/"
	}

	return &RedashClient{
		BaseURL:    baseURL,
		APIKey:     apiKey,
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
	}, nil
}

// makeRequest makes an HTTP request to the Redash API
func (c *RedashClient) makeRequest(method, endpoint string, body interface{}) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, c.BaseURL+endpoint, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Key "+c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	return resp, nil
}

// GetQueries gets all queries with pagination
func (c *RedashClient) GetQueries(page, pageSize int, search string) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf("/api/queries?page=%d&page_size=%d", page, pageSize)
	if search != "" {
		endpoint += "&q=" + search
	}

	resp, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil
}

// GetQuery gets a specific query by ID
func (c *RedashClient) GetQuery(queryID int) (*RedashQuery, error) {
	endpoint := fmt.Sprintf("/api/queries/%d", queryID)

	resp, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	var query RedashQuery
	if err := json.NewDecoder(resp.Body).Decode(&query); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &query, nil
}

// CreateQuery creates a new query
func (c *RedashClient) CreateQuery(queryData CreateQueryRequest) (*RedashQuery, error) {
	resp, err := c.makeRequest("POST", "/api/queries", queryData)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var query RedashQuery
	if err := json.NewDecoder(resp.Body).Decode(&query); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &query, nil
}

// UpdateQuery updates an existing query
func (c *RedashClient) UpdateQuery(queryID int, updateData UpdateQueryRequest) (*RedashQuery, error) {
	endpoint := fmt.Sprintf("/api/queries/%d", queryID)

	resp, err := c.makeRequest("POST", endpoint, updateData)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var query RedashQuery
	if err := json.NewDecoder(resp.Body).Decode(&query); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &query, nil
}

// ArchiveQuery archives a query
func (c *RedashClient) ArchiveQuery(queryID int) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf("/api/queries/%d", queryID)
	updateData := UpdateQueryRequest{
		IsArchived: &[]bool{true}[0],
	}

	resp, err := c.makeRequest("POST", endpoint, updateData)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return map[string]interface{}{"success": true}, nil
}

// GetDataSources gets all data sources
func (c *RedashClient) GetDataSources() ([]interface{}, error) {
	resp, err := c.makeRequest("GET", "/api/data_sources", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	var dataSources []interface{}
	if err := json.NewDecoder(resp.Body).Decode(&dataSources); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return dataSources, nil
}

// GetQueryResult gets cached query results by result ID
func (c *RedashClient) GetQueryResult(resultID int) (*RedashQueryResult, error) {
	endpoint := fmt.Sprintf("/api/query_results/%d", resultID)

	resp, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	// Read the full response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// The API returns the result wrapped in a query_result object
	var wrappedResult struct {
		QueryResult RedashQueryResult `json:"query_result"`
	}
	if err := json.Unmarshal(bodyBytes, &wrappedResult); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &wrappedResult.QueryResult, nil
}

// ExecuteQueryFresh executes a query with fresh data, skipping any cached results
func (c *RedashClient) ExecuteQueryFresh(queryID int, parameters map[string]interface{}) (*RedashQueryResult, error) {
	// Force fresh execution by directly calling the results endpoint
	endpoint := fmt.Sprintf("/api/queries/%d/results", queryID)

	var body interface{}
	if parameters != nil {
		body = map[string]interface{}{
			"parameters": parameters,
		}
	}

	resp, err := c.makeRequest("POST", endpoint, body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body first
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check if this is a job response (either 202 Accepted or 200 with job)
	var jobResponse struct {
		Job struct {
			ID     string `json:"id"`
			Status int    `json:"status"`
		} `json:"job"`
	}

	isJobResponse := false
	if resp.StatusCode == http.StatusAccepted {
		isJobResponse = true
	} else if resp.StatusCode == http.StatusOK {
		// Check if this is a job response with status 200
		if err := json.Unmarshal(bodyBytes, &jobResponse); err == nil && jobResponse.Job.ID != "" {
			isJobResponse = true
		}
	}

	if isJobResponse {
		// Re-parse the job response if we haven't already
		if jobResponse.Job.ID == "" {
			if err := json.Unmarshal(bodyBytes, &jobResponse); err != nil {
				return nil, fmt.Errorf("failed to decode job response: %w", err)
			}
		}

		// Poll for results with improved logic
		result, err := c.pollQueryResults(jobResponse.Job.ID)
		if err != nil {
			return nil, fmt.Errorf("fresh query execution failed: %w", err)
		}
		return result, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// Check if the response is wrapped like cached results
	var wrappedResult struct {
		QueryResult RedashQueryResult `json:"query_result"`
	}
	if err := json.Unmarshal(bodyBytes, &wrappedResult); err == nil && wrappedResult.QueryResult.ID != 0 {
		return &wrappedResult.QueryResult, nil
	}

	// Try direct parsing
	var result RedashQueryResult
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// pollQueryResults polls for query execution results with improved logic
func (c *RedashClient) pollQueryResults(jobID string) (*RedashQueryResult, error) {
	timeout := time.After(60 * time.Second)   // Reduced timeout to 1 minute
	ticker := time.NewTicker(1 * time.Second) // Reduced polling interval for faster response
	defer ticker.Stop()

	var queryID int // Store the query ID for fallback
	pollCount := 0  // Track number of polls

	for {
		select {
		case <-timeout:
			return nil, fmt.Errorf("query execution timeout after 60 seconds (polled %d times)", pollCount)
		case <-ticker.C:
			pollCount++
			endpoint := fmt.Sprintf("/api/jobs/%s", jobID)
			resp, err := c.makeRequest("GET", endpoint, nil)
			if err != nil {
				// Log error but continue polling
				if pollCount > 30 { // After 30 seconds of errors, give up
					return nil, fmt.Errorf("too many network errors during polling: %w", err)
				}
				continue
			}

			bodyBytes, err := io.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				continue
			}

			// Check for HTTP errors
			if resp.StatusCode != http.StatusOK {
				if pollCount > 10 { // Give up after 10 failed status codes
					return nil, fmt.Errorf("job polling failed with status %d: %s", resp.StatusCode, string(bodyBytes))
				}
				continue
			}

			// Parse as generic JSON first to handle different response formats
			var genericResponse map[string]interface{}
			if err := json.Unmarshal(bodyBytes, &genericResponse); err != nil {
				if pollCount > 20 {
					return nil, fmt.Errorf("failed to parse job response after %d attempts: %w", pollCount, err)
				}
				continue
			}

			// Extract job information from generic response
			jobData, ok := genericResponse["job"].(map[string]interface{})
			if !ok {
				if pollCount > 20 {
					return nil, fmt.Errorf("no job data found in response after %d attempts", pollCount)
				}
				continue
			}

			// Get job status
			status, ok := jobData["status"].(float64)
			if !ok {
				continue
			}

			// Store query ID for potential fallback
			if queryIDFloat, ok := jobData["query_id"].(float64); ok {
				queryID = int(queryIDFloat)
			}

			switch int(status) {
			case 3: // Completed
				// Try to get result from query_result field
				if queryResult, ok := jobData["query_result"].(map[string]interface{}); ok {
					// Parse the query result
					resultBytes, _ := json.Marshal(queryResult)
					var result RedashQueryResult
					if err := json.Unmarshal(resultBytes, &result); err == nil {
						return &result, nil
					}
				}

				// Try to get result from result field (could be object or ID)
				if resultData, ok := jobData["result"]; ok {
					// Check if it's a full object
					if resultObj, ok := resultData.(map[string]interface{}); ok {
						resultBytes, _ := json.Marshal(resultObj)
						var result RedashQueryResult
						if err := json.Unmarshal(resultBytes, &result); err == nil {
							return &result, nil
						}
					}
					// Check if it's a result ID (number)
					if resultID, ok := resultData.(float64); ok && resultID > 0 {
						if result, err := c.GetQueryResult(int(resultID)); err == nil {
							return result, nil
						}
					}
				}

				// Try to get query_result_id
				if queryResultID, ok := jobData["query_result_id"].(float64); ok && queryResultID > 0 {
					if result, err := c.GetQueryResult(int(queryResultID)); err == nil {
						return result, nil
					}
				}

				// Check if there's a top-level query_result in the response
				if queryResult, ok := genericResponse["query_result"].(map[string]interface{}); ok {
					resultBytes, _ := json.Marshal(queryResult)
					var result RedashQueryResult
					if err := json.Unmarshal(resultBytes, &result); err == nil {
						return &result, nil
					}
				}

				// If still no result but we have a query ID, try to get the latest result
				if queryID > 0 {
					if query, err := c.GetQuery(queryID); err == nil && query.LatestQueryDataID > 0 {
						if result, err := c.GetQueryResult(query.LatestQueryDataID); err == nil {
							return result, nil
						}
					}
				}

				// If all else fails, return an error
				return nil, fmt.Errorf("query completed but no result found in job response after %d polls", pollCount)

			case 4: // Failed
				errorMsg := "query execution failed"
				if errorStr, ok := jobData["error"].(string); ok && errorStr != "" {
					errorMsg = fmt.Sprintf("query execution failed: %s", errorStr)
				}
				return nil, fmt.Errorf(errorMsg)

			case 1, 2: // Pending or Started
				// Check if we already have a result ID even if job is still running
				if queryResultID, ok := jobData["query_result_id"].(float64); ok && queryResultID > 0 {
					if result, err := c.GetQueryResult(int(queryResultID)); err == nil {
						return result, nil
					}
				}

				// Prevent infinite polling - if we've been polling for too long, try fallback
				if pollCount > 45 && queryID > 0 { // After 45 seconds, try fallback
					if query, err := c.GetQuery(queryID); err == nil && query.LatestQueryDataID > 0 {
						if result, err := c.GetQueryResult(query.LatestQueryDataID); err == nil {
							return result, nil
						}
					}
				}
				continue // Keep polling

			default:
				// Unknown status - if we've been polling too long, give up
				if pollCount > 30 {
					return nil, fmt.Errorf("unknown job status %d after %d polls", int(status), pollCount)
				}
				continue // Keep polling for unknown statuses
			}
		}
	}
}

// GetDashboards gets all dashboards with pagination
func (c *RedashClient) GetDashboards(page, pageSize int) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf("/api/dashboards?page=%d&page_size=%d", page, pageSize)

	resp, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil
}

// GetDashboard gets a specific dashboard by ID
func (c *RedashClient) GetDashboard(dashboardID int) (*RedashDashboard, error) {
	endpoint := fmt.Sprintf("/api/dashboards/%d", dashboardID)

	resp, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	var dashboard RedashDashboard
	if err := json.NewDecoder(resp.Body).Decode(&dashboard); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &dashboard, nil
}

// GetVisualization gets a specific visualization by ID
func (c *RedashClient) GetVisualization(visualizationID int) (*RedashVisualization, error) {
	endpoint := fmt.Sprintf("/api/visualizations/%d", visualizationID)

	resp, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	var visualization RedashVisualization
	if err := json.NewDecoder(resp.Body).Decode(&visualization); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &visualization, nil
}
