package core

import "crypto/md5"
import "encoding/hex"

type ID string

// In bytes
const SafeTrimSize = 4

func GenID(s string) ID {
	return ID(GenString(s))
}

func GenString(s string) string {
	hasher := md5.New()
	hasher.Write([]byte(s))
	return hex.EncodeToString(hasher.Sum([]byte{}))[0 : SafeTrimSize*2]
}
