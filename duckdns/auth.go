package duckdns

//Auth structure
type Auth struct {
	Token string
}

//AuthService structure
type AuthService struct {
	client *Client
	auth   Auth
}

//SetToken function
func (s *AuthService) SetToken(token string) {
	s.auth.Token = token
}

//GetToken function
func (s *AuthService) GetToken() string {
	return s.auth.Token
}
