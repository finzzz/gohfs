package main

import(
	"fmt"
	"log"
	"net"
	"net/http"
	"io/ioutil"
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
var host string
var port string
var current_dir string

type Templ struct {
	Dir 		string
	IP			string
	Port		string
	Items		[]Item
	NItems		string
}

type Item struct {
	Name		string
	Type		string
	Size		string
	RawSize		string
	ModTime		string
	RawModTime	string
}

func main(){
	flag.StringVar(&host, "host", "0.0.0.0", "Host")
	flag.StringVar(&port, "port", "8080", "Port")
	flag.StringVar(&current_dir, "dir", ".", "Directory to serve")
	flag.Parse()

	http_handler := http.HandlerFunc(catch)
    http_server := &http.Server{
            Addr:           host + ":" + port,
            Handler:        http_handler,
	}
	
	fmt.Printf("Serving HTTP on %s port %s (http://%s:%s/) ...\n", host, port, host, port)
	log.Fatal(http_server.ListenAndServe()) 
}

func catch(w http.ResponseWriter, req *http.Request){
	path := req.URL.Path

	if req.Method != "POST" {
		log.Printf("From: %s - %s %s", req.RemoteAddr, req.Method, path[1:])
	}else{
		file, fileHeader, err := req.FormFile("file")
		fsize, fbytes := parseSize(fileHeader.Size)
		log.Printf("From: %s - %s %s  filename: %s  size: %g %s", req.RemoteAddr, req.Method, req.URL, fileHeader.Filename, fsize, fbytes)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()

		// read content
		fileBytes, _ := ioutil.ReadAll(file)

		// write to file
		err = ioutil.WriteFile( current_dir + path + fileHeader.Filename, fileBytes, 0644)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Fprintf(w, "<script>alert(\"Upload Success!\")</script>")
	}

	if string(path[len(path)-1]) == "/"{
		w.Header().Set("Content-Type", "text/html; charset=UTF-8")
		ls(w, req, current_dir + path)
	}else {
		http.ServeFile(w, req, current_dir + path)
	}
	
}

func ls(w http.ResponseWriter, req *http.Request, dir string){
    // check permission
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(w, "<h1>permission denied</h1>")
		log.Printf("permission denied : %s", dir)
		return
	}

    // get host ip
	var ip string
	if strings.Split(req.Host,":")[0] != "0.0.0.0" {
		ip = strings.Split(req.Host,":")[0]
	} else {
		ip = getIP()
	}

	index, _ := template.ParseFS(static, "static/*")

	templ := Templ{
        Dir	    : dir[:len(dir)-1],
        IP      : ip,
        Port    : port,
        NItems  : strconv.Itoa(len(files)),
    }

    for _, file := range files {
        tmp := Item{
            Name: file.Name(),
            ModTime: file.ModTime().Format(time.RFC1123),
            RawModTime: file.ModTime().Format(time.RFC3339),
        }

		if file.IsDir() {
            tmp.Type = "Directory"
            tmp.Size = "--"
            tmp.RawSize = "-1"
		} else {
            fsize, suffix := parseSize(file.Size())

            tmp.Type = "File"
            tmp.Size = strconv.FormatFloat(fsize, 'f', 1, 64) + " " + suffix
            tmp.RawSize = strconv.FormatInt(file.Size(), 10)
        }

        templ.Items = append(templ.Items, tmp)
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

func getIP() (string) {
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
