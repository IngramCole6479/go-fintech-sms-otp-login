package main

import "fmt"

type loginRequest struct{ Phone, Code string }

func authorizeRiskAction(req loginRequest, verify func(string, string) (bool, error)) (string, error) {
	if req.Phone == "" || req.Code == "" {
		return "", fmt.Errorf("phone and code are required")
	}
	verified, err := verify(req.Phone, req.Code)
	if err != nil {
		return "", err
	}
	if !verified {
		return "", fmt.Errorf("phone verification failed")
	}
	return "high-risk action authorized", nil
}
