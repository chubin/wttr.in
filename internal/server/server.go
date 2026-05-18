package server

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
	stdlog "log"
	"net/http"
	"path"
	"runtime/debug"
	"strings"
	"time"

	"github.com/chubin/wttr.in/internal/assets"
	"github.com/chubin/wttr.in/internal/logging"
	"github.com/chubin/wttr.in/internal/weather"
)

const logLineStart = "LOG_LINE_START "

var ErrNoServersConfigured = errors.New("no servers configured")

// certMap maps domain name -> certificate
type certMap map[string]*tls.Certificate

func suppressMessages() []string {
	return []string{
		"error reading preface from client",
		"TLS handshake error from",
		"URL query contains semicolon, which is no longer a supported separator",
		"connection error: PROTOCOL_ERROR",
	}
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	data, err := assets.GetFile("share/static/favicon.ico")
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "image/x-icon")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

// staticFilesHandler serves embedded files under /files/
func staticFilesHandler(w http.ResponseWriter, r *http.Request) {
	// Remove /files/ prefix
	filePath := strings.TrimPrefix(r.URL.Path, "/files/")
	if filePath == "" || filePath == "/" {
		http.NotFound(w, r)
		return
	}

	// Clean the path to prevent directory traversal
	filePath = path.Clean(filePath)
	if strings.HasPrefix(filePath, "..") {
		http.NotFound(w, r)
		return
	}

	embedPath := "share/static/" + filePath
	data, err := assets.GetFile(embedPath)
	if err != nil {
		log.Println(err)
		http.NotFound(w, r)
		return
	}

	// Set proper Content-Type
	contentType := "application/octet-stream"
	switch path.Ext(filePath) {
	case ".css":
		contentType = "text/css; charset=utf-8"
	case ".js":
		contentType = "application/javascript; charset=utf-8"
	case ".png":
		contentType = "image/png"
	case ".ico":
		contentType = "image/x-icon"
	case ".html":
		contentType = "text/html; charset=utf-8"
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

// panicRecovery is a middleware that recovers from panics in handlers
func panicRecovery(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				stack := debug.Stack()
				log.Printf("PANIC in handler %s %s: %v\n%s", r.Method, r.URL.Path, rec, stack)
				fmt.Fprintf(log.Writer(), "PANIC recovered: %v\n%s\n", rec, stack)

				if !wroteHeader(w) {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}
		}()
		next(w, r)
	}
}

func wroteHeader(w http.ResponseWriter) bool {
	type headerWriter interface {
		WroteHeader() bool
	}
	if hw, ok := w.(headerWriter); ok {
		return hw.WroteHeader()
	}
	return false
}

func serveHTTP(mux *http.ServeMux, port int, logFile io.Writer, errs chan<- error) {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		ErrorLog:     stdlog.New(logFile, logLineStart, stdlog.LstdFlags),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  1 * time.Second,
		Handler:      mux,
	}
	errs <- srv.ListenAndServe()
}

// loadCertificates loads all certificates and returns a map + default cert
func loadCertificates(conf *Config) (certMap, *tls.Certificate, error) {
	certs := make(certMap)
	var defaultCert *tls.Certificate

	// Load new multi-certificate config
	for _, tc := range conf.TLSCerts {
		cert, err := tls.LoadX509KeyPair(tc.CertFile, tc.KeyFile)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to load certificate for domain %s: %w", tc.Domain, err)
		}
		certs[tc.Domain] = &cert

		if defaultCert == nil {
			defaultCert = &cert
		}
	}

	// Backward compatibility: legacy single certificate
	if len(certs) == 0 && conf.HasLegacyCert() {
		cert, err := tls.LoadX509KeyPair(conf.TLSCertFile, conf.TLSKeyFile)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to load legacy certificate: %w", err)
		}
		defaultCert = &cert
		certs[""] = defaultCert // empty key = default
	}

	return certs, defaultCert, nil
}

// serveHTTPS starts HTTPS server with support for multiple certificates via SNI
func serveHTTPS(mux *http.ServeMux, port int, conf *Config, logFile io.Writer, errs chan<- error) {
	certs, defaultCert, err := loadCertificates(conf)
	if err != nil {
		errs <- fmt.Errorf("TLS certificate loading failed: %w", err)
		return
	}

	tlsConfig := &tls.Config{
		// This function is called for every TLS handshake (SNI)
		GetCertificate: func(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
			domain := hello.ServerName

			// Exact match
			if cert, ok := certs[domain]; ok {
				return cert, nil
			}

			// Wildcard support (*.example.com)
			for d, cert := range certs {
				if strings.HasPrefix(d, "*.") && strings.HasSuffix(domain, d[1:]) {
					return cert, nil
				}
			}

			// Fallback to default certificate
			if defaultCert != nil {
				return defaultCert, nil
			}

			return nil, fmt.Errorf("no certificate configured for domain: %s", domain)
		},

		MinVersion: tls.VersionTLS13,
		// You can uncomment and customize if needed:
		// CipherSuites: []uint16{
		// 	tls.TLS_CHACHA20_POLY1305_SHA256,
		// 	tls.TLS_AES_256_GCM_SHA384,
		// 	tls.TLS_AES_128_GCM_SHA256,
		// },
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		ErrorLog:     stdlog.New(logFile, logLineStart, stdlog.LstdFlags),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 20 * time.Second,
		IdleTimeout:  1 * time.Second,
		TLSConfig:    tlsConfig,
		Handler:      mux,
	}

	// We pass empty strings because we manage certificates via TLSConfig.GetCertificate
	errs <- srv.ListenAndServeTLS("", "")
}

func Serve(conf *Config, logConf *logging.Config, ws *weather.WeatherService) error {
	var (
		mux = http.NewServeMux()

		logger = logging.NewRequestLogger(
			logConf.AccessLog,
			time.Duration(logConf.Interval)*time.Second,
		)

		errs = make(chan error, 1)

		numberOfServers int

		errorsLog = logging.NewLogSuppressor(
			logConf.ErrorsLog,
			suppressMessages(),
			logLineStart,
		)
	)

	if err := errorsLog.Open(); err != nil {
		return err
	}

	// Register handlers with panic recovery
	mux.HandleFunc("/", panicRecovery(mainHandler(ws, logger)))
	mux.HandleFunc("/favicon.ico", panicRecovery(faviconHandler))
	mux.HandleFunc("/files/", panicRecovery(staticFilesHandler))

	if conf.PortHTTP != 0 {
		go serveHTTP(mux, conf.PortHTTP, errorsLog, errs)
		numberOfServers++
	}

	if conf.PortHTTPS != 0 {
		go serveHTTPS(mux, conf.PortHTTPS, conf, errorsLog, errs)
		numberOfServers++
	}

	if numberOfServers == 0 {
		return ErrNoServersConfigured
	}

	return <-errs // block until one server fails
}

func mainHandler(ws *weather.WeatherService, logger *logging.RequestLogger) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := logger.Log(r); err != nil {
			log.Println(err)
		}
		ws.WeatherHandler(w, r)
	}
}
