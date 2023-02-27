package util

import (
	"bytes"
	"crypto/sha1"
	"encoding/gob"
	"encoding/hex"
)

func HashFromStruct(v interface{}) (*string, error) {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	hash := sha1.New()
	hash.Write(buffer.Bytes())
	sha1_hash := hex.EncodeToString(hash.Sum(nil))
	return &sha1_hash, nil
}
