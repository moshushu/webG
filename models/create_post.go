package models

import "time"

// 引用一个Go内存对齐的概念
// 将同一类型的tag放到一块，有利于缩小内存

// 创建帖子的模型
type CreatePost struct {
	// `json:"id,sting"`,tag结构体中添加string来告诉json包从字符串中解析相应字段的数据
	// 简单来说：请求中带着id字段的数值过来，但由于可能数值过大，
	// 超过了json的int类型的最大值（负2的53次方减1，到2的53次方减1）
	// 因此需要将id的数值转成string类型来避免id数值失真（丢失）
	// 在id字段的tag中添加string，来告诉json，在解析请求中的id字段时需要去字符串中
	// 找到对应的id字段，再将其转化成int64类型存储再结构体的ID字段中
	ID          int64     `json:"id,string" db:"post_id"`
	AuthorID    int64     `json:"author_id" db:"author_id"`
	CommunityID int64     `json:"community_id" db:"community_id" binding:"required"`
	Content     string    `json:"content" db:"content" binding:"required"`
	Title       string    `json:"title" db:"title" binding:"required"`
	Status      int32     `json:"status" db:"status"`
	CreateTime  time.Time `json:"create_time" db:"create_time"`
}

// ApiPostDetail 帖子详情接口的结构体
type ApiPostDetail struct {
	// 作者名称
	AuthorName       string                    `json:"author_name"`
	*CreatePost                                // 嵌入帖子结构体
	*CommunityDetail `json:"community_detail"` // 嵌入社区分类信息
}
