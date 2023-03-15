package models

import (
	"fmt"

	"github.com/SonataUtopia/IM/utils"
	"gorm.io/gorm"
)

// 人员关系
type Contact struct {
	gorm.Model
	OwnerId  uint
	TargetId uint
	Type     int //关系类型	1.好友	2.群组	3.
	Desc     string
}

func (table *Contact) TableName() string {
	return "contact"
}

func SearchFriend(userId uint) []UserBasic {
	contacts := make([]Contact, 0)
	objIds := make([]uint64, 0)
	utils.DB.Where("owner_id = ? and type = 1", userId).Find(&contacts)
	for _, v := range contacts {
		fmt.Println(v)
		objIds = append(objIds, uint64(v.TargetId))
	}
	users := make([]UserBasic, 0)
	utils.DB.Where("id in ?", objIds).Find(&users)
	return users
}

// 添加好友
func AddFriend(userId uint, targetName string) (int, string) {
	fmt.Println("add friend start")
	user := UserBasic{}
	if targetName != "" {
		user = FindUserByName(targetName)
		if user.Name != "" {
			if userId == user.ID {
				return -1, "不能添加自己"
			}
			contactCheck := Contact{}
			utils.DB.Where("ownerId = ? and target_id = ? and type = 1", userId, user.ID).Find(&contactCheck)
			if contactCheck.ID != 0 {
				return -1, "不能重复添加好友"
			}

			tx := utils.DB.Begin()

			//发生任何异常都将回滚
			defer func() {
				if r := recover(); r != nil {
					tx.Rollback()
				}
			}()

			contact := Contact{
				OwnerId:  userId,
				TargetId: user.ID,
				Type:     1,
			}
			if err := utils.DB.Create(&contact).Error; err != nil {
				tx.Rollback()
				return -1, "添加好友失败"
			}
			contactAnother := Contact{
				OwnerId:  user.ID,
				TargetId: userId,
				Type:     1,
			}
			if err := utils.DB.Create(&contactAnother).Error; err != nil {
				tx.Rollback()
				return -1, "添加好友失败"
			}
			tx.Commit()
			return 0, "添加好友成功"
		}
		return -1, "添加好友失败"
	}
	return -1, "好友姓名不能为空"
}

func SearchUserByGroupId(communityId uint) []uint {
	contacts := make([]Contact, 0)
	objIds := make([]uint, 0)
	utils.DB.Where("target_id = ? and type=2", communityId).Find(&contacts)
	for _, v := range contacts {
		objIds = append(objIds, uint(v.OwnerId))
	}
	return objIds
}
