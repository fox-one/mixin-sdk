package utils

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
)

// GZipDecode gzip decode
func GZipDecode(binary []byte) ([]byte, error) {
	byteReader := bytes.NewReader(binary)
	gReader, err := gzip.NewReader(byteReader)

	if err != nil {
		return []byte{}, err
	}

	return ioutil.ReadAll(gReader)
}

// GZipEncode gzip encode
func GZipEncode(p []byte) ([]byte, error) {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	if _, err := gz.Write(p); err != nil {
		return nil, err
	}
	gz.Flush()
	gz.Close()

	return buf.Bytes(), nil
}
