package web

import "github.com/gin-gonic/gin"

func ErrorHandler(c *gin.Context) {
	c.Next()

	if len(c.Errors) == 0 {
		return
	}

	var code int
	var template string

	last := c.Errors.Last().Err

	switch e := last.(type) {
	case *HttpError:
		code = e.code
		template = e.template
	default:
		code = 500
		template = "error.html.tmpl"
	}

	c.Negotiate(code, gin.Negotiate{
		Offered:  []string{gin.MIMEJSON, gin.MIMEHTML, gin.MIMEPlain},
		HTMLName: template,
		Data:     c.Errors,
	})

	c.Abort()
}

type HttpError struct {
	msg      string
	code     int
	template string
}

func (e *HttpError) Error() string {
	return e.msg
}

func (e *HttpError) Code() int {
	return e.code
}

func (e *HttpError) Template() string {
	return e.template
}

func NotFoundError(msg string) *HttpError {
	return &HttpError{
		msg,
		404,
		"error404.html.tmpl",
	}
}
