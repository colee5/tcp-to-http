package request

import (
	"errors"
	"io"
	"log"
	"module-lol/types"
	"regexp"
	"strings"
)

func NewRequest() *Request {
	return &Request{
		Headers: make(map[string]string),
		Body:    []byte{},
		State:   types.StateInitialized,
	}
}

const bufferSize = 8

func RequestFromReader(reader io.Reader) (*types.Request, error) {
	var readToIndex = 0

	buf := make([]byte, bufferSize)

	// Create a new request parser
	req := NewRequest()

	n, err := reader.Read(buf[readToIndex:])
	if err != nil {
		log.Fatal(err)
	}

	readToIndex += n
	for req.State != types.StateDone {
		if readToIndex >= len(buf) {
			// Grow buffer - create new slice
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}

		if err == io.EOF {
			break
		}

		bytesConsumed, err := req.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}

		if bytesConsumed > 0 {
			copy(buf, buf[bytesConsumed:readToIndex])
			readToIndex -= bytesConsumed
		} else {
			n, err := reader.Read(buf[readToIndex:])
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, err
			}
			readToIndex += n
		}
	}

	// Convert local Request back to types.Request and return
	return (*types.Request)(req), nil

}

func parseRequestLine(data string) (*types.RequestLine, int, error) {

	crlfIndex := strings.Index(data, "\r\n")
	if crlfIndex == -1 {
		return nil, 0, nil
	}

	requestLine := data[:crlfIndex]
	bytesConsumed := crlfIndex + 2

	parts := strings.Split(requestLine, " ")

	if len(parts) != 3 {
		return nil, 0, errors.New("invalid request line format")
	}

	method := parts[0]
	uri := parts[1]
	version := parts[2]

	// Validate method contains only uppercase letters
	validMethod := regexp.MustCompile(`^[A-Z]+$`)
	if !validMethod.MatchString(method) {
		return nil, 0, errors.New("invalid HTTP method")
	}

	// Validate HTTP version format and extract version number
	if version != "HTTP/1.1" {
		return nil, 0, errors.New("unsupported HTTP version")
	}

	// Extract just the version number from the version
	httpVersion := version[5:]

	return &types.RequestLine{
		Method:        method,
		RequestTarget: uri,
		HttpVersion:   httpVersion,
	}, bytesConsumed, nil
}

type Request types.Request

func (r *Request) parse(data []byte) (int, error) {
	if r.State == types.StateInitialized {
		requestLine, bytesParsed, err := parseRequestLine(string(data))
		if err != nil {
			return 0, err
		}

		if bytesParsed == 0 {
			return 0, nil
		}

		r.RequestLine = *requestLine
		r.State = types.StateDone
		return bytesParsed, nil
	}

	if r.State == types.StateDone {
		return 0, errors.New("error: trying to read data in a done state")
	}

	return 0, errors.New("error: unknown state")
}
