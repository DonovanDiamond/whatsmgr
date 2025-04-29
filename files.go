package whatsmgr

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"os"
)

func (conn *Connection) hashFile(data []byte) string {
	hash := sha1.New()
	hash.Write(data)
	return fmt.Sprintf("%d-%s", len(data), hex.EncodeToString(hash.Sum(nil)))
}

func (conn *Connection) writeFileIfNotExists(path string, data []byte) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}
		return err
	}
	defer f.Close()
	_, err = f.Write(data)
	return err
}
