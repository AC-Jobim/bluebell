package redis

import (
	"bluebell/models"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"strconv"
	"time"
)

/**
 * @Author: zhaobin
 * @Description TODO
 * @Date: 2022-04-11 18:12
 */

// getIDsForKey 根据key返回zset集合中的ids
func getIDsForKey(key string, page, size int64) ([]string, error) {
	// 2.确定查询的索引起始点
	start := (page - 1) * size
	end := start + size - 1
	// 3. ZREVRANGE 按分数从大到小的顺序查询指定数量的元素
	return rdb.ZRevRange(key, start, end).Result()
}

// GetPostIDsInOrder 根据order（时间或者分数）去获取帖子列表ids
func GetPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	// 从redis获取id
	// 1.根据用户请求中携带的order参数确定要查询的redis key
	key := getRedisKey(KeyPostTimeZset)
	zap.L().Debug("GetPostIDsInOrder", zap.Any("p.Order", p.Order))
	if p.Order == models.OrderScore {
		zap.L().Info("根据得分顺序排序帖子列表")
		key = getRedisKey(KeyPostScoreZset)
	} else {
		zap.L().Info("根据时间顺序排序帖子列表")
	}

	// 存在的话就直接根据key查询ids
	return getIDsForKey(key, p.Page, p.Size)
}

// GetCommunityPostIDsInOrder 根据order（时间或者分数）去获取某社区帖子列表ids
// p.CommunityID 为该社区的id
func GetCommunityPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	orderKey := getRedisKey(KeyPostTimeZset)
	if p.Order == models.OrderScore {
		orderKey = getRedisKey(KeyPostScoreZset) //如果order为 models.OrderScore帖子列表按照分数排序
	}
	// 使用 zinterstore 把分区的帖子set与帖子分数的 zset 生成一个新的zset
	// 针对新的zset 按之前的逻辑取数据
	// 社区set集合中的key
	cKey := getRedisKey(KeyCommunitySetPrefix + strconv.Itoa(int(p.CommunityID)))

	// 新生成的key
	key := orderKey + ":" + strconv.Itoa(int(p.CommunityID))

	// 利用缓存key减少zinterstore的执行次数，首先查询缓存中是否存在
	if rdb.Exists(key).Val() < 1 {
		// 缓存若不存在，需要计算
		pipeline := rdb.Pipeline()
		pipeline.ZInterStore(key, redis.ZStore{Aggregate: "MAX"}, cKey, orderKey) //计算最大值
		pipeline.Expire(key, 60*time.Second)                                      // 设置超时时间60s
		_, err := pipeline.Exec()
		if err != nil {
			zap.L().Debug("GetCommunityPostIDsInOrder", zap.Error(err))
			return nil, err
		}
	}
	// 存在的话就直接根据key查询ids
	return getIDsForKey(key, p.Page, p.Size)
}

// GetPostVoteData 根据ids查询每篇帖子的投赞成票的数据
func GetPostVoteData(ids []string) (data []int64, err error) {
	// 使用pipeline一次发送多条命令，减少RTT
	pipeline := rdb.Pipeline()
	for _, id := range ids {
		key := getRedisKey(KeyPostVotedHashPrefix + id)
		pipeline.HVals(key)
	}
	cmders, _ := pipeline.Exec()
	if err != nil {
		return nil, err
	}
	data = make([]int64, 0, len(cmders))
	for _, cmder := range cmders {
		vales := cmder.(*redis.StringSliceCmd).Val()
		var count int64 = 0
		for _, val := range vales {
			if val == "1" {
				count++
			} else if val == "-1" {
				count--
			}
		}
		data = append(data, count)
	}
	return
}
