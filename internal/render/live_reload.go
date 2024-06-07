package render

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

type LiveRealod struct {
	token string
	store []*chan string
}

func NewLiveReloadSever(token string) *LiveRealod {
	lr := &LiveRealod{
		store: make([]*chan string, 0),
		token: token,
	}
	return lr
}

func (lr *LiveRealod) removeChannel(ch *chan string) {
	pos := -1
	storeLen := len(lr.store)
	for i, msgChan := range lr.store {
		if ch == msgChan {
			pos = i
		}
	}

	if pos == -1 {
		return
	}
	lr.store[pos] = lr.store[storeLen-1]
	lr.store = lr.store[:storeLen-1]
	slog.Debug("Connection remains", "count", len(lr.store))
}

func (lr *LiveRealod) Broadcast(msg string) {
	for _, ch := range lr.store {
		*ch <- msg
	}
}

func (lr *LiveRealod) CreateLiveReloadHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		lr.ExecuteLiveReloadHandler(w, r)
	}
}

func (lr *LiveRealod) ExecuteLiveReloadHandler(w http.ResponseWriter, r *http.Request) {
	ch := make(chan string)
	lr.store = append(lr.store, &ch)

	slog.Debug("Client connected", "count", len(lr.store))

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	defer func() {
		close(ch)
		lr.removeChannel(&ch)
		slog.Debug("Client closed connection")
	}()

	flusher, ok := w.(http.Flusher)
	if !ok {
		slog.Error("Could not init http.Flusher")
		os.Exit(1)
	}

	_, err := fmt.Fprintf(w, "data: %s\n\n", lr.token)
	if err != nil {
		slog.Error("Could write SSE message", "error", err.Error())
		os.Exit(1)
	}
	flusher.Flush()

	for {
		select {
		case message := <-ch:
			_, err := fmt.Fprintf(w, "data: %s\n\n", message)
			if err != nil {
				slog.Error("Could write SSE message", "error", err.Error())
				os.Exit(1)
				return
			}
			flusher.Flush()
		case <-r.Context().Done():
			slog.Debug("Client closed connection")
			return
		}
	}
}
