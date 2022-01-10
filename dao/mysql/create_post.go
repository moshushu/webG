package mysql

import (
	"web1/models"
)

// CreatePost 将创建的帖子写入数据库
func CreatePost(p *models.CreatePost) (err error) {
	sqlStr := `insert into post(
		post_id, title, content, author_id, community_id)
		values (?, ?, ?, ?, ?)
		`
	_, err = db.Exec(sqlStr, p.ID, p.Title, p.Content, p.AuthorID, p.CommunityID)
	return
}

// GetPostById 根据id获取帖子的详细信息
func GetPostById(id int64) (data *models.CreatePost, err error) {
	data = new(models.CreatePost)
	sqlStr := `select 
		post_id, title, content, author_id, community_id, create_time
		from post
		where post_id = ?
	`
	err = db.Get(data, sqlStr, id)
	return
}

// GetPostList 获取帖子列表
// 分页功能，page和size，page表示第几页，size表示多少条
func GetPostList(page, size int64) (posts []*models.CreatePost, err error) {
	// limit，表示查询从？到？的记录
	sqlStr := `select post_id, title, content, author_id, community_id, create_time from post limit ?,?`
	posts = make([]*models.CreatePost, 0, 2) // 不要写成make([]*models.Post, 2)
	err = db.Select(&posts, sqlStr, (page-1)*size, size)
	return
}
