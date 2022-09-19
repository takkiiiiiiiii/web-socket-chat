package main

import (
	// "bytes"
	"io"
	"io/ioutil"
	// "mime/multipart"
	"net/http"
	"path"
)

// uploaderHandler expects two fields to be posted, userid and avatarFile.
func uploaderHandler(w http.ResponseWriter, req *http.Request) {
	// body := &bytes.Buffer{}
	// mw := multipart.NewWriter(body)
	// req.Header.Add("Content-Type", mw.FormDataContentType())
	userID := req.FormValue("userid")
	file, header, err := req.FormFile("avatarFile")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	filename := path.Join("avatars", userID+path.Ext(header.Filename))
	err = ioutil.WriteFile(filename, data, 0777)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, "Successful")
}
