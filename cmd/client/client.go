package main

import (
	"github.com/skdiver33/gophkeeper/internal/client"
	"github.com/skdiver33/gophkeeper/logger"
)

func main() {
	logger.Log.Info("Hello", "123", "456")
	client.Run()
}
