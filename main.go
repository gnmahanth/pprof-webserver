package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"path"
	"runtime"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetOutput(os.Stdout)
	log.SetReportCaller(true)
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.TextFormatter{
		DisableLevelTruncation: true,
		PadLevelText:           true,
		FullTimestamp:          true,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			// return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", path.Base(f.File), f.Line)
			return "", fmt.Sprintf(" %s:%d", path.Base(f.File), f.Line)
		},
	})

	// command line arguments
	storePath := flag.String(
		"storage",
		"data",
		"path to directory containing profile files",
	)
	listenPort := flag.String(
		"port",
		"8080",
		"server listen port",
	)
	debugMode := flag.Bool(
		"debug",
		false,
		"enable debug logs",
	)

	flag.Parse()

	if *debugMode {
		log.SetLevel(log.DebugLevel)
	}

	mux := http.NewServeMux()

	// local file storage handlers
	storage := &fileServer{directory: *storePath}
	mux.HandleFunc("/", storage.handleIndex)
	mux.HandleFunc("/pprof/", storage.handleFile)
	mux.HandleFunc("/upload", storage.handleUpload)
	mux.HandleFunc("/remove/", storage.handleRemove)

	log.Infof("starting server at port %s", *listenPort)
	if err := http.ListenAndServe(net.JoinHostPort("0.0.0.0", *listenPort), mux); err != nil {
		log.Fatal(err)
	}
}
