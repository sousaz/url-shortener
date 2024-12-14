package socket

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/sousaz/urlshortener/database"
	"github.com/sousaz/urlshortener/utils"
)

var cache = make(map[string]*utils.Node)
var list = utils.List{}

func HandleConnection(conn net.Conn) {
	// É empilhada então só vai funcionar depois que esta função terminar
	defer conn.Close()

	request, err := parseRequest(conn)
	if err != nil {
		fmt.Fprintf(conn, "HTTP/1.1 400 Bad Request\r\n\r\nError processing request.")
		return
	}

	routeRequest(request, conn)
}

type Request struct {
	Method  string
	Path    string
	Headers map[string]string
	Body    map[string]interface{}
}

func parseRequest(conn net.Conn) (*Request, error) {
	// O bufio mantém um ponteiro onde parou a leitura anterior ou seja ele vai continuar apartir dali
	reader := bufio.NewReader(conn)
	var buffer bytes.Buffer

	// Pega (GET / HTTP/1.1) essa parte
	statusLine, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	buffer.WriteString(statusLine)

	// Mesma coisa que um split, vai devolver um slice de substrings separadas por um espaço
	parts := strings.Fields(statusLine)
	if len(parts) < 3 {
		return nil, fmt.Errorf("linha de status inválida: %s", statusLine)
	}
	method, path := parts[0], parts[1]

	headers := make(map[string]string)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		if line == "\r\n" {
			break
		}

		headerParts := strings.SplitN(line, ":", 2)
		if len(headerParts) >= 2 {
			headers[strings.TrimSpace(headerParts[0])] = strings.TrimSpace(headerParts[1])
		}
		buffer.WriteString(line)
	}

	var body map[string]interface{}
	if contentLength, ok := headers["content-length"]; ok {
		// Atoi é a mesma coisa que um ParseInt
		length, err := strconv.Atoi(contentLength)
		if err != nil {
			return nil, fmt.Errorf("content-length inválido: %s", contentLength)
		}

		bodyBuffer := make([]byte, length)
		_, err = reader.Read(bodyBuffer)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(bodyBuffer, &body)
		if err != nil {
			fmt.Printf("Erro ao decodificar JSON %q", err)
		}
	}

	return &Request{
		Method:  method,
		Path:    path,
		Headers: headers,
		Body:    body,
	}, nil
}

func routeRequest(req *Request, conn net.Conn) {
	switch {
	case req.Method == "GET":
		id := strings.Split(req.Path, "/")[1]
		cacheResponse, err := list.Get(id, cache)
		if err != nil {
			res, err := database.GetUrl(id)
			if err != nil {
				fmt.Printf("Error: %q", err)
			}
			list.Add(*res, cache, id)
			redirect(conn, *res)
			return
		}
		redirect(conn, *cacheResponse)
	case req.Method == "POST" && req.Path == "/submit":
		res, err := database.AddUrl(req.Body)
		if err != nil {
			respond(conn, 400, "Bad request", req.Body)
		}
		list.Add(req.Body["original"].(string), cache, *res)
		respond(conn, 200, "Requisição POST recebida", req.Headers["Host"]+"/"+*res)
	default:
		respond(conn, 404, "404 - Endpoint não encontrado", nil)
	}
}

func respond(conn net.Conn, statusCode int, message string, response interface{}) {
	responseData := map[string]interface{}{
		"status":   statusCode,
		"message":  message,
		"response": response,
	}

	var buffer bytes.Buffer
	err := json.NewEncoder(&buffer).Encode(responseData)
	if err != nil {
		fmt.Println("Error coding JSON")
		return
	}

	fmt.Fprintf(conn, "HTTP/1.1 %d OK\r\n", statusCode)
	fmt.Fprintf(conn, "Content-Type: application/json\r\n")
	fmt.Fprintf(conn, "content-length: %d\r\n", buffer.Len())
	fmt.Fprintf(conn, "\r\n")

	conn.Write(buffer.Bytes())
}

func redirect(conn net.Conn, location string) {
	response := fmt.Sprintf(
		"HTTP/1.1 302 Found\r\n"+
			"Location: %s\rz\n"+
			"Content-Length: 0\r\n"+
			"\r\n",
		location,
	)

	conn.Write([]byte(response))
}
