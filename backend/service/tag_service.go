package service

import (
	"errors"
	"keep_coding_blog/db"
	"keep_coding_blog/models"
)

// TagService 标签服务结构体
type TagService struct{}

// CreateTag 创建标签 (insert)
func (s *TagService) CreateTag(name string) (*models.Tag, error) {
	// 验证数据合法性
	if name == "" {
		return nil, errors.New("tag name cannot be empty")
	}
	if len(name) > 50 {
		return nil, errors.New("tag name must be less than 50 characters")
	}

	// 使用事务处理
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 检查标签是否已存在
	var existingTag models.Tag
	if err := tx.Where("name = ?", name).First(&existingTag).Error; err == nil {
		tx.Rollback()
		return nil, errors.New("tag already exists")
	}

	// 创建标签
	tag := &models.Tag{Name: name}

	if err := tx.Create(tag).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return tag, tx.Commit().Error
}

// GetTag 获取单个标签 (select)
func (s *TagService) GetTag(id uint) (*models.Tag, error) {
	var tag models.Tag
	if err := db.DB.First(&tag, id).Error; err != nil {
		return nil, err
	}
	return &tag, nil
}

// GetAllTags 获取所有标签 (select)
func (s *TagService) GetAllTags() ([]models.Tag, error) {
	var tags []models.Tag
	if err := db.DB.Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

// UpdateTag 更新标签 (update)
func (s *TagService) UpdateTag(id uint, name string) (*models.Tag, error) {
	// 验证数据合法性
	if name == "" {
		return nil, errors.New("tag name cannot be empty")
	}
	if len(name) > 50 {
		return nil, errors.New("tag name must be less than 50 characters")
	}

	// 使用事务处理
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 检查标签是否存在
	var existingTag models.Tag
	if err := tx.Where("id = ?", id).First(&existingTag).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("tag not found")
	}

	// 更新标签
	existingTag.Name = name

	if err := tx.Save(&existingTag).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return &existingTag, tx.Commit().Error
}

// DeleteTag 删除标签 (delete)
func (s *TagService) DeleteTag(id uint) error {
	// 验证数据合法性
	if id == 0 {
		return errors.New("tag id cannot be 0")
	}

	// 使用事务处理
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 检查标签是否存在
	var existingTag models.Tag
	if err := tx.Where("id = ?", id).First(&existingTag).Error; err != nil {
		tx.Rollback()
		return errors.New("tag not found")
	}

	// 删除标签
	if err := tx.Delete(&existingTag).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
