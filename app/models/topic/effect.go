package topic

import (
	"gin_bbs/database"

	"github.com/lexkong/log"
)

// Create -
func (t *Topic) Create() (err error) {
	if err = database.DB.Create(&t).Error; err != nil {
		log.Warnf("topic 创建失败: %v", err)
		return err
	}

	return nil
}

// Update -
func (t *Topic) Update() (err error) {
	if err = database.DB.Save(&t).Error; err != nil {
		log.Warnf("topic 更新失败: %v", err)
		return err
	}

	return nil
}

// Delete -
func Delete(id int) (err error) {
	topic := &Topic{}
	topic.ID = uint(id)

	if err = database.DB.Delete(&topic).Error; err != nil {
		log.Warnf("topic 删除失败: %v", err)
		return err
	}

	return nil
}
