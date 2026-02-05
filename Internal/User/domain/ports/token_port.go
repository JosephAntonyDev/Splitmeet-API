package ports

type TokenManager interface {
	GenerateToken(userId int, email string, name string) (string, error)
	ValidateToken(token string) (bool, map[string]interface{}, error)
}