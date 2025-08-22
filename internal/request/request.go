package request

import (
	"errors"
	"io"
	"log"
	"module-lol/types"
	"regexp"
	"strings"
)

func RequestFromReader(reader io.Reader) (*types.Request, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		log.Fatal(err)
	}

	stringifiedData := string(data)

	//
	requestLine, err := parseRequestLine(stringifiedData)
	if err != nil {
		return nil, err
	}

	request := &types.Request{
		RequestLine: *requestLine,
		Headers:     make(map[string]string),
		Body:        []byte{},
	}

	return request, nil

}

func parseRequestLine(data string) (*types.RequestLine, error) {
	lines := strings.Split(data, "\r\n")
	if len(lines) == 0 {
		return nil, errors.New("no request line found")
	}

	requestLine := lines[0]
	parts := strings.Split(requestLine, " ")

	if len(parts) != 3 {
		return nil, errors.New("invalid request line format")
	}

	method := parts[0]
	uri := parts[1]
	version := parts[2]

	// Validate method contains only uppercase letters
	validMethod := regexp.MustCompile(`^[A-Z]+$`)
	if !validMethod.MatchString(method) {
		return nil, errors.New("invalid HTTP method")
	}

	// Validate HTTP version format and extract version number
	if version != "HTTP/1.1" {
		return nil, errors.New("unsupported HTTP version")
	}

	// Extract just the version number from the version
	httpVersion := version[5:]

	return &types.RequestLine{
		Method:        method,
		RequestTarget: uri,
		HttpVersion:   httpVersion,
	}, nil
}
