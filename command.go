package main

type ZbotCommandHandler interface {
	nextCommand()
}

type ZbotCommand struct {
	pattern string
	next ZbotCommandHandler
}

