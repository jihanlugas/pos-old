package controller

import (
	"errors"
	"github.com/jihanlugas/pos/config"
	"github.com/jihanlugas/pos/cryption"
	"github.com/jihanlugas/pos/db"
	"github.com/jihanlugas/pos/model"
	"github.com/jihanlugas/pos/request"
	"github.com/jihanlugas/pos/response"
	"github.com/jihanlugas/pos/utils"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type User struct{}

func UserComposer() User {
	return User{}
}

// SignIn Sign In user
// @Summary Sign in a user
// @Tags Authentication
// @Accept json
// @Produce json
// @Param req body request.Signin true "json req body"
// @Success      200  {object}	response.Response
// @Failure      500  {object}  response.Response
// @Router /sign-in [post]
func (h User) SignIn(c echo.Context) error {
	var err error
	var user model.User
	var usercompany model.Usercompany

	req := new(request.Signin)
	if err = c.Bind(req); err != nil {
		return err
	}

	if err = c.Validate(req); err != nil {
		return response.Error(http.StatusBadRequest, "error validation", response.ValidationError(err)).SendJSON(c)
	}

	conn, closeConn := db.GetConnection()
	defer closeConn()

	if utils.IsValidEmail(req.Username) {
		user.Email = req.Username
		err = conn.Where("email = ? ", user.Email).First(&user).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return response.Error(http.StatusBadRequest, "invalid username or password", response.Payload{}).SendJSON(c)
			}
			errorInternal(c, err)
		}
	} else {
		user.Username = req.Username
		err = conn.Where("username = ? ", user.Username).First(&user).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return response.Error(http.StatusBadRequest, "invalid username or password", response.Payload{}).SendJSON(c)
			}
			errorInternal(c, err)
		}
	}

	if !user.Active {
		return response.Error(http.StatusBadRequest, "user not active", response.Payload{}).SendJSON(c)
	}

	err = cryption.CheckAES64(req.Passwd, user.Passwd)
	if err != nil {
		return response.Error(http.StatusBadRequest, "invalid username or password", response.Payload{}).SendJSON(c)
	}

	err = conn.Where("user_id = ? ", user.ID).Where("default_company = ? ", true).First(&usercompany).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.Error(http.StatusBadRequest, "data not found", response.Payload{}).SendJSON(c)
		}
		errorInternal(c, err)
	}

	tx := conn.Begin()

	now := time.Now()
	user.LastLoginDt = &now
	user.UpdateDt = now
	tx.Save(&user)

	err = tx.Commit().Error
	if err != nil {
		errorInternal(c, err)
	}

	expiredAt := time.Now().Add(time.Hour * time.Duration(config.AuthTokenExpiredHour))
	token, err := CreateToken(user.ID, user.RoleID, usercompany.CompanyID, user.PassVersion, expiredAt)
	if err != nil {
		return response.Error(http.StatusBadRequest, "Failed generate token", response.Payload{}).SendJSON(c)
	}

	return response.Success(http.StatusOK, "success", response.Payload{
		"token": token,
	}).SendJSON(c)
}

func (h User) SignOut(c echo.Context) error {
	return response.Success(http.StatusOK, "success", response.Payload{}).SendJSON(c)
}