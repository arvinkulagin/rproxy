# RProxy
RProxy is a simple reverse proxy server writen in Golang. It can replace substring in body of backend server responses. Replacing work only for responses with Content-Type header set in "text/plain" and "text/html".
## Get
```
go get github.com/arvinkulagin/rproxy
```
## Get with git
```
mkdir $GOPATH/src/github.com/arvinkulagin
cd $GOPATH/src/github.com/arvinkulagin
git clone https://github.com/arvinkulagin/rproxy.git
```
## Build
```
cd $GOPATH/src/github.com/arvinkulagin/rproxy
make build
```
## Build docker image
```
cd $GOPATH/src/github.com/arvinkulagin/rproxy
make build-image
```
## Usage (binary)
```
rproxy -config rproxy.conf
```
## Usage (docker)
```
docker run -d --net "host" -v "/path/to/rproxy/donfig/dir:/config" rproxy
```
## Configuration
```
address     - interface and port to listening incoming requests in
target      - address of target server
pattern     - regexp pattern to match substrings for replacing
replacement - replacement substring
```