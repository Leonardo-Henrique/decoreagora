package ports

type AiHandler interface {
	GenerateImage(imageData []byte, imageType, instructions string) ([]byte, error)
}
