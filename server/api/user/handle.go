package user

import (
	"server/api"
)

func init() {
	api.SetRouterRegister(func(router *api.RouterGroup) {
		userRouteGroup := router.Group("/user")
		//用户相关加载
		userRouteGroup.StdGET("getUserList", DoGetUserInfoList)
	})
}