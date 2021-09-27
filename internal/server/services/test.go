package services

import "context"

var (
	TestIdentitySrv identityServer
	TestAuthSrv     authServer
	TestMovieSrv    movieServer
)

type TestConfig struct {
	Server  interface{}
	URL     string
	Body    interface{}
	Context context.Context
}
