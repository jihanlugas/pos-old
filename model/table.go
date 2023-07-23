package model

import (
	"time"
)

type User struct {
	ID          string     `gorm:"primaryKey"`
	RoleID      string     `gorm:"not null"`
	Email       string     `gorm:"not null"`
	Username    string     `gorm:"not null"`
	NoHp        string     `gorm:"not null"`
	Fullname    string     `gorm:"not null"`
	Passwd      string     `gorm:"not null"`
	PassVersion int        `gorm:"not null"`
	Active      bool       `gorm:"not null"`
	PhotoID     string     `gorm:"not null"`
	LastLoginDt *time.Time `gorm:"null"`
	CreateBy    string     `gorm:"not null"`
	CreateDt    time.Time  `gorm:"not null"`
	UpdateBy    string     `gorm:"not null"`
	UpdateDt    time.Time  `gorm:"not null"`
	DeleteBy    string     `gorm:"not null"`
	DeleteDt    *time.Time `gorm:"null"`
}

type Company struct {
	ID       string     `gorm:"primaryKey"`
	Name     string     `gorm:"not null"`
	CreateBy string     `gorm:"not null"`
	CreateDt time.Time  `gorm:"not null"`
	UpdateBy string     `gorm:"not null"`
	UpdateDt time.Time  `gorm:"not null"`
	DeleteBy string     `gorm:"not null"`
	DeleteDt *time.Time `gorm:"null"`
}

type Usercompany struct {
	ID               string     `gorm:"primaryKey"`
	UserID           string     `gorm:"not null"`
	CompanyID        string     `gorm:"not null"`
	IsDefaultCompany bool       `gorm:"not null"`
	IsCreator        bool       `gorm:"not null"`
	CreateBy         string     `gorm:"not null"`
	CreateDt         time.Time  `gorm:"not null"`
	UpdateBy         string     `gorm:"not null"`
	UpdateDt         time.Time  `gorm:"not null"`
	DeleteBy         string     `gorm:"not null"`
	DeleteDt         *time.Time `gorm:"null"`
}
