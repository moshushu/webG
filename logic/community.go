package logic

import (
	"web1/dao/mysql"
	"web1/models"
)

func GetCommunityList() ([]*models.Community, error) {
	// 查找数据库 查找到所有的community并返回
	return mysql.GetCommunityList()
}

func GetCommunityDetail(id int64) (*models.CommunityDetail, error) {
	// 根据id查找社区的内容
	return mysql.GetCommunityDetailByID(id)
}
