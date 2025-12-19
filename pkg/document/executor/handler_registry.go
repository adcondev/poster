package executor

// HandlerRegistry mantiene el registro de handlers
type HandlerRegistry struct {
	handlers map[string]CommandHandler
}

// NewRegistry crea un nuevo registro
func NewRegistry() *HandlerRegistry {
	return &HandlerRegistry{
		handlers: make(map[string]CommandHandler),
	}
}

// Register registra un handler
func (r *HandlerRegistry) Register(cmdType string, handler CommandHandler) {
	r.handlers[cmdType] = handler
}

// Get obtiene un handler
func (r *HandlerRegistry) Get(cmdType string) (CommandHandler, bool) {
	h, ok := r.handlers[cmdType]
	return h, ok
}

// List retorna todos los tipos registrados
func (r *HandlerRegistry) List() []string {
	types := make([]string, 0, len(r.handlers))
	for t := range r.handlers {
		types = append(types, t)
	}
	return types
}
