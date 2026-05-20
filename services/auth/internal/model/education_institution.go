// internal/model/education_institution.go
package model

// EducationInstitution представляет модель учебного заведения
type EducationInstitution struct {
	ID        int64  `db:"id" json:"id"`
	Title     string `db:"title" json:"title"`
	CreatedAt string `db:"created_at" json:"created_at"`
	UpdatedAt string `db:"updated_at" json:"updated_at"`
}
