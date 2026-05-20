// internal/model/user.go
package model

// User представляет модель пользователя в системе
type User struct {
	ID           int64  `db:"id" json:"id"`
	FirstName    string `db:"first_name" json:"first_name"`
	LastName     string `db:"last_name" json:"last_name"`
	MiddleName   *string `db:"middle_name" json:"middle_name,omitempty"`
	PhoneNumber  *string `db:"phone_number" json:"phone_number,omitempty"`
	Email        string `db:"email" json:"email"` // Уникальный
	DateOfBirth  *string `db:"date_of_birth" json:"date_of_birth,omitempty"`
	VKProfileURL *string `db:"vk_profile_url" json:"vk_profile_url,omitempty"`
	PasswordHash string `db:"password_hash" json:"-"`
	CreatedAt    string `db:"created_at" json:"created_at"`
	UpdatedAt    string `db:"updated_at" json:"updated_at"`
}
