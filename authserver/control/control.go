package control

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"encoding/hex"
	"fmt"

	"github.com/EnsurityTechnologies/logger"
	"github.com/EnsurityTechnologies/thincvirtual/authserver/storage"
)

const (
	ReadSectorCmd uint8 = iota
	WriteSectorCmd
	EraseSectorCmd
	EraseAllSectorsCmd
)

const (
	CommandSuccess uint8 = iota
	InvalidCommandErr
	CommadFailedErr
)

type Control struct {
	log      logger.Logger
	s        storage.Storage
	initMode bool
	flag     bool
}

type Command struct {
	Command      uint8
	UserID       []byte
	SectorNumber uint32
	Address      uint32
	DataLength   uint32
	Data         []byte
}

func NewControl(log logger.Logger, s storage.Storage, im bool) (*Control, error) {
	c := &Control{
		log:      log,
		s:        s,
		initMode: im,
	}
	return c, nil
}

func MarshalCommand(c *Command) []byte {
	ul := len(c.UserID)
	dl := 14 + ul + int(c.DataLength)
	d := make([]byte, dl)
	d[0] = c.Command
	d[1] = byte(ul)
	copy(d[2:], []byte(c.UserID))
	binary.BigEndian.PutUint32(d[2+ul:], c.SectorNumber)
	binary.BigEndian.PutUint32(d[2+ul+4:], c.Address)
	binary.BigEndian.PutUint32(d[2+ul+8:], c.DataLength)
	if len(c.Data) > 0 {
		copy(d[2+ul+12:], c.Data)
	}
	return d
}

func UnMarshalCommand(d []byte) *Command {
	l := len(d)
	c := &Command{
		Command: d[0],
	}
	ul := d[1]
	dl := l - int(ul) - 14
	c.UserID = make([]byte, ul)
	copy(c.UserID, d[2:2+ul])
	c.SectorNumber = binary.BigEndian.Uint32(d[2+ul:])
	c.Address = binary.BigEndian.Uint32(d[2+ul+4:])
	c.DataLength = binary.BigEndian.Uint32(d[2+ul+8:])
	if dl > 0 {
		c.Data = make([]byte, dl)
		copy(c.Data, d[2+ul+12:])
	}
	return c
}

func EncryptData(data []byte, key []byte, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("could not create new cipher: %v", err)
	}
	encData := make([]byte, len(data))
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(encData, data)
	return encData, nil
}

func DecryptData(data []byte, key []byte, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("could not create new cipher: %v", err)
	}
	decData := make([]byte, len(data))
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(decData, data)
	return decData, nil
}

func (c *Control) HandleCommand(data []byte) ([]byte, uint8) {
	cmd := UnMarshalCommand(data)
	if cmd == nil {
		return nil, InvalidCommandErr
	}
	userID := cmd.UserID
	if cmd.SectorNumber < 4 {
		userID = make([]byte, 32)
		for i := range userID {
			userID[i] = 0xFF
		}
	}
	if userID[0] == 0 || userID[0] == 0xFF {
		c.flag = false
		//c.log.Debug("Accessing", "userID", hex.EncodeToString(userID))
	}

	if cmd.Address == 5120 || cmd.Address == 16384 || cmd.Address == 10240 {
		c.log.Debug("Accessing", "address", cmd.Address, "userID", hex.EncodeToString(userID))
	}

	switch cmd.Command {
	case ReadSectorCmd:
		if !c.flag && cmd.Address == 943104 {
			c.flag = true
			c.log.Debug("Accessing", "userID", hex.EncodeToString(userID))
		}
		rd, err := c.s.ReadData(hex.EncodeToString(userID), cmd.Address, cmd.DataLength)
		if err != nil {
			return nil, CommadFailedErr
		}
		cmd.Data = make([]byte, len(rd))
		copy(cmd.Data, rd)
	case WriteSectorCmd:
		if cmd.Address != 0 || c.initMode {
			err := c.s.WriteData(hex.EncodeToString(userID), cmd.Address, cmd.Data)
			if err != nil {
				return nil, CommadFailedErr
			}
		} else {
			c.log.Error("MF Setting write")
		}
	case EraseSectorCmd:
		if cmd.SectorNumber >= 4 {
			err := c.s.EraseSector(hex.EncodeToString(userID), cmd.SectorNumber)
			if err != nil {
				return nil, CommadFailedErr
			}
		}
	case EraseAllSectorsCmd:
		err := c.s.EraseFlash(hex.EncodeToString(cmd.UserID))
		if err != nil {
			return nil, CommadFailedErr
		}
	}
	wd := MarshalCommand(cmd)
	if wd == nil {
		return nil, CommadFailedErr
	}
	return wd, CommandSuccess
}

func (c *Control) RemoveUser(userID string) error {
	return c.s.EraseFlash(userID)
}
