package utils

import (
    "archive/zip"
    "log"
    "os"
)

func ZipWrite(s string) (string) {
	tmpZip := RandStr(12) + ".zip"

    outFile, err := os.Create(tmpZip)
    if err != nil {
        log.Println(err)
    }
    defer outFile.Close()

    w := zip.NewWriter(outFile)

	if IsDirPath(s) {
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