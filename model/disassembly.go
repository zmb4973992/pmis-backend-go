package model

import (
	"fmt"
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
	//删除相关的子表记录
	//先find，再delete，可以激活相关的钩子函数
	fmt.Println("拆解id：", d.ID)
	var progresses []Progress
	err = tx.Where("disassembly_id = ?", d.ID).
		Find(&progresses).Delete(&progresses).Error

	if err != nil {
		return err
	}

	return nil
}
