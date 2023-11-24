package security

import "fmt"

type InMemoryUserRepository struct {
	db map[string]User
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		db: make(map[string]User),
	}
}

func (r *InMemoryUserRepository) GetUser(t Token) (*User, error) {
	user, ok := r.db[t.AccountID]
	if !ok {
		return nil, fmt.Errorf("user not found")
	}

	return &user, nil
}

func (r *InMemoryUserRepository) GetUserByID(id string) (*User, error) {
	user, ok := r.db[id]
	if !ok {
		return nil, fmt.Errorf("user not found")
	}

	return &user, nil
}

func (r *InMemoryUserRepository) Save(u User) error {
	r.db[u.AccountID] = u
	return nil
}

type InMemoryTokenRepository struct {
	db map[string]ActivationToken
}

func NewInMemoryTokenRepository() *InMemoryTokenRepository {
	return &InMemoryTokenRepository{
		db: make(map[string]ActivationToken),
	}
}

func (i InMemoryTokenRepository) SaveToken(t ActivationToken) error {
	i.db[t.Token] = t

	return nil
}

func (i InMemoryTokenRepository) FindActivationToken(token string) (*ActivationToken, error) {
	t, ok := i.db[token]
	if !ok {
		return nil, fmt.Errorf("token not found")
	}

	return &t, nil
}
