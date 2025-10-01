package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/TheAmirhosssein/cool-password-manage/config"
	"github.com/TheAmirhosssein/cool-password-manage/internal/infrastructure/server"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	conf := config.GetConfig()
	if err := server.Run(ctx, conf); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
