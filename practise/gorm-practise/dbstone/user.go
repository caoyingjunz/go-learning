package dbstone

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"

	"go-learning/practise/gorm-practise/models"
)

var (
	RecordNotUpdated = errors.New("record not updated")
)

type UserInterface interface {
	Get(ctx context.Context, name string, age int) (*models.User, error)
	Create(ctx context.Context, user *models.User) error
	Update(ctx context.Context, name string, resourceVersion int64, updates map[string]interface{}) error
	Delete(ctx context.Context, uid int64) error

	List(ctx context.Context, name string) ([]models.User, error)

	OptimisticUpdate(ctx context.Context, name string, resourceVersion int64, updates map[string]interface{}) error
}

func NewUserDB() UserInterface {
	return &UserDB{
		dbstone: DB,
	}
}

type UserDB struct {
	dbstone *gorm.DB
}

func (u *UserDB) Get(ctx context.Context, name string, age int) (*models.User, error) {
	var user models.User
	// Frist 获取第一个
	// Find 获取满足条件，如果只有一个返回，返回最后一个
	if err := u.dbstone.Where("name = ? and age = ?", name, age).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UserDB) List(ctx context.Context, name string) ([]models.User, error) {
	var users []models.User
	if err := u.dbstone.Where("name = ?", name).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (u *UserDB) Create(ctx context.Context, user *models.User) error {
	return u.dbstone.Create(&user).Error
}

func (u *UserDB) Update(ctx context.Context, name string, resourceVersion int64, updates map[string]interface{}) error {
	updates["resource_version"] = resourceVersion + 1
	updates["gmt_modified"] = time.Now()

	return u.dbstone.Model(&models.User{}).
		Where("name = ? and resource_version = ?", name, resourceVersion).
		Updates(updates).Error
}

func (u *UserDB) Delete(ctx context.Context, uid int64) error {
	return u.dbstone.Where("id = ?", uid).Delete(&models.User{}).Error
}

func (u *UserDB) GetRawUsers(names []string) (user []models.User, err error) {
	// 如果需要，强制使用索引
	if err = u.dbstone.Raw("select * from user force index(idx_user_name_age) where name in (?)", names).Scan(&user).Error; err != nil {
		return
	}
	return
}

// OptimisticUpdate 自定义乐观锁
func (u *UserDB) OptimisticUpdate(ctx context.Context, name string, resourceVersion int64, updates map[string]interface{}) error {
	updates["resource_version"] = resourceVersion + 1
	updates["gmt_modified"] = time.Now()

	uc := u.dbstone.Model(&models.User{}).
		Where("name = ? and resource_version = ?", name, resourceVersion).Update(updates)

	if err := uc.Error; err != nil {
		return err
	}
	// RowsAffected 为 0 的时候，说明未匹配到更新
	if uc.RowsAffected == 0 {
		return RecordNotUpdated
	}

	fmt.Println(uc.RowsAffected)
	return nil
}
