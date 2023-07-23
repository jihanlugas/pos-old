package router

import (
	"encoding/json"
	"fmt"
	"github.com/jihanlugas/pos/config"
	"github.com/jihanlugas/pos/controller"
	"github.com/jihanlugas/pos/response"
	"github.com/labstack/echo/v4"
	"net/http"

	_ "github.com/jihanlugas/pos/docs"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func Init() *echo.Echo {

	router := websiteRouter()

	userController := controller.UserComposer()

	router.GET("/swg/*", echoSwagger.WrapHandler)

	router.GET("/", controller.Ping)
	router.POST("/sign-in", userController.SignIn)
	router.GET("/sign-out", userController.SignOut)
	//router.GET("/refresh-token", userController.RefreshToken, checkToken)

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
