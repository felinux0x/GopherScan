package probes

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

// HTTPProbe é um probe para detectar servidores HTTP.
type HTTPProbe struct{}

func (p *HTTPProbe) Name() string {
	return "http"
}

func (p *HTTPProbe) Run(conn net.Conn, timeout time.Duration) (*ProbeResult, error) {
	// Envia uma requisição HEAD simples.
	request := "HEAD / HTTP/1.1\r\nHost: %s\r\nUser-Agent: PentScan\r\n\r\n"
	host := strings.Split(conn.RemoteAddr().String(), ":")[0]

	conn.SetWriteDeadline(time.Now().Add(timeout))
	_, err := fmt.Fprintf(conn, request, host)
	if err != nil {
		return nil, err
	}

	conn.SetReadDeadline(time.Now().Add(timeout))
	reader := bufio.NewReader(conn)
	responseLine, err := reader.ReadString('\n')
	if err != nil {
		// Não é necessariamente um erro, pode ser apenas um serviço que não respondeu.
		return nil, nil
	}

	// Verifica se a resposta parece ser HTTP.
	if strings.HasPrefix(responseLine, "HTTP/") {
		return &ProbeResult{
			ServiceName: p.Name(),
			Banner:      strings.TrimSpace(responseLine),
		},
		nil
	}

	// Não é um servidor HTTP.
	return nil, nil
}
