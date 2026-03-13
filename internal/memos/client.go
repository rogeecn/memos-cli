package memos

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Client struct {
	baseURL     string
	apiKey      string
	adminAPIKey string
	httpClient  *http.Client
}

type ListMemosParams struct {
	Filter    string
	PageSize  int
	PageToken string
}

type UpdateMemoPayload struct {
	Content    string `json:"content,omitempty"`
	Visibility string `json:"visibility,omitempty"`
}

type Memo struct {
	Name       string `json:"name"`
	Content    string `json:"content"`
	Visibility string `json:"visibility,omitempty"`
	CreateTime string `json:"createTime,omitempty"`
}

type MemoPayload struct {
	Content    string `json:"content"`
	Visibility string `json:"visibility,omitempty"`
}

type ListMemosResponse struct {
	Memos         []Memo `json:"memos"`
	NextPageToken string `json:"nextPageToken,omitempty"`
}

type ListUsersResponse struct {
	Users []map[string]any `json:"users"`
}

func NewClient(baseURL, apiKey, adminAPIKey string) *Client {
	return &Client{
		baseURL:     strings.TrimRight(baseURL, "/"),
		apiKey:      apiKey,
		adminAPIKey: adminAPIKey,
		httpClient:  &http.Client{},
	}
}

func (c *Client) ListMemos(params ListMemosParams) (ListMemosResponse, error) {
	query := url.Values{}
	if params.Filter != "" {
		query.Set("filter", params.Filter)
	}
	if params.PageSize > 0 {
		query.Set("pageSize", strconv.Itoa(params.PageSize))
	}
	if params.PageToken != "" {
		query.Set("pageToken", params.PageToken)
	}

	var response ListMemosResponse
	err := c.doJSON(http.MethodGet, "/api/v1/memos", query, nil, c.apiKey, &response)
	return response, err
}

func (c *Client) ListUsers() (ListUsersResponse, error) {
	if strings.TrimSpace(c.adminAPIKey) == "" {
		return ListUsersResponse{}, errors.New("admin api key is required")
	}

	var response ListUsersResponse
	err := c.doJSON(http.MethodGet, "/api/v1/users", nil, nil, c.adminAPIKey, &response)
	return response, err
}

func (c *Client) GetMemo(memoID string) (Memo, error) {
	var response Memo
	err := c.doJSON(http.MethodGet, "/api/v1/"+normalizeMemoName(memoID), nil, nil, c.apiKey, &response)
	return response, err
}

func (c *Client) CreateMemo(payload MemoPayload) (Memo, error) {
	var response Memo
	err := c.doJSON(http.MethodPost, "/api/v1/memos", nil, payload, c.apiKey, &response)
	return response, err
}

func (c *Client) UpdateMemo(memoID string, payload UpdateMemoPayload) (Memo, error) {
	var response Memo
	err := c.doJSON(http.MethodPatch, "/api/v1/"+normalizeMemoName(memoID), nil, payload, c.apiKey, &response)
	return response, err
}

func (c *Client) DeleteMemo(memoID string) error {
	return c.doJSON(http.MethodDelete, "/api/v1/"+normalizeMemoName(memoID), nil, nil, c.apiKey, nil)
}

func (c *Client) CreateComment(memoID string, payload MemoPayload) (Memo, error) {
	var response Memo
	path := fmt.Sprintf("/api/v1/%s/comments", normalizeMemoName(memoID))
	err := c.doJSON(http.MethodPost, path, nil, payload, c.apiKey, &response)
	return response, err
}

func (c *Client) doJSON(method, path string, query url.Values, body any, token string, target any) error {
	requestURL := c.baseURL + path
	if len(query) > 0 {
		requestURL += "?" + query.Encode()
	}

	var bodyReader io.Reader
	if body != nil {
		payload, err := json.Marshal(body)
		if err != nil {
			return err
		}
		bodyReader = bytes.NewReader(payload)
	}

	req, err := http.NewRequest(method, requestURL, bodyReader)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if strings.TrimSpace(token) != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("memos api error: %s", strings.TrimSpace(string(body)))
	}

	if target == nil {
		return nil
	}
	return json.NewDecoder(resp.Body).Decode(target)
}

func normalizeMemoName(memoID string) string {
	memoID = strings.TrimSpace(memoID)
	memoID = strings.TrimPrefix(memoID, "memos/")
	return "memos/" + memoID
}
