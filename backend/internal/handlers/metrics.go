package handlers

import (
	"bufio"
	"log"
	"net"
	"net/http"
	"net/url"
	"runtime/debug"

	"github.com/gofiber/fiber/v2"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"mylittleprice/internal/utils"
)

// MetricsHandler handles Prometheus metrics endpoint
type MetricsHandler struct {
	// Pre-created handler to avoid recreating on every request
	handler http.Handler
}

// NewMetricsHandler creates a new metrics handler
func NewMetricsHandler() *MetricsHandler {
	// Create the Prometheus handler once at initialization
	// This is more efficient and avoids potential issues with recreating handlers
	promHandler := promhttp.Handler()

	return &MetricsHandler{
		handler: promHandler,
	}
}

// fiberResponseWriter adapts fiber.Ctx to http.ResponseWriter
type fiberResponseWriter struct {
	ctx        *fiber.Ctx
	header     http.Header
	statusCode int
	written    bool
}

func newFiberResponseWriter(ctx *fiber.Ctx) *fiberResponseWriter {
	return &fiberResponseWriter{
		ctx:    ctx,
		header: make(http.Header),
	}
}

func (w *fiberResponseWriter) Header() http.Header {
	return w.header
}

func (w *fiberResponseWriter) Write(b []byte) (int, error) {
	if !w.written {
		w.written = true
		if w.statusCode == 0 {
			w.statusCode = http.StatusOK
		}

		// Write status code
		w.ctx.Status(w.statusCode)

		// Copy headers from our http.Header to fiber response
		for key, values := range w.header {
			for _, value := range values {
				w.ctx.Set(key, value)
			}
		}
	}
	return w.ctx.Write(b)
}

func (w *fiberResponseWriter) WriteHeader(statusCode int) {
	if !w.written {
		w.statusCode = statusCode
	}
}

// Hijack implements http.Hijacker
func (w *fiberResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return nil, nil, http.ErrNotSupported
}

// GetMetrics returns Prometheus metrics
// @Summary Get Prometheus metrics
// @Description Returns Prometheus metrics in text format
// @Tags monitoring
// @Produce plain
// @Success 200 {string} string "Prometheus metrics"
// @Router /metrics [get]
func (h *MetricsHandler) GetMetrics(c *fiber.Ctx) error {
	utils.LogDebug(c.Context(), "📊 Metrics requested")

	// Wrap handler with panic recovery and detailed logging
	defer func() {
		if r := recover(); r != nil {
			stack := debug.Stack()
			log.Printf("❌ Panic in metrics handler: %v\nStack trace:\n%s", r, string(stack))
			c.Status(fiber.StatusInternalServerError).SendString("Error collecting metrics")
		}
	}()

	// Create our custom ResponseWriter that wraps fiber.Ctx
	w := newFiberResponseWriter(c)

	// Create http.Request from fiber.Ctx
	// Convert fasthttp.URI to standard url.URL
	uri := c.Context().URI()
	path := string(uri.Path())
	if len(uri.QueryString()) > 0 {
		path = path + "?" + string(uri.QueryString())
	}

	req := &http.Request{
		Method:     c.Method(),
		URL:        &url.URL{Path: path},
		Header:     make(http.Header),
		RemoteAddr: c.IP(),
	}

	// Copy headers from fiber to http.Request
	c.Request().Header.VisitAll(func(key, value []byte) {
		req.Header.Add(string(key), string(value))
	})

	log.Printf("🔍 Calling Prometheus handler")

	// Call Prometheus handler with proper error recovery
	func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("❌ Panic inside Prometheus handler: %v", r)
				log.Printf("Stack: %s", debug.Stack())
			}
		}()
		h.handler.ServeHTTP(w, req)
	}()

	log.Printf("✅ Prometheus handler completed, statusCode: %d, headers: %v", w.statusCode, w.header)

	return nil
}
