package writer

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"strconv"

	"github.com/user/pentscan/internal/types"
	"go.uber.org/zap"
)

// ResultWriter é a interface para escrever resultados de varredura.
type ResultWriter interface {
	WriteHeader()
	Write(result types.ScanResult)
	Close()
}

// ----------------------------------------
// JSON Writer
// ----------------------------------------

type JSONWriter struct {
	encoder *json.Encoder
	logger  *zap.Logger
	file    io.Closer
}

func NewJSONWriter(w io.Writer, logger *zap.Logger) (*JSONWriter, error) {
	jw := &JSONWriter{
		encoder: json.NewEncoder(w),
		logger:  logger,
	}
	if f, ok := w.(io.Closer); ok {
		jw.file = f
	}
	return jw, nil
}

func (jw *JSONWriter) WriteHeader() { /* JSON não precisa de header */ }

func (jw *JSONWriter) Write(result types.ScanResult) {
	if err := jw.encoder.Encode(result); err != nil {
		jw.logger.Error("Failed to write JSON output", zap.Error(err))
	}
}

func (jw *JSONWriter) Close() {
	if jw.file != nil {
		if err := jw.file.Close(); err != nil {
			jw.logger.Error("Failed to close output file", zap.Error(err))
		}
	}
}

// ----------------------------------------
// CSV Writer
// ----------------------------------------

type CSVWriter struct {
	writer *csv.Writer
	logger *zap.Logger
	file   io.Closer
}

func NewCSVWriter(w io.Writer, logger *zap.Logger) (*CSVWriter, error) {
	cw := &CSVWriter{
		writer: csv.NewWriter(w),
		logger: logger,
	}
	if f, ok := w.(io.Closer); ok {
		cw.file = f
	}
	return cw, nil
}

func (cw *CSVWriter) WriteHeader() {
	header := []string{"host", "port", "status", "service_name", "banner", "error"}
	if err := cw.writer.Write(header); err != nil {
		cw.logger.Error("Failed to write CSV header", zap.Error(err))
	}
}

func (cw *CSVWriter) Write(result types.ScanResult) {
	record := []string{
		result.Target.Host,
		strconv.Itoa(result.Target.Port),
		result.Status.String(),
		result.ServiceName,
		result.Banner,
		result.Error,
	}
	if err := cw.writer.Write(record); err != nil {
		cw.logger.Error("Failed to write CSV record", zap.Error(err))
	}
}

func (cw *CSVWriter) Close() {
	cw.writer.Flush()
	if err := cw.writer.Error(); err != nil {
		cw.logger.Error("Failed to flush CSV writer", zap.Error(err))
	}
	if cw.file != nil {
		if err := cw.file.Close(); err != nil {
			cw.logger.Error("Failed to close output file", zap.Error(err))
		}
	}
}
