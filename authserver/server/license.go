package server

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"io/ioutil"
	"time"

	"github.com/EnsurityTechnologies/enscrypt"
	"github.com/EnsurityTechnologies/thincvirtual/authserver/storage"
	"github.com/jaypipes/ghw"
)

const (
	LicenseKey  string = "license@3563Aw2930*373622nKlaseX230497/>aqwofkgh"
	LicenseText string = "license20294857KSjfughShkdksWl;sldlf9347823.?ald;ALskd12@3Dla0a0"
)

func GetSystemHash() string {
	block, err := ghw.Block()
	if err != nil {
		return ""
	}
	h := sha256.New()
	baseboard, err := ghw.Baseboard()
	if err != nil {
		return ""
	}
	h.Write([]byte(baseboard.Vendor))
	h.Write([]byte(baseboard.SerialNumber))
	for _, disk := range block.Disks {
		h.Write([]byte(disk.SerialNumber))
	}
	product, err := ghw.Product()
	if err != nil {
		return ""
	}
	h.Write([]byte(product.Name))
	h.Write([]byte(product.SerialNumber))
	h.Write([]byte(product.Vendor))
	h.Write([]byte(product.UUID))

	net, err := ghw.Network()
	if err != nil {
		return ""
	}

	for _, nic := range net.NICs {
		if !nic.IsVirtual {
			h.Write([]byte(nic.MacAddress))
		}
	}

	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}

func GenerateLicense(hs string, key string) string {
	// h := sha256.New()
	// h.Write([]byte(hs))
	// h.Write([]byte(LicenseText))
	// hb := h.Sum(nil)
	// text := hex.EncodeToString(hb)
	eb, err := enscrypt.Seal(generateKey(hs), []byte(key))
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(eb)
}

func generateKey(hs string) string {
	h := sha256.New()
	h.Write([]byte(hs))
	h.Write([]byte(LicenseKey))
	hb := h.Sum(nil)
	key := hex.EncodeToString(hb)
	return key
}

func (s *Server) checkLicense() string {
	hs := GetSystemHash()
	s.log.Info("System hash : " + hs)
	rb, err := ioutil.ReadFile("license.txt")
	if err != nil {
		key := RandString(32)
		lic := GenerateLicense(hs, key)
		err := ioutil.WriteFile("license.txt", []byte(lic), 0660)
		if err != nil {
			s.log.Error("failed to create license file", "err", err)
			return ""
		}
		return key
	}
	db, err := base64.StdEncoding.DecodeString(string(rb))
	if err != nil {
		s.log.Error("invalid license file", "err", err)
		return ""
	}
	d, err := enscrypt.UnSeal(generateKey(hs), db)
	if err != nil {
		s.log.Error("invalid license file", "err", err)
		return ""
	}
	return string(d)
}

func (s *Server) initLicense() error {
	var l License
	err := s.ts.Read(LicenseTable, "license", &l)
	if err == storage.ErrDec {
		s.log.Error("invalid license key, failed to initialize license", "err", err)
		return err
	}
	if err != nil {
		l = License{
			MaxNumberUsers:   20,
			MaxNumberDevices: 20,
			NumberUsers:      0,
			NumberDevices:    0,
			ExpiryTime:       time.Now().AddDate(0, 2, 0),
		}
		err = s.ts.Write(LicenseTable, "license", l)
		if err != nil {
			s.log.Error("failed to init license", "err", err)
			return err
		}
	}
	return nil
}

func (s *Server) getLicense() (*License, error) {
	var l License
	s.ll.Lock()
	err := s.ts.Read(LicenseTable, "license", &l)
	if err != nil {
		s.log.Error("failed to get license", "err", err)
		s.ll.Unlock()
		return nil, err
	}
	return &l, nil
}

func (s *Server) releaseLicense() {
	s.ll.Unlock()
}

func (s *Server) updateLicense(l *License) error {
	err := s.ts.Write(LicenseTable, "license", *l)
	if err != nil {
		s.log.Error("failed to update license", "err", err)
		return err
	}
	return nil
}
