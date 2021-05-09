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
)

func (h HandlerObj) Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" && utils.IsDirPath(string(r.RequestURI)){
		h.uploadHandler(w, r)
	} else {
		log.Printf("From: %s - %s %s", r.RemoteAddr, r.Method, r.URL)
	}

	if strings.HasPrefix(r.RequestURI, h.Config.ZipPath) {
		z := utils.ZipWrite(h.Config.Dir + strings.TrimPrefix(r.RequestURI, h.Config.ZipPath))

		w.Header().Set("Content-Disposition", "attachment; filename=" + utils.Basename(r.RequestURI) + ".zip")
		http.ServeFile(w, r, z)
		_ = os.Remove(z)
		return
	}

	if utils.IsDirPath(string(r.RequestURI)) {
		h.listingHandler(w, r)
		return
	}

	if strings.HasPrefix(r.RequestURI, h.Config.WebPath) {
		fname := strings.Trim(strings.TrimPrefix(r.RequestURI, h.Config.WebPath),"/")

		if strings.HasSuffix(fname, ".css") {
			w.Header().Add("content-type","text/css; charset=utf-8")
		} else if strings.HasSuffix(fname, ".js") {
			w.Header().Add("content-type","text/javascript; charset=utf-8")
		}

		data, err := h.Config.Web.ReadFile(fname)
		if err != nil {
			log.Println(err)
			return
		}

		fmt.Fprintln(w, string(data))

		return
	}

	http.ServeFile(w, r, h.Config.Dir + r.RequestURI)	
}

func (h HandlerObj) uploadHandler(w http.ResponseWriter, r *http.Request){
	file, fileHeader, err := r.FormFile("file")
	fsize, fbytes := utils.ParseSize(fileHeader.Size)
	log.Printf("From: %s - %s %s  filename: %s  size: %g %s", r.RemoteAddr, r.Method, r.URL, fileHeader.Filename, fsize, fbytes)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fileBytes, _ := io.ReadAll(file) // read content
	err = os.WriteFile( h.Config.Dir + r.RequestURI + fileHeader.Filename, fileBytes, 0644) // write to file
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintln(w, `<script>alert("Upload Success!")</script>`)
}

func (h HandlerObj) listingHandler(w http.ResponseWriter, r *http.Request){
	dir := h.Config.Dir + r.RequestURI

	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(w, "<h1>permission denied</h1>")
		log.Printf("permission denied : %s", dir)
		return
	}

	index, err := template.ParseFS(h.Config.Web, "template/index.html")
	if err != nil {
		log.Println(err)
	}

	templ := web.Templ{
        IP      : utils.GetIP(strings.Split(r.Host,":")[0]),
        Port    : h.Config.Port,
		WebPath	: h.Config.WebPath,
    }

	if ! h.Config.Hide {
		templ.Dir = dir[:len(dir)-1]
		templ.NItems = strconv.Itoa(len(files))

		for _, file := range files {
			info,_ := file.Info()
			templ.Items = append(templ.Items, web.ParseItem(info, h.Config.ZipPath))
		}
	}

    err = index.Execute(w, templ)
    if err != nil {
        log.Fatal(err)
    }
}