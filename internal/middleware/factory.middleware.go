package middleware

import "github.com/gin-gonic/gin"

type Factory struct {
}

func NewFactory() *Factory {
    return &Factory{}
}

func (f *Factory) Auth() gin.HandlerFunc {
    return Auth()
}

func (f *Factory) RequirePermission(permission string) gin.HandlerFunc {
    return RequirePermission(permission)
}