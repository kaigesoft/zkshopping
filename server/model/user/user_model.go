package user

import (
	"server/model"
	"time"
	"fmt"
)

/**
CREATE TABLE tbUserInfo (
  id int(11) NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  userName varchar(256) NOT NULL DEFAULT '' COMMENT '用户名',
  userGroupId int(32) NOT NULL DEFAULT 0 COMMENT '用户角色关联id',
  userDesc varchar(1024) NOT NULL DEFAULT 0 COMMENT '描述信息',
  status int(10) NOT NULL DEFAULT 1 COMMENT '记录状态:1表示正常，0表示删除',
  statusTime datetime NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '状态变更时间',
  creator varchar(64) NOT NULL DEFAULT '' COMMENT '创建人',
  createTime datetime NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '创建时间@now',
  updater varchar(64) NOT NULL DEFAULT '' COMMENT '更新人',
  updateTime datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  watcher varchar(256) NOT NULL DEFAULT 'all' COMMENT '关注人',
  PRIMARY KEY (id),
  UNIQUE KEY UNIQUEID(userName,userGroupId,statusTime)
);
JSGEN **/

var _ = time.Now

// UserInfoConnection UserInfo连接类型
type UserInfoConnection func() model.DBConnect

// DefaultUserInfoConnection DefaultUserInfo默认连接
var DefaultUserInfoConnection UserInfoConnection

// UserInfo UserInfo值类型
type UserInfo struct {
	ID int `json:"id"`
	UserName string `json:"userName"`
	UserGroupID int `json:"userGroupId"`
	UserDesc string `json:"userDesc"`
	Status int `json:"status"`
	StatusTime time.Time `json:"statusTime"`
	Creator string `json:"creator"`
	CreateTime time.Time `json:"createTime"`
	Updater string `json:"updater"`
	UpdateTime time.Time `json:"updateTime"`
	Watcher string `json:"watcher"`
}

// Add 插入UserInfo
func (c UserInfoConnection) Add(model *UserInfo) (int64, error) {
	sqlStr := "INSERT INTO `tbUserInfo` (`userName`, `userGroupId`, `userDesc`, `status`, `statusTime`, `creator`, `createTime`, `updater`, `updateTime`, `watcher`) VALUES(?, ?, ?, ?, ?, ?, now(), ?, now(), ?)"
	result, err := c().Exec(sqlStr, model.UserName, model.UserGroupID, model.UserDesc, model.Status, model.StatusTime, model.Creator, model.Updater, model.Watcher)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// AddUserInfo 插入UserInfo
func AddUserInfo(model *UserInfo) (int64, error) {
	return DefaultUserInfoConnection.Add(model)
}

// Find 查询UserInfo
func (c UserInfoConnection) Find(condition string, args ...interface{}) ([]*UserInfo, error) {
	sqlStr := "SELECT `id`, `userName`, `userGroupId`, `userDesc`, `status`, `statusTime`, `creator`, `createTime`, `updater`, `updateTime`, `watcher` FROM `tbUserInfo`"
	if len(condition) > 0 {
		sqlStr = sqlStr + " WHERE " + condition
	}
	results := make([]*UserInfo, 0)

	stmt, err := c().Prepare(sqlStr)
	if err != nil {
		return results, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(args...)
	if err != nil {
		return results, err
	}

	defer rows.Close()
	for rows.Next() {
		model := UserInfo{}
		values := []interface{}{
			new(interface{}),
			new(interface{}),
			new(interface{}),
			new(interface{}),
			new(interface{}),
			new(interface{}),
			new(interface{}),
			new(interface{}),
			new(interface{}),
			new(interface{}),
			new(interface{}),
		}
		rows.Scan(values...)
		if *(values[0].(*interface{})) != nil {
			tmp := int((*(values[0].(*interface{}))).(int64))
			model.ID = tmp
		}
		if *(values[1].(*interface{})) != nil {
			tmp := string((*(values[1].(*interface{}))).([]uint8))
			model.UserName = tmp
		}
		if *(values[2].(*interface{})) != nil {
			tmp := int((*(values[2].(*interface{}))).(int64))
			model.UserGroupID = tmp
		}
		if *(values[3].(*interface{})) != nil {
			tmp := string((*(values[3].(*interface{}))).([]uint8))
			model.UserDesc = tmp
		}
		if *(values[4].(*interface{})) != nil {
			tmp := int((*(values[4].(*interface{}))).(int64))
			model.Status = tmp
		}
		if *(values[5].(*interface{})) != nil {
			tmp := (*(values[5].(*interface{}))).(time.Time)
			model.StatusTime = tmp
		}
		if *(values[6].(*interface{})) != nil {
			tmp := string((*(values[6].(*interface{}))).([]uint8))
			model.Creator = tmp
		}
		if *(values[7].(*interface{})) != nil {
			tmp := (*(values[7].(*interface{}))).(time.Time)
			model.CreateTime = tmp
		}
		if *(values[8].(*interface{})) != nil {
			tmp := string((*(values[8].(*interface{}))).([]uint8))
			model.Updater = tmp
		}
		if *(values[9].(*interface{})) != nil {
			tmp := (*(values[9].(*interface{}))).(time.Time)
			model.UpdateTime = tmp
		}
		if *(values[10].(*interface{})) != nil {
			tmp := string((*(values[10].(*interface{}))).([]uint8))
			model.Watcher = tmp
		}
		results = append(results, &model)
	}
	return results, nil
}

// FindUserInfo 查询UserInfo
func FindUserInfo(condition string, args ...interface{}) ([]*UserInfo, error) {
	return DefaultUserInfoConnection.Find(condition, args...)
}
// PagedQuery 分页查询UserInfo
func (c UserInfoConnection) PagedQuery(condition string, pageSize uint, page uint, args ...interface{}) (totalCount uint, rows []*UserInfo, err error) {
	sqlStr := "SELECT COUNT(1) as cnt FROM `tbUserInfo`"
	if len(condition) > 0 {
		sqlStr = sqlStr + " WHERE " + condition
	}

	cr := c().QueryRow(sqlStr, args...)

	err = cr.Scan(&totalCount)
	if err != nil {
		return 0, nil, err
	}
	if page > 0 {
		page = page - 1
	}
	offset := page * pageSize
	if totalCount <= offset {
		return totalCount, []*UserInfo{}, nil
	}

	if len(condition) == 0 {
		condition = fmt.Sprintf("1=1")
	}
	condition = condition + fmt.Sprintf(" LIMIT %d OFFSET %d", pageSize, offset)
	rows, err = c.Find(condition, args...)
	return
}

// UserInfoPagedQuery 分页查询UserInfo
func UserInfoPagedQuery(condition string, pageSize uint, page uint, args ...interface{}) (totalCount uint, rows []*UserInfo, err error) {
	return DefaultUserInfoConnection.PagedQuery(condition, pageSize, page, args...)
}

// Get 获取UserInfo
func (c UserInfoConnection) Get(condition string, args ...interface{}) (*UserInfo, error) {
	results, err := c.Find(condition, args...)

	if err != nil {
		return nil, err
	}

	if len(results) > 0 {
		return results[0], nil
	}

	return nil, nil
}


// GetUserInfo 获取UserInfo
func GetUserInfo(condition string, args ...interface{}) (*UserInfo, error) {
	return DefaultUserInfoConnection.Get(condition, args...)
}

// Update 更新UserInfo
func (c UserInfoConnection) Update(model *UserInfo) (int64, error) {
	sqlStr := "UPDATE `tbUserInfo` SET `userName` = ?, `userGroupId` = ?, `userDesc` = ?, `status` = ?, `statusTime` = ?, `creator` = ?, `updater` = ?, `updateTime` = now(), `watcher` = ? WHERE `id` = ?"
	result, err := c().Exec(sqlStr, model.UserName, model.UserGroupID, model.UserDesc, model.Status, model.StatusTime, model.Creator, model.Updater, model.Watcher, model.ID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// UpdateUserInfo 更新UserInfo
func UpdateUserInfo(model *UserInfo) (int64, error) {
	return DefaultUserInfoConnection.Update(model)
}
