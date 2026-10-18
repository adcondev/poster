package connection

import "io"

// Connector define la interfaz para cualquier tipo de conexión con la impresora
type Connector interface {
	io.WriteCloser // Write([]byte) (int, error) y Close() error
}
