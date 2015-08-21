package main

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"strconv"
)

const (
	RequestPayload          = '1'
	ResponsePayload         = '2'
	ReplayedResponsePayload = '3'
)

func uuid() []byte {
	b := make([]byte, 20)
	rand.Read(b)

	uuid := make([]byte, 40)
	hex.Encode(uuid, b)

	return uuid
}

var payloadSeparator = "\n🐵🙈🙉\n"

func payloadScanner(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	if i := bytes.Index(data, []byte(payloadSeparator)); i >= 0 {
		// We have a full newline-terminated line.
		return i + len([]byte(payloadSeparator)), data[0:i], nil
	}

	if atEOF {
		return len(data), data, nil
	}
	return 0, nil, nil
}

// Timing is request start or round-trip time, depending on payloadType
func payloadHeader(payloadType byte, uuid []byte, timing int64) (header []byte) {
	sTime := strconv.FormatInt(timing, 10)

	//Example:
	//  3 f45590522cd1838b4a0d5c5aab80b77929dea3b3 1231\n
	// `+ 1` indicates space characters or end of line
	header = make([]byte, 1+1+len(uuid)+1+len(sTime)+1)
	header[0] = payloadType
	header[1] = ' '
	header[2+len(uuid)] = ' '
	header[len(header)-1] = '\n'

	copy(header[2:], uuid)
	copy(header[3+len(uuid):], sTime)

	return header
}

func payloadBody(payload []byte) []byte {
	headerSize := bytes.IndexByte(payload, '\n')
	return payload[headerSize+1:]
}

func payloadMeta(payload []byte) [][]byte {
	headerSize := bytes.IndexByte(payload, '\n')
	return bytes.Split(payload[:headerSize], []byte{' '})
}

func isOriginPayload(payload []byte) bool {
	switch payload[0] {
	case RequestPayload, ResponsePayload:
		return true
	default:
		return false
	}
}

func isRequestPayload(payload []byte) bool {
	return payload[0] == RequestPayload
}
