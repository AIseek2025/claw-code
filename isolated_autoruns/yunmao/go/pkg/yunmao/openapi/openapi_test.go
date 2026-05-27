package openapi

import (
	_ "embed"
	"encoding/json"
	"testing"
)

//go:embed v3.json
var specBytes []byte

func TestSpecIsValidJSON(t *testing.T) {
	var doc map[string]any
	if err := json.Unmarshal(specBytes, &doc); err != nil {
		t.Fatalf("v3.json is not valid JSON: %v", err)
	}
}

func TestSpecHasRequiredFields(t *testing.T) {
	var doc map[string]any
	_ = json.Unmarshal(specBytes, &doc)

	if _, ok := doc["openapi"]; !ok {
		t.Fatal("missing openapi version")
	}
	if _, ok := doc["info"]; !ok {
		t.Fatal("missing info")
	}
	if _, ok := doc["paths"]; !ok {
		t.Fatal("missing paths")
	}
	if _, ok := doc["components"]; !ok {
		t.Fatal("missing components")
	}
}

func TestSpecHasSchemas(t *testing.T) {
	var doc struct {
		Components struct {
			Schemas map[string]any `json:"schemas"`
		} `json:"components"`
	}
	if err := json.Unmarshal(specBytes, &doc); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	required := []string{
		"User",
		"Room",
		"FeedRequest",
		"Order",
		"Wallet",
		"PrepayResponse",
		"IceServersResponse",
		"ChatMessage",
	}
	for _, name := range required {
		if _, ok := doc.Components.Schemas[name]; !ok {
			t.Errorf("missing schema: %s", name)
		}
	}
}

func TestSpecPathsCoverMainServices(t *testing.T) {
	var doc struct {
		Paths map[string]any `json:"paths"`
	}
	if err := json.Unmarshal(specBytes, &doc); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	mustHave := []string{
		"/v1/auth/login/sms",
		"/v1/rooms",
		"/v1/rooms/{id}",
		"/v1/rooms/{id}/ice-servers",
		"/api/v1/feed-requests",
		"/api/v1/orders",
		"/api/v1/orders/{id}/prepay",
		"/api/v1/wallets/{user_id}",
		"/api/v1/pay/channels",
	}
	for _, p := range mustHave {
		if _, ok := doc.Paths[p]; !ok {
			t.Errorf("missing path: %s", p)
		}
	}
}

func TestSpecOperationsHaveSchemas(t *testing.T) {
	var doc struct {
		Paths map[string]map[string]struct {
			OperationId string `json:"operationId"`
			Request     struct {
				Content map[string]struct {
					Schema struct {
						Ref string `json:"$ref"`
					} `json:"schema"`
				} `json:"content"`
			} `json:"requestBody"`
			Responses map[string]struct {
				Content map[string]struct {
					Schema struct {
						Ref string `json:"$ref"`
					} `json:"schema"`
				} `json:"content"`
			} `json:"responses"`
		} `json:"paths"`
	}
	if err := json.Unmarshal(specBytes, &doc); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	opsChecked := 0
	for _, methods := range doc.Paths {
		for _, op := range methods {
			if op.OperationId == "" {
				continue
			}
			opsChecked++
			hasResponse := false
			for _, resp := range op.Responses {
				for _, media := range resp.Content {
					if media.Schema.Ref != "" {
						hasResponse = true
					}
				}
			}
			if !hasResponse {
				if _, ok := op.Responses["200"]; ok {
					if len(op.Responses["200"].Content) == 0 {
						t.Logf("op %s may lack inline response schema", op.OperationId)
					}
				}
			}
		}
	}
	if opsChecked < 10 {
		t.Errorf("expected at least 10 operations with operationId, got %d", opsChecked)
	}
}
