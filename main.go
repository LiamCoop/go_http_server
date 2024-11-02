package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
	"strconv"
)

type Header map[string][]string

type HTTPRequestFormat struct {
	Method   string
	Path     string
	Protocol string
	Headers  Header
	Body     io.Reader
}

type HTTPResponseFormat struct {
	Protocol string
	Status string
	Headers  Header
	Body     io.Reader
}

func main() {
	// implement a TCP listener
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Printf("%v", err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("%v", err)
		}
		go handleConnection(conn)
	}
}

// accepts connection, parses request
func handleConnection(conn net.Conn) (error) {
    defer conn.Close()

    reader := bufio.NewReader(conn)

	requestLine, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("error reading request line: %v", err)
	}
    requestLine = strings.TrimSpace(requestLine)
    parts := strings.Split(requestLine, " ")
    if len(parts) < 3 {
        return fmt.Errorf("invalid request line")
    }

    method, path, protocol := parts[0], parts[1], parts[2]

    // parse headers
    headers := make(map[string][]string)
    for {
        line, err := reader.ReadString('\n')
        if err != nil {
            return fmt.Errorf("error reading headers: %v", err)
        }
        line = strings.TrimSpace(line)
        if line == "" {
            // empty line signifies the end of headers
            break
        }

        headerParts := strings.Split(line, ": ")
        if len(headerParts) != 2 {
            return fmt.Errorf("invalid header line: %s", line)
        }
        key, value := headerParts[0], headerParts[1]
        headers[key] = append(headers[key], value)
    }

    // step 3: parse body
    var body io.Reader 
    if contentLengthValues, ok := headers["Content-Length"]; ok && len(contentLengthValues) > 0 {
        contentLength, err := strconv.Atoi(contentLengthValues[0])
        if err != nil {
            return fmt.Errorf("invalid content-length: %v", err)
        }

        body = io.LimitReader(reader, int64(contentLength))
    }

    request := &HTTPRequestFormat{
        Method: method,
        Path: path,
        Protocol: protocol,
        Headers: headers,
        Body: body,
    }

    responseStruct := &HTTPResponseFormat{
        Protocol: request.Protocol,
        Status: "200 OK",
        Headers: make(map[string][]string),
        Body: request.Body,
    }

    responseStr, err := createResponse(responseStruct)
    if err != nil {
		fmt.Printf("error building response", err)
    }

	_, err = conn.Write([]byte(responseStr))
	if err != nil {
		fmt.Printf("Error writing response: %v\n", err)
	}

    return nil
}

func createResponse(response *HTTPResponseFormat) (string, error) {
    var sb strings.Builder

    sb.WriteString(response.Protocol + " " + response.Status + "\r\n")

    // attach all headers, content length to be dealt with later
    // as part of dealing with the body itself.
    for key, values := range response.Headers {
        for _, value := range values {
            sb.WriteString(key + ": " + value + "\r\n")
        }
    }

    // determine length, attach that header, also attach the body itself
    contentLength := 0
    if response.Body != nil {
        bodyBytes, err := io.ReadAll(response.Body)
        if err != nil {
            return "", fmt.Errorf("failed to read response body: %v", err)
        }
        contentLength = len(bodyBytes)
        // Write body content after headers
        sb.WriteString("Content-Length: " + strconv.Itoa(contentLength) + "\r\n")
        sb.WriteString("\r\n")
        sb.Write(bodyBytes) 
    } else {
        // no body, so double newline after headers
        sb.WriteString("\r\n")
    }

    st := sb.String()
    return st, nil
}

const postRequest = `POST /api/v1/users HTTP/1.1
Host: example.com
Content-Type: application/json
Content-Length: 55

{"name": "John Doe", "email": "john@example.com"}
`

const request = `POST /api/v1/users HTTP/1.1
Host: example.com
Content-Type: application/json
Content-Length: 34

{"name": "John", "email": "john@example.com"}`

