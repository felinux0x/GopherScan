package writer

import (
	"fmt"
	"io"
	"strings"

	"github.com/user/pentscan/internal/types"
	"go.uber.org/zap"
)

// ANSI Color codes
const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorGray   = "\033[90m"
)

// TXTWriter formata os resultados em texto plano colorido.
type TXTWriter struct {
	output io.Writer
	logger *zap.Logger
}

// NewTXTWriter cria um novo writer de texto.
func NewTXTWriter(output io.Writer, logger *zap.Logger) (*TXTWriter, error) {
	return &TXTWriter{
		output: output,
		logger: logger,
	}, nil
}

// WriteHeader não faz nada para este formato de TXT.
func (w *TXTWriter) WriteHeader() {}

// Write formata e escreve um único resultado de varredura.
func (w *TXTWriter) Write(result types.ScanResult) {
	// Não exibir portas fechadas para manter a saída limpa, a menos que haja um erro.
	if result.Status == types.StatusClosed && result.Error == "" {
		return
	}

	var statusSymbol, statusColor string

	switch result.Status {
	case types.StatusOpen:
		statusSymbol = "[+]"
		statusColor = colorGreen
	case types.StatusClosed:
		statusSymbol = "[-]"
		statusColor = colorRed
	case types.StatusFiltered:
		statusSymbol = "[?]"
		statusColor = colorYellow
	default:
		statusSymbol = "[ ]"
		statusColor = colorGray
	}

	targetStr := fmt.Sprintf("%-21s", result.Target.String()) // Pad para 21 chars (e.g. 255.255.255.255:65535)
	statusStr := fmt.Sprintf("%-8s", result.Status.String())

	line := fmt.Sprintf("%s %s %s %s",
		statusColor+statusSymbol+colorReset,
		targetStr,
		statusColor+statusStr+colorReset,
		colorBlue+result.ServiceName+colorReset,
	)

	if result.Banner != "" {
		cleanBanner := strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(result.Banner, "\n", " "), "\r", ""))
		line += fmt.Sprintf(" %s%s%s", colorYellow, cleanBanner, colorReset)
	} else if result.Error != "" {
		line += fmt.Sprintf(" %s%s%s", colorGray, result.Error, colorReset)
	}

	fmt.Fprintln(w.output, line)
}

// Close não faz nada para o formato TXT.
func (w *TXTWriter) Close() {}
