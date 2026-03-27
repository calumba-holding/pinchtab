package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/pinchtab/pinchtab/internal/bridge"
	"github.com/pinchtab/pinchtab/internal/bridge/observe"
	"github.com/pinchtab/pinchtab/internal/httpx"
)

// HandleNetworkExport exports captured network data in a registered format (HAR, NDJSON, etc.).
//
// @Endpoint GET /network/export
// @Description Export captured network entries in HAR 1.2, NDJSON, or other registered formats
//
// @Param tabId   string query Tab ID (optional, uses current tab)
// @Param format  string query Export format: "har" (default), "ndjson", or any registered format
// @Param output  string query "file" to save to disk (optional)
// @Param path    string query Filename when output=file (optional, auto-generated if omitted)
// @Param body    string query "true" to include response bodies (can be slow)
// @Param filter  string query URL pattern filter
// @Param method  string query HTTP method filter
// @Param status  string query Status code range filter (e.g. "4xx")
// @Param type    string query Resource type filter
// @Param limit   string query Maximum entries to export
//
// @Response 200 application/har+json|application/x-ndjson  Exported data (streamed)
// @Response 200 application/json                           File save result when output=file
// @Response 400 application/json                           Invalid format or parameters
// @Response 500 application/json                           Export error
func (h *Handlers) HandleNetworkExport(w http.ResponseWriter, r *http.Request) {
	if err := h.ensureChrome(); err != nil {
		httpx.Error(w, 500, fmt.Errorf("chrome initialization: %w", err))
		return
	}

	tabID := r.URL.Query().Get("tabId")
	tabCtx, resolvedTabID, err := h.tabContext(r, tabID)
	if err != nil {
		httpx.Error(w, 404, err)
		return
	}
	if _, ok := h.enforceCurrentTabDomainPolicy(w, r, tabCtx, resolvedTabID); !ok {
		return
	}

	// Resolve format from registry
	formatName := r.URL.Query().Get("format")
	if formatName == "" {
		formatName = "har"
	}
	factory := observe.GetFormat(formatName)
	if factory == nil {
		httpx.JSON(w, 400, map[string]any{
			"code":      "unknown_format",
			"error":     fmt.Sprintf("unknown export format %q", formatName),
			"available": observe.ListFormats(),
		})
		return
	}

	// Get network entries
	nm := h.Bridge.NetworkMonitor()
	if nm == nil {
		// No monitor — export empty
		enc := factory("PinchTab", h.version())
		w.Header().Set("Content-Type", enc.ContentType())
		_ = enc.Start(w)
		_ = enc.Finish()
		return
	}

	bufferSize := parseBufferSize(r)
	buf := nm.GetBuffer(resolvedTabID)
	if buf == nil {
		if err := nm.StartCaptureWithSize(tabCtx, resolvedTabID, bufferSize); err != nil {
			httpx.Error(w, 500, fmt.Errorf("start capture: %w", err))
			return
		}
		buf = nm.GetBuffer(resolvedTabID)
	}

	filter := parseNetworkFilter(r)
	entries := buf.List(filter)
	if filter.Limit > 0 && len(entries) > filter.Limit {
		entries = entries[len(entries)-filter.Limit:]
	}

	includeBody := r.URL.Query().Get("body") == "true"
	output := r.URL.Query().Get("output")

	// Convert entries
	exportEntries := make([]observe.ExportEntry, 0, len(entries))
	for _, entry := range entries {
		var body string
		var b64 bool
		if includeBody && entry.Finished && !entry.Failed {
			body, b64, _ = nm.GetResponseBody(tabCtx, entry.RequestID)
		}
		exportEntries = append(exportEntries, observe.NetworkEntryToExport(entry, body, b64))
	}

	enc := factory("PinchTab", h.version())

	if output == "file" {
		if err := h.writeExportFile(w, r, enc, exportEntries, formatName); err != nil {
			httpx.Error(w, 500, fmt.Errorf("write file: %w", err))
		}
		return
	}

	// Stream to response
	w.Header().Set("Content-Type", enc.ContentType())
	if err := enc.Start(w); err != nil {
		return
	}
	for _, entry := range exportEntries {
		if err := enc.Encode(entry); err != nil {
			return
		}
	}
	_ = enc.Finish()
}

func (h *Handlers) writeExportFile(
	w http.ResponseWriter,
	r *http.Request,
	enc observe.ExportEncoder,
	entries []observe.ExportEntry,
	formatName string,
) error {
	userPath := r.URL.Query().Get("path")
	if userPath == "" {
		ts := time.Now().Format("20060102-150405")
		userPath = fmt.Sprintf("network-%s%s", ts, enc.FileExtension())
	}

	dir := filepath.Join(h.Config.StateDir, "exports")
	if err := os.MkdirAll(dir, 0750); err != nil {
		return fmt.Errorf("create dir: %w", err)
	}

	finalPath := filepath.Join(dir, filepath.Base(userPath))
	tmpPath := finalPath + ".tmp"

	f, err := os.OpenFile(tmpPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}

	if err := enc.Start(f); err != nil {
		f.Close()
		os.Remove(tmpPath)
		return err
	}
	for _, entry := range entries {
		if err := enc.Encode(entry); err != nil {
			f.Close()
			os.Remove(tmpPath)
			return err
		}
	}
	if err := enc.Finish(); err != nil {
		f.Close()
		os.Remove(tmpPath)
		return err
	}
	if err := f.Close(); err != nil {
		os.Remove(tmpPath)
		return err
	}

	if err := os.Rename(tmpPath, finalPath); err != nil {
		os.Remove(tmpPath)
		return err
	}

	httpx.JSON(w, 200, map[string]any{
		"path":    finalPath,
		"entries": len(entries),
		"format":  formatName,
	})
	return nil
}

// HandleNetworkExportStream streams network entries to a file as they arrive.
//
// @Endpoint GET /network/export/stream
// @Description Live capture: write entries to file as they are captured
//
// @Param tabId   string query Tab ID
// @Param format  string query Export format (default: har)
// @Param path    string query Output filename (required)
// @Param body    string query "true" to include response bodies
// @Param filter... (same as HandleNetworkExport)
//
// @Response 200 text/event-stream  SSE progress events
func (h *Handlers) HandleNetworkExportStream(w http.ResponseWriter, r *http.Request) {
	if err := h.ensureChrome(); err != nil {
		httpx.Error(w, 500, fmt.Errorf("chrome initialization: %w", err))
		return
	}

	userPath := r.URL.Query().Get("path")
	if userPath == "" {
		httpx.Error(w, 400, fmt.Errorf("path required for streaming export"))
		return
	}

	tabID := r.URL.Query().Get("tabId")
	tabCtx, resolvedTabID, err := h.tabContext(r, tabID)
	if err != nil {
		httpx.Error(w, 404, err)
		return
	}
	if _, ok := h.enforceCurrentTabDomainPolicy(w, r, tabCtx, resolvedTabID); !ok {
		return
	}

	formatName := r.URL.Query().Get("format")
	if formatName == "" {
		formatName = "har"
	}
	factory := observe.GetFormat(formatName)
	if factory == nil {
		httpx.JSON(w, 400, map[string]any{
			"code":      "unknown_format",
			"error":     fmt.Sprintf("unknown format %q", formatName),
			"available": observe.ListFormats(),
		})
		return
	}

	nm := h.Bridge.NetworkMonitor()
	if nm == nil {
		httpx.Error(w, 500, fmt.Errorf("network monitor not available"))
		return
	}

	bufferSize := parseBufferSize(r)
	buf := nm.GetBuffer(resolvedTabID)
	if buf == nil {
		if err := nm.StartCaptureWithSize(tabCtx, resolvedTabID, bufferSize); err != nil {
			httpx.Error(w, 500, fmt.Errorf("start capture: %w", err))
			return
		}
		buf = nm.GetBuffer(resolvedTabID)
	}

	includeBody := r.URL.Query().Get("body") == "true"
	filter := parseNetworkFilter(r)

	// Open file
	dir := filepath.Join(h.Config.StateDir, "exports")
	if err := os.MkdirAll(dir, 0750); err != nil {
		httpx.Error(w, 500, fmt.Errorf("create dir: %w", err))
		return
	}
	finalPath := filepath.Join(dir, filepath.Base(userPath))
	f, err := os.OpenFile(finalPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		httpx.Error(w, 500, fmt.Errorf("create file: %w", err))
		return
	}

	enc := factory("PinchTab", h.version())
	if err := enc.Start(f); err != nil {
		f.Close()
		httpx.Error(w, 500, fmt.Errorf("start encoder: %w", err))
		return
	}

	// Subscribe for live entries
	subID, ch := buf.Subscribe()
	defer buf.Unsubscribe(subID)

	// SSE setup
	flusher, ok := w.(http.Flusher)
	if !ok {
		enc.Finish()
		f.Close()
		httpx.Error(w, 500, fmt.Errorf("streaming not supported"))
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	count := 0
	keepalive := time.NewTicker(15 * time.Second)
	defer keepalive.Stop()

	for {
		select {
		case <-r.Context().Done():
			_ = enc.Finish()
			f.Close()
			fmt.Fprintf(w, "event: done\ndata: {\"entries\":%d,\"path\":%q}\n\n", count, finalPath)
			flusher.Flush()
			return

		case entry, ok := <-ch:
			if !ok {
				_ = enc.Finish()
				f.Close()
				fmt.Fprintf(w, "event: done\ndata: {\"entries\":%d,\"path\":%q}\n\n", count, finalPath)
				flusher.Flush()
				return
			}
			if !filter.Match(entry) {
				continue
			}
			var body string
			var b64 bool
			if includeBody && entry.Finished && !entry.Failed {
				body, b64, _ = nm.GetResponseBody(tabCtx, entry.RequestID)
			}
			export := observe.NetworkEntryToExport(entry, body, b64)
			if err := enc.Encode(export); err != nil {
				_ = enc.Finish()
				f.Close()
				return
			}
			count++
			fmt.Fprintf(w, "event: export\ndata: {\"entries\":%d,\"url\":%q}\n\n", count, truncateURL(entry.URL))
			flusher.Flush()

		case <-keepalive.C:
			fmt.Fprintf(w, ": keepalive\n\n")
			flusher.Flush()
		}
	}
}

// HandleTabNetworkExport handles GET /tabs/{id}/network/export.
func (h *Handlers) HandleTabNetworkExport(w http.ResponseWriter, r *http.Request) {
	tabID := r.PathValue("id")
	if tabID == "" {
		httpx.Error(w, 400, fmt.Errorf("tab id required"))
		return
	}
	q := r.URL.Query()
	q.Set("tabId", tabID)
	req := r.Clone(r.Context())
	u := *r.URL
	u.RawQuery = q.Encode()
	req.URL = &u
	h.HandleNetworkExport(w, req)
}

// HandleTabNetworkExportStream handles GET /tabs/{id}/network/export/stream.
func (h *Handlers) HandleTabNetworkExportStream(w http.ResponseWriter, r *http.Request) {
	tabID := r.PathValue("id")
	if tabID == "" {
		httpx.Error(w, 400, fmt.Errorf("tab id required"))
		return
	}
	q := r.URL.Query()
	q.Set("tabId", tabID)
	req := r.Clone(r.Context())
	u := *r.URL
	u.RawQuery = q.Encode()
	req.URL = &u
	h.HandleNetworkExportStream(w, req)
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func parseNetworkFilter(r *http.Request) bridge.NetworkFilter {
	f := bridge.NetworkFilter{
		URLPattern:   r.URL.Query().Get("filter"),
		Method:       r.URL.Query().Get("method"),
		StatusRange:  r.URL.Query().Get("status"),
		ResourceType: r.URL.Query().Get("type"),
	}
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			f.Limit = n
		}
	}
	return f
}

func (h *Handlers) version() string {
	return "dev"
}

func truncateURL(u string) string {
	if len(u) > 120 {
		return u[:120]
	}
	return u
}
