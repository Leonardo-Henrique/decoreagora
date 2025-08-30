package ai

import (
	"context"
	"fmt"
	"log"

	"github.com/Leonardo-Henrique/decoreagora/app/core/config"
	"google.golang.org/genai"
)

type Gemini struct {
	genclient *genai.Client
	ctx       context.Context
}

func NewGemini(ctx context.Context) *Gemini {
	client, err := genai.NewClient(
		ctx,
		&genai.ClientConfig{
			APIKey: config.C.GOOGLE_GEMINI_API_KEY,
		},
	)
	if err != nil {
		log.Fatal("could not start gemini client")
	}
	return &Gemini{
		genclient: client,
		ctx:       ctx,
	}
}

func (g *Gemini) GenerateImage(imageData []byte, imageType, instructions string) ([]byte, error) {
	parts := []*genai.Part{
		genai.NewPartFromText(instructions),
		&genai.Part{
			InlineData: &genai.Blob{
				MIMEType: imageType,
				Data:     imageData,
			},
		},
	}

	contents := []*genai.Content{
		genai.NewContentFromParts(parts, genai.RoleUser),
	}

	result, err := g.genclient.Models.GenerateContent(
		g.ctx,
		"gemini-2.5-flash-image-preview",
		contents,
		nil,
	)
	if err != nil {
		return nil, err
	}

	for _, part := range result.Candidates[0].Content.Parts {
		if part.Text != "" {
			fmt.Println(part.Text)
		} else if part.InlineData != nil {
			imageBytes := part.InlineData.Data
			//outputFilename := "generated.png"
			//_ = os.WriteFile(outputFilename, imageBytes, 0644)
			return imageBytes, nil
		}

	}

	return nil, nil

}
