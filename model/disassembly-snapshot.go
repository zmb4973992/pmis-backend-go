package model

import (
	"gorm.io/gorm"
	"time"
)

// deprecated
type DisassemblySnapshot struct {
	BaseModel
	Name               *string    //名称
	Date               *time.Time `gorm:"type:datetime"` //添加记录的日期
	IDWithDate         *string    //带日期的拆解情况id
	ProjectID          *int       //项目id，外键
	SuperiorIDWithDate *string    //带日期的上级id
	Level              *int       //层级
	Weight             *float64   //权重
	//WorkProgresses     []Progress `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// TableName 修改表名
func (*DisassemblySnapshot) TableName() string {
	return "disassembly_snapshot"
}

func (d *DisassemblySnapshot) BeforeDelete(tx *gorm.DB) error {
	if d.ID > 0 {
		//如果有删除人的id，则记录下来
		if d.Deleter != nil && *d.Deleter > 0 {
			err := tx.Model(&DisassemblySnapshot{}).Where("id = ?", d.ID).
				Update("deleter", d.Deleter).Error
			if err != nil {
				return err
			}
		}
		//删除相关的子表记录
		err = tx.Model(&Progress{}).Where("disassembly_id = ?", d.ID).
			Updates(map[string]any{
				"deleted_at": time.Now(),
				"deleter":    d.Deleter,
			}).Error
		if err != nil {
			return err
		}
	}
	return nil
}
