package logic

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/models"
	"bluebell/pkg/snowflake"
	"go.uber.org/zap"
)

/**
 * @Author: zhaobin
 * @Description TODO
 * @Date: 2022-04-08 16:39
 */

func CreatePost(p *models.Post) (err error) {
	//生成postID
	p.ID = int64(snowflake.GenID())
	//保存到数据库
	err = mysql.CreatePost(p)
	if err != nil {
		return
	}
	err = redis.CreatePost(p.ID, p.CommunityID)
	return
}

// GetPostById 根据id查看帖子详情数据
func GetPostById(pid int64) (data *models.ApiPostDetail, error error) {

	//查询数据的基础信息
	post, err := mysql.GetPostById(pid)
	if err != nil {
		zap.L().Error("mysql.GetPostById(pid) failed", zap.Int64("pid", pid), zap.Error(err))
		return
	}
	user, err := mysql.GetUserById(post.AuthorID)
	if err != nil {
		zap.L().Error("mysql.GetUserById(post.AuthorID) failed", zap.Int64("author_id", post.AuthorID),
			zap.Error(err))
		return
	}

	//根据社区ID查询社区详情信息
	community, err := mysql.GetCommunityDetailByID(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityDetailByID(post.AuthorID) failed", zap.Int64("community_id", post.CommunityID),
			zap.Error(err))
		return
	}
	data = &models.ApiPostDetail{
		AuthorName:      user.Username,
		Post:            post,
		CommunityDetail: community,
	}
	return
}

// GetPostList 获取帖子数据，单纯的分页查询
func GetPostList(page int64, size int64) (data []*models.ApiPostDetail, err error) {
	posts, err := mysql.GetPostList(page, size)
	if err != nil {
		return nil, err
	}
	data = make([]*models.ApiPostDetail, 0, len(posts))
	for _, post := range posts {
		// 根据作者id查询作者信息
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(err))
			continue
		}
		// 根据社区id查询社区详细信息
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
				zap.Int64("community_id", post.CommunityID),
				zap.Error(err))
			continue
		}
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postDetail)
	}

	return
}

// GetPostListNew 进行逻辑判断，判断是否根据社区查询帖子列表
func GetPostListNew(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	// 根据请求参数的不同，执行不同的逻辑
	if p.CommunityID == 0 {
		// 查所有
		zap.L().Info("查询所有的帖子列表")
		data, err = GetPostListByTimeOrScore(p)
	} else {
		// 根据社区id查询
		zap.L().Info("根据社区id进行查询", zap.Any("社区id", p.CommunityID))
		data, err = GetCommunityPostList(p)
	}

	if err != nil {
		zap.L().Error("GetPostListNew failed", zap.Error(err))
		return nil, err
	}
	return
}

// GetPostListByTimeOrScore 查询帖子列表，可以根据时间顺序或者得分顺序
func GetPostListByTimeOrScore(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	// 1.去redis查询id列表
	ids, err := redis.GetPostIDsInOrder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDsInOrder(p) return 0 data")
		return
	}

	// 2.根据id去MySQL数据库查询帖子详细信息，返回的数据还要按照我给定的顺序返回
	return getPostListDetailByIds(ids)
}

// GetCommunityPostList 根据社区查询帖子列表
func GetCommunityPostList(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	// 1.去redis中查找该社区的 id 列表
	ids, err := redis.GetCommunityPostIDsInOrder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetCommunityPostIDsInOrder(p) return 0 data")
		return
	}
	// 2.根据id去MySQL数据库查询帖子详细信息，返回的数据还要按照我给定的顺序返回
	return getPostListDetailByIds(ids)
}

// getPostListDetailByIds 根据id去MySQL数据库查询帖子详细信息，返回的数据还要按照我给定的顺序返回
func getPostListDetailByIds(ids []string) (data []*models.ApiPostDetail, err error) {

	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}

	// 提前查询好每篇帖子的投票数
	voteData, err := redis.GetPostVoteData(ids)

	if err != nil {
		return
	}

	// 将帖子的作者及分区信息查询出来填充到帖子中
	for idx, post := range posts {
		// 根据作者id查询作者信息
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(err))
			continue
		}
		// 根据社区id查询社区详细信息
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
				zap.Int64("community_id", post.CommunityID),
				zap.Error(err))
			continue
		}
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			VoteNum:         voteData[idx],
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postDetail)
	}
	return
}
