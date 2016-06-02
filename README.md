# RadiusAuthProxy

Small reverse proxy in Golang that authenticates users against a Radius server.

When a user is not yet authenticated it will ask for credentials using HTTP basic authentication. 
The credentials will be checked against a radius server. If access is allowed a cookie will be set for
subsequent requests.

### Installation

```sh
go get github.com/tpjg/RadiusAuthProxy
```

### Usage

```sh
go build
./RadiusAuthProxy -h
Usage of ./RadiusAuthProxy:
  -backend string
    	The backend to proxy to (default "http://127.0.0.1:80")
  -bind string
    	Address to bind (default ":8888")
  -cookie string
    	Name of cookie (default "RADAUTH")
  -radport string
    	Radius server port (default "1812")
  -radsecret string
    	Radius shared secret (default "testing123")
  -radserver string
    	Radius server IP address or hostname (default "127.0.0.1")
```

### Example

```sh
./RadiusAuthProxy -bind=":8080" -backend="http://myserver.com" -radserver="192.168.1.1"
```

### Notes

Does not work for servers that redirect to another URL (including changing to https). 
Also a backend with "https" will only work if the reverse proxy is using https, currently not supported directly - so use nginx/haproxy or another SSL frontend.
