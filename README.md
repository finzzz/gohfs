# GoHFS
Feature Rich Simple HTTP Server

![](ss.png)
  
# Installation
```bash
wget -O gohfs https://github.com/finzzz/gohfs/raw/master/releases/gohfs-linux-amd64
chmod +x gohfs
```

# Getting started
```bash
# running in current directory
./gohfs

# specifying parameters
./gohfs -host 127.0.0.1 -port 8081 -dir /tmp 

# disable directory listing
./gohfs -hide

# authentication
./gohfs -user gopher -pass gopher # raw password
./gohfs -user gopher -hpass 9cc1ee455a3363ffc504f40006f70d0c8276648a5d3eb3f9524e94d1b7a83aef # sha256 hashed

# getting help
./gohfs -h
```

# TODO
- [x] disable directory listing option
- [ ] https support
- [x] add authentication
- [ ] verbose logging
- [x] download file/folder as zip
