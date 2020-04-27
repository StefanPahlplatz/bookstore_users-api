package app

import (
	"github.com/StefanPahlplatz/bookstore_utils-go/logger"
	"github.com/gin-gonic/gin"
)

var (
	router = gin.Default()
)

func StartApplication() {
	mapUrls()

	logger.Info("Starting application...")
	router.Run(":8081")
}
