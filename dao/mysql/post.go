package mysql

import (
	"database/sql"
	"strings"
	"web_app/models"

	"github.com/jmoiron/sqlx"

	"go.uber.org/zap"
)

func CreatePost(p *models.Post) (err error) {
	sqlStr := `insert into post (post_id, title, content, author_id, community_id)
			  values(?,?,?,?,?)`
	_, err = db.Exec(sqlStr, p.ID, p.Title, p.Content, p.AuthorID, p.CommunityID)
	return err
}

// 根据id查询单个帖子数据
func GetPostByID(pid int64) (postDetail *models.Post, err error) {
	sqlStr := `select post_id, title, content,author_id, community_id, create_time 
			   from post where post_id = ?`
	postDetail = new(models.Post)
	if err := db.Get(postDetail, sqlStr, pid); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("There is no Postdetail in db")
			err = ErrorInvalidPid
		}
	}
	return
}

// 获取帖子列表
func GetPostList(page, size int64) (posts []*models.Post, err error) {
	sqlStr := `select post_id, title, content,author_id, community_id, create_time 
			   from post 
			   ORDER BY create_time
			   DESC
			   limit ?,?
			   `
	posts = make([]*models.Post, 0, 2)
	if err := db.Select(&posts, sqlStr, (page-1)*size, size); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("There is no Postlist in db")
			err = ErrorInvalidPid
		}
	}
	return
}

// 根据给定的id列表查询帖子数据
func GetPostListByIds(idlist []string) (postList []*models.Post, err error) {
	sqlStr := `select post_id, title, content, author_id, community_id, create_time, update_time
			   from post
			   where post_id in(?)
			   order by  FIND_IN_SET(post_id,?)`
	//sql.In返回带？bindvar的查询语句，使用rebind()重新綁定
	query, args, err := sqlx.In(sqlStr, idlist, strings.Join(idlist, ","))
	if err != nil {
		return nil, err
	}
	query = db.Rebind(query)
	err = db.Select(&postList, query, args...)
	return
}
