package model

type User struct {
	Base
	Email     string `json:"email" gorm:"tinytext;unique"`
	Firstname string `json:"firstname" gorm:"tinytext"`
	Lastname  string `json:"lastname" gorm:"tinytext"`
	Role      string `json:"role" gorm:"tinytext"`
}
