package main

import "net/http"

type RedirectingFileSystem struct {
http.FileSystem
}

func (fs RedirectingFileSystem) Open(name string) (http.File, error) {
file, err := fs.FileSystem.Open(name)
if err != nil {
file, err = fs.FileSystem.Open("index.html")
}

if err != nil {
return nil, err
}
return file, nil
}