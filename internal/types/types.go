package types

import "fmt"

// Target representa um único alvo de varredura (IP/Host + Porta).
type Target struct {
	Host string
	Port int
}

func (t Target) String() string {
	return fmt.Sprintf("%s:%d", t.Host, t.Port)
}

// ScanStatus representa o estado de uma porta.
type ScanStatus int

const (
	StatusUnknown ScanStatus = iota
	StatusOpen
	StatusClosed
	StatusFiltered // ou Timeout
)

func (s ScanStatus) String() string {
	switch s {
	case StatusOpen:
		return "open"
	case StatusClosed:
		return "closed"
	case StatusFiltered:
		return "filtered"
	default:
		return "unknown"
	}
}

// ScanResult contém o resultado da varredura para um único alvo.
type ScanResult struct {
	SchemaVersion string     `json:"schema_version"`
	Target        Target     `json:"target"`
	Status        ScanStatus `json:"status"`
	ServiceName   string     `json:"service_name,omitempty"`
	Banner        string     `json:"banner,omitempty"`
	Error         string     `json:"error,omitempty"`
}