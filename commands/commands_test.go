package command

var user = User{
	Username: "ssalvatori",
	Ident: "stefano",
	Host: "localhost",
}

type FakeCommand struct {
	Next HandlerCommand
}
func (handler *FakeCommand) ProcessText(text string, user User) string {
	return "Fake OK"
}
