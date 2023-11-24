package security

type Token struct {
	AccessToken string `json:"access_token"`
	AccountID   string `json:"account_id"`
}

type User struct {
	AccountID     string `json:"account_id"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
}

type TokenProvider interface {
	GetToken(code string) (*Token, error)
}

type UserRepository interface {
	GetUser(Token) (*User, error)
	GetUserByID(string) (*User, error)
	Save(User) error
}

type Validator struct {
	users  UserRepository
	client ExternalUserProvider
}

type ExternalUserProvider interface {
	GetUser(token Token) (*User, error)
}

func NewValidator(users UserRepository, c ExternalUserProvider) *Validator {
	return &Validator{users, c}
}

func (v *Validator) HasEmailVerified(t *Token) bool {
	if t == nil {
		return false
	}

	user, err := v.client.GetUser(*t)
	if user == nil || err != nil {
		return false
	}

	if user.EmailVerified {
		return true
	}

	user, err = v.users.GetUser(*t)
	if err != nil {
		return false
	}

	return user.EmailVerified
}
