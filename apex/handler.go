package apex

import (
	"encoding/json"
	"net"

	"github.com/apex/log"
)

// Handler converts Apex logs into logd UDP packets and transports them to a local logd server.
type Handler struct {
	conn net.Conn
}

// New creates a new handler.
func New() (*Handler, error) {
	conn, err := net.Dial("udp", ":9044")
	if err != nil {
		return nil, err
	}

	return &Handler{
		conn: conn,
	}, nil
}

// HandleLog by converting to UDP and sending to logd.
func (h *Handler) HandleLog(e *log.Entry) error {
	b, err := json.Marshal(e)
	if err != nil {
		return err
	}
	_, err = h.conn.Write(b)
	return err
}
