package entity

type Account struct {
	Id              string           `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserId          string           `json:"userId" gorm:"type:uuid;index"`
	AccountDevices  []AccountDevices `json:"accountDevices" gorm:"foreignkey:AccountID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	AccountSettings *AccountSettings `json:"accountSettings" gorm:"foreignkey:AccountID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type AccountDevices struct {
	Id         string `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	AccountID  string `json:"AccountID" gorm:"type:uuid;index"`
	Name       string `json:"name"`
	OS         string `json:"os"`
	MacAddress string `json:"macAddress"`
	Active     bool   `json:"active"`
}

type AccountSettings struct {
	Id        string `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	AccountID string `json:"AccountID" gorm:"type:uuid;index"`
	Language  string `json:"language"`
}
