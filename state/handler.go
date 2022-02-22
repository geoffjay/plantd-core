package main

type Handler struct {
	running   bool
	cancelled bool
	callbacks map[string]HandlerCallback
}

type HandlerCallback interface {
	Execute(msgBody string) ([]byte, error)
}

func NewHandler() *Handler {
	return &Handler{
		running:   false,
		cancelled: false,
		callbacks: make(map[string]HandlerCallback),
	}
}

func (h *Handler) AddCallback(name string, callback HandlerCallback) {
	h.callbacks[name] = callback
}
