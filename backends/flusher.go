package backends

import "io"
import "sync"
import "time"
import "net/http"

type flushableWriter interface {
	io.Writer
	http.Flusher
}

type WriterFlusher struct {
	destination   flushableWriter
	flushInterval time.Duration
	lock          sync.Mutex
	done          chan struct{}
}

func IsFlushable(deatination io.Writer) bool {
	_, ok := destination.(flushableWriter)
	return ok
}

func WrapWriterAsFlushable(destination io.Writer, flushInterval time.Duration) *WriterFlusher {
	return &WriterFlusher{
		destination:   destination,
		flushInterval: flushInterval,
		done:          make(chan struct{}),
	}
}

func (wr *WriterFlusher) Write(p []byte) (int, error) {
	wr.lock.Lock()
	defer wr.lock.Unlock()

	return wr.destination.Write(p)
}

func (wr *WriterFlusher) Start() {
	go func() {
		flushTick := time.NewTicker(wr.flushInterval)
		defer flushTick.Stop()

		for {
			select {
			case <-wr.done:
				return
			case <-flushTick.C:
				wr.lock.Lock()
				wr.destination.Flush()
				wr.lock.Unlock()
			}
		}
	}()
}

func (wr *WriterFlusher) Stop() {
	wr.done <- struct{}{}
}
