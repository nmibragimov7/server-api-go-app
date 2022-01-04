package helpers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
)

func checkRole(c *gin.Context, userRole string) (err error) {
	role := c.GetString("role")
	err = nil
	if role != userRole {
		err = errors.New("unauthorized to access this resource")
		return err
	}
	return err
}

func Access(c *gin.Context, id string) (err error) {
	role := c.GetString("role")
	uid := c.GetString("_id")

	fmt.Println("role" + role)
	fmt.Println("uid" + uid)

	err = nil
	if role == "user" && uid != id {
		fmt.Println("error")
		err = errors.New("unauthorized to access this resource")
		return err
	}
	err = checkRole(c, role)
	return err
}
