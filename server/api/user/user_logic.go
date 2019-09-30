package user

import (
	"server/api"
	"server/config"
	"server/model/user"
)

//DoGetUserInfoList
func DoGetUserInfoList(c *api.Context) (int,string,interface{}) {
	id := c.Query("id")
	rtList,err := user.FindUserInfo("id = ?",id)
	if err != nil {
		return config.DatabaseError,c.Error(err).Error(),nil
	}
	return 0,"",rtList
}
