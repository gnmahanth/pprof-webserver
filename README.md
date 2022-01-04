# pprof-webserver
simple wrapper for [google/pprof](https://github.com/google/pprof) to run as webserver 

- left column shows list of profiles stored on the server 
- new profiles can be uploaded using dialog box on the top left 
- clicking on red "x" deletes profile file from server

![Screenshot](Screenshot.png?raw=true)

## Getting Started
- clone this repo 
- run *__make dev__* to start server
- default server port is 8080 and is accessible at url http://localhost:8080 in any browser

## Command line options
```bash
$ ./pprof-webserver --help
Usage of ./pprof-webserver:
  -debug
    	enable debug logs
  -port string
    	server listen port (default "8080")
  -storage string
    	path to directory containing profile files (default "data")
```

## Build
- binary 
```bash
make static
```
- docker image
```bash
docker build -t pprof-server:v1 .
```

## Run
- command line
```bash
$ ./pprof-webserver -port 8080
```
- docker image
```bash
docker run -it -p 8080:8080 pprof-server:v1
```

## Security 
Current implementation does not have kind of auth implementation, if you are deploying to the internet use a reverse proxy to secure the endpoint.
