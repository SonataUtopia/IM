package models

import (
	"fmt"

	"github.com/SonataUtopia/IM/utils"
	"gorm.io/gorm"
)

type Community struct {
	gorm.Model
	Name    string
	OwnerId uint
	Icon    string
	Desc    string
}

func CreateCommunity(community Community) (int, string) {
	if len(community.Name) == 0 {
		return -1, "群名称不能为空"
	}
	if community.OwnerId == 0 {
		return -1, "请先登录"
	}

	if err := utils.DB.Create(&community).Error; err != nil {
		fmt.Println(err)
		return -1, "创建失败"
	}
	return 0, "创建成功"
}

func LoadCommunity(ownerId uint) ([]*Community, string) {
	contacts := make([]Contact, 0)
	objIds := make([]uint64, 0)
	utils.DB.Where("owner_id = ? and type=2", ownerId).Find(&contacts)
	for _, v := range contacts {
		objIds = append(objIds, uint64(v.TargetId))
	}

	data := make([]*Community, 10)
	utils.DB.Where("id in ?", objIds).Find(&data)
	for _, v := range data {
		fmt.Println(v)
	}
	// utils.DB.Where()
	return data, "加载成功"
}
