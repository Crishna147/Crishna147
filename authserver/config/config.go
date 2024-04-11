package config

import "github.com/EnsurityTechnologies/config"

type Config struct {
	ServerConfig config.Config `json:"server_config"`
	License      string        `json:"license"`
	Type         string        `json:"type"`
	Address      string        `json:"address"`
	Port         string        `json:"port"`
	InitMode     bool          `json:"init_mode"`
	SecureStream bool          `json:"secure_stream"`
	SectorSize   uint32        `json:"sector_size"`
	Client       bool          `json:"client"`
}
