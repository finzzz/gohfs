package config

import(
	"flag"
	"log"
	"embed"
)

type Config struct {
	Host			string
	Port			string
	Dir				string
	Scheme			string
	CertPem			string
	KeyPem			string
	User			string
	Pass			string
	HashedPass		string
	WebPath			string
	ZipPath			string
	Hide			bool
	Web				embed.FS
}

func ParseConf(config *Config) {
	var tls bool

	// basic params
	flag.StringVar(&config.Host, "host", "0.0.0.0", "Host")
	flag.StringVar(&config.Port, "port", "8080", "Port")
	flag.StringVar(&config.Dir, "dir", ".", "Directory to serve")

	// tls
	flag.BoolVar(&tls, "tls", false, "Enable HTTPS")
	flag.StringVar(&config.CertPem, "cert", "", "Public certificate")
	flag.StringVar(&config.KeyPem, "key", "", "Private certificate")

	// auth
	flag.StringVar(&config.User, "user", "admin", "Username (default: admin)")
	flag.StringVar(&config.Pass, "pass", "", "Password")
	flag.StringVar(&config.HashedPass, "hpass", "", "Hashed Password (sha-256)")

	// api path
	flag.StringVar(&config.WebPath, "webpath", "/gohfs-web", "UI Path")
	flag.StringVar(&config.ZipPath, "zippath", "/gohfs-zip", "Zip Path")
	
	// disable feature
	flag.BoolVar(&config.Hide, "hide", false, "Disable Listing")

	flag.Parse()

	if tls {
		(*config).Scheme = "https"
	} else {
		(*config).Scheme = "http"
	}
}

func VerifyConf(config Config) {
	if config.Pass != "" && config.HashedPass != "" {
		log.Fatal(`Can only define either "Password" or "Hashed Password"`)
	}

	if config.Scheme == "https" && (config.CertPem == "" || config.KeyPem == "") {
		log.Fatal(`Must specify both "-cert" and "-key" if HTTPS is enabled`)
	}
}