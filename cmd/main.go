package main

import (
	"fmt"

	"github.com/DavydAbbasov/trecker_bot/config"
)

func main() {
	config.LoadConfig()
	fmt.Println("installed token",config.TelegramToken)
}
