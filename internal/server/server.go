package server

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
	stdlog "log"
	"net/http"
	"time"

	"github.com/chubin/wttr.go/internal/logging"
	"github.com/chubin/wttr.go/internal/weather"
)

type Config struct {
	PortHTTP    int    `yaml:"portHttp"`
	PortHTTPS   int    `yaml:"portHttps"`
	TLSCertFile string `yaml:"tlsCertFile"`
	TLSKeyFile  string `yaml:"tlsKeyFile"`
}

const logLineStart = "LOG_LINE_START "

var ErrNoServersConfigured = errors.New("no servers configured")

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

func serveHTTPS(mux *http.ServeMux, port int, certFile, keyFile string, logFile io.Writer, errs chan<- error) {
	tlsConfig := &tls.Config{
		// CipherSuites: []uint16{
		// 	tls.TLS_CHACHA20_POLY1305_SHA256,
		// 	tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		// 	tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		// },
		// MinVersion: tls.VersionTLS13,
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
	errs <- srv.ListenAndServeTLS(certFile, keyFile)
}

func Serve(conf *Config, logConf *logging.Config, ws *weather.WeatherService) error {
	var (
		// mux is main HTTP/HTTP requests multiplexer.
		mux = http.NewServeMux()

		// logger is optimized requests logger.
		logger = logging.NewRequestLogger(
			logConf.AccessLog,
			time.Duration(logConf.Interval)*time.Second)

		// errs is the servers errors channel.
		errs = make(chan error, 1)

		// numberOfServers started. If 0, exit.
		numberOfServers int

		errorsLog = logging.NewLogSuppressor(
			logConf.ErrorsLog,
			suppressMessages(),
			logLineStart,
		)

		err error
	)

	err = errorsLog.Open()
	if err != nil {
		return err
	}

	mux.HandleFunc("/", mainHandler(ws, logger))

	if conf.PortHTTP != 0 {
		go serveHTTP(mux, conf.PortHTTP, errorsLog, errs)
		numberOfServers++
	}
	if conf.PortHTTPS != 0 {
		go serveHTTPS(mux, conf.PortHTTPS, conf.TLSCertFile, conf.TLSKeyFile, errorsLog, errs)
		numberOfServers++
	}
	if numberOfServers == 0 {
		return ErrNoServersConfigured
	}

	return <-errs // block until one of the servers writes an error
}

func mainHandler(
	ws *weather.WeatherService,
	logger *logging.RequestLogger,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := logger.Log(r); err != nil {
			log.Println(err)
		}
		ws.WeatherHandler(w, r)

		// response, err := rp.ProcessRequest(r)
		// if err != nil {
		// 	log.Println(err)

		// 	return
		// }
		// if response.StatusCode == 0 {
		// 	log.Println("status code 0", response)

		// 	return
		// }

		// copyHeader(w.Header(), response.Header)
		// w.Header().Set("Access-Control-Allow-Origin", "*")
		// w.WriteHeader(response.StatusCode)
		// _, err = w.Write(response.Body)
		// if err != nil {
		// 	log.Println(err)
		// }
	}
}
