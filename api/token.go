package api

import (
	"context"
	"fmt"
	"io"
)

type TokenInterface interface {
	RefreshToken(ctx context.Context) error
}

func (az *AzureTTSClient) RefreshToken(ctx context.Context) error {
	req, err := az.newTokenRequest(ctx, "POST", az.TokenRefreshURL, nil)
	if err != nil {
		return err
	}

	resp, err := az.performReq(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	az.AccessToken = string(body)
	return nil
}
