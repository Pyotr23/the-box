package main

import (
	"log/slog"

	"github.com/Pyotr23/the-box/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		slog.Error(err.Error())
	}
}
