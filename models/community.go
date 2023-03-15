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

	tx := utils.DB.Begin()

	//发生任何异常都将回滚
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := utils.DB.Create(&community).Error; err != nil {
		fmt.Println(err)
		tx.Rollback()
		return -1, "创建失败"
	}

	utils.DB.Where("name = ?", community.Name).Find(&community)
	if community.Name == "" {
		tx.Rollback()
		return -1, "创建失败"
	}

	contact := Contact{
		OwnerId:  community.OwnerId,
		TargetId: community.ID,
		Type:     2,
	}
	if err := utils.DB.Create(&contact).Error; err != nil {
		fmt.Println(err)
		tx.Rollback()
		return -1, "创建失败"
	}

	tx.Commit()

	return 0, "创建成功"
}

func JoinGroup(userId uint, comId string) (int, string) {
	contact := Contact{
		OwnerId: userId,
		Type:    2,
	}
	community := Community{}

	utils.DB.Where("id = ? or name = ?", comId, comId).Find(&community)
	if community.Name == "" {
		return -1, "没有找到群"
	}
	utils.DB.Where("owner_id = ? and target_id = ? and type = 2 ", userId, comId).Find(&contact)
	if !contact.CreatedAt.IsZero() {
		return -1, "已加过此群"
	} else {
		contact.TargetId = community.ID
		utils.DB.Create(&contact)
		return 0, "加群成功"
	}
}

func LoadCommunity(ownerId uint) ([]*Community, string) {
	contacts := make([]Contact, 0)
	objIds := make([]uint64, 0)
	utils.DB.Where("owner_id = ? and type = 2", ownerId).Find(&contacts)
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
