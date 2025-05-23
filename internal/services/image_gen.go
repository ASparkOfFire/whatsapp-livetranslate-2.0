package services

import "context"

type ImageGenerator interface {
	GenerateImage(ctx context.Context, prompt string) ([]byte, error)
}
