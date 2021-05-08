package main

import(
	"fmt"
	"log"
	"net"
	"net/http"
	"io"
	"io/fs"
	"os"
	"flag"
	"math"
	"strings"
    "strconv"
	"time"
	"embed"
	"html/template"
)

//go:embed static/*
var static embed.FS
var config Config

type Config struct {
	Host			string
	Port			string
	Dir				string
	Hide			bool
}

type Templ struct {
	Dir 			string
	IP				string
	Port			string
	Items			[]Item
	NItems			string
}

type Item struct {
	Name			string
	Type			string
	Size			string
	RawSize			string
	ModTime			string
	RawModTime		string
}

func main(){
	flag.StringVar(&config.Host, "host", "0.0.0.0", "Host")
	flag.StringVar(&config.Port, "port", "8080", "Port")
	flag.StringVar(&config.Dir, "dir", ".", "Directory to serve")
	flag.BoolVar(&config.Hide, "hide", false, "Disable Listing")
	flag.Parse()

	http_handler := http.HandlerFunc(handler)
    http_server := &http.Server{
            Addr:           config.Host + ":" + config.Port,
            Handler:        http_handler,
	}
	
	fmt.Printf("Serving HTTP on %s port %s (http://%s:%s/) ...\n", config.Host, config.Port, config.Host, config.Port)
	log.Fatal(http_server.ListenAndServe()) 
}

func handler(w http.ResponseWriter, r *http.Request){
	if r.Method == "POST" {
		if isDirPath(string(r.RequestURI)) {
			uploadHandler(w, r)
		} else {
			fmt.Println("Invalid")
		}
	} else {
		log.Printf("From: %s - %s %s", r.RemoteAddr, r.Method, r.URL)
	}

	if isDirPath(string(r.RequestURI)) {
		listingHandler(w, r)
		return
	}
	
	http.ServeFile(w, r, config.Dir + r.RequestURI)
}

func uploadHandler(w http.ResponseWriter, r *http.Request){
	file, fileHeader, err := r.FormFile("file")
	fsize, fbytes := parseSize(fileHeader.Size)
	log.Printf("From: %s - %s %s  filename: %s  size: %g %s", r.RemoteAddr, r.Method, r.URL, fileHeader.Filename, fsize, fbytes)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fileBytes, _ := io.ReadAll(file) // read content
	err = os.WriteFile( config.Dir + r.RequestURI + fileHeader.Filename, fileBytes, 0644) // write to file
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintln(w, `<script>alert("Upload Success!")</script>`)
}

func listingHandler(w http.ResponseWriter, r *http.Request){
	dir := config.Dir + r.RequestURI

	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(w, "<h1>permission denied</h1>")
		log.Printf("permission denied : %s", dir)
		return
	}

	index, _ := template.ParseFS(static, "static/*")

	templ := Templ{
        IP      : getIP(strings.Split(r.Host,":")[0]),
        Port    : config.Port,
    }

	if ! config.Hide {
		templ.Dir = dir[:len(dir)-1]
		templ.NItems = strconv.Itoa(len(files))

		for _, file := range files {
			info,_ := file.Info()
			templ.Items = append(templ.Items, parseItem(info))
		}
	}

    err = index.ExecuteTemplate(w, "index", templ)
    if err != nil {
        fmt.Println(err)
    }
}

func parseSize(s int64) (float64, string){
	suffix := []string{"B","KB","MB","GB","TB"}
	var i int

	val := float64(s)
	for i=0;i<len(suffix);i++ {
		if val < 1024 {
			break
		}
		val = val / 1024
	}

	val = math.Round(val*100)/100

	return val, suffix[i]
}

func parseItem(info fs.FileInfo) Item {
	tmp := Item{
		Name: info.Name(),
		ModTime: info.ModTime().Format(time.RFC1123),
		RawModTime: info.ModTime().Format(time.RFC3339),
	}

	if info.IsDir() {
		tmp.Type = "Directory"
		tmp.Size = "--"
		tmp.RawSize = "-1"
	} else {
		fsize, suffix := parseSize(info.Size())

		tmp.Type = "File"
		tmp.Size = strconv.FormatFloat(fsize, 'f', 1, 64) + " " + suffix
		tmp.RawSize = strconv.FormatInt(info.Size(), 10)
	}

	return tmp
}

func getIP(host string) (string) {
	if host != "0.0.0.0" {
		return host
	}

    netInterfaceAddresses, _ := net.InterfaceAddrs()

    for _, netInterfaceAddress := range netInterfaceAddresses {
        networkIp, ok := netInterfaceAddress.(*net.IPNet)
        if ok && !networkIp.IP.IsLoopback() && networkIp.IP.To4() != nil {
            return networkIp.IP.String()
        }
    }

    return "IP"
}

func isDirPath(s string) (bool) {
	if string(s[len(s)-1]) == "/" {
		return true
	}

	return false
}