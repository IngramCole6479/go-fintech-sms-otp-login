package main

import "testing"

func TestAuthorizeRiskAction(t *testing.T) {
	tests := []struct {
		name     string
		req      loginRequest
		verified bool
		wantErr  bool
	}{
		{"verified customer", loginRequest{"+14155550100", "123456"}, true, false},
		{"wrong code", loginRequest{"+14155550100", "000000"}, false, true},
		{"missing code", loginRequest{"+14155550100", ""}, true, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := authorizeRiskAction(tt.req, func(string, string) (bool, error) { return tt.verified, nil })
			if (err != nil) != tt.wantErr {
				t.Fatalf("error=%v wantErr=%v", err, tt.wantErr)
			}
			if !tt.wantErr && got != "high-risk action authorized" {
				t.Fatalf("unexpected result %q", got)
			}
		})
	}
}
