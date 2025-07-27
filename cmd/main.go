package main

import (
	"context"
	"github.com/folivorra/ziper/internal/config"
	"log/slog"
	"os"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})

	cfg := config.NewConfig()

}
