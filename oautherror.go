package jwttransport

type JWTAuthError struct {
	prefix string
	msg    string
}

func (oe JWTAuthError) Error() string {
	return "JWTAuthError: " + oe.prefix + ": " + oe.msg
}
