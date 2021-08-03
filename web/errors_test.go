package web

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestErrorHandler(t *testing.T) {
	app := gin.Default()
	app.Use(ErrorHandler)
	app.GET("/", func(c *gin.Context) {
		c.Error(errors.New("error message"))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	app.ServeHTTP(w, req)

	assert.Equal(t, 500, w.Code)
	assert.Contains(t, w.Body.String(), "error message")
}

func TestErrorHandlerContentNegotiation(t *testing.T) {
	app := gin.Default()
	app.HTMLRender = NewLayoutRender(templatesFS, "templates/*.tmpl")
	app.Use(ErrorHandler)
	app.GET("/", func(c *gin.Context) {
		c.Error(errors.New("error message"))
		c.Error(errors.New("2nd error message"))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Accept", "text/html")

	app.ServeHTTP(w, req)

	assert.Equal(t, 500, w.Code)
	assert.Contains(t, w.Body.String(), "Ooops")
	assert.Contains(t, w.Body.String(), "error message</br>")
	assert.Contains(t, w.Body.String(), "2nd error message</br>")
}

func TestErrorHandlerWithHttpError(t *testing.T) {
	app := gin.Default()
	app.Use(ErrorHandler)
	app.GET("/", func(c *gin.Context) {
		c.Error(NotFoundError("error message"))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	app.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)
	assert.Contains(t, w.Body.String(), "error message")
}
