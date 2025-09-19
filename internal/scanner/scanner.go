package scanner

import (
	"bufio"
	"net"
	"strings"
	time "time"

	"github.com/user/pentscan/internal/probes"
	"github.com/user/pentscan/internal/types"
)

// Scanner é a interface que todos os tipos de scanner (Connect, SYN, etc.) devem implementar.
type Scanner interface {
	Scan(target types.Target, timeout time.Duration) types.ScanResult
	Close()
}

// ----------------------------------------
// Connect Scanner
// ----------------------------------------

var registeredProbes []probes.Probe

func init() {
	registeredProbes = append(registeredProbes, &probes.HTTPProbe{})
	registeredProbes = append(registeredProbes, &probes.SSHProbe{})
}

// ConnectScanner implementa a varredura usando net.Dial (TCP Connect).
type ConnectScanner struct{}

// NewConnectScanner cria um novo scanner de conexão.
func NewConnectScanner() (*ConnectScanner, error) {
	return &ConnectScanner{}, nil
}

// Scan realiza a varredura TCP Connect e o probing de serviço.
func (cs *ConnectScanner) Scan(target types.Target, timeout time.Duration) types.ScanResult {
	result := types.ScanResult{
		SchemaVersion: "1.1",
		Target:        target,
		Status:        types.StatusUnknown,
	}

	conn, err := net.DialTimeout("tcp", target.String(), timeout)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			result.Status = types.StatusFiltered
		} else if opErr, ok := err.(*net.OpError); ok {
			if strings.Contains(opErr.Err.Error(), "connection refused") {
				result.Status = types.StatusClosed
			} else {
				result.Status = types.StatusFiltered
				result.Error = err.Error()
			}
		} else {
			result.Status = types.StatusFiltered
			result.Error = err.Error()
		}
		return result
	}
	defer conn.Close()

	result.Status = types.StatusOpen

	for _, probe := range registeredProbes {
		// Re-estabelece a conexão para cada probe para garantir um estado limpo.
		// Isso é um pouco ineficiente, mas mais confiável.
		probeConn, err := net.DialTimeout("tcp", target.String(), timeout)
		if err != nil {
			continue
		}
		defer probeConn.Close()

		probeResult, _ := probe.Run(probeConn, timeout)
		if probeResult != nil {
			result.ServiceName = probeResult.ServiceName
			result.Banner = probeResult.Banner
			return result
		}
	}

	// Fallback para banner genérico se nenhum probe funcionar.
	fallbackConn, err := net.DialTimeout("tcp", target.String(), timeout)
	if err == nil {
		defer fallbackConn.Close()
		fallbackConn.SetReadDeadline(time.Now().Add(timeout))
		reader := bufio.NewReader(fallbackConn)
		banner, err := reader.ReadString('\n')
		if err == nil {
			result.Banner = strings.TrimSpace(banner)
		}
	}

	return result
}

// Close não faz nada para o ConnectScanner.
func (cs *ConnectScanner) Close() {}