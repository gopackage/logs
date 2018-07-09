package apex

import (
	"encoding/json"
	"net"

	"github.com/apex/log"
	"github.com/gopackage/logs"
)

// Handler converts Apex logs into logd UDP packets and transports them to a local logd server.
type Handler struct {
	conn  net.Conn
	buf   chan *log.Entry
	stats logs.Stats
}

// New creates a new handler.
func New() (*Handler, error) {
	return NewHandler(":9044", make(chan *log.Entry, 1000), nil)
}

// NewHandler creates a new handler with specific settings
func NewHandler(target string, buffer chan *log.Entry, stats logs.Stats) (*Handler, error) {
	conn, err := net.Dial("udp", target)
	if err != nil {
		return nil, err
	}

	h := &Handler{
		conn:  conn,
		stats: stats,
	}
	h.Buffer(buffer)
	return h, nil
}

// Stats configures the statistics handler for the logger. If set to nil (the default)
// the logger does not report stats.
func (h *Handler) Stats(stats logs.Stats) {
	h.stats = stats
}

// Buffer configures the internal buffer to use a provided channel. If the channel
// is nil, the method returns the current buffer channel.
func (h *Handler) Buffer(buf chan *log.Entry) chan *log.Entry {
	if buf != nil {
		if h.buf != nil {
			close(h.buf) // Close the current buffer
		}
		h.buf = buf
		go h.start() // Start a go routine to process the buffer
	}
	return h.buf
}

// HandleLog by converting to UDP and sending to logd.
func (h *Handler) HandleLog(e *log.Entry) error {
	h.buf <- e
	return nil
}

// start the handler proessing the current buffer (blocks until the buffer is closed).
func (h *Handler) start() error {
	buf := h.buf
	for {
		e, ok := <-buf
		if ok {
			b, err := json.Marshal(e)
			if err != nil {
				return err
			}
			size, err := h.conn.Write(b)
			if err != nil {
				return err
			}
			if h.stats != nil {
				h.stats.Count("logs.sent", 1)
				h.stats.Count("logs.size", size)
			}
		} else {
			return nil
		}
	}
}
