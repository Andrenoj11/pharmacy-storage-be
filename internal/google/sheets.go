package google

import (
	"context"
	"fmt"
	"os"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func NewSheetsService(ctx context.Context, credentialsFile string) (*sheets.Service, error) {
	if credentialsFile == "" {
		return nil, fmt.Errorf("google credentials file is required")
	}

	credentialsJSON, err := os.ReadFile(credentialsFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read google credentials file: %w", err)
	}

	config, err := google.JWTConfigFromJSON(credentialsJSON, sheets.SpreadsheetsScope)
	if err != nil {
		return nil, fmt.Errorf("failed to parse google credentials file: %w", err)
	}

	client := config.Client(ctx)

	service, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("failed to create sheets service: %w", err)
	}

	return service, nil
}
