package engine

import (
	"context"
	"sync"
	time "time"

	"github.com/user/pentscan/internal/metrics"
	"github.com/user/pentscan/internal/scanner"
	"github.com/user/pentscan/internal/types"
	"github.com/user/pentscan/internal/writer"
	"go.uber.org/zap"
)

// Engine coordena a varredura.
type Engine struct {
	logger  *zap.Logger
	writer  writer.ResultWriter
	scanner scanner.Scanner
	workers int
	rate    int
	timeout time.Duration
}

// New cria um novo motor de varredura.
func New(logger *zap.Logger, writer writer.ResultWriter, scanner scanner.Scanner, workers, rate int, timeout time.Duration) *Engine {
	return &Engine{
		logger:  logger,
		writer:  writer,
		scanner: scanner,
		workers: workers,
		rate:    rate,
		timeout: timeout,
	}
}

// Run inicia a varredura.
func (e *Engine) Run(ctx context.Context, targets []types.Target) {
	var wg sync.WaitGroup
	targetsChan := make(chan types.Target, e.workers)

	go func() {
		defer close(targetsChan)
		var ticker *time.Ticker
		if e.rate > 0 {
			ticker = time.NewTicker(time.Second / time.Duration(e.rate))
			defer ticker.Stop()
		}
		for _, target := range targets {
			if ticker != nil {
				select {
				case <-ticker.C:
				case <-ctx.Done():
					return
				}
			}
			select {
			case targetsChan <- target:
			case <-ctx.Done():
				return
			}
		}
	}()

	for i := 0; i < e.workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case target, ok := <-targetsChan:
					if !ok {
						return
					}
					metrics.TargetsTotal.Inc() // Incrementa o total de alvos processados
					result := e.scanner.Scan(target, e.timeout)

					switch result.Status {
					case types.StatusOpen:
						metrics.PortsOpen.Inc()
					case types.StatusClosed:
						metrics.PortsClosed.Inc()
					case types.StatusFiltered:
						metrics.PortsFiltered.Inc()
					}

					if result.Status == types.StatusOpen || (result.Status == types.StatusClosed && result.Error != "") || result.Status == types.StatusFiltered {
						e.writer.Write(result)
					}
				case <-ctx.Done():
					return
				}
			}
		}()
	}

	wg.Wait()
}
