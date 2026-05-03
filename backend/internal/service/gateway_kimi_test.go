//go:build unit

package service

import (
	"context"
	"testing"
)

func TestKimiAPIKeyUsesAnthropicAPIKeyToken(t *testing.T) {
	svc := newMiniMaxGatewayTestService()
	account := &Account{
		ID:       1,
		Platform: PlatformKimi,
		Type:     AccountTypeAPIKey,
		Credentials: map[string]any{
			"api_key": "sk-kimi",
		},
	}

	token, tokenType, err := svc.GetAccessToken(context.Background(), account)
	if err != nil {
		t.Fatalf("GetAccessToken returned error: %v", err)
	}
	if token != "sk-kimi" {
		t.Fatalf("token = %q, want %q", token, "sk-kimi")
	}
	if tokenType != "apikey" {
		t.Fatalf("tokenType = %q, want %q", tokenType, "apikey")
	}
}

func TestKimiBuildUpstreamRequest(t *testing.T) {
	svc := newMiniMaxGatewayTestService()
	account := &Account{
		ID:          1,
		Platform:    PlatformKimi,
		Type:        AccountTypeAPIKey,
		Credentials: map[string]any{},
	}

	req, err := svc.buildUpstreamRequest(
		context.Background(),
		nil,
		account,
		[]byte(`{"model":"kimi-for-coding","messages":[{"role":"user","content":"hi"}],"max_tokens":16}`),
		"sk-kimi",
		"apikey",
		"kimi-for-coding",
		false,
		false,
	)
	if err != nil {
		t.Fatalf("buildUpstreamRequest returned error: %v", err)
	}

	if got := req.URL.String(); got != kimiAPIURL {
		t.Fatalf("url = %q, want %q", got, kimiAPIURL)
	}
	if got := getHeaderRaw(req.Header, "x-api-key"); got != "sk-kimi" {
		t.Fatalf("x-api-key = %q, want %q", got, "sk-kimi")
	}
	if got := getHeaderRaw(req.Header, "authorization"); got != "" {
		t.Fatalf("authorization = %q, want empty", got)
	}
}

func TestKimiBuildCountTokensRequest(t *testing.T) {
	svc := newMiniMaxGatewayTestService()
	account := &Account{
		ID:          1,
		Platform:    PlatformKimi,
		Type:        AccountTypeAPIKey,
		Credentials: map[string]any{},
	}

	req, err := svc.buildCountTokensRequest(
		context.Background(),
		nil,
		account,
		[]byte(`{"model":"kimi-for-coding","messages":[{"role":"user","content":"hi"}]}`),
		"sk-kimi",
		"apikey",
		"kimi-for-coding",
		false,
	)
	if err != nil {
		t.Fatalf("buildCountTokensRequest returned error: %v", err)
	}

	if got := req.URL.String(); got != kimiCountTokensURL {
		t.Fatalf("url = %q, want %q", got, kimiCountTokensURL)
	}
	if got := getHeaderRaw(req.Header, "x-api-key"); got != "sk-kimi" {
		t.Fatalf("x-api-key = %q, want %q", got, "sk-kimi")
	}
	if got := getHeaderRaw(req.Header, "authorization"); got != "" {
		t.Fatalf("authorization = %q, want empty", got)
	}
}

func TestKimiDefaultModelMapping(t *testing.T) {
	account := &Account{
		Platform:    PlatformKimi,
		Type:        AccountTypeAPIKey,
		Credentials: map[string]any{},
	}

	if got := account.GetMappedModel("claude-sonnet-4-5"); got != "kimi-for-coding" {
		t.Fatalf("mapped model = %q, want %q", got, "kimi-for-coding")
	}
	if got := account.GetMappedModel("kimi-for-coding"); got != "kimi-for-coding" {
		t.Fatalf("mapped model = %q, want %q", got, "kimi-for-coding")
	}
}
