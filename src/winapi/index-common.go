package winapi

import (
	"bytes"
	"encoding/binary"
	"io"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func ByteToString(data []byte, coding string) (string, error) {
	buf := bytes.TrimRightFunc(data, func(r rune) bool {
		// 去年结尾的 \0 字符
		return r == 0
	})
	if len(buf) == 0 {
		return "", nil
	}

	if coding == "gbk" {
		data, err := io.ReadAll(transform.NewReader(
			bytes.NewBuffer(buf),
			simplifiedchinese.GBK.NewDecoder(),
		))
		if err != nil {
			return "", err
		}

		return string(data), nil
	} else {
		return string(buf), nil
	}
}

func ByteToUInt16 (data []byte) uint16 {
	buf := []byte{data[0], data[1]}
	return binary.LittleEndian.Uint16(buf)
}

func ByteToUInt32 (data []byte) uint32 {
	buf := []byte{data[0], data[1], data[2], data[3]}
	return binary.LittleEndian.Uint32(buf)
}

func ByteToUInt64 (data []byte) uint64 {
	buf := []byte{data[0], data[1], data[2], data[3], data[4], data[5], data[6], data[7]}
	return binary.LittleEndian.Uint64(buf)
}