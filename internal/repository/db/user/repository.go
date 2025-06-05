package user

type Repository interface {
	GetUserRepository() error
}
