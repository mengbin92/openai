package server

import (
	"os"

	"github.com/mengbin92/openai/log"
)

var (
	logger = log.DefaultLogger().Sugar()
	token  = os.Getenv("TOKEN")
)
