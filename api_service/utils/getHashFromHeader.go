package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"io"
	"net/http"
)

func GetHashFromHeader(header http.Header) string{
	//logger.Println("++++++++++", header)
	digest := header.Get("Digest")
	if len(digest) == 0 || digest[:8] != "SHA-256=" {
		return ""
	}

	return digest[8:]
}

func CalculateHash(r io.Reader) string {
	h := sha256.New()
	io.Copy(h, r)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}