package logic

import (
	"web_app/dao/mysql"
	"web_app/dao/redis"
	"web_app/models"
	"web_app/pkg/snowflake"

	"go.uber.org/zap"
)

func CreatePost(p *models.Post) (err error) {
	//1.生成post id
	p.ID = snowflake.GenID()

	// //2.保存到数据库
	err = mysql.CreatePost(p)
	if err != nil {
		return err
	}
	err = redis.CreatePost(p.ID, p.CommunityID)
	//3.返回
	return
}

// GetPostByID根据帖子id获取帖子详情
func GetPostByID(pid int64) (data *models.ApiPostDetail, err error) {
	//查询并组合接口想用的数据
	post, err := mysql.GetPostByID(pid)
	if err != nil {
		zap.L().Error("mysql.GetPostByID(pid) failed",
			zap.Int64("pid", pid),
			zap.Error(err))
		return
	}
	//根据作者id查询作者信息
	user, err := mysql.GetUserByID(post.AuthorID)
	if err != nil {
		zap.L().Error("mysql.GetUserByID failed",
			zap.Int64("author_id", post.AuthorID),
			zap.Error(err))
		return
	}
	votedata, err := redis.GetPostVoteNum(pid)
	//根据社区id社区详细信息
	community, err := mysql.GetCommunityDetailByID(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityDetailByID failed",
			zap.Int64("community_id", community.ID),
			zap.Error(err))
		return
	}
	//数据拼接
	data = &models.ApiPostDetail{
		AuthorName:      user.Username,
		Post:            post,
		CommunityDetail: community,
		VoteNumber:      votedata,
	}
	return
}

// 获取帖子列表
func GetPostList(page, size int64) (data []*models.ApiPostDetail, err error) {
	posts, err := mysql.GetPostList(page, size)
	if err != nil {
		return nil, err
	}
	//有多少帖子就有多少帖子详情
	data = make([]*models.ApiPostDetail, 0, len(posts))
	for _, post := range posts {
		//根据作者id查询作者信息
		user, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserByID failed",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(err))
			continue
		}
		//根据社区id社区详细信息
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID failed",
				zap.Int64("community_id", community.ID),
				zap.Error(err))
			continue
		}

		postdetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			Post:            post,
			CommunityDetail: community,
		}

		data = append(data, postdetail)
	}
	return
}

func GetPosts(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	//去redis查询帖子id列表
	//posts, err := mysql.GetPostList(page, size)
	//if err != nil {
	//	return nil, err
	//}
	//去redis查询id列表
	idlist, err := redis.GetPostIDsInOrder(p)
	if err != nil {
		return
	}
	if len(idlist) == 0 {
		zap.L().Warn(" redis.GetPostIDsInOrder(p) return no idlist")
		return
	}
	zap.L().Debug("GetPosts", zap.Any("idlist", idlist))
	//去mysql根据id查询帖子详情
	//返回的数据按照给定id的顺序返回
	posts, err := mysql.GetPostListByIds(idlist)
	if err != nil {
		return
	}
	//提前查询好每篇帖子的投票数
	votedata, err := redis.GetPostVoteData(idlist)
	if err != nil {
		return
	}

	//将帖子的作者id及分区信息查询出来填充到帖子中
	zap.L().Debug("GetPosts", zap.Any("posts", posts))
	data = make([]*models.ApiPostDetail, 0, len(posts))
	for idx, post := range posts {
		//根据作者id查询作者信息
		user, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserByID failed",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(err))
			continue
		}
		//根据社区id社区详细信息
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID failed",
				zap.Int64("community_id", community.ID),
				zap.Error(err))
			continue
		}

		postdetail := &models.ApiPostDetail{
			VoteNumber:      votedata[idx],
			AuthorName:      user.Username,
			Post:            post,
			CommunityDetail: community,
		}

		data = append(data, postdetail)
	}
	return
}

func GetCommunityPostList(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	//去redis查询id列表
	idlist, err := redis.GetCommunityPostIDsInOrder(p)
	if err != nil {
		return
	}
	if len(idlist) == 0 {
		zap.L().Warn(" redis.GetCommunityPostIDsInOrder(p) return no idlist")
		return
	}
	zap.L().Debug("GetCommunityPostIDsInOrder", zap.Any("idlist", idlist))
	//去mysql根据id查询帖子详情
	//返回的数据按照给定id的顺序返回
	posts, err := mysql.GetPostListByIds(idlist)
	if err != nil {
		return
	}
	//提前查询好每篇帖子的投票数
	votedata, err := redis.GetPostVoteData(idlist)
	if err != nil {
		return
	}

	//将帖子的作者id及分区信息查询出来填充到帖子中
	zap.L().Debug("GetPosts", zap.Any("posts", posts))
	data = make([]*models.ApiPostDetail, 0, len(posts))
	for idx, post := range posts {
		//根据作者id查询作者信息
		user, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserByID failed",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(err))
			continue
		}
		//根据社区id社区详细信息
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID failed",
				zap.Int64("community_id", community.ID),
				zap.Error(err))
			continue
		}

		postdetail := &models.ApiPostDetail{
			VoteNumber:      votedata[idx],
			AuthorName:      user.Username,
			Post:            post,
			CommunityDetail: community,
		}

		data = append(data, postdetail)
	}
	return
}

// 将两个查询逻辑合二为一
func GetPostListNew(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	//根据请求参数的不同调用不同方法
	if p.CommunityID == 0 {
		//查所有的帖子
		data, err = GetPosts(p)
	} else {
		//根据社区id查询
		data, err = GetCommunityPostList(p)
	}
	if err != nil {
		zap.L().Error("GetPostListNew failed", zap.Error(err))
		return nil, err
	}
	return
}
