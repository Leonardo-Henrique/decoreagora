package models

import (
	"mime/multipart"
	"time"
)

type ImageUploadRequest struct {
	Image       *multipart.FileHeader `form:"image"`
	Description string                `form:"description"`
}

type ImageUploadResponse struct {
	Success  bool   `json:"success"`
	Message  string `json:"message"`
	ImageURL string `json:"image_url,omitempty"`
	Filename string `json:"filename,omitempty"`
}

type EditedImageResponse struct {
	PublicKey              string    `json:"public_key"`
	OriginalImageBucketKey string    `json:"original_image_key"`
	EditedImageBucketKey   string    `json:"edited_image_key"`
	Prompt                 string    `json:"prompt"`
	CreatedAt              time.Time `json:"created_at"`
}

type CreditsResponse struct {
	CurrentQtd int `json:"current_qtd"`
}
