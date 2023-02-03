package zapx

import (
	"strings"
	"sync"
	"time"

	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

type loggedEntry struct {
	zapcore.Entry
	context []zapcore.Field
}
type cachedLogs struct {
	mu   sync.RWMutex
	logs []loggedEntry
}

func (p *cachedLogs) add(log loggedEntry) {
	p.mu.Lock()
	p.logs = append(p.logs, log)
	p.mu.Unlock()
}

func (p *cachedLogs) clone() *cachedLogs {
	clog := &cachedLogs{
		logs: p.logs[:len(p.logs):len(p.logs)],
	}
	return clog
}

var _cachedCorePool = sync.Pool{
	New: func() interface{} {
		core, _ := newCachedCore(encoder, sink, level)
		return core
	},
}

func getCachedCore() *cachedCore {
	return _cachedCorePool.Get().(*cachedCore)
}

func putCachedCore(core *cachedCore) {
	core.clog.logs = core.clog.logs[:0]
	core.context = core.context[:0]
	_cachedCorePool.Put(core)
}

type cachedCore struct {
	zapcore.LevelEnabler
	enc     zapcore.Encoder
	out     zapcore.WriteSyncer
	clog    *cachedLogs
	context []zapcore.Field
}

func newCachedCore(enc zapcore.Encoder, ws zapcore.WriteSyncer, enab zapcore.LevelEnabler) (zapcore.Core, *cachedLogs) {
	ol := &cachedLogs{}
	return &cachedCore{
		LevelEnabler: enab,
		clog:         ol,
		enc:          enc,
		out:          ws,
	}, ol
}

// Check implement zapcore.Core func Check
func (p *cachedCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if p.Enabled(ent.Level) {
		return ce.AddCore(ent, p)
	}
	return ce
}

// With implement zapcore.Core func With
func (p *cachedCore) With(fields []zapcore.Field) zapcore.Core {
	return &cachedCore{
		LevelEnabler: p.LevelEnabler,
		enc:          p.enc.Clone(),
		out:          p.out,
		clog:         p.clog.clone(),
		context:      append(p.context[:len(p.context):len(p.context)], fields...),
	}
}

// Write implement zapcore.Core func Write
func (p *cachedCore) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	if ent.Message == globalKeyPrefix {
		p.context = append(p.context, fields...)
		return nil
	}
	p.clog.add(loggedEntry{ent, fields})
	return nil
}

// Sync implement zapcore.Core func Sync
func (p *cachedCore) Sync() error {
	defer putCachedCore(p)
	p.appendAll()
	return nil
}

func (p *cachedCore) appendAll() {
	buff := GetBufPool()
	defer buff.Free()
	buff.AppendByte('{')
	p.appendEntrys(buff)
	p.appendGlobalKV(buff)
	buff.AppendByte('}')
	buff.AppendByte('\n')
	_, _ = p.out.Write(buff.Bytes())
}

func (p *cachedCore) appendEntrys(buf *buffer.Buffer) {
	var skips map[string]int8
	if len(p.clog.logs) > 0 {
		skips = make(map[string]int8, len(p.clog.logs))
	}
	for _, log := range p.clog.logs {
		if strings.Contains(log.Message, ".") {
			ms := strings.Split(log.Message, ".")
			if _, ok := skips[ms[0]]; ok {
				continue
			}
			skips[ms[0]] = 0
			p.appendKey(ms[0], buf)
			buf.AppendByte('{')
			for _, llog := range p.clog.logs {
				if strings.HasPrefix(llog.Message, ms[0]+".") {
					mms := strings.Split(llog.Message, ".")
					llog.Message = mms[1]
					p.appendEntry(llog, buf)
				}
			}
			bs := buf.Bytes()
			bs[len(bs)-1] = '}'
			bs = append(bs, ',')
			buf.Reset()
			_, _ = buf.Write(bs)
		} else {
			p.appendEntry(log, buf)
		}
	}
}

func (p *cachedCore) appendEntry(entry loggedEntry, buf *buffer.Buffer) {
	buff, _ := p.enc.EncodeEntry(entry.Entry, entry.context)
	defer buff.Free()
	p.appendKey(entry.Message, buf)
	_, _ = buf.Write(buff.Bytes())
}

func (p *cachedCore) appendKey(key string, buf *buffer.Buffer) {
	buf.AppendByte('"')
	buf.AppendString(key)
	buf.AppendByte('"')
	buf.AppendByte(':')
}

func (p *cachedCore) appendGlobalKV(buf *buffer.Buffer) {
	buff, _ := p.enc.EncodeEntry(zapcore.Entry{Time: time.Now()}, p.context)
	defer buff.Free()
	gb := buff.Bytes()
	gb = gb[1 : len(gb)-2]
	_, _ = buf.Write(gb)
}
