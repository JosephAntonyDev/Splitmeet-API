package app

import (
	"github.com/JosephAntonyDev/splitmeet-api/internal/user/domain/entities"
	"github.com/JosephAntonyDev/splitmeet-api/internal/user/domain/repository"
)

type SearchUsers struct {
	userRepo repository.UserRepository
}

func NewSearchUsers(userRepo repository.UserRepository) *SearchUsers {
	return &SearchUsers{userRepo: userRepo}
}

type UserSearchResult struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone,omitempty"`
}

func (s *SearchUsers) Execute(query string, limit int) ([]UserSearchResult, error) {
	if limit <= 0 || limit > 20 {
		limit = 10
	}

	users, err := s.userRepo.SearchByUsername(query, limit)
	if err != nil {
		return nil, err
	}

	results := make([]UserSearchResult, len(users))
	for i, u := range users {
		results[i] = mapToSearchResult(u)
	}

	return results, nil
}

func mapToSearchResult(u entities.User) UserSearchResult {
	return UserSearchResult{
		ID:       u.ID,
		Username: u.Username,
		Name:     u.Name,
		Email:    u.Email,
		Phone:    u.Phone,
	}
}
