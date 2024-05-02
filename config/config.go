package config

import (
	"github.com/joho/godotenv"
	"os"
	"strings"
)

var (
	Bind       string
	Port       string
	ProxyUrl   string
	GatewayUrl string
)

func init() {
	_ = godotenv.Load()
	Bind = os.Getenv("BIND")
	if Bind == "" {
		Bind = "0.0.0.0"
	}

	Port = os.Getenv("PORT")
	if Port == "" {
		Port = "8484"
	}

	ProxyUrl = os.Getenv("PROXY")

	GatewayUrl = os.Getenv("GATEWAY_URL")
	if GatewayUrl == "" {
		GatewayUrl = "https://chat.openai.com"
	} else {
		GatewayUrl = strings.TrimRight(GatewayUrl, "/")
	}
}
