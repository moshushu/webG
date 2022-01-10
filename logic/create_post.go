package logic

import (
	"web1/dao/mysql"
	"web1/models"
	"web1/pkg/snowflake"

	"go.uber.org/zap"
)

// CreatePost 创建帖子
func CreatePost(p *models.CreatePost) (err error) {
	// 1、生成post ID
	// 利用雪花算法生成 post ID
	p.ID = snowflake.GenID()
	// 2、保存到数据库
	return mysql.CreatePost(p)
}

// GetPostById 根据id获取帖子的详细信息
func GetPostById(id int64) (data *models.ApiPostDetail, err error) {
	// 查询帖子的信息
	post, err := mysql.GetPostById(id)
	if err != nil {
		zap.L().Error("mysql.GetPostById(id) failed", zap.Error(err))
		return
	}
	// 根据作者id查询作者信息
	// 将post.AuthorID传入查询出username，作者的名字
	user, err := mysql.GetUserById(post.AuthorID)
	if err != nil {
		zap.L().Error("mysql.GetUserById(post.AuthorID) failed", zap.Error(err))
		return
	}
	//  根据社区的id查询社区的详细信息（帖子属于哪个分区）
	community, err := mysql.GetCommunityDetailByID(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityDetailByID(post.CommunityID) failed", zap.Error(err))
		return
	}
	// 将查出来的所有数据整合到一起（接口数据拼接），ApiPostDetail类型返回
	data = &models.ApiPostDetail{
		AuthorName:      user.Username, // 返回作者的姓名
		CreatePost:      post,          // 返回帖子的信息
		CommunityDetail: community,     // 返回帖子的分区
	}
	return
}

// GetPostList 获取帖子列表
func GetPostList(page, size int64) (data []*models.ApiPostDetail, err error) {
	// 1、查找获取帖子的信息
	posts, err := mysql.GetPostList(page, size)
	if err != nil {
		zap.L().Error("mysql.GetPostList(page,size) failed", zap.Error(err))
		// 数据库查询出错
		return nil, err
	}
	// 对data进行初始化（申请内存），切片的容量由查询到的记录数量来决定len(posts)
	// data 用来存放查询到的帖子记录
	data = make([]*models.ApiPostDetail, 0, len(posts))
	// 由于数据由多条，利用循环将记录分别添加到data切片中
	for _, post := range posts {
		// 根据作者id查询作者信息
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById failed", zap.Error(err))
			// 这里不使用return，是因为不能因为一个查不到就停止所有的查询
			// 直接查询下一条
			continue
		}
		// 根据社区id查询社区详细信息
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID(post.CommunityID) failed", zap.Error(err))
			// 这里不使用return，是因为不能因为一个查不到就停止所有的查询
			// 直接查询下一条
			continue
		}
		// 将得到的信息进行拼接
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			CreatePost:      post,
			CommunityDetail: community,
		}
		// 将数据添加到data切片中
		data = append(data, postDetail)
	}
	return
}
