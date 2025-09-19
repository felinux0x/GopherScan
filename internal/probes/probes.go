package probes

import (
	"net"
	"time"
)

// ProbeResult contém os dados de um banner grab bem-sucedido.
type ProbeResult struct {
	// ServiceName é o nome do serviço identificado (ex: "http", "ssh").
	ServiceName string
	// Banner são os dados brutos ou processados retornados pelo alvo.
	Banner string
}

// Probe é a interface que todos os módulos de banner devem implementar.
type Probe interface {
	// Name retorna o nome do probe.
	Name() string
	// Run executa o probe em uma conexão já estabelecida.
	// Retorna um resultado se o serviço for identificado com sucesso.
	// Retorna (nil, nil) se o serviço não corresponder a este probe.
	// Retorna (nil, error) em caso de um erro inesperado de rede.
	Run(conn net.Conn, timeout time.Duration) (*ProbeResult, error)
}
