package data

import (
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

// Define a custom ErrRecordNotFound error. We'll return this from our Get() method when
// looking up a movie that doesn't exist in our database.

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

// Create a Models struct which wraps the MovieModel. We'll add other models to this,
// like a UserModel and PermissionModel, as our build progresses.
type Models struct {
	// Set the Movies field to be an interface containing the methods that both the
	// 'real' model and mock model need to support.
	Movies interface {
		Insert(movie *Movie) error
		Get(id int64) (*Movie, error)
		Update(movie *Movie) error
		Delete(id int64) error
		GetAll(title string, genres []string, filters Filters) ([]*Movie, Metadata, error)
	}
	Users interface {
		Insert(user *User) error
		GetByEmail(email string) (*User, error)
		Update(user *User) error
		GetForToken(tokenScope, tokenPlaintext string) (*User, error)
	}
	Tokens interface {
		Insert(token *Token) error
		New(userID int64, ttl time.Duration, scope string) (*Token, error)
		DeleteAllForUser(scope string, userID int64) error
	}
	Permissions interface {
		GetAllForUser(userID int64) (Permissions, error)
		AddForUser(userID int64, codes ...string) error
	}
}

// For ease of use, we also add a New() method which returns a Models struct containing
// the initialized MovieModel.
func NewModels(pool *pgxpool.Pool) Models {
	return Models{
		Movies:      MovieModel{pool: pool},
		Users:       UserModel{pool: pool},
		Tokens:      TokenModel{pool: pool},
		Permissions: PermissionModel{pool: pool},
	}
}
