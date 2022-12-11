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
	geoip "github.com/chubin/wttr.in/internal/geo/ip"
	geoloc "github.com/chubin/wttr.in/internal/geo/location"
	"github.com/chubin/wttr.in/internal/logging"
	"github.com/chubin/wttr.in/internal/processor"
)

var cli struct {
	ConfigCheck bool   `name:"config-check" help:"Check configuration"`
	ConfigDump  bool   `name:"config-dump" help:"Dump configuration"`
	GeoResolve  string `name:"geo-resolve" help:"Resolve location"`

	ConfigFile string `name:"config-file" arg:"" optional:"" help:"Name of configuration file"`

	ConvertGeoIPCache       bool `name:"convert-geo-ip-cache" help:"Convert Geo IP data cache to SQlite"`
	ConvertGeoLocationCache bool `name:"convert-geo-location-cache" help:"Convert Geo Location data cache to SQlite"`
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
		ErrorLog:     log.New(logFile, logLineStart, log.LstdFlags),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 20 * time.Second,
		IdleTimeout:  1 * time.Second,
		TLSConfig:    tlsConfig,
		Handler:      mux,
	}
	errs <- srv.ListenAndServeTLS(certFile, keyFile)
}

func serve(conf *config.Config) error {
	var (
		// mux is main HTTP/HTTP requests multiplexer.
		mux *http.ServeMux = http.NewServeMux()

		// logger is optimized requests logger.
		logger *logging.RequestLogger

		rp *processor.RequestProcessor

		// errs is the servers errors channel.
		errs chan error = make(chan error, 1)

		// numberOfServers started. If 0, exit.
		numberOfServers int

		errorsLog *logging.LogSuppressor

		err error
	)

	// logger is optimized requests logger.
	logger = logging.NewRequestLogger(
		conf.Logging.AccessLog,
		time.Duration(conf.Logging.Interval)*time.Second)

	errorsLog = logging.NewLogSuppressor(
		conf.Logging.ErrorsLog,
		[]string{
			"error reading preface from client",
			"TLS handshake error from",
		},
		logLineStart,
	)

	rp, err = processor.NewRequestProcessor(conf)
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

	if conf.Server.PortHTTP != 0 {
		go serveHTTP(mux, conf.Server.PortHTTP, errorsLog, errs)
		numberOfServers++
	}
	if conf.Server.PortHTTPS != 0 {
		go serveHTTPS(mux, conf.Server.PortHTTPS, conf.Server.TLSCertFile, conf.Server.TLSKeyFile, errorsLog, errs)
		numberOfServers++
	}
	if numberOfServers == 0 {
		return errors.New("no servers configured")
	}
	return <-errs // block until one of the servers writes an error
}

func main() {
	var (
		conf *config.Config
		err  error
	)

	ctx := kong.Parse(&cli)

	if cli.ConfigFile != "" {
		conf, err = config.Load(cli.ConfigFile)
		if err != nil {
			log.Fatalf("reading config from %s: %s\n", cli.ConfigFile, err)
		}
	} else {
		conf = config.Default()
	}

	if cli.ConfigDump {
		fmt.Print(string(conf.Dump()))
	}

	if cli.ConfigCheck || cli.ConfigDump {
		return
	}

	if cli.ConvertGeoIPCache {
		geoIPCache, err := geoip.NewCache(conf)
		if err != nil {
			ctx.FatalIfErrorf(err)
		}
		ctx.FatalIfErrorf(geoIPCache.ConvertCache())
		return
	}

	if cli.ConvertGeoLocationCache {
		geoLocCache, err := geoloc.NewCache(conf)
		if err != nil {
			ctx.FatalIfErrorf(err)
		}
		ctx.FatalIfErrorf(geoLocCache.ConvertCache())
		return
	}

	if cli.GeoResolve != "" {
		sr := geoloc.NewSearcher(conf)
		loc, err := sr.Search(cli.GeoResolve)
		ctx.FatalIfErrorf(err)
		if loc != nil {
			fmt.Println(*loc)

		}
	}

	err = serve(conf)
	ctx.FatalIfErrorf(err)
}
