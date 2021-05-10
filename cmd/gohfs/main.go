package main

import(
	"fmt"
	"log"
	"net/http"

	"gohfs/web"
	"gohfs/internal/config"
	"gohfs/internal/handler"
)

func main(){
	var cfg config.Config
	
	web.Embed(&cfg)
	config.ParseConf(&cfg)
	config.VerifyConf(cfg)

	handlerObj := handler.HandlerObj{Config: cfg}
	http_handler := http.HandlerFunc(handlerObj.AuthHandler(handlerObj.Handler))

    http_server := &http.Server{
            Addr:           cfg.Host + ":" + cfg.Port,
            Handler:        http_handler,
	}
	
	fmt.Printf("Serving on %s://%s:%s/ ...\n", cfg.Scheme, cfg.Host, cfg.Port)

	if cfg.Scheme == "https" {
		log.Fatal(http_server.ListenAndServeTLS(cfg.CertPem, cfg.KeyPem)) 
	} else {
		log.Fatal(http_server.ListenAndServe()) 
	}
}