package mysql

import (
	"database/sql"
	"web_app/models"

	"go.uber.org/zap"
)

// 获取community数据
func GetCommunityList() (communitydata []*models.Community, err error) {
	sqlStr := "select community_id, community_name from community"
	if err := db.Select(&communitydata, sqlStr); err != nil {
		if err == sql.ErrNoRows {
			//没有查到数据
			zap.L().Warn("There is no community in db")
			err = nil
		}
	}
	return
}

// 根据ID查询社区详情
func GetCommunityDetailByID(id int64) (communitydetail *models.CommunityDetail, err error) {
	communitydetail = new(models.CommunityDetail)
	sqlStr := `select community_id, community_name, introduction, create_time 
			   from community where community_id = ?`
	if err := db.Get(communitydetail, sqlStr, id); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("There is no communitydetail in db")
			err = ErrorInvalidID

		}
	}
	//查数据库，查找到所有的community并返回
	return communitydetail, err
}
