package client

import (
	"github.com/EnsurityTechnologies/ensweb"
	"github.com/EnsurityTechnologies/logger"
	"github.com/EnsurityTechnologies/thincvirtual/authserver/config"
	"github.com/EnsurityTechnologies/uuid"
)

type Client struct {
	ensweb.Client
	secret string
	cfg    *config.Config
	ss     string
	log    logger.Logger
	uuid   string
}

func NewClient(cfg *config.Config, log logger.Logger) (*Client, error) {
	var err error
	c := &Client{
		cfg: cfg,
		log: log,
	}
	c.Client, err = ensweb.NewClient(&cfg.ServerConfig, log.Named("ensclient"), ensweb.SetClientTokenHelper("token.txt"), ensweb.EnableClientSecureAPI(cfg.License))
	if err != nil {
		log.Error("failed to create ensclient")
		return nil, err
	}
	// c.spub, err = c.GetPublicKey()
	// if err != nil {
	// 	log.Error("failed to get public key", "err", err)
	// 	return nil, err
	// }
	// log.Info("gnerating shared secret")
	// c.ss, err = c.getSharedSecret(c.spub)
	// if err != nil {
	// 	log.Error("failed to get shared secret key", "err", err)
	// 	return nil, err
	// }
	id := uuid.New()
	c.uuid = id.String()
	return c, nil
}
