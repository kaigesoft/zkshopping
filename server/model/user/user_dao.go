package user

import (
	"server/config"
	"server/model"
)

func init() {
	DefaultUserInfoConnection = UserInfoConnection(model.GetDBToDBConnect(config.GetDBConnect("db")))
}