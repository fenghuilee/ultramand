package server

const (
	NotAuthorized = `HTTP/1.0 401 Not Authorized
WWW-Authenticate: Basic realm="ultraman"
Content-Length: 23

Authorization required
`

	NotFound = `HTTP/1.0 404 Not Found
Content-Length: %d

Tunnel %s not found
`

	BadRequest = `HTTP/1.0 400 Bad Request
Content-Length: 12

Bad Request
`
)

type CliOptions struct {
	domain    string
	http      string
	webSocket string
	ssdb      string
	logTo     string
	logLevel  string
}
