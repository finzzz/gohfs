### Running in Docker

This project can be run in a Docker container, it exposes only port 8080 (default). It is advised 
to run this configuration through a reverse proxy providing SSL if the service will be exposed over the internet, or use the integrated feature and provide gohfs with a valid certificate.

Password authentification is disabled per default. To enable it set the required environment variables. You only have to set either `GOHFS_PASSWORD` or `GOHFS_PASSWORD_RAW`.

An minimal example run would be

    # with docker run
    cd <PROJECT_ROOT>
    docker build . -t finzzz/gohfs:latest
    docker run -p '8080:8080' -v './data:/data' -e PGID=1000 -e PUID=1000 -name gohfs finzzz/gohfs:latest
    
    # with docker-compose
    cd <PROJECT_ROOT>
    docker-compose build
    docker-compose up -d
    
You can then connect to `localhost:8080` or `<server_ip>:8080` to access the application.
The application will then serve the data from the /data directory, if permissions are set correct.

Supportet environment variables
-----
| Name               | Description                                                                                |
| ------------------ | ------------------------------------------------------------------------------------------ |
| **Application**    |                                                                                            |
| TZ                 | Time Zone for the container. Is shown in the Webinterface (https://en.wikipedia.org/wiki/List_of_tz_database_time_zones) |
| PUID               | Unix user ID to run the container as                                                       |
| PGID               | Unix group ID to run the container as                                                      |
|                    |                                                                                            |
| GOHFS_HOST         | IP address on which the webserver listens *(default: 0.0.0.0)*                             |
| GOHFS_PORT         | Port on which the webserver listens *(default: 8080)*                                      |
| GOHFS_USER         | Username for authentification *(default: disabled)*                                        |
| GOHFS_PASSWORD     | Sha256 hashed password for authentification *(default: disabled/empty)*                    |
| GOHFS_PASSWORD_RAW | Raw (cleartext) password for authentification *(default: disabled/empty)*                  |
| GOHFS_MAXUP        | Maximum upload size in Bytes *(default -1)*                                                |
| GOHFS_ARGS         | Additional gohfs command-line arguments *(e.g. `"-dz -du -dl"`)*                           |


### Contribution

Like this project? Want to contribute? Awesome! Feel free to open some pull requests or just open an issue.

### Changelog

Detailed changes for each release are documented in the [release notes](https://github.com/finzzz/gohfs/releases).
