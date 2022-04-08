package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
)

/**
 * @Author: zhaobin
 * @Description 查询社区的相关数据
 * @Date: 2022-04-08 15:02
 */

func GetCommunityList() (data []*models.Community, err error) {
	return mysql.GetCommunityList()
}

func GetCommunityDetail(id int64) (data *models.CommunityDetail, err error) {
	return mysql.GetCommunityDetailByID(id)
}
