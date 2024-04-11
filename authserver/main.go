package main

import (
	"flag"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/EnsurityTechnologies/apiconfig"
	"github.com/EnsurityTechnologies/logger"
	"github.com/EnsurityTechnologies/thincvirtual/authserver/config"
	"github.com/EnsurityTechnologies/thincvirtual/authserver/server"
)

const (
	APIConfigKey string = "TestKey@Prod#$erver^2021"
)

func main() {
	var debugMode bool
	var initMode bool
	flag.BoolVar(&debugMode, "d", false, "Enable debug")
	flag.BoolVar(&initMode, "i", false, "Init mode")
	flag.Parse()
	level := logger.Info
	if debugMode {
		level = logger.Debug
	}
	fp, err := os.OpenFile("log.txt",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	logOptions := &logger.LoggerOptions{
		Name:   "Main",
		Level:  level,
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

	cfg.InitMode = initMode

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
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM)
	signal.Notify(ch, syscall.SIGINT)
	select {
	case <-ch:
	}
	s.Shutdown()
	log.Info("Shutting down...")
}
