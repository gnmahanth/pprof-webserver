package main

import (
	_ "embed"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"text/template"
	"time"

	"github.com/google/pprof/driver"
	"github.com/google/pprof/profile"
	log "github.com/sirupsen/logrus"
)

//go:embed templates/index.html
var index string

var indexTemplate = template.Must(template.New("index").Parse(index))

type fileFetcher struct {
	directory string
	id        string
}

func (f fileFetcher) Fetch(src string, duration, timeout time.Duration) (*profile.Profile, string, error) {
	var p *profile.Profile
	var err error

	fin := path.Join(f.directory, f.id)
	in, err := os.Open(fin)
	if err != nil {
		log.Errorf("error opening file %s, %s", fin, err)
		return p, "", err
	}
	defer in.Close()

	p, err = profile.Parse(in)
	if err != nil {
		log.Errorf("error parsing file %s, %s", fin, err)
		return p, "", err
	}
	return p, "", nil
}

type fileServer struct {
	directory string
}

func listfiles(d string) []string {
	if _, err := os.Stat(d); err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(d, 0775); err != nil {
				log.Error(err)
			}
		} else {
			log.Error(err)
		}
	}
	l := []string{}
	files, err := os.ReadDir(d)
	if err != nil {
		log.Error(err)
		return l
	}
	for _, f := range files {
		l = append(l, f.Name())
	}
	return l
}

func parsePath(reqPath string, trimPath string) (string, string) {
	parts := strings.Split(path.Clean(strings.TrimPrefix(reqPath, trimPath)), "/")
	if len(parts) < 1 {
		return "", ""
	} else if len(parts) == 1 {
		return parts[0], "/"
	}
	return parts[0], "/" + strings.Join(parts[1:], "/")
}

func (s *fileServer) handleFile(w http.ResponseWriter, r *http.Request) {
	file, rPath := parsePath(r.URL.Path, "/pprof/")
	log.Debugf("req path: %s file: %s path: %s", r.URL.Path, file, rPath)

	flagset := &pprofFlags{
		FlagSet: flag.NewFlagSet("pprof", flag.ContinueOnError),
		args: []string{
			"--symbolize", "none",
			"--http", "localhost:0",
			"",
		},
	}

	fetcher := &fileFetcher{
		directory: s.directory,
		id:        file,
	}

	server := func(args *driver.HTTPServerArgs) error {
		handler, ok := args.Handlers[rPath]
		if !ok {
			return fmt.Errorf("unknown endpoint %s\n", rPath)
		}
		handler.ServeHTTP(w, r)
		return nil
	}

	opt := &driver.Options{
		Flagset:    flagset,
		HTTPServer: server,
		Fetch:      fetcher,
		UI:         &noUI{},
	}
	if err := driver.PProf(opt); err != nil {
		log.Error(err)
		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			log.Error(err)
		}

	}
}

func (s *fileServer) handleUpload(w http.ResponseWriter, r *http.Request) {
	// Maximum upload of 10 MB files
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Error(err)
		return
	}

	file, handler, err := r.FormFile("inputfile")
	if err != nil {
		log.Errorf("error reading file from %s", err)
		return
	}

	defer file.Close()

	log.Infof("uploaded file: %s size: %dkb", handler.Filename, handler.Size/1024)

	savefile := path.Join(s.directory, handler.Filename)
	dst, err := os.Create(savefile)
	if err != nil {
		log.Errorf("error creating save file %s, %s", savefile, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		log.Errorf("error copying data to %s, %s", savefile, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func (s *fileServer) handleRemove(w http.ResponseWriter, r *http.Request) {
	file, _ := parsePath(r.URL.Path, "/remove/")
	// Validate: disallow path separators or ".." in file name
	if strings.Contains(file, "/") || strings.Contains(file, "\\") || strings.Contains(file, "..") {
		http.Error(w, "Invalid file name", http.StatusBadRequest)
		log.Warnf("attempted deletion with invalid filename: %q", file)
		return
	}
	removefile := path.Join(s.directory, file)
	log.Infof("removing file %s", removefile)
	if err := os.Remove(removefile); err != nil {
		log.Errorf("error removing file %s", err)
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func (s *fileServer) handleIndex(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	err := indexTemplate.Execute(w, struct{ Files []string }{Files: listfiles(s.directory)})
	if err != nil {
		log.Error(err)
	}
}
