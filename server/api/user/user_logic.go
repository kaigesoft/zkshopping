package user

import (
	"server/api"
	"server/config"
	"server/model/user"
)

//DoGetUserInfoList
func DoGetUserInfoList(c *api.Context) (int, string, interface{}) {
	par := map[string]*api.ParamConstruct{
		"id":        {FieldName: "id", DefaultValue: nil, CheckValue: nil, Need: false, Link: "and", Symbol: "="},
		"userName":  {FieldName: "userName", DefaultValue: nil, CheckValue: nil, Need: false, Link: "and", Symbol: "="},
		"searchKey": {FieldName: "id|userGroupId|userName|userDesc", Need: false, Link: "and"},
		"status":    {FieldName: "status", DefaultValue: 1, CheckValue: []interface{}{0, 1}, Need: false, Link: "and", Symbol: "="},
		"orderBy":   {DefaultValue: "id|desc", Need: false},
	}
	condition, args, err := c.GetConditionByParam(par)
	page, pageSize, _, err := c.GetPager()
	total, rtList, err := user.UserInfoPagedQuery(condition, pageSize, page, args...)
	if err != nil {
		return config.DatabaseError, c.Error(err).Error(), nil
	}
	ret := make([]interface{},0)
	for _, u := range rtList {
		item := api.H{}
		item["id"] = u.ID
		item["userName"] = u.UserName
		item["userDesc"] = u.UserDesc
		item["userGroupID"] = u.UserGroupID
		item["status"] = u.Status
		item["watcher"] = u.Watcher
		item["creator"] = u.Creator
		item["createTime"] = u.CreateTime.Format(api.Layout)
		item["updater"] = u.Updater
		item["updateTime"] = u.UpdateTime.Format(api.Layout)
		ret = append(ret,&item)
	}
	return 0, "", api.H{
		"count":    total,
		"dataList": ret,
	}
}

//DoSaveUserInfo
func DoSaveUserInfo(c *api.Context) (int, string, interface{}) {
	reqData := user.UserInfo{}
	if err := c.ShouldBind(&reqData); err != nil {
		return config.IllegalArgument, c.Error(err).Error(), nil
	}
	if reqData.ID > 0 {
		rowsAffected, err := user.UpdateUserInfo(&reqData)
		if err != nil {
			return config.DatabaseError, c.Error(err).Error(), nil
		}
		return 0, "", rowsAffected
	}
	id, err := user.AddUserInfo(&reqData)
	if err != nil {
		return config.DatabaseError, c.Error(err).Error(), nil
	}
	return 0, "", id
}
