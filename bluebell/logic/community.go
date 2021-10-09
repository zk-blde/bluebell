package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
)

func GetCommunityList()([]*models.Community, error){
	// 查找数据库所有额Community并返回
	return mysql.GetCommunityList()

}

func GetCommunityDetail(id int64)(community *models.CommunityDetail, err error){
	return mysql.GetCommunityDetailByID(id)
}
