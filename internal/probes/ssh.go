package probes

import (
	"bufio"
	"net"
	"strings"
	time "time"
)

// SSHProbe é um probe para detectar servidores SSH.
type SSHProbe struct{}

func (p *SSHProbe) Name() string {
	return "ssh"
}

func (p *SSHProbe) Run(conn net.Conn, timeout time.Duration) (*ProbeResult, error) {
	conn.SetReadDeadline(time.Now().Add(timeout))
	reader := bufio.NewReader(conn)
	banner, err := reader.ReadString('\n')
	if err != nil {
		return nil, nil // O serviço não enviou um banner a tempo.
	}

	// Verifica se o banner se parece com um banner SSH.
	if strings.HasPrefix(banner, "SSH-") {
		return &ProbeResult{
			ServiceName: p.Name(),
			Banner:      strings.TrimSpace(banner),
		},
		nil
	}

	return nil, nil
}
