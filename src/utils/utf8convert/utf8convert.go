package utf8convert

import (
	"bytes"
	"io/ioutil"

	"golang.org/x/net/html/charset"
)

func ConvrtToUTF8(str string, origEncoding string) string {
	strBytes := []byte(str)
	byteReader := bytes.NewReader(strBytes)
	reader, _ := charset.NewReaderLabel(origEncoding, byteReader)
	strBytes, _ = ioutil.ReadAll(reader)
	return string(strBytes)
}
