package main

import (
	"github.com/TheAmirhosssein/cool-password-manage/config"
	"github.com/TheAmirhosssein/cool-password-manage/internal/infrastructrue/server"
)

func main() {
	conf := config.GetConfig()
	server.Run(conf)
}
