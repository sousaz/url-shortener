package socket

import (
	"net"
	"testing"
)

type mockConn struct {
	net.Conn
	data []byte
}

// Sobrescreve a função Write da net
func (m *mockConn) Write(b []byte) (n int, err error) {
	m.data = append(m.data, b...)
	return len(b), nil
}

func TestRedirect(t *testing.T) {
	tests := []struct {
		location string
		expected string
	}{
		{
			location: "http://example.com",
			expected: "HTTP/1.1 302 Found\r\nLocation: http://example.com\rz\nContent-Length: 0\r\n\r\n",
		},
		{
			location: "https://example.com/path",
			expected: "HTTP/1.1 302 Found\r\nLocation: https://example.com/path\rz\nContent-Length: 0\r\n\r\n",
		},
	}

	for _, test := range tests {
		conn := &mockConn{}
		redirect(conn, test.location)

		if string(conn.data) != test.expected {
			t.Errorf("expected %q, got %q", test.expected, string(conn.data))
		}
	}
}

func TestRespond(t *testing.T) {
	tests := []struct {
		statusCode int
		message    string
		response   interface{}
		expected   string
	}{
		{
			statusCode: 200,
			message:    "Tudo certo!",
			response:   "http://localhost",
			expected:   "HTTP/1.1 200 OK\r\nContent-Type: application/json\r\ncontent-length: 69\r\n\r\n{\"message\":\"Tudo certo!\",\"response\":\"http://localhost\",\"status\":200}\n",
		},
		{
			statusCode: 200,
			message:    "Tudo certo!",
			response:   nil,
			expected:   "HTTP/1.1 200 OK\r\nContent-Type: application/json\r\ncontent-length: 55\r\n\r\n{\"message\":\"Tudo certo!\",\"response\":null,\"status\":200}\n",
		},
	}

	for _, test := range tests {
		conn := &mockConn{}
		respond(conn, test.statusCode, test.message, test.response)
		if string(conn.data) != test.expected {
			t.Errorf("expected %q, got %q", test.expected, string(conn.data))
		}
	}
}
