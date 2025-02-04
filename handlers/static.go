package handlers

import (
	"net/http"
	"os"
)

type justFilesFilesystem struct {
	fs http.FileSystem
}

type neuteredReaddirFile struct {
	http.File
}

func (fs justFilesFilesystem) Open(name string) (http.File, error) {
	f, err := fs.fs.Open(name)
	if err != nil {
		return nil, err
	}
	return neuteredReaddirFile{f}, nil
}

func (f neuteredReaddirFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}

func HandleStatic() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			nfs := justFilesFilesystem{http.Dir("static")}
			fs := http.FileServer(nfs)
			fs.ServeHTTP(w, r)
		},
	)
}
