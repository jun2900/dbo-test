// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameUser = "users"

// User mapped from table <users>
type User struct {
	ID       int32  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Email    string `gorm:"column:email;not null" json:"email"`
	Password string `gorm:"column:password;not null" json:"password"`
}

// TableName User's table name
func (*User) TableName() string {
	return TableNameUser
}
