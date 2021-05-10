package handler

import (
	"fmt"
	"log"
	"io"
	"os"
	"net/http"
	"html/template"
	"strings"
	"strconv"

	"gohfs/web"
	"gohfs/internal/utils"
	"gohfs/internal/logger"
)

func (h HandlerObj) Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" && utils.IsDirPath(string(r.RequestURI)){
		h.uploadHandler(w, r)
	}

	if strings.HasPrefix(r.RequestURI, h.Config.WebPath) {
		h.staticHandler(w, r)
		return
	}

	log.Printf("From: %s - %s %s", r.RemoteAddr, r.Method, r.URL)

	if strings.HasPrefix(r.RequestURI, h.Config.ZipPath) {
		h.zipHandler(w, r)
		return
	}

	if utils.IsDirPath(string(r.RequestURI)) {
		h.listingHandler(w, r)
		return
	}

	http.ServeFile(w, r, h.Config.Dir + r.RequestURI)	
}

func (h HandlerObj) uploadHandler(w http.ResponseWriter, r *http.Request){
	file, fileHeader, err := r.FormFile("file")
	if logger.LogErr("uploadHandler", err) {
		return
	}
	defer file.Close()

	fsize, fbytes := utils.ParseSize(fileHeader.Size)
	log.Printf("From: %s - %s %s  filename: %s  size: %g %s", r.RemoteAddr, r.Method, r.URL, fileHeader.Filename, fsize, fbytes)
	if logger.LogErr("uploadHandler", err) {
		return
	}
	
	fileBytes, err := io.ReadAll(file) // read content
	err = os.WriteFile( h.Config.Dir + r.RequestURI + fileHeader.Filename, fileBytes, 0644) // write to file
	if logger.LogErr("uploadHandler", err) {
		return
	}

	fmt.Fprintln(w, `<script>alert("Upload Success!")</script>`)
}

func (h HandlerObj) listingHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	dir := h.Config.Dir + r.RequestURI

	files, err := os.ReadDir(dir)
	if logger.LogErr("listingHandler (ReadDir)", err) {
		return
	}

	index, err := template.ParseFS(h.Config.Web, "template/index.html")
	if logger.LogErr("listingHandler (parse template)", err) {
		return
	}

	templ := web.Templ{
		Scheme	: h.Config.Scheme,
        IP      : utils.GetIP(strings.Split(r.Host,":")[0]),
        Port    : h.Config.Port,
		WebPath	: h.Config.WebPath,
		ZipPath	: h.Config.ZipPath,
    }

	if ! h.Config.Hide {
		templ.Dir = dir[:len(dir)-1]
		templ.NItems = strconv.Itoa(len(files))

		for _, file := range files {
			info,_ := file.Info()
			templ.Items = append(templ.Items, web.ParseItem(info))
		}
	}

    err = index.Execute(w, templ)
	if logger.LogErr("listingHandler (execute template)", err) {
		return
	}
}

func (h HandlerObj) zipHandler(w http.ResponseWriter, r *http.Request) {
	z := utils.ZipWrite(h.Config.Dir + strings.TrimPrefix(r.RequestURI, h.Config.ZipPath))

	w.Header().Set("Content-Disposition", "attachment; filename=" + utils.Basename(r.RequestURI) + ".zip")
	http.ServeFile(w, r, z)
	_ = os.Remove(z)

	return
}

func (h HandlerObj) staticHandler(w http.ResponseWriter, r *http.Request) {
	fname := strings.Trim(strings.TrimPrefix(r.RequestURI, h.Config.WebPath),"/")

	if strings.HasSuffix(fname, ".css") {
		w.Header().Add("content-type","text/css; charset=utf-8")
	} else if strings.HasSuffix(fname, ".js") {
		w.Header().Add("content-type","text/javascript; charset=utf-8")
	} else if strings.HasSuffix(fname, ".svg") {
		w.Header().Add("content-type","image/svg+xml; charset=utf-8")
	}

	data, err := h.Config.Web.ReadFile(fname)
	if logger.LogErr("staticHandler", err) {
		return
	}

	fmt.Fprintln(w, string(data))
}