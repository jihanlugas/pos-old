package router

import (
	"encoding/json"
	"fmt"
	"github.com/jihanlugas/pos-old/config"
	"github.com/jihanlugas/pos-old/constant"
	"github.com/jihanlugas/pos-old/controller"
	"github.com/jihanlugas/pos-old/db"
	"github.com/jihanlugas/pos-old/model"
	"github.com/jihanlugas/pos-old/response"
	"github.com/labstack/echo/v4"
	"net/http"

	_ "github.com/jihanlugas/pos-old/docs"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func Init() *echo.Echo {
	router := websiteRouter()
	checkToken := checkTokenMiddleware()

	userController := controller.UserComposer()

	router.GET("/swg/*", echoSwagger.WrapHandler)

	router.GET("/", controller.Ping)
	router.POST("/sign-in", userController.SignIn)
	router.GET("/sign-out", userController.SignOut)
	router.GET("/refresh-token", userController.RefreshToken, checkToken)

	return router

}

func httpErrorHandler(err error, c echo.Context) {
	var errorResponse *response.Response
	code := http.StatusInternalServerError
	switch e := err.(type) {
	case *echo.HTTPError:
		// Handle pada saat URL yang di request tidak ada. atau ada kesalahan server.
		code = e.Code
		errorResponse = &response.Response{
			Status:  false,
			Message: fmt.Sprintf("%v", e.Message),
			Payload: map[string]interface{}{},
			Code:    code,
		}
	case *response.Response:
		errorResponse = e
	default:
		// Handle error dari panic
		code = http.StatusInternalServerError
		if config.Debug {
			errorResponse = &response.Response{
				Status:  false,
				Message: err.Error(),
				Payload: map[string]interface{}{},
				Code:    http.StatusInternalServerError,
			}
		} else {
			errorResponse = &response.Response{
				Status:  false,
				Message: "Internal server error",
				Payload: map[string]interface{}{},
				Code:    http.StatusInternalServerError,
			}
		}
	}

	js, err := json.Marshal(errorResponse)
	if err == nil {
		_ = c.Blob(code, echo.MIMEApplicationJSONCharsetUTF8, js)
	} else {
		b := []byte("{error: true, message: \"unresolved error\"}")
		_ = c.Blob(code, echo.MIMEApplicationJSONCharsetUTF8, b)
	}
}

func checkTokenMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var err error

			userLogin, err := controller.ExtractClaims(c.Request().Header.Get(config.HeaderAuthName))
			if err != nil {
				return response.ErrorForce(http.StatusUnauthorized, err.Error(), response.Payload{}).SendJSON(c)
			}

			conn, closeConn := db.GetConnection()
			defer closeConn()

			var user model.User
			err = conn.Where("id = ? ", userLogin.UserID).First(&user).Error
			if err != nil {
				return response.ErrorForce(http.StatusUnauthorized, "Token Expired!", response.Payload{}).SendJSON(c)
			}

			if user.PassVersion != userLogin.PassVersion {
				return response.ErrorForce(http.StatusUnauthorized, "Token Expired~", response.Payload{}).SendJSON(c)
			}

			c.Set(constant.TokenUserContext, userLogin)
			return next(c)
		}
	}
}
