//go:build unit

package service

import (
	"context"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
)

func newMiniMaxGatewayTestService() *GatewayService {
	return &GatewayService{
		cfg: &config.Config{
			Security: config.SecurityConfig{
				URLAllowlist: config.URLAllowlistConfig{Enabled: false},
			},
		},
	}
}

func TestMiniMaxAPIKeyUsesBearerToken(t *testing.T) {
	svc := newMiniMaxGatewayTestService()
	account := &Account{
		ID:       1,
		Platform: PlatformMiniMax,
		Type:     AccountTypeAPIKey,
		Credentials: map[string]any{
			"api_key": "sk-minimax",
		},
	}

	token, tokenType, err := svc.GetAccessToken(context.Background(), account)
	if err != nil {
		t.Fatalf("GetAccessToken returned error: %v", err)
	}
	if token != "sk-minimax" {
		t.Fatalf("token = %q, want %q", token, "sk-minimax")
	}
	if tokenType != "bearer" {
		t.Fatalf("tokenType = %q, want %q", tokenType, "bearer")
	}
}

func TestMiniMaxBuildUpstreamRequest(t *testing.T) {
	svc := newMiniMaxGatewayTestService()
	account := &Account{
		ID:          1,
		Platform:    PlatformMiniMax,
		Type:        AccountTypeAPIKey,
		Credentials: map[string]any{},
	}

	req, err := svc.buildUpstreamRequest(
		context.Background(),
		nil,
		account,
		[]byte(`{"model":"MiniMax-M2.7","messages":[{"role":"user","content":"hi"}],"max_tokens":16}`),
		"sk-minimax",
		"bearer",
		"MiniMax-M2.7",
		false,
		false,
	)
	if err != nil {
		t.Fatalf("buildUpstreamRequest returned error: %v", err)
	}

	if got := req.URL.String(); got != miniMaxAPIURL {
		t.Fatalf("url = %q, want %q", got, miniMaxAPIURL)
	}
	if got := getHeaderRaw(req.Header, "authorization"); got != "Bearer sk-minimax" {
		t.Fatalf("authorization = %q, want %q", got, "Bearer sk-minimax")
	}
	if got := getHeaderRaw(req.Header, "x-api-key"); got != "" {
		t.Fatalf("x-api-key = %q, want empty", got)
	}
}

func TestMiniMaxBuildCountTokensRequest(t *testing.T) {
	svc := newMiniMaxGatewayTestService()
	account := &Account{
		ID:          1,
		Platform:    PlatformMiniMax,
		Type:        AccountTypeAPIKey,
		Credentials: map[string]any{},
	}

	req, err := svc.buildCountTokensRequest(
		context.Background(),
		nil,
		account,
		[]byte(`{"model":"MiniMax-M2.7","messages":[{"role":"user","content":"hi"}]}`),
		"sk-minimax",
		"bearer",
		"MiniMax-M2.7",
		false,
	)
	if err != nil {
		t.Fatalf("buildCountTokensRequest returned error: %v", err)
	}

	if got := req.URL.String(); got != miniMaxCountTokensURL {
		t.Fatalf("url = %q, want %q", got, miniMaxCountTokensURL)
	}
	if got := getHeaderRaw(req.Header, "authorization"); got != "Bearer sk-minimax" {
		t.Fatalf("authorization = %q, want %q", got, "Bearer sk-minimax")
	}
	if got := getHeaderRaw(req.Header, "x-api-key"); got != "" {
		t.Fatalf("x-api-key = %q, want empty", got)
	}
}

func TestMiniMaxDefaultModelMapping(t *testing.T) {
	account := &Account{
		Platform:    PlatformMiniMax,
		Type:        AccountTypeAPIKey,
		Credentials: map[string]any{},
	}

	if got := account.GetMappedModel("claude-sonnet-4-5"); got != "MiniMax-M2.7" {
		t.Fatalf("mapped model = %q, want %q", got, "MiniMax-M2.7")
	}
	if got := account.GetMappedModel("MiniMax-M2.7-highspeed"); got != "MiniMax-M2.7-highspeed" {
		t.Fatalf("mapped model = %q, want %q", got, "MiniMax-M2.7-highspeed")
	}
}
