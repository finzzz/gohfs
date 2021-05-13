package handler

import (
	"fmt"
	"log"
	"io"
	"os"
	"net/http"
	"net/url"
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

	if strings.HasPrefix(r.RequestURI, h.Config.SHA1Path) {
		h.sha1Handler(w, r)
		return
	}

	if utils.IsDirPath(string(r.RequestURI)) {
		h.listingHandler(w, r)
		return
	}

	link, err := url.PathUnescape(r.RequestURI)
	if logger.LogErr("Handler", err) {
		return
	}

	http.ServeFile(w, r, h.Config.Dir + link)
}

func (h HandlerObj) uploadHandler(w http.ResponseWriter, r *http.Request){
	if h.Config.DisableUp {
		fmt.Fprintln(w, `<script>alert("Upload is disabled!")</script>`)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if h.Config.MaxUpload != -1 && int(fileHeader.Size) >= h.Config.MaxUpload {
		fmt.Fprintln(w, `<script>alert("Upload rejected, maximum upload size exceeded!")</script>`)
		return
	}

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

	link, err := url.PathUnescape(r.RequestURI)
	if logger.LogErr("Handler", err) {
		return
	}
	dir := h.Config.Dir + link

	files, err := os.ReadDir(dir)
	if logger.LogErr("listingHandler (ReadDir)", err) {
		return
	}

	index, err := template.ParseFS(h.Config.Web, "template/index.html")
	if logger.LogErr("listingHandler (parse template)", err) {
		return
	}

	templ := web.Templ{
		Scheme		: h.Config.Scheme,
        IP      	: utils.GetIP(strings.Split(r.Host,":")[0]),
        Port    	: h.Config.Port,
		WebPath		: h.Config.WebPath,
		ZipPath		: h.Config.ZipPath,
		SHA1Path	: h.Config.SHA1Path,
    }

	if ! h.Config.DisableListing {
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
	if h.Config.DisableZip {
		fmt.Fprintln(w, `<script>alert("Zip is disabled!"); history.back()</script>`)
		return
	}

	link, err := url.PathUnescape(r.RequestURI)
	if logger.LogErr("zipHandler", err) {
		return
	}

	link = h.Config.Dir + strings.TrimPrefix(link, h.Config.ZipPath)

	// handle dir not exist
	if ! utils.IsDirExist(link) {
		return
	}

	// handle dir not containing trailing slash
	link = utils.CleanDirPath(link)

	z := utils.ZipWrite(link, utils.CleanDirPath(h.Config.ZipTemp))

	w.Header().Set("Content-Disposition", "attachment; filename=" + utils.Basename(link) + ".zip")
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

func (h HandlerObj) sha1Handler(w http.ResponseWriter, r *http.Request) {
	link, err := url.PathUnescape(r.RequestURI)
	if logger.LogErr("sha1Handler", err) {
		return
	}

	sha1 := utils.SHA1(h.Config.Dir + strings.TrimPrefix(link, h.Config.SHA1Path))
	fmt.Fprintln(w, `<script>alert("` + sha1 + `"); history.back()</script>`)
}