package connection

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

const (
	// DefaultTCPPort is the standard port for network thermal printers (RAW 9100).
	DefaultTCPPort = "9100"

	// DefaultDialTimeout is the default connection timeout.
	DefaultDialTimeout = 5 * time.Second

	// DefaultWriteTimeout is the default timeout for write operations.
	DefaultWriteTimeout = 10 * time.Second
)

// TCPOption configures a TCPConnector.
type TCPOption func(*TCPConnector)

// WithDialTimeout sets the network dial timeout.
func WithDialTimeout(d time.Duration) TCPOption {
	return func(c *TCPConnector) {
		c.dialTimeout = d
	}
}

// WithWriteTimeout sets the network write deadline timeout.
func WithWriteTimeout(d time.Duration) TCPOption {
	return func(c *TCPConnector) {
		c.writeTimeout = d
	}
}

// TCPConnector implements Connector over a raw TCP socket connection.
type TCPConnector struct {
	conn         net.Conn
	address      string
	dialTimeout  time.Duration
	writeTimeout time.Duration
	mu           sync.Mutex
	closed       bool
}

// NewTCPConnector connects to a network thermal printer at the given address.
// If address does not include a port, DefaultTCPPort (9100) is appended.
func NewTCPConnector(address string, opts ...TCPOption) (*TCPConnector, error) {
	trimmed := strings.TrimSpace(address)
	if trimmed == "" {
		return nil, errors.New("tcp address cannot be empty")
	}

	if !strings.Contains(trimmed, ":") {
		trimmed = trimmed + ":" + DefaultTCPPort
	}

	c := &TCPConnector{
		address:      trimmed,
		dialTimeout:  DefaultDialTimeout,
		writeTimeout: DefaultWriteTimeout,
	}

	for _, opt := range opts {
		if opt != nil {
			opt(c)
		}
	}

	dialer := &net.Dialer{
		Timeout: c.dialTimeout,
	}
	conn, err := dialer.DialContext(context.Background(), "tcp", c.address)
	if err != nil {
		return nil, fmt.Errorf("dialing tcp %s: %w", c.address, err)
	}

	if tcpConn, ok := conn.(*net.TCPConn); ok {
		_ = tcpConn.SetKeepAlive(true)
	}

	c.conn = conn
	return c, nil
}

// Write sends raw ESC/POS bytes over the TCP connection.
func (c *TCPConnector) Write(data []byte) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return 0, errors.New("tcp connector is closed")
	}

	if len(data) == 0 {
		return 0, nil
	}

	if c.writeTimeout > 0 {
		_ = c.conn.SetWriteDeadline(time.Now().Add(c.writeTimeout))
	}

	n, err := c.conn.Write(data)
	if err != nil {
		return n, fmt.Errorf("tcp write to %s: %w", c.address, err)
	}
	return n, nil
}

// Close terminates the TCP connection. It is idempotent.
func (c *TCPConnector) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return nil
	}
	c.closed = true

	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			return fmt.Errorf("closing tcp connection to %s: %w", c.address, err)
		}
	}
	return nil
}

// Address returns the normalized TCP target address.
func (c *TCPConnector) Address() string {
	return c.address
}

// Ensure TCPConnector implements Connector.
var _ Connector = (*TCPConnector)(nil)
