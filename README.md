# GoHFS
Feature-rich HTTP File Server

![](https://raw.githubusercontent.com/finzzz/images/master/gohfs.png)

# Features and roadmap
- UI
    - [x] Show QR Link
    - [x] SHA1 checksum
    - [x] Command line cheatsheet (curl, wget, PS)
    - [ ] Hot reload
    - [ ] Relative last modified
    - [ ] Regex filtering
- Functionality
    - [ ] Send message to backend
    - [ ] Global search using fzf
- Upload
    - [x] Single file upload
    - [ ] Folder upload
    - [ ] Multi file upload
    - [ ] Limit upload size
- Download
    - [x] as ZIP
    - [ ] Base64 encode
    - [ ] Multi file download
- Security
    - [x] HTTPS
    - [x] Basic Auth
        - [x] Can store as hashed password
    - [ ] Regex listing
- Options
    - [x] Disable directory listing
    - [x] Disable upload
    - [x] Disable zip
    - [ ] Specify temporary zip folder
- Others
    - [x] [Docker support](docker/README.md)
    - [ ] Log to file
    - [x] Compress binary
    - [ ] Minify JS on build
    - [ ] Show version
    - [x] Specify ip, port, dir
    - [ ] Replace linuxserver.io docker base image

# Getting started
```bash
# running in current directory
./gohfs

# specifying parameters
./gohfs -host 127.0.0.1 -port 8081 -dir /tmp 

# https
./gohfs -tls -cert selfsigned.cert -key selfsigned.key

# disable directory listing
./gohfs -dl

# authentication
./gohfs -user gopher -pass gopher # raw password
./gohfs -user gopher -hpass 9cc1ee455a3363ffc504f40006f70d0c8276648a5d3eb3f9524e94d1b7a83aef # sha256 hashed

# getting help
./gohfs -h
```

# Full Usage
```bash
$ ./gohfs -h
Usage of ./gohfs:
  -cert string
        Public certificate
  -dir string
        Directory to serve (default ".")
  -dl
        Disable Listing
  -du
        Disable Upload
  -dz
        Disable Zip
  -host string
        Host (default "0.0.0.0")
  -hpass string
        Hashed Password (sha-256)
  -key string
        Private certificate
  -pass string
        Password
  -port string
        Port (default "8080")
  -sha1path string
        SHA1 API (default "/gohfs-sha1")
  -tls
        Enable HTTPS
  -user string
        Username (default "admin")
  -webpath string
        UI API (default "/gohfs-web")
  -zippath string
        Zip API (default "/gohfs-zip")
```

# Contribution
Like this project? Want to contribute? Awesome! Feel free to open some pull requests or just open an issue.

# Changelog
Detailed changes for each release are documented in the [release notes](https://github.com/finzzz/gohfs/releases).
