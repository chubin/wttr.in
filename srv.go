package main

import (
	"crypto/tls"
	"fmt"
	"io"
	stdlog "log"
	"net/http"
	"time"

	"github.com/alecthomas/kong"
	log "github.com/sirupsen/logrus"

	"github.com/chubin/wttr.in/internal/config"
	geoip "github.com/chubin/wttr.in/internal/geo/ip"
	geoloc "github.com/chubin/wttr.in/internal/geo/location"
	"github.com/chubin/wttr.in/internal/logging"
	"github.com/chubin/wttr.in/internal/processor"
	"github.com/chubin/wttr.in/internal/types"
	"github.com/chubin/wttr.in/internal/view/v1"
)

//nolint:gochecknoglobals
var cli struct {
	ConfigFile string `name:"config-file" arg:"" optional:"" help:"Name of configuration file"`

	ConfigCheck             bool   `name:"config-check" help:"Check configuration"`
	ConfigDump              bool   `name:"config-dump" help:"Dump configuration"`
	ConvertGeoIPCache       bool   `name:"convert-geo-ip-cache" help:"Convert Geo IP data cache to SQlite"`
	ConvertGeoLocationCache bool   `name:"convert-geo-location-cache" help:"Convert Geo Location data cache to SQlite"`
	GeoResolve              string `name:"geo-resolve" help:"Resolve location"`
	LogLevel                string `name:"log-level" short:"l" help:"Show log messages with level" default:"info"`

	V1 struct v1.Configuration
}

const logLineStart = "LOG_LINE_START "

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

func serve(conf *config.Config) error {
	var (
		// mux is main HTTP/HTTP requests multiplexer.
		mux = http.NewServeMux()

		// logger is optimized requests logger.
		logger = logging.NewRequestLogger(
			conf.Logging.AccessLog,
			time.Duration(conf.Logging.Interval)*time.Second)

		rp *processor.RequestProcessor

		// errs is the servers errors channel.
		errs = make(chan error, 1)

		// numberOfServers started. If 0, exit.
		numberOfServers int

		errorsLog = logging.NewLogSuppressor(
			conf.Logging.ErrorsLog,
			suppressMessages(),
			logLineStart,
		)

		err error
	)

	rp, err = processor.NewRequestProcessor(conf)
	if err != nil {
		return fmt.Errorf("log processor initialization: %w", err)
	}

	err = errorsLog.Open()
	if err != nil {
		return err
	}

	err = rp.Start()
	if err != nil {
		return err
	}

	mux.HandleFunc("/", mainHandler(rp, logger))

	if conf.Server.PortHTTP != 0 {
		go serveHTTP(mux, conf.Server.PortHTTP, errorsLog, errs)
		numberOfServers++
	}
	if conf.Server.PortHTTPS != 0 {
		go serveHTTPS(mux, conf.Server.PortHTTPS, conf.Server.TLSCertFile, conf.Server.TLSKeyFile, errorsLog, errs)
		numberOfServers++
	}
	if numberOfServers == 0 {
		return types.ErrNoServersConfigured
	}

	return <-errs // block until one of the servers writes an error
}

func mainHandler(
	rp *processor.RequestProcessor,
	logger *logging.RequestLogger,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := logger.Log(r); err != nil {
			log.Println(err)
		}

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
		_, err = w.Write(response.Body)
		if err != nil {
			log.Println(err)
		}
	}
}

func main() {
	var (
		conf *config.Config
		err  error
	)

	ctx := kong.Parse(&cli)
	ctx.FatalIfErrorf(setLogLevel(cli.LogLevel))

	if cli.ConfigFile != "" {
		conf, err = config.Load(cli.ConfigFile)
		if err != nil {
			log.Fatalf("reading config from %s: %s\n", cli.ConfigFile, err)
		}
	} else {
		conf = config.Default()
	}

	if cli.ConfigDump {
		//nolint:forbidigo
		fmt.Print(string(conf.Dump()))

		return
	}

	if cli.ConfigCheck {
		return
	}

	switch {
	case cli.ConvertGeoIPCache:
		ctx.FatalIfErrorf(convertGeoIPCache(conf))
	case cli.ConvertGeoLocationCache:
		ctx.FatalIfErrorf(convertGeoLocationCache(conf))
	case cli.GeoResolve != "":
		sr := geoloc.NewSearcher(conf)
		loc, err := sr.Search(cli.GeoResolve)
		ctx.FatalIfErrorf(err)
		if loc != nil {
			//nolint:forbidigo
			fmt.Println(*loc)
		}
	default:
		err = serve(conf)
		ctx.FatalIfErrorf(err)
	}
}

func convertGeoIPCache(conf *config.Config) error {
	geoIPCache, err := geoip.NewCache(conf)
	if err != nil {
		return err
	}

	return geoIPCache.ConvertCache()
}

func convertGeoLocationCache(conf *config.Config) error {
	geoLocCache, err := geoloc.NewCache(conf)
	if err != nil {
		return err
	}

	return geoLocCache.ConvertCache(false)
}

func setLogLevel(logLevel string) error {
	parsedLevel, err := log.ParseLevel(logLevel)
	if err != nil {
		return err
	}
	log.SetLevel(parsedLevel)

	return nil
}
