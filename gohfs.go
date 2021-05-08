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
	"math/rand"
	"strings"
    "strconv"
	"crypto/sha256"
	"time"
	"embed"
	"html/template"
	"archive/zip"
)

//go:embed static/*
var static embed.FS
var config Config

type Config struct {
	Host			string
	Port			string
	Dir				string
	User			string
	Pass			string
	HashedPass		string
	ZipPath			string
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
	ZipURL			string
	RawSize			string
	ModTime			string
	RawModTime		string
}

func main(){
	flag.StringVar(&config.Host, "host", "0.0.0.0", "Host")
	flag.StringVar(&config.Port, "port", "8080", "Port")
	flag.StringVar(&config.User, "user", "admin", "Username")
	flag.StringVar(&config.Pass, "pass", "", "Password")
	flag.StringVar(&config.HashedPass, "hpass", "", "Hashed Password (sha-256)")
	flag.StringVar(&config.ZipPath, "zippath", "/gohfs-zip", "Directory to serve")
	flag.StringVar(&config.Dir, "dir", ".", "Directory to serve")
	flag.BoolVar(&config.Hide, "hide", false, "Disable Listing")
	flag.Parse()

	if config.Pass != "" && config.HashedPass != "" {
		log.Fatal(`Can only define either "Password" or "Hashed Password"`)
	}

	http_handler := http.HandlerFunc(authHandler(handler))
    http_server := &http.Server{
            Addr:           config.Host + ":" + config.Port,
            Handler:        http_handler,
	}
	
	fmt.Printf("Serving HTTP on %s port %s (http://%s:%s/) ...\n", config.Host, config.Port, config.Host, config.Port)
	log.Fatal(http_server.ListenAndServe()) 
}

func authHandler(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, _ := r.BasicAuth()
		if !checkAuth(user, pass) {
			w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}

		f(w, r)
	}
}

func checkAuth(user, pass string) (bool) {
	var p, ipass string

	if config.Pass == "" && config.HashedPass == "" {
		return true // doesn't have auth
	}

	if config.Pass != "" {
		p = config.Pass
		ipass = pass
	} else {
		p = config.HashedPass
		ipass = fmt.Sprintf("%x", sha256.Sum256([]byte(pass)))
	}

	if strings.Compare(config.User, user) == 0 && strings.Compare(p, ipass) == 0 {
		return true
	}

	return false
}

func handler(w http.ResponseWriter, r *http.Request){
	if r.Method == "POST" && isDirPath(string(r.RequestURI)){
		uploadHandler(w, r)
	} else {
		log.Printf("From: %s - %s %s", r.RemoteAddr, r.Method, r.URL)
	}

	if strings.HasPrefix(r.RequestURI, config.ZipPath) {
		z := zipWrite(config.Dir + strings.TrimPrefix(r.RequestURI, config.ZipPath))

		w.Header().Set("Content-Disposition", "attachment; filename=" + basename(r.RequestURI) + ".zip")
		http.ServeFile(w, r, z)
		_ = os.Remove(z)
		return
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
		ZipURL: config.ZipPath + "/" + info.Name(),
		ModTime: info.ModTime().Format(time.RFC1123),
		RawModTime: info.ModTime().Format(time.RFC3339),
	}

	if info.IsDir() {
		tmp.Type = "Directory"
		tmp.Size = "--"
		tmp.RawSize = "-1"
		tmp.ZipURL += "/"
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

func randStr(n int) string {
	rand.Seed(time.Now().UnixNano())

	letters := []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
    b := make([]rune, n)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}

func zipWrite(s string) (string) {
	tmpZip := randStr(12) + ".zip"

    outFile, err := os.Create(tmpZip)
    if err != nil {
        log.Println(err)
    }
    defer outFile.Close()

    w := zip.NewWriter(outFile)

	if isDirPath(s) {
		zipAdd(w, s, "")
	} else {
		zipFile(w, s)
	}

    if err != nil {
        log.Println(err)
    }

    err = w.Close()
    if err != nil {
        log.Println(err)
    }

	return tmpZip
}

func zipAdd(w *zip.Writer, basePath, baseInZip string) {
    files, err := os.ReadDir(basePath)
    if err != nil {
        log.Println(err)
    }

    for _, file := range files {
        if !file.IsDir() {
            zipFile(w, basePath + file.Name())
        } else if file.IsDir(){
            newBase := basePath + file.Name() + "/"
            zipAdd(w, newBase, baseInZip  + file.Name() + "/")
        }
    }
}

func zipFile(w *zip.Writer, path string) {
	dat, err := os.ReadFile(path)
	if err != nil {
		log.Println(err)
	}

	f, err := w.Create(path)
	if err != nil {
		log.Println(err)
	}

	_, err = f.Write(dat)
	if err != nil {
		log.Println(err)
	}
}

func basename(s string) (string){
	splits := strings.Split(s, "/")

	if isDirPath(s) {
		return splits[len(splits)-2]
	}
	
	return splits[len(splits)-1]
}