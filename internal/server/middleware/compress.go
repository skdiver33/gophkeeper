package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

type gzipWriter struct {
	http.ResponseWriter
}

func (w gzipWriter) Write(b []byte) (int, error) {
	//use compression if size>100 bytes
	if len(b) > 100 {
		gz, err := gzip.NewWriterLevel(w.ResponseWriter, gzip.BestSpeed)
		if err != nil {
			slog.Error("create gzip compressor for response", "error", err)
			return w.ResponseWriter.Write(b)
		}
		w.Header().Set("Content-Encoding", "gzip")
		count, err := gz.Write(b)
		gz.Close()
		return count, err
	}
	return w.ResponseWriter.Write(b)
}

func GzipHandle(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Encoding") == "gzip" {
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				slog.Error("create gzip archivator", "error", err.Error())
				return
			}
			decompressBody, err := io.ReadAll(gz)
			if err != nil {
				slog.Error("decompress body", "error", err.Error())
				return
			}
			gz.Close()
			r.Body = io.NopCloser(bytes.NewReader(decompressBody))
			r.ContentLength = int64(len(decompressBody))
		}

		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}
		next.ServeHTTP(gzipWriter{ResponseWriter: w}, r)
	})
}
