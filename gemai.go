package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Leonardo-Henrique/decoreagora/app/core/config"
	"google.golang.org/genai"
)

func CreateImage() {

	ctx := context.Background()
	genclient, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: config.C.GOOGLE_GEMINI_API_KEY,
	})
	if err != nil {
		log.Fatal(err)
	}

	imagePath := "./test.jpeg"
	imgData, err := os.ReadFile(imagePath)
	if err != nil {
		log.Fatal(err)
	}

	parts := []*genai.Part{
		genai.NewPartFromText("Put a couple sit down in the sofa, they are watching tv"),
		&genai.Part{
			InlineData: &genai.Blob{
				MIMEType: "image/png",
				Data:     imgData,
			},
		},
	}

	contents := []*genai.Content{
		genai.NewContentFromParts(parts, genai.RoleUser),
	}

	result, err := genclient.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash-image-preview",
		contents,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	for _, part := range result.Candidates[0].Content.Parts {
		if part.Text != "" {
			fmt.Println(part.Text)
		} else if part.InlineData != nil {
			imageBytes := part.InlineData.Data
			outputFilename := "generated.png"
			_ = os.WriteFile(outputFilename, imageBytes, 0644)
		}

	}

}
