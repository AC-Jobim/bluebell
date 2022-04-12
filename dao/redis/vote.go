package redis

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"math"
	"strconv"
	"time"
)

/**
 * @Author: zhaobin
 * @Description TODO
 * @Date: 2022-04-11 13:58
 */

const (
	oneWeekInSeconds = 7 * 24 * 3600
	scorePerVote     = 432 // 每一票值多少分
)

var (
	ErrVoteTimeExpire = errors.New("投票时间已过")
	ErrVoteRepeated   = errors.New("不允许重复投票")
)

func VoteForPost(userID, postID string, value float64) error {
	// 1. 判断投票限制
	// 去redis取帖子发布时间
	postTime := rdb.ZScore(getRedisKey(KeyPostTimeZset), postID).Val()
	if float64(time.Now().Unix())-postTime > oneWeekInSeconds {
		return ErrVoteTimeExpire
	}

	// 2和3需要放到一个pipeline事务中操作

	// 2. 更新贴子的分数
	// 先查当前用户给当前帖子的投票记录
	ovStr := rdb.HGet(getRedisKey(KeyPostVotedHashPrefix+postID), userID).Val()
	fmt.Println()

	ov, err := strconv.ParseFloat(ovStr, 64)
	if err != nil {
		zap.L().Error("vote.go", zap.Error(err))
	}
	// 更新：如果这一次投票的值和之前保存的值一致，就提示不允许重复投票
	fmt.Println("value:", value, "ov:", ov)
	if value == ov {
		return ErrVoteRepeated
	}
	var op float64
	if value > ov {
		// 分数要加
		op = 1
	} else {
		// 分数要减
		op = -1
	}
	diff := math.Abs(ov - value) // 计算两次投票的差值

	//开启一个事务
	pipeline := rdb.TxPipeline()
	pipeline.ZIncrBy(getRedisKey(KeyPostScoreZset), op*diff*scorePerVote, postID)

	// 3. 记录用户为该贴子投票的数据
	if value == 0 {
		// 删除投票记录
		pipeline.HDel(getRedisKey(KeyPostVotedHashPrefix+postID), userID)
	} else {
		// 增加该用户的投票记录
		pipeline.HSet(getRedisKey(KeyPostVotedHashPrefix+postID), userID, value)
	}
	_, err = pipeline.Exec()
	return err
}

func CreatePost(postId int64) (err error) {
	pipeline := rdb.TxPipeline()
	// 帖子时间
	pipeline.ZAdd(getRedisKey(KeyPostTimeZset), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postId,
	})
	//帖子分数
	pipeline.ZAdd(getRedisKey(KeyPostScoreZset), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postId,
	})
	_, err = pipeline.Exec()
	return
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
