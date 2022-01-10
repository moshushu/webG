package mysql

import (
	"database/sql"
	"web1/models"

	"go.uber.org/zap"
)

// GetCommunityList 查询数据库内容
func GetCommunityList() (communityList []*models.Community, err error) {
	sqlStr := "select community_id, community_name from community"
	if err := db.Select(&communityList, sqlStr); err != nil {
		// 如果查询的是空的，没有记录，应该返回一个空的社区/列表
		if err == sql.ErrNoRows {
			zap.L().Warn("there is no community in db")
			err = nil
		}
	}
	return
}

// GetCommunityDetailByID
// models.CommunityDetail，在models中定义一个结构体用来存储查询到的信息
func GetCommunityDetailByID(id int64) (communityDetail *models.CommunityDetail, err error) {
	communityDetail = new(models.CommunityDetail)
	sqlStr := `select community_id, community_name, introduction, create_time from community where community_id=?`
	if err := db.Get(communityDetail, sqlStr, id); err != nil {
		// 如果查询的是空的，没有记录
		if err == sql.ErrNoRows {
			// 无效id
			err = ErrorInvalidID
			return nil, err
		}
	}
	return
}
