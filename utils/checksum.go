package utils

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
)

func ComputeMd5Checksum(data []byte) []byte {
	h := md5.New()
	h.Write(data)
	return h.Sum(nil)
}

func ComputeMd5ChecksumHex(data []byte) string {
	return hex.EncodeToString(ComputeMd5Checksum(data))
}

func ComputeSha256Checksum(data []byte) []byte {
	h := sha256.New()
	h.Write(data)
	return h.Sum(nil)
}

func ComputeSha256ChecksumHex(data []byte) string {
	return hex.EncodeToString(ComputeSha256Checksum(data))
}

func ComputeHmacSha256Sign(key, msg []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(msg)

	return mac.Sum(nil)
}

func ComputeHmacSha256SignHex(key, msg []byte) string {
	return hex.EncodeToString(ComputeHmacSha256Sign(key, msg))
}
