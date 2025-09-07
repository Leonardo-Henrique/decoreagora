package controllers

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/Leonardo-Henrique/decoreagora/app/core/config"
	"github.com/Leonardo-Henrique/decoreagora/app/core/config/logger"
	"github.com/Leonardo-Henrique/decoreagora/app/core/models"
	"github.com/Leonardo-Henrique/decoreagora/app/core/usecases"
	"github.com/Leonardo-Henrique/decoreagora/app/core/utils"
	"github.com/gofiber/fiber/v2"
)

type ImageController struct {
	imageUC     usecases.ImagesUsecase
	creditsUC   usecases.CreditsUsecase
	maxFileSize int64
}

func NewImageController(creditsUC usecases.CreditsUsecase, imageUC usecases.ImagesUsecase) *ImageController {
	return &ImageController{
		creditsUC:   creditsUC,
		imageUC:     imageUC,
		maxFileSize: 10 * 1024 * 1024,
	}
}

func (i *ImageController) CreateNewImage(c *fiber.Ctx) error {
	logger.Logging.Info("Starting CreateNewImage controller")
	ctx := context.TODO()

	form, err := c.MultipartForm()
	if err != nil {
		logger.Logging.Error("Error when getting MultipartForm", err)
		return c.Status(fiber.StatusBadRequest).JSON(models.ImageUploadResponse{
			Success: false,
			Message: "Failed to parse form data",
		})
	}

	files := form.File["image"]
	if len(files) == 0 {
		logger.Logging.Error("No file provided in image field", nil)
		return c.Status(fiber.StatusBadRequest).JSON(models.ImageUploadResponse{
			Success: false,
			Message: "No image file provided",
		})
	}

	imageFile := files[0]
	filename := imageFile.Filename
	extension := filepath.Ext(filename)

	if imageFile.Size > i.maxFileSize {
		logger.Logging.Error("Image size is greater than maxSize", nil)
		return c.Status(fiber.StatusBadRequest).JSON(models.ImageUploadResponse{
			Success: false,
			Message: fmt.Sprintf("File size exceeds maximum limit of %d bytes", i.maxFileSize),
		})
	}

	descriptions := form.Value["description"]
	if len(descriptions) == 0 {
		logger.Logging.Error("No image description prompt provided", nil)
		return c.Status(fiber.StatusBadRequest).JSON(models.ImageUploadResponse{
			Success: false,
			Message: "Description is required",
		})
	}

	description := descriptions[0]
	if strings.TrimSpace(description) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ImageUploadResponse{
			Success: false,
			Message: "Description cannot be empty",
		})
	}

	if len(strings.Fields(description)) < 3 {
		return c.Status(fiber.StatusBadRequest).JSON(models.ImageUploadResponse{
			Success: false,
			Message: "Small decription, try adding more details",
		})
	}

	prompt := fmt.Sprintf("Generate a high-quality image based on this description. Do not ask clarifying questions. Description: %s", description)

	if !i.isValidImageType(imageFile.Header.Get("Content-Type")) {
		logger.Logging.Error("User is trying to upload a not allowed image format", nil)
		return c.Status(fiber.StatusBadRequest).JSON(models.ImageUploadResponse{
			Success: false,
			Message: "Invalid image type. Only JPEG, PNG, and WebP are allowed",
		})
	}

	file, err := imageFile.Open()
	if err != nil {
		logger.Logging.Error("Error when opening image:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.ImageUploadResponse{
			Success: false,
			Message: "Failed to open uploaded file",
		})
	}
	defer file.Close()

	logger.Logging.Info("Saving original Image to S3")
	originalImageKey, err := i.imageUC.SaveImage(ctx, file, extension)
	if err != nil {
		logger.Logging.Error("Error when saving image in S3:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.ImageUploadResponse{
			Success: false,
			Message: "Failed to save img",
		})
	}
	logger.Logging.Info("Original image saved in S3")

	creation_date := time.Now()

	logger.Logging.Info("Registering new image generation entry in Database")
	originalImagePublicKey, err := i.imageUC.RegisterImage(originalImageKey, description, utils.GetCurrentUserID(c), creation_date)
	if err != nil {
		logger.Logging.Error("Error when creating new image generation registry:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.ImageUploadResponse{
			Success: false,
			Message: "Failed to register image",
		})
	}

	//editedFileBucketKey := "images/a5999"

	logger.Logging.Info("Sending image to AI edition")
	editedFile, err := i.imageUC.EditWithAI(file, prompt)
	if err != nil {
		logger.Logging.Error("Error editing with AI:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.ImageUploadResponse{
			Success: false,
			Message: "Failed to edit image with AI",
		})
	}

	logger.Logging.Info("Starting decrementing credit")
	if err := i.creditsUC.DecrementCredit(utils.GetCurrentUserID(c)); err != nil {
		logger.Logging.Error("Error when decrementing user credit:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.ImageUploadResponse{
			Success: false,
			Message: "Failed to update user funds",
		})
	}

	editedFileBucketKey, err := i.imageUC.SaveImage(ctx, bytes.NewReader(editedFile), extension)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ImageUploadResponse{
			Success: false,
			Message: "Failed to save IA image result",
		})
	}

	logger.Logging.Info("Registering the generated image key in database")
	err = i.imageUC.FinishImageEdition(editedFileBucketKey, originalImagePublicKey)
	if err != nil {
		logger.Logging.Error("Error when registering generatedimage key:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.ImageUploadResponse{
			Success: false,
			Message: "Failed to finish AI process",
		})
	}

	logger.Logging.Info("Signing original image")
	originalImageKey, err = i.imageUC.SignURL(
		ctx,
		config.C.AWS_IMAGES_BUCKET_NAME,
		originalImageKey,
	)
	if err != nil {
		logger.Logging.Error("Error when signing original image key:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.ImageUploadResponse{
			Success: false,
			Message: "Can sign object 1",
		})
	}

	logger.Logging.Info("Signing generated image")
	editedFileBucketKey, err = i.imageUC.SignURL(
		ctx,
		config.C.AWS_IMAGES_BUCKET_NAME,
		editedFileBucketKey,
	)
	if err != nil {
		logger.Logging.Error("Error when signing generated image key:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.ImageUploadResponse{
			Success: false,
			Message: "Can sign object 2",
		})
	}

	resp := models.EditedImageResponse{
		PublicKey:              originalImagePublicKey,
		OriginalImageBucketKey: originalImageKey,
		EditedImageBucketKey:   editedFileBucketKey,
		CreatedAt:              creation_date,
	}

	logger.Logging.Info("Signing EditedImageResponse to frontend")

	return c.Status(http.StatusCreated).JSON(resp)
}

func (i *ImageController) isValidImageType(contentType string) bool {
	validTypes := []string{
		"image/jpeg",
		"image/png",
		"image/webp",
		"image/jpg",
	}

	for _, validType := range validTypes {
		if contentType == validType {
			return true
		}
	}
	return false
}

func (i *ImageController) GetUserImages(c *fiber.Ctx) error {
	logger.Logging.Info("Starting GetUserImages controller")
	userID := utils.GetCurrentUserID(c)
	if userID == 0 {
		logger.Logging.Error("There is no userID in c:", nil)
		return c.Status(fiber.StatusBadRequest).JSON(models.ImageUploadResponse{
			Success: false,
			Message: fmt.Sprintf("Couldnt identify user"),
		})
	}

	logger.Logging.Info("Getting user generated images")
	images, err := i.imageUC.GetUserImages(userID)
	if err != nil {
		logger.Logging.Error("Error when getting user generated images", err)
		return c.Status(fiber.StatusBadRequest).JSON(models.ImageUploadResponse{
			Success: false,
			Message: fmt.Sprintf("error when retrieving user images"),
		})
	}

	logger.Logging.Info("Sending user images to frontend")
	return c.Status(http.StatusOK).JSON(images)
}

func (i *ImageController) GetUserCredits(c *fiber.Ctx) error {
	userID := utils.GetCurrentUserID(c)
	if userID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(models.ImageUploadResponse{
			Success: false,
			Message: fmt.Sprintf("Couldnt identify user"),
		})
	}
	credits, err := i.creditsUC.GetUserCredits(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ImageUploadResponse{
			Success: false,
			Message: fmt.Sprintf("error when retrieving user credits"),
		})
	}
	return c.Status(http.StatusOK).JSON(models.CreditsResponse{
		CurrentQtd: credits,
	})
}
