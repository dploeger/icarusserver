package endpoints

import "github.com/gin-gonic/gin"

type Endpoint interface {
	Register(group *gin.RouterGroup) error
}
