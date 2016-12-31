package command

type FakeCommand struct {
	Next HandlerCommand
}
func (handler *FakeCommand) ProcessText(text string) string {
	return "Fake OK"
}
