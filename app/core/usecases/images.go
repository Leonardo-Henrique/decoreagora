package usecases

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"time"

	"github.com/Leonardo-Henrique/decoreagora/app/core/config"
	"github.com/Leonardo-Henrique/decoreagora/app/core/config/logger"
	"github.com/Leonardo-Henrique/decoreagora/app/core/models"
	"github.com/Leonardo-Henrique/decoreagora/app/core/ports"
	"github.com/google/uuid"
)

type ImagesUsecase struct {
	db ports.Database
	s3 ports.ImagesHandler
	ai ports.AiHandler
}

func NewImagesUsecase(db ports.Database, s3 ports.ImagesHandler, ai ports.AiHandler) *ImagesUsecase {
	return &ImagesUsecase{
		db: db,
		s3: s3,
		ai: ai,
	}
}

func (i *ImagesUsecase) SaveImage(ctx context.Context, image io.Reader, extension string) (string, error) {

	imageKey := fmt.Sprintf("images/%s"+extension,
		uuid.New().String()[:5],
	)

	if err := i.s3.SaveImage(
		ctx,
		image,
		imageKey,
		config.C.AWS_IMAGES_BUCKET_NAME,
	); err != nil {
		log.Println(err)
		return "", errors.New("we couldnt save the image")
	}

	return imageKey, nil
}

func (i *ImagesUsecase) RegisterImage(imageKey, prompt_description string, userID int, date time.Time) (string, error) {
	if imageKey == "" || prompt_description == "" || userID == 0 {
		return "", errors.New("required image metadata not found")
	}
	public_id := uuid.NewString()
	_, err := i.db.CreateNewImageEntry(public_id, imageKey, prompt_description, userID, date)
	if err != nil {
		return "", err
	}
	return public_id, nil
}

func (i *ImagesUsecase) EditWithAI(original_file multipart.File, prompt string) ([]byte, error) {
	imgData, err := io.ReadAll(original_file)
	if err != nil {
		return nil, err
	}
	newImg, err := i.ai.GenerateImage(imgData, "image/png", prompt)
	if err != nil {
		return nil, err
	}
	return newImg, nil
}

func (i *ImagesUsecase) FinishImageEdition(generatedFileBucketKey, originalFilePublicKey string) error {
	return i.db.FinishImageEdition(generatedFileBucketKey, originalFilePublicKey)
}

func (i *ImagesUsecase) SignURL(ctx context.Context, bucket string, key string) (string, error) {
	expire := time.Minute * 15
	return i.s3.SignURL(ctx, bucket, key, expire)
}

func (i *ImagesUsecase) GetUserImages(userID int) ([]models.EditedImageResponse, error) {
	images, err := i.db.GetUserImages(userID)
	if err != nil {
		return nil, err
	}

	expire := time.Minute * 40

	for j := 0; j < len(images); j++ {
		images[j].OriginalImageBucketKey, err = i.s3.SignURL(
			context.TODO(),
			config.C.AWS_IMAGES_BUCKET_NAME,
			images[j].OriginalImageBucketKey,
			expire,
		)
		if err != nil {
			logger.Logging.Error("Error when signing image at user generated images", err)
		}

		images[j].EditedImageBucketKey, err = i.s3.SignURL(
			context.TODO(),
			config.C.AWS_IMAGES_BUCKET_NAME,
			images[j].EditedImageBucketKey,
			expire,
		)
		if err != nil {
			logger.Logging.Error("Error when signing image at user generated images", err)
		}
	}

	return images, nil
}
