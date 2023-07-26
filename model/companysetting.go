package model

import (
	"gorm.io/gorm"
	"time"
)

func (u *Companysettings) BeforeCreate(tx *gorm.DB) (err error) {
	now := time.Now()

	u.CreateDt = now
	u.UpdateDt = now
	return
}

func (u *Companysettings) BeforeUpdate(tx *gorm.DB) (err error) {
	if u.DeleteDt == nil {
		now := time.Now()
		u.UpdateDt = now
	}
	return
}
