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
    - [x] Can disable directory listing
    - [ ] Can disable upload
    - [ ] Can disable zip
    - [ ] Regex listing
- Others
    - [ ] Log to file
    - [x] Compress binary
    - [ ] Minify JS on build
    - [ ] Show version
    - [x] Specify ip, port, dir

# Getting started
```bash
# running in current directory
./gohfs

# specifying parameters
./gohfs -host 127.0.0.1 -port 8081 -dir /tmp 

# https
./gohfs -tls -cert selfsigned.cert -key selfsigned.key

# disable directory listing
./gohfs -hide

# authentication
./gohfs -user gopher -pass gopher # raw password
./gohfs -user gopher -hpass 9cc1ee455a3363ffc504f40006f70d0c8276648a5d3eb3f9524e94d1b7a83aef # sha256 hashed

# getting help
./gohfs -h
```