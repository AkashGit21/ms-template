package services

var (
	TestIdentitySrv identityServer
	TestAuthSrv     authServer
	TestMovieSrv    movieServer
)

type TestConfig struct {
	Server interface{}
	URL    string
	Body   interface{}
}
