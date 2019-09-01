package model

import "github.com/gin-gonic/gin"

type itemList struct {
	ID     int
	Sort   int
	Name   string
	Status string
	Forums []forumsList
}

type forumsList struct {
	ID        int
	Fgroup    int
	Sort      int
	Name      string
	ShowName  string
	Msg       string
	Interval  int
	CreatedAt string
	UpdateAt  string
	Status    string
}

func GetList(c *gin.Context) {
	flist := []forumsList{
		forumsList{
			ID:        3,
			Fgroup:    4,
			Sort:      1,
			Name:      "每日好货3",
			ShowName:  "",
			Msg:       "msg",
			Interval:  1,
			CreatedAt: "2018-07-20 15:49:28",
			UpdateAt:  "2018-07-20 15:49:28",
			Status:    "n",
		},
	}
	rtList := itemList{
		ID:     1,
		Sort:   1,
		Name:   "每日好货D",
		Status: "n",
	}
	rtList.Forums = flist
	c.JSON(0, rtList)
}