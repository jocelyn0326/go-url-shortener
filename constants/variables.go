package constants

import (
	"context"
	"os"
)

var (
	WebServerPort = os.Getenv("WebServerPort")
	BaseUrl       = os.Getenv("BaseUrl")
	RedisAddr     = os.Getenv("RedisHost") + ":" + os.Getenv("RedisPort")
	Ctx           = context.TODO()
)
