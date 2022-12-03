package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/alecthomas/kong"

	"github.com/chubin/wttr.in/internal/config"
	"github.com/chubin/wttr.in/internal/logging"
	"github.com/chubin/wttr.in/internal/processor"
)

var cli struct {
	ConfigCheck bool   `name:"config-check" help:"Check configuration"`
	ConfigDump  bool   `name:"config-dump" help:"Dump configuration"`
	ConfigFile  string `name:"config-file" arg:"" optional:"" help:"Name of configuration file"`
}

const logLineStart = "LOG_LINE_START "

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
		ErrorLog:     log.New(logFile, logLineStart, log.LstdFlags),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  1 * time.Second,
		Handler:      mux,
	}
	errs <- srv.ListenAndServe()
}

func serveHTTPS(mux *http.ServeMux, port int, logFile io.Writer, errs chan<- error) {
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
		ErrorLog:     log.New(logFile, logLineStart, log.LstdFlags),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 20 * time.Second,
		IdleTimeout:  1 * time.Second,
		TLSConfig:    tlsConfig,
		Handler:      mux,
	}
	errs <- srv.ListenAndServeTLS(config.Conf.Server.TLSCertFile, config.Conf.Server.TLSKeyFile)
}

func main() {
	var (
		// mux is main HTTP/HTTP requests multiplexer.
		mux *http.ServeMux = http.NewServeMux()

		// logger is optimized requests logger.
		logger *logging.RequestLogger = logging.NewRequestLogger(
			config.Conf.Logging.AccessLog,
			time.Duration(config.Conf.Logging.Interval)*time.Second)

		rp *processor.RequestProcessor

		// errs is the servers errors channel.
		errs chan error = make(chan error, 1)

		// numberOfServers started. If 0, exit.
		numberOfServers int

		errorsLog *logging.LogSuppressor = logging.NewLogSuppressor(
			config.Conf.Logging.ErrorsLog,
			[]string{
				"error reading preface from client",
				"TLS handshake error from",
			},
			logLineStart,
		)

		err error
	)

	rp, err = processor.NewRequestProcessor(config.Conf)
	if err != nil {
		log.Fatalln("log processor initialization:", err)
	}

	err = errorsLog.Open()
	if err != nil {
		log.Fatalln("errors log:", err)
	}

	rp.Start()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := logger.Log(r); err != nil {
			log.Println(err)
		}
		// printStat()
		response, err := rp.ProcessRequest(r)
		if err != nil {
			log.Println(err)
			return
		}
		if response.StatusCode == 0 {
			log.Println("status code 0", response)
			return
		}

		copyHeader(w.Header(), response.Header)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(response.StatusCode)
		w.Write(response.Body)
	})

	if config.Conf.Server.PortHTTP != 0 {
		go serveHTTP(mux, config.Conf.Server.PortHTTP, errorsLog, errs)
		numberOfServers++
	}
	if config.Conf.Server.PortHTTPS != 0 {
		go serveHTTPS(mux, config.Conf.Server.PortHTTPS, errorsLog, errs)
		numberOfServers++
	}
	if numberOfServers == 0 {
		log.Println("no servers configured; exiting")
		return
	}
	log.Fatal(<-errs) // block until one of the servers writes an error
}
