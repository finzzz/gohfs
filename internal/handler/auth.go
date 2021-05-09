package handler

import (
	"fmt"
    "net/http"
	"strings"
	"crypto/sha256"

	"gohfs/internal/config"
)

func (h HandlerObj) AuthHandler(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		iuser, ipass, _ := r.BasicAuth()
		if ! checkAuth(iuser, ipass, h.Config) {
			w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}else {
			f(w,r)
		}
	}
}

func checkAuth(iuser, ipass string, cfg config.Config) (bool) {
	pass := cfg.Pass

	if cfg.Pass == "" && cfg.HashedPass == "" {
		return true // doesn't have auth
	}

	if cfg.HashedPass != "" {
		pass = cfg.HashedPass
		ipass = fmt.Sprintf("%x", sha256.Sum256([]byte(ipass)))
	}

	if strings.Compare(cfg.User, iuser) == 0 && strings.Compare(pass, ipass) == 0 {
		return true
	}

	return false
}