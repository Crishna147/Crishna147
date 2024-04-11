package server

import (
	"crypto/ecdh"
	"crypto/rand"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/EnsurityTechnologies/enscrypt"
	"github.com/EnsurityTechnologies/ensweb"
	"github.com/EnsurityTechnologies/logger"
	"github.com/EnsurityTechnologies/thincvirtual/authserver/config"
	"github.com/EnsurityTechnologies/thincvirtual/authserver/docs"
	"github.com/EnsurityTechnologies/thincvirtual/authserver/storage"
	"github.com/EnsurityTechnologies/thincvirtual/authserver/stream"
	"github.com/patrickmn/go-cache"
)

const (
	HeaderLength int    = 4
	FrameHeader  uint16 = 0x5354
)

const (
	RequestID    string = "requestid"
	PublicKeyID  string = "publickey"
	LicenseKeyID string = "licensekey"
)
const (
	APIGetPublicKey   string = "/api/getpublickey"
	APIAdminLogin     string = "/api/adminlogin"
	APIAddDevice      string = "/api/adddevice"
	APIBlockDevice    string = "/api/blockdevice"
	APIUnBlockDevice  string = "/api/unblockdevice"
	APIResetUser      string = "/api/resetuser"
	APISetNumFingers  string = "/api/setnumfingers"
	APISetPinRequired string = "/api/setpinrequired"
	APIUserLogin      string = "/api/userlogin"
	APIDeviceStatus   string = "/api/devicestatus"
	APIFingerEnrolled string = "/api/fingerenrolled"
	APIKeyAdded       string = "/api/keyadded"
)

const (
	DefaultExpiration = 5 * time.Minute
	PurgeTime         = 10 * time.Minute
)

const (
	MaxAllowedFingers int = 5
)

type Server struct {
	ensweb.Server
	cfg      *config.Config
	log      logger.Logger
	l        net.Listener
	shutDown bool
	ll       sync.Mutex
	initMode bool
	s        storage.Storage
	ts       storage.TableStorage
	users    *cache.Cache
}

type HandlerFunc func(req *ensweb.Request, dr *DataRequest) *ensweb.Result

func NewServer(cfg *config.Config, log logger.Logger) (*Server, error) {
	var st storage.Storage
	var err error
	st, err = storage.NewFileStorage("./profiles/", 1024, 100, cfg.InitMode)
	if err != nil {
		log.Error("failed to create file storage")
		return nil, err
	}

	key, err := ecdh.P256().GenerateKey(rand.Reader)
	if err != nil {
		log.Error("failed to generate private key")
		return nil, err
	}
	s := &Server{
		cfg:      cfg,
		log:      log.Named("Server"),
		s:        st,
		initMode: cfg.InitMode,
		users:    cache.New(DefaultExpiration, PurgeTime),
	}
	s.Server, err = ensweb.NewServer(&cfg.ServerConfig, nil, log.Named("apiserver"), ensweb.EnableSecureAPI(key, s.cfg.License))
	if err != nil {
		s.log.Error("failed to create api server", "err", err)
		return nil, err
	}
	secret := s.checkLicense()
	if secret == "" {
		s.log.Error("invalid license shutting down server")
		return nil, fmt.Errorf("invalid license")
	}
	ts, err := storage.NewTableFileStorage("./tables/", secret, log.Named("table"))
	if err != nil {
		log.Error("failed to create table file storage")
		return nil, err
	}
	s.ts = ts
	err = s.initLicense()
	if err != nil {
		s.log.Error("failed to initialize license", "err", err)
		return nil, err
	}
	u, err := s.getUser("authadmin")
	if err != nil || u.UserName != "authadmin" {
		u = &User{
			UserName:       "authadmin",
			UserID:         "0000-0000-0000-0000",
			Password:       enscrypt.HashPassword("admin@123$", 3, 1, 1000),
			Role:           AdminRole,
			IsActive:       true,
			AllowedFingers: 5,
		}
		s.log.Info("creating authadmin user")
		err = s.putUser(u)
		if err != nil {
			log.Error("failed to create authadmin", "err", err)
			return nil, err
		}
	}
	url := s.GetServerURL()
	if strings.HasPrefix(url, "http://") {
		url = strings.TrimLeft(url, "http://")
	} else if strings.HasPrefix(url, "https://") {
		url = strings.TrimLeft(url, "https://")
	}
	docs.SwaggerInfo.Title = "THINC VIRTUAL"
	docs.SwaggerInfo.Description = "This is THINC VIRTUAL  server framework"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = url
	docs.SwaggerInfo.BasePath = ""
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	s.EnableSWagger("http://localhost:12501")
	s.SetDebugMode()
	s.SetShutdown(s.closeAgentService)
	s.RegisterRoutes()
	
	go s.Start()
	if err != nil {
		log.Error("failed to create authadmin", "err", err)
		return nil, err
	}
	//s.Regi
	return s, nil
}

func (s *Server) RegisterRoutes() {
	s.AddRoute(APIAdminLogin, "POST", s.adminLogin)
	s.AddRoute(APIAddDevice, "POST", s.BasicAuthHandle(&BearerToken{}, s.addDevice, nil, nil))
	s.AddRoute(APIBlockDevice, "POST", s.BasicAuthHandle(&BearerToken{}, s.blockDevice, nil, nil))
	s.AddRoute(APIUnBlockDevice, "POST", s.BasicAuthHandle(&BearerToken{}, s.unblockDevice, nil, nil))
	s.AddRoute(APIResetUser, "POST", s.BasicAuthHandle(&BearerToken{}, s.resetUser, nil, nil))
	s.AddRoute(APISetNumFingers, "POST", s.BasicAuthHandle(&BearerToken{}, s.setNumFingers, nil, nil))
	s.AddRoute(APISetPinRequired, "POST", s.BasicAuthHandle(&BearerToken{}, s.setPinRequired, nil, nil))
	s.AddRoute(APIUserLogin, "POST", s.userLogin)
	//s.AddRoute(APIDeviceStatus, "POST", s.BasicAuthHandle(&BearerToken{}, s.deviceStatus, nil, nil))
	s.AddRoute(APIDeviceStatus, "POST", s.deviceStatus)
	s.AddRoute(APIFingerEnrolled, "POST", s.BasicAuthHandle(&BearerToken{}, s.fingerEnrolled, nil, nil))
	s.AddRoute(APIKeyAdded, "POST", s.BasicAuthHandle(&BearerToken{}, s.keyAdded, nil, nil))

}

func (s *Server) getPrivilege(role string) int {
	switch role {
	case AdminRole:
		return 1
	case UserRole:
		return 0
	default:
		return -1
	}
}

func (s *Server) authorize(role string) ensweb.AuthFunc {
	return ensweb.AuthFunc(func(req *ensweb.Request) bool {
		bt, ok := req.ClientToken.Model.(BearerToken)
		if !ok {
			s.log.Error("invalid bearer token, failed to autorize")
			return false
		}
		if s.getPrivilege(bt.Role) >= s.getPrivilege(role) {
			s.log.Error("permission denied", "exp", role, "recv", bt.Role)
			return false
		}
		return true
	})
}

func (s *Server) getHeader(req *ensweb.Request, key string, ss string) string {
	v := req.Headers.Get(key)
	if v == "" {
		s.log.Error("header is not present", "key", key)
		return ""
	}
	ds, err := DecryptData(ss, v)
	if err != nil {
		s.log.Error("failed to get header", "err", err)
		return ""
	}
	return ds
}

func (s *Server) Listen() {
	var err error
	s.l, err = net.Listen(s.cfg.Type, s.cfg.Address+":"+s.cfg.Port)
	if err != nil {
		s.log.Error("failed to listen", "err", err)
		return
	}
	s.log.Info("Server started", "address", s.l.Addr().String())
	for {
		s.log.Debug("Waiting for the stream")
		c, err := s.l.Accept()
		if err != nil {
			if !s.shutDown {
				s.log.Error("Failed to acccept connection", "err", err)
				return
			}
		}
		if c == nil {
			return
		}
		go s.handleConnection(c)
	}
}

func (s *Server) handleConnection(c net.Conn) {
	defer c.Close()
	st, err := stream.NewStream(c, s.cfg, s.log, false)
	if err != nil {
		s.log.Error("failed to open stream", "err", err)
		return
	}
	s.log.Debug("Stream connected")
	for {
		rf, err := st.RecvFrame()
		if err != nil {
			s.log.Error("failed to recv frame", "err", err)
			return
		}
		hd, status := s.HandleCommand(rf.Data)
		wr := &stream.Frame{
			Command: rf.Command,
			Status:  status,
			Length:  uint16(len(hd)),
		}
		if len(hd) != 0 {
			wr.Data = hd
		}
		err = st.SendFrame(wr)
		if err != nil {
			s.log.Error("failed to send frame", "err", err)
			return
		}
	}
}

func (s *Server) closeAgentService() error {
	s.shutDown = true
	return s.l.Close()
}
