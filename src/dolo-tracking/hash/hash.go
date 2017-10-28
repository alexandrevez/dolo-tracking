package hash

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
)

// Sha256String returns a SHA-256 hash (as hex string) by concatenating `partList` strings
func Sha256String(partList ...string) string {
	merge := ""
	for _, part := range partList {
		merge += part
	}

	hash := sha256.Sum256([]byte(merge))
	return hex.EncodeToString(hash[:32])
}

// MD5String returns a MD5 hash (as hex string) by concatenating `partList` strings
func MD5String(partList ...string) string {
	merge := ""
	for _, part := range partList {
		merge += part
	}

	hash := md5.Sum([]byte(merge))
	return hex.EncodeToString(hash[:16])
}
