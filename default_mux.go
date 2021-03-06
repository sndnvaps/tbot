package tbot

import "regexp"

// DefaultMux is a default multiplexer
// Supports parametrized commands
type DefaultMux struct {
	handlers       Handlers
	defaultHandler *Handler
}

// NewDefaultMux creates new DefaultMux
func NewDefaultMux() Mux {
	return &DefaultMux{handlers: make(Handlers)}
}

// Handlers returns list of handlers currently in mux
func (dm *DefaultMux) Handlers() Handlers {
	return dm.handlers
}

// DefaultHandler returns default handler, nil if it's not set
func (dm *DefaultMux) DefaultHandler() *Handler {
	return dm.defaultHandler
}

// Mux takes message content and returns corresponding handler
// and parsed vars from message
func (dm *DefaultMux) Mux(path string) (*Handler, MessageVars) {
	for _, handler := range dm.handlers {
		re := regexp.MustCompile(handler.pattern)
		matches := re.FindStringSubmatch(path)

		if len(matches) > 0 {
			messageData := make(map[string]string)
			matches := matches[1:]
			for i, match := range matches {
				messageData[handler.variables[i]] = match
			}
			return handler, messageData
		}
	}
	return dm.defaultHandler, nil
}

// HandleFunc adds new handler function to mux
func (dm *DefaultMux) HandleFunc(path string, handler HandlerFunction, description ...string) {
	dm.handlers[path] = NewHandler(handler, path, description...)
}

// Handle is a shortcut for HandleFunc to reply just with text
func (dm *DefaultMux) Handle(path string, reply string, description ...string) {
	f := func(m Message) {
		m.Reply(reply)
	}
	dm.HandleFunc(path, f, description...)
}

// HandleDefault adds new default handler, when nothing matches with message
func (dm *DefaultMux) HandleDefault(handler HandlerFunction, description ...string) {
	dm.defaultHandler = NewHandler(handler, "", description...)
}
