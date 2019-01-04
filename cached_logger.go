package zapx

import (
	"sync"

	"go.uber.org/zap"
)

// CachedLogger cached logger
type CachedLogger struct {
	*zap.Logger
}

// Add add kv pair in root
func (p *CachedLogger) Add(fields ...zap.Field) {
	p.Info(globalKeyPrefix, fields...)
}

// Flush flush cached log
func (p *CachedLogger) Flush(wg *sync.WaitGroup) {
	if wg != nil {
		wg.Wait()
	}
	p.Sync()
}

// CachedSugar get sugared logger from cached logger
func (p *CachedLogger) CachedSugar() *CachedSugaredLogger {
	return &CachedSugaredLogger{SugaredLogger: p.Sugar()}
}

// CachedSugaredLogger cached logger struct
type CachedSugaredLogger struct {
	*zap.SugaredLogger
}

// Flush flush cached sugar log
func (p *CachedSugaredLogger) Flush(wg *sync.WaitGroup) {
	if wg != nil {
		wg.Wait()
	}
	p.Sync()
}

// DeCachedSugar get cached logger from sugared logger
func (p *CachedSugaredLogger) DeCachedSugar() *CachedLogger {
	return &CachedLogger{Logger: p.SugaredLogger.Desugar()}
}
