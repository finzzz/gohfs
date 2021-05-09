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
	User			string
	Pass			string
	HashedPass		string
	WebPath			string
	ZipPath			string
	Hide			bool
	Web				embed.FS
}

func ParseConf(config *Config) {
	flag.StringVar(&config.Host, "host", "0.0.0.0", "Host")
	flag.StringVar(&config.Port, "port", "8080", "Port")
	flag.StringVar(&config.User, "user", "admin", "Username (default: admin)")
	flag.StringVar(&config.Pass, "pass", "", "Password")
	flag.StringVar(&config.HashedPass, "hpass", "", "Hashed Password (sha-256)")
	flag.StringVar(&config.WebPath, "webpath", "/gohfs-web", "UI Path")
	flag.StringVar(&config.ZipPath, "zippath", "/gohfs-zip", "Zip Path")
	flag.StringVar(&config.Dir, "dir", ".", "Directory to serve")
	flag.BoolVar(&config.Hide, "hide", false, "Disable Listing")
	flag.Parse()
}

func VerifyConf(config Config) {
	if config.Pass != "" && config.HashedPass != "" {
		log.Fatal(`Can only define either "Password" or "Hashed Password"`)
	}
}