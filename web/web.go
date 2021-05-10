package web

import (
	"embed"
	"io/fs"
	"time"
	"strconv"

	"gohfs/internal/config"
	"gohfs/internal/utils"
)

//go:embed static/* template/*
var web embed.FS

type Templ struct {
	Dir 			string
	Scheme			string
	IP				string
	Port			string
	Items			[]Item
	NItems			string
	WebPath			string
	ZipPath			string
}

type Item struct {
	Name			string
	Type			string
	Size			string
	RawSize			string
	ModTime			string
	SHA1			string
	RawModTime		string
}

func Embed(cfg *config.Config){
	(*cfg).Web = web
}

func ParseItem(info fs.FileInfo) Item {
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
		fsize, suffix := utils.ParseSize(info.Size())

		tmp.Type = "File"
		tmp.Size = strconv.FormatFloat(fsize, 'f', 1, 64) + " " + suffix
		tmp.RawSize = strconv.FormatInt(info.Size(), 10)
	}

	return tmp
}