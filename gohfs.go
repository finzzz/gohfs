package main

import(
	"fmt"
	"log"
	"net/http"
	"io/ioutil"
	"flag"
	"time"
)

var host string
var port string
var current_dir string

func main(){
	flag.StringVar(&host, "host", "127.0.0.1", "Host")
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

	log.Printf("From: %s - %s %s", req.RemoteAddr, req.Method, current_dir + path)

	if req.Method == "POST" {
		file, fileHeader, err := req.FormFile("file")
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
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(w, "<h1>permission denied</h1>")
		log.Printf("permission denied : %s", dir)
		return
	}

	fmt.Fprintf(w, "<h2>Directory listing for %s</h2><hr>", dir)
	fmt.Fprintf(w, "<form ENCTYPE=multipart/form-data method=post><input name=file type=\"file\"/><input type=submit value=\"upload\"/></form><hr>")
	fmt.Fprintf(w, "<table>")

    for _, file := range files {
		file_name := file.Name()
		if file.IsDir() {
			file_name += "/"
		}

		fmt.Fprintf(w, "<tr>")
		fmt.Fprintf(w, "<td><a href=\"%s\">%s<a></td>", file_name, file_name)

		if ! file.IsDir(){
			fmt.Fprintf(w, "<td align=\"right\">%d B</td>", file.Size())
			fmt.Fprintf(w, "<td align=\"right\">%s</td>", file.ModTime().Format(time.RFC1123))
		}

		fmt.Fprintf(w, "</tr>")
	}
	
	fmt.Fprintf(w, "</table>")
}