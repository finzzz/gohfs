package config

import(
	"flag"
	"log"
	"embed"

	"gohfs/internal/utils"
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
	ZipTemp			string
	SHA1Path		string
	DisableListing	bool
	DisableZip		bool
	DisableUp		bool
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
	flag.StringVar(&config.User, "user", "admin", "Username")
	flag.StringVar(&config.Pass, "pass", "", "Password")
	flag.StringVar(&config.HashedPass, "hpass", "", "Hashed Password (sha-256)")

	// api path
	flag.StringVar(&config.WebPath, "webpath", "/gohfs-web", "UI API")
	flag.StringVar(&config.ZipPath, "zippath", "/gohfs-zip", "Zip API")
	flag.StringVar(&config.SHA1Path, "sha1path", "/gohfs-sha1", "SHA1 API")
	
	// disable feature
	flag.BoolVar(&config.DisableListing, "dl", false, "Disable Listing")
	flag.BoolVar(&config.DisableZip, "dz", false, "Disable Zip")
	flag.BoolVar(&config.DisableUp, "du", false, "Disable Upload")

	// others
	flag.StringVar(&config.ZipTemp, "ziptemp", ".", "Temporary Zip Folder")

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

	if ! utils.IsDirExist(config.ZipTemp) {
		log.Fatal(`Temporary Zip Folder doesn't exist`)
	}
}