package service

import (
	"errors"
	"keep_coding_blog/db"
	"keep_coding_blog/models"
)

// PostService 文章服务结构体
type PostService struct{}

// CreatePost 创建文章 (insert)
func (s *PostService) CreatePost(title, content string, userID uint, tagNames []string) (*models.Post, error) {
	// 验证数据合法性
	if title == "" || content == "" || userID == 0 {
		return nil, errors.New("title, content, and userID cannot be empty")
	}
	if len(title) > 200 {
		return nil, errors.New("title cannot be longer than 200 characters")
	}

	// 开始事务
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 检查文章标题是否已存在
	var existingPost models.Post
	if err := tx.Where("title = ?", title).First(&existingPost).Error; err == nil {
		tx.Rollback()
		return nil, errors.New("title already exists")
	}

	// 创建文章
	post := &models.Post{
		Title:   title,
		Content: content,
		UserID:  userID,
	}

	if err := tx.Create(post).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 处理标签
	if len(tagNames) > 0 {
		for _, tagName := range tagNames {
			var tag models.Tag
			// 查找或创建标签
			if err := tx.Where("name = ?", tagName).FirstOrCreate(&tag, models.Tag{Name: tagName}).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
			// 关联标签到文章
			if err := tx.Model(post).Association("Tags").Append(&tag); err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	// 重新加载文章信息，包括关联的标签
	if err := tx.Preload("Tags").Preload("User").First(post, post.ID).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return post, tx.Commit().Error
}

// GetPost 获取单个文章 (select)
func (s *PostService) GetPost(id uint) (*models.Post, error) {
	var post models.Post
	if err := db.DB.Preload("Tags").Preload("User").First(&post, id).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

// GetAllPosts 获取文章列表 (select)
func (s *PostService) GetAllPosts(page, pageSize int) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64

	// 获取总数
	if err := db.DB.Model(&models.Post{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	if err := db.DB.Preload("Tags").Preload("User").
		Offset(offset).Limit(pageSize).
		Order("created_at DESC").
		Find(&posts).Error; err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

// UpdatePost 更新文章 (update)
func (s *PostService) UpdatePost(id uint, userID uint, title, content string, tagNames []string) (*models.Post, error) {
	// 验证数据合法性
	if id == 0 || userID == 0 || title == "" || content == "" {
		return nil, errors.New("id, userID, title and content cannot be empty")
	}
	if len(title) > 200 {
		return nil, errors.New("title cannot be longer than 200 characters")
	}

	// 开始事务
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 检查文章是否存在
	var post models.Post
	if err := tx.First(&post, id).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 检查是否是文章作者
	if post.UserID != userID {
		tx.Rollback()
		return nil, errors.New("unauthorized to update this post")
	}

	// 检查标签是否存在
	if len(tagNames) > 0 {
		for _, tagName := range tagNames {
			var tag models.Tag
			if err := tx.Where("name = ?", tagName).First(&tag).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	// 更新文章基本信息
	post.Title = title
	post.Content = content
	if err := tx.Save(&post).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 清除原有标签关联
	if err := tx.Model(&post).Association("Tags").Clear(); err != nil {
		tx.Rollback()
		return nil, err
	}

	// 添加新标签
	if len(tagNames) > 0 {
		for _, tagName := range tagNames {
			var tag models.Tag
			if err := tx.Where("name = ?", tagName).FirstOrCreate(&tag, models.Tag{Name: tagName}).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
			if err := tx.Model(&post).Association("Tags").Append(&tag); err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	// 重新加载文章信息
	if err := tx.Preload("Tags").Preload("User").First(&post, post.ID).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return &post, tx.Commit().Error
}

// DeletePost 删除文章 (delete)
func (s *PostService) DeletePost(id uint, userID uint) error {
	// 验证数据合法性
	if id == 0 || userID == 0 {
		return errors.New("id, userID cannot be empty")
	}

	// 开始事务
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 检查文章是否存在
	var post models.Post
	if err := tx.First(&post, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 检查是否是文章作者
	if post.UserID != userID {
		tx.Rollback()
		return errors.New("unauthorized to delete this post")
	}

	// 清除标签关联
	if err := tx.Model(&post).Association("Tags").Clear(); err != nil {
		tx.Rollback()
		return err
	}

	// 删除文章
	if err := tx.Delete(&post).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// GetPostComments 获取文章的所有评论 (select)
func (s *PostService) GetPostComments(postID uint) ([]models.Comment, int64, error) {
	var comments []models.Comment
	var total int64

	// 获取评论
	if err := db.DB.Where("post_id = ?", postID).Find(&comments).Error; err != nil {
		return nil, 0, errors.New("failed to get post comments")
	}

	// 获取评论总数
	if err := db.DB.Model(&models.Comment{}).Where("post_id = ?", postID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return comments, total, nil
}

/*
// SearchPosts 搜索文章 (select)
func (s *PostService) SearchPosts(query string, tags []string, startTime, endTime *time.Time, page, pageSize int) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64

	db := db.DB.Model(&models.Post{})

	// 标题搜索
	if query != "" {
		db = db.Where("title ILIKE ?", "%"+query+"%")
	}

	// 标签搜索
	if len(tags) > 0 {
		db = db.Joins("JOIN post_tags ON posts.id = post_tags.post_id").
			Joins("JOIN tags ON post_tags.tag_id = tags.id").
			Where("tags.name IN ?", tags)
	}

	// 时间范围搜索
	if startTime != nil {
		db = db.Where("created_at >= ?", startTime)
	}
	if endTime != nil {
		db = db.Where("created_at <= ?", endTime)
	}

	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	if err := db.Preload("Tags").
		Preload("User").
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&posts).Error; err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

// GetPostTags 获取文章标签 (select)
func (s *PostService) GetPostTags(postID uint) ([]models.Tag, error) {
	var tags []models.Tag
	if err := db.DB.Joins("JOIN post_tags ON posts.id = post_tags.post_id").
		Joins("JOIN tags ON post_tags.tag_id = tags.id").
		Where("posts.id = ?", postID).
		Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

*/
