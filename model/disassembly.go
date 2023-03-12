package model

import (
	"gorm.io/gorm"
)

type Disassembly struct {
	BasicModel
	Name       *string  //名称
	ProjectID  *int     //项目id，外键
	SuperiorID *int     //上级id
	Level      *int     //层级
	Weight     *float64 //权重
}

// TableName 修改表名
func (*Disassembly) TableName() string {
	return "disassembly"
}

func (d *Disassembly) BeforeDelete(tx *gorm.DB) error {
	if d.ID > 0 {
		//如果有删除人的id，则记录下来
		//if d.Deleter != nil && *d.Deleter > 0 {
		//	err := tx.Model(&Disassembly{}).Where("id = ?", d.ID).
		//		Update("deleter", d.Deleter).Error
		//	if err != nil {
		//		return err
		//	}
		//}
		//删除相关的子表记录
		//err = tx.Model(&Progress{}).Where("disassembly_id = ?", d.ID).
		//	Updates(map[string]any{
		//		"deleted_at": time.Now(),
		//		"deleter":    d.Deleter,
		//	}).Error
		//if err != nil {
		//	return err
		//}
	}
	return nil
}
