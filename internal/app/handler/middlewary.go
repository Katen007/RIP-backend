package handler

import "github.com/gin-gonic/gin"

const CtxUserID = "user_id"

func LabFixedUser() gin.HandlerFunc {
	return func(c *gin.Context) { c.Set(CtxUserID, 1); c.Next() }
}
func getUserID(c *gin.Context) int {
	v, _ := c.Get(CtxUserID)
	if id, ok := v.(int); ok {
		return id
	}
	return 1
}
