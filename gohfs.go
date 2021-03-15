package main

// TODO
// auth??

import(
	"fmt"
	"log"
	"net"
	"net/http"
	"io/ioutil"
	"flag"
	"math"
	"strings"
	"time"
)

var host string
var port string
var current_dir string

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
	// header
	fmt.Fprintf(w, `<!DOCTYPE html>
<html>
	<head>
		<style>
			body {
				font-family: arial, sans-serif;
				font-weight: 300;
				line-height: 1.625;
				color: #555;
				margin-left: 20px;
			}
			
			table {
				border-collapse: separate;
				border-spacing: 0;
				text-align: left;
				box-sizing: content-box;
				border-top: 1px solid #dee2e6;
				border-bottom: 1px solid #111;
				width: calc(100%% - 40px);
			}
			
			a {
				color: #007bff;
				text-decoration: none;
			}

			th {
				height: 20px;
				border-bottom: 1px solid #111;
				border-top: 1px solid #dee2e6;
				padding-left: 20px;
			}

			th:hover span {
				display:none;
			}

			th:hover:before {
				content: "sort";
			}

			td {
				height: 30px;
				border-top: 1px solid #dee2e6;
				padding-left: 20px;
			}

			code {
				background-color: #efefef;
				padding-left: 0.5rem;
				padding-right: 0.5rem;
				padding-top: 0.15rem;
				padding-bottom: 0.15rem;
				border-radius: 0.3125rem;
				margin-right: 5px;
			}

			#upload {
				padding-bottom: 40px;
			}

			.sort {
				display:none;
			}
		</style>
	</head>
	<body>
		<h2>Directory : %s</h2>
	`, dir[:len(dir)-1])

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(w, "<h1>permission denied</h1>")
		log.Printf("permission denied : %s", dir)
		return
	}

	var ip string
	if strings.Split(req.Host,":")[0] != "0.0.0.0" {
		ip = strings.Split(req.Host,":")[0]
	} else {
		ip = getIP()
	}

	// upload
	fmt.Fprintf(w, `
		<div id=upload>
			<div style="float: left;">
				<form ENCTYPE=multipart/form-data id=upload-form method=post onclick=submitForm()>
					<input name=file type="file"/>
					<input type=submit value="upload"/>
				</form>
			</div>
			<div style="float: right; margin-right: 40px">
				<code id=uploadcmd>curl -F 'file=@uploadthis.txt' %s:%s</code>
				<button onclick=copyText()>copy</button>
			</div>
		</div>
	`, ip, port)

	// n items
	fmt.Fprintf(w, `
		<div style="float: right; margin-right: 40px">
			%d items
		</div>
	`, len(files))

	// table
	fmt.Fprintf(w, `
		<table id=filetable>
			<tr>
				<th onclick=sortTable(0)><span>Name</span></th>
				<th onclick=sortTable(1)><span>Type</span></th>
				<th onclick=sortTable(3,"num")><span>Size</span></th>
				<th onclick=sortTable(5,"date")><span>Last Modified</span></th>
			</tr>
	`)

    for _, file := range files {
		file_name := file.Name()
		if file.IsDir() {
			file_name += "/"
		}

		fmt.Fprintf(w, `
			<tr>
				<td><a href="%s">%s</a></td>`, file_name, file_name)

		if ! file.IsDir(){
			fsize, suffix := parseSize(file.Size())
			fmt.Fprintf(w, `
				<td>File</td>
				<td>%g %s</td>
				<td style="display: none;">%d</td>`, fsize, suffix, file.Size())
		}else {
			fmt.Fprintf(w, `
				<td>Directory</td>
				<td>--</td>
				<td style="display: none;">-1</td>`)
		}

		fmt.Fprintf(w, `
				<td>%s</td>
				<td style="display: none;">%s</td>
			</tr>
			`, file.ModTime().Format(time.RFC1123), file.ModTime().Format(time.RFC3339))
	}
	
	fmt.Fprintf(w, `
		</table>
		<p class=sort id=th_0>1</p>
		<p class=sort id=th_1>1</p>
		<p class=sort id=th_3>1</p>
		<p class=sort id=th_5>1</p>
	<script>
		function submitForm() {
			var f = document.getElementsByName('upload-form')[0];
			f.submit();
			f.reset();
		}

		function copyText() {
			const tmp = document.createElement('textarea');
			tmp.value = (document.getElementById("uploadcmd")).innerHTML;
			document.body.appendChild(tmp);
			tmp.select();
			document.execCommand('copy');
			document.body.removeChild(tmp);
		}

		function sortTable(idx, type) {
			var table, rows, switching, i, x, y, shouldSwitch, sortorder;
			table = document.getElementById("filetable");
			sortorder = document.getElementById("th_" + idx);
			switching = true;
			while (switching) {
				switching = false;
				rows = table.rows;
				for (i = 1; i < (rows.length - 1); i++) {
					shouldSwitch = false;
					x = rows[i].getElementsByTagName("TD")[idx];
					y = rows[i + 1].getElementsByTagName("TD")[idx];

					// comparison here
					if (type == "num") {
						if (Number(x.innerHTML) > Number(y.innerHTML) && sortorder.innerHTML > 0) {
							shouldSwitch = true;
							break;
						}
						
						if (Number(x.innerHTML) < Number(y.innerHTML) && sortorder.innerHTML < 0) {
							shouldSwitch = true;
							break;
						}
					} else if (type == "date"){
						x = new Date(x.innerHTML);
						y = new Date(y.innerHTML);
						if (x > y && sortorder.innerHTML > 0) {
							shouldSwitch = true;
							break;
						}
						
						if (x < y && sortorder.innerHTML < 0) {
							shouldSwitch = true;
							break;
						}
					} else {
						if (x.innerHTML.toLowerCase() > y.innerHTML.toLowerCase() && sortorder.innerHTML > 0) {
							shouldSwitch = true;
							break;
						}

						if (x.innerHTML.toLowerCase() < y.innerHTML.toLowerCase() && sortorder.innerHTML < 0) {
							shouldSwitch = true;
							break;
						}
					}
				}
				if (shouldSwitch) {
					rows[i].parentNode.insertBefore(rows[i + 1], rows[i]);
					switching = true;
				}
			}
			sortorder.innerHTML = Number(sortorder.innerHTML) * -1;
		}
	</script>
	</body>
</html>`)
		
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
