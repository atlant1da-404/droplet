package entity

type User struct {
	Id       string `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
