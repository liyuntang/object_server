package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"io"
	"net/http"
)

func GetHashFromHeader(header http.Header) string{
	return header.Get("Digest")

}

func CalculateHash(r io.Reader) string {
	h := sha256.New()
	io.Copy(h, r)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}