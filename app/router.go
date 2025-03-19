package main

import (
	"io/fs"
	"net/http"
	"path/filepath"
	"runtime"

	"app/api"
	appView "app/view"

	"github.com/go-chi/chi/v5"
)

type fileSystemWrapper struct {
	fs http.FileSystem
}

func (fswrapper fileSystemWrapper) Open(path string) (http.File, error) {
	f, err := fswrapper.fs.Open(path)
	if err != nil {
		return nil, err
	}
	s, err := f.Stat()
	if err != nil {
		defer f.Close()
		return nil, err
	}
	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		f2, err := fswrapper.fs.Open(index)
		if err == nil {
			defer f2.Close()
		} else {
			defer f.Close()
			err := fs.ErrNotExist
			return nil, err
		}
	}
	return f, nil
}

func Route(r *chi.Mux) {
	//static
	_, filename, _, _ := runtime.Caller(0)
	filename = filepath.Join(filepath.Dir(filename), "/static")
	fswrapper := fileSystemWrapper{http.Dir(filename)}
	fs := http.FileServer(fswrapper)
	r.Handle("/*", fs)

	//error
	errorRouter := chi.NewRouter()
	errorRouter.Get("/500", appView.ServerError)
	r.Mount("/error/", errorRouter)

	//api
	//api := &ApiInjector{}
	apiRouter := chi.NewRouter()
	apiRouter.Get("/", api.Index)
	apiRouter.Get("/data", api.Data)
	r.Mount("/api/", apiRouter)

}
