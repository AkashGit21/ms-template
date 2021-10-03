package services

var (
	TestIdentitySrv *identityServer
	TestAuthSrv     *authServer
	TestMovieSrv    *movieServer
)

type TestCase struct {
	name        string
	args        interface{}
	expected    interface{}
	expectedErr string
}
