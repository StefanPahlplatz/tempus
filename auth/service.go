package auth

type Service interface {
	Login(username string, password string) error
}
