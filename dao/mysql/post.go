package mysql

import "bluebell/models"

/**
 * @Author: zhaobin
 * @Description TODO
 * @Date: 2022-04-08 16:40
 */

func CreatePost(p *models.Post) error {
	sqlStr := "insert into post(post_id,title,content,author_id,community_id) values(?,?,?,?,?)"
	_, err := db.Exec(sqlStr, p.ID, p.Title, p.Content, p.AuthorID, p.CommunityID)
	return err
}

func GetPostById(pid int64) (p *models.Post, err error) {
	p = new(models.Post)
	sqlStr := "select post_id,title,content,author_id,community_id,create_time from post where post_id=?"
	err = db.Get(p, sqlStr, pid)

	return
}

// GetPostList 查询帖子列表函数
func GetPostList(page int64, size int64) (posts []*models.Post, err error) {
	sqlStr := `select 
	post_id, title, content, author_id, community_id, create_time
	from post
	ORDER BY create_time
	DESC
	limit ?,?
	`
	posts = make([]*models.Post, 0, 2)
	err = db.Select(&posts, sqlStr, (page-1)*size, size)
	return
}
