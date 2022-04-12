package redis

/**
 * @Author: zhaobin
 * @Description TODO
 * @Date: 2022-04-09 11:06
 */

const (
	KeyPrefix        = "bluebell:"
	KeyPostTimeZset  = "zset:post:time"  // 帖子及发帖时间
	KeyPostScoreZset = "zset:post:score" // 帖子及分数
	// KeyPostVotedHashPrefix 记录用户及投票类型，参数是post_id。其中-1表示反对票，0表示不投票，1表示赞成票
	KeyPostVotedHashPrefix = "hash:post:voted:"
	KeyCommunitySetPrefix  = "set:community:" // set; 保存每个分区下帖子的id
)

func getRedisKey(key string) string {
	return KeyPrefix + key
}
