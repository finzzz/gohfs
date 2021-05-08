# GoHFS
Golang implementation of simple HTTP server with upload feature.  

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
./gohfs -host 127.0.0.1 -port 8081 -dir /tmp -hide

# getting help
./gohfs -h
```

# TODO
- [x] disable directory listing option
- [ ] https support
- [ ] add authentication
- [ ] verbose logging
- [ ] compressed folder download
