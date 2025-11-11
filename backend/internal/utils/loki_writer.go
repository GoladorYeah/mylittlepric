package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// LokiWriter sends logs to Grafana Loki
type LokiWriter struct {
	endpoint   string
	labels     map[string]string
	client     *http.Client
	buffer     []lokiEntry
	bufferLock sync.Mutex
	batchSize  int
	flushTimer *time.Timer
	ctx        context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup
}

type lokiEntry struct {
	timestamp time.Time
	line      string
}

type lokiStream struct {
	Stream map[string]string `json:"stream"`
	Values [][]string        `json:"values"`
}

type lokiRequest struct {
	Streams []lokiStream `json:"streams"`
}

// NewLokiWriter creates a new Loki writer
func NewLokiWriter(endpoint string, labels map[string]string) *LokiWriter {
	ctx, cancel := context.WithCancel(context.Background())

	lw := &LokiWriter{
		endpoint:   endpoint,
		labels:     labels,
		client:     &http.Client{Timeout: 10 * time.Second},
		buffer:     make([]lokiEntry, 0, 100),
		batchSize:  100,
		ctx:        ctx,
		cancel:     cancel,
		flushTimer: time.NewTimer(5 * time.Second),
	}

	// Start background flusher
	lw.wg.Add(1)
	go lw.backgroundFlusher()

	return lw
}

// Write implements io.Writer interface
func (lw *LokiWriter) Write(p []byte) (n int, err error) {
	if lw == nil || lw.endpoint == "" {
		// If Loki is not configured, just return success
		return len(p), nil
	}

	entry := lokiEntry{
		timestamp: time.Now(),
		line:      string(p),
	}

	lw.bufferLock.Lock()
	lw.buffer = append(lw.buffer, entry)
	shouldFlush := len(lw.buffer) >= lw.batchSize
	lw.bufferLock.Unlock()

	if shouldFlush {
		go lw.flush()
	}

	return len(p), nil
}

// backgroundFlusher periodically flushes the buffer
func (lw *LokiWriter) backgroundFlusher() {
	defer lw.wg.Done()

	for {
		select {
		case <-lw.ctx.Done():
			// Final flush before exit
			lw.flush()
			return
		case <-lw.flushTimer.C:
			lw.flush()
			lw.flushTimer.Reset(5 * time.Second)
		}
	}
}

// flush sends buffered logs to Loki
func (lw *LokiWriter) flush() {
	lw.bufferLock.Lock()
	if len(lw.buffer) == 0 {
		lw.bufferLock.Unlock()
		return
	}

	// Copy buffer and clear it
	entries := make([]lokiEntry, len(lw.buffer))
	copy(entries, lw.buffer)
	lw.buffer = lw.buffer[:0]
	lw.bufferLock.Unlock()

	// Convert entries to Loki format
	values := make([][]string, len(entries))
	for i, entry := range entries {
		// Loki expects timestamp in nanoseconds as string
		ts := fmt.Sprintf("%d", entry.timestamp.UnixNano())
		values[i] = []string{ts, entry.line}
	}

	stream := lokiStream{
		Stream: lw.labels,
		Values: values,
	}

	req := lokiRequest{
		Streams: []lokiStream{stream},
	}

	// Send to Loki
	if err := lw.sendToLoki(req); err != nil {
		// Log error but don't fail - we don't want logging to break the app
		fmt.Printf("Failed to send logs to Loki: %v\n", err)
	}
}

// sendToLoki sends the request to Loki API
func (lw *LokiWriter) sendToLoki(req lokiRequest) error {
	data, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(lw.ctx, "POST", lw.endpoint, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := lw.client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("loki returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// Close closes the Loki writer and flushes remaining logs
func (lw *LokiWriter) Close() error {
	if lw == nil {
		return nil
	}

	// Stop the background flusher
	lw.cancel()

	// Wait for background flusher to finish
	lw.wg.Wait()

	return nil
}
