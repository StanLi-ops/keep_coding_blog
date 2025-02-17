package service

import (
	"errors"
	"keep_coding_blog/db"
	"keep_coding_blog/models"
	"time"
)

// PostService 文章服务结构体
type PostService struct{}

// CreatePost 创建文章
func (s *PostService) CreatePost(title, content string, userID uint, tagNames []string) (*models.Post, error) {
	// 创建文章
	post := &models.Post{
		Title:   title,
		Content: content,
		UserID:  userID,
	}

	// 开启事务
	tx := db.DB.Begin()

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

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// 重新加载文章信息，包括关联的标签
	if err := db.DB.Preload("Tags").Preload("User").First(post, post.ID).Error; err != nil {
		return nil, err
	}

	return post, nil
}

// GetPost 获取单个文章
func (s *PostService) GetPost(id uint) (*models.Post, error) {
	var post models.Post
	if err := db.DB.Preload("Tags").Preload("User").First(&post, id).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

// GetPosts 获取文章列表
func (s *PostService) GetPosts(page, pageSize int) ([]models.Post, int64, error) {
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

// UpdatePost 更新文章
func (s *PostService) UpdatePost(id uint, userID uint, title, content string, tagNames []string) (*models.Post, error) {
	var post models.Post
	if err := db.DB.First(&post, id).Error; err != nil {
		return nil, err
	}

	// 检查是否是文章作者
	if post.UserID != userID {
		return nil, errors.New("unauthorized to update this post")
	}

	// 开启事务
	tx := db.DB.Begin()

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

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// 重新加载文章信息
	if err := db.DB.Preload("Tags").Preload("User").First(&post, post.ID).Error; err != nil {
		return nil, err
	}

	return &post, nil
}

// DeletePost 删除文章
func (s *PostService) DeletePost(id uint, userID uint) error {
	var post models.Post
	if err := db.DB.First(&post, id).Error; err != nil {
		return err
	}

	// 检查是否是文章作者
	if post.UserID != userID {
		return errors.New("unauthorized to delete this post")
	}

	// 开启事务
	tx := db.DB.Begin()

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

// SearchPosts 搜索文章
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

// GetAllTags 获取所有标签
func (s *PostService) GetAllTags() ([]models.Tag, error) {
	var tags []models.Tag
	if err := db.DB.Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}
