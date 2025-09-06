package ports

type EmailHandler interface {
	SendWithoutPreStyle(to, subject, body string) error
}
