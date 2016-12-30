package command

type HandlerCommand interface {
	ProcessText(text string) string
}