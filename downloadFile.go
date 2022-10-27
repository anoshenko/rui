package rui

import (
	"bytes"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type downloadFile struct {
	filename string
	path     string
	data     []byte
}

var currentDownloadId = int(rand.Int31())
var downloadFiles = map[string]downloadFile{}

func (session *sessionData) startDownload(file downloadFile) {
	currentDownloadId++
	id := strconv.Itoa(currentDownloadId)
	downloadFiles[id] = file
	session.runScript(fmt.Sprintf(`startDowndload("%s", "%s")`, id, file.filename))
}

func serveDownloadFile(id string, w http.ResponseWriter, r *http.Request) bool {
	if file, ok := downloadFiles[id]; ok {
		delete(downloadFiles, id)
		if file.data != nil {
			http.ServeContent(w, r, file.filename, time.Now(), bytes.NewReader(file.data))
			return true
		} else if _, err := os.Stat(file.path); err == nil {
			http.ServeFile(w, r, file.path)
			return true
		}
	}
	return false
}

// DownloadFile starts downloading the file on the client side.
func (session *sessionData) DownloadFile(path string) {
	if _, err := os.Stat(path); err != nil {
		ErrorLog(err.Error())
		return
	}

	_, filename := filepath.Split(path)
	session.startDownload(downloadFile{
		filename: filename,
		path:     path,
		data:     nil,
	})
}

// DownloadFileData starts downloading the file on the client side. Arguments specify the name of the downloaded file and its contents
func (session *sessionData) DownloadFileData(filename string, data []byte) {
	if data == nil {
		ErrorLog("Invalid download data. Must be not nil.")
		return
	}

	session.startDownload(downloadFile{
		filename: filename,
		path:     "",
		data:     data,
	})
}
