package main

import (
	"petProject/internal/app"
)

const configDirectoryPath = "config"

func main() {
	app.Run(configDirectoryPath)
}
