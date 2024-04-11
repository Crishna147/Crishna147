package server

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"encoding/hex"
	"fmt"
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

type Command struct {
	Command      uint8
	UserID       []byte
	SectorNumber uint32
	Address      uint32
	DataLength   uint32
	Data         []byte
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

func EncryptControlData(data []byte, key []byte, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("could not create new cipher: %v", err)
	}
	encData := make([]byte, len(data))
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(encData, data)
	return encData, nil
}



func DecryptControlData(data []byte, key []byte, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("could not create new cipher: %v", err)
	}
	decData := make([]byte, len(data))
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(decData, data)
	return decData, nil
}

func (s *Server) HandleCommand(data []byte) ([]byte, uint8) {
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

	uid := hex.EncodeToString(userID)

	u, err := s.getUserByID(uid)
	if err != nil {
		s.log.Error("user does not exist", "uid", uid)
		return nil, CommadFailedErr
	}
	if !u.IsActive {
		s.log.Error("user is not active", "username", u.UserName)
		return nil, CommadFailedErr
	}

	switch cmd.Command {
	case ReadSectorCmd:
		rd, err := s.s.ReadData(uid, cmd.Address, cmd.DataLength)
		if err != nil {
			return nil, CommadFailedErr
		}
		cmd.Data = make([]byte, len(rd))
		copy(cmd.Data, rd)
	case WriteSectorCmd:
		if cmd.Address != 0 || s.initMode {
			err := s.s.WriteData(uid, cmd.Address, cmd.Data)
			if err != nil {
				return nil, CommadFailedErr
			}
		} else {
			s.log.Error("MF Setting write")
		}
	case EraseSectorCmd:
		if cmd.SectorNumber >= 4 {
			err := s.s.EraseSector(uid, cmd.SectorNumber)
			if err != nil {
				return nil, CommadFailedErr
			}
		}
	case EraseAllSectorsCmd:
		err := s.s.EraseFlash(uid)
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

func (s *Server) RemoveUser(userID string) error {
	return s.s.EraseFlash(userID)
}
