package main

import (
	"github.com/samsoft00/social-authentication-in-Golang/setup"
	"go.uber.org/fx"
)

func main() {
	options := setup.GetOptions()
	app := fx.New(options...)
	app.Run()
}
