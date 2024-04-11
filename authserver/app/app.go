package app

import (
	"io"
	"os"

	"github.com/EnsurityTechnologies/apiconfig"
	"github.com/EnsurityTechnologies/logger"
	"github.com/EnsurityTechnologies/thincvirtual/authserver/config"
	"github.com/EnsurityTechnologies/thincvirtual/authserver/server"
)

const (
	APIConfigKey string = "TestKey@Prod#$erver^2021"
)

func Run() {
	fp, err := os.OpenFile("log.txt",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	logOptions := &logger.LoggerOptions{
		Name:   "Main",
		Level:  logger.Info,
		Color:  []logger.ColorOption{logger.AutoColor, logger.ColorOff},
		Output: []io.Writer{logger.DefaultOutput, fp},
	}

	log := logger.New(logOptions)

	var cfg config.Config

	err = apiconfig.LoadAPIConfig("api_config.json", APIConfigKey, &cfg)

	if err != nil {
		log.Error("Failed to parse api config file")
		return
	}

	cfg.InitMode = false

	// cfg := &server.Config{
	// 	Type:    "tcp",
	// 	Address: "localhost",
	// 	Port:    "12500",
	// 	Config: stream.Config{
	// 		SecureStream: true,
	// 	},
	// }
	s, err := server.NewServer(&cfg, log)
	if err != nil {
		log.Error("faield to setup new server", "err", err)
		return
	}
	go s.Listen()
}
