package mysql

import (
	"bluebell/models"
	"database/sql"
	"go.uber.org/zap"
)

/**
 * @Author: zhaobin
 * @Description TODO
 * @Date: 2022-04-08 15:04
 */

func GetCommunityList() (data []*models.Community, err error) {
	sqlStr := "select community_id,community_name from community"
	if err = db.Select(&data, sqlStr); err != nil {
		// 查询为空
		if err == sql.ErrNoRows {
			zap.L().Warn("there is no community in db")
			err = nil
		}
	}
	return
}

func GetCommunityDetailByID(id int64) (community *models.CommunityDetail, err error) {
	community = new(models.CommunityDetail)
	sqlStr := "select community_id, community_name,introduction,create_time from community where community_id=?"
	if err = db.Get(community, sqlStr, 1); err != nil {
		if err == sql.ErrNoRows {
			err = ErrorInvalidID
		}
	}
	return
}
