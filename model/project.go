package model

import (
	"gorm.io/gorm"
	"time"
)

type Project struct {
	BaseModel
	Code               *string
	Name               *string
	Country            *int //见dictionary_item
	Type               *int //见dictionary_item
	SegmentedType      *int //细分的项目类型，见dictionary_item
	Amount             *float64
	Currency           *int //见dictionary_item
	ExchangeRate       *float64
	Status             *int       //见dictionary_item
	OurSignatory       *int       //我方签约主体，见dictionary_item
	ConstructionPeriod *int       //工期，天
	SigningDate        *time.Time `gorm:"type:date"` //签约日期
	EffectiveDate      *time.Time `gorm:"type:date"` //生效日期
	CommissioningDate  *time.Time `gorm:"type:date"` //调试日期
	DepartmentID       *int       //见department
	RelatedPartyID     *int       //见related_party
	Content            *string    //工作内容
}

// TableName 将表名改为project
func (*Project) TableName() string {
	return "project"
}

func (d *Project) BeforeDelete(tx *gorm.DB) error {
	if d.ID > 0 {
		//如果有删除人的id，则记录下来
		if d.Deleter != nil && *d.Deleter > 0 {
			err := tx.Model(&Project{}).Where("id = ?", d.ID).
				Update("deleter", d.Deleter).Error
			if err != nil {
				return err
			}
		}
		//删除相关的子表记录
		err = tx.Model(&Disassembly{}).Where("project_id = ?", d.ID).
			Updates(map[string]any{
				"deleted_at": time.Now(),
				"deleter":    d.Deleter,
			}).Error
		if err != nil {
			return err
		}

		err = tx.Model(&PlannedReceiptAndPayment{}).Where("project_id = ?", d.ID).
			Updates(map[string]any{
				"deleted_at": time.Now(),
				"deleter":    d.Deleter,
			}).Error
		if err != nil {
			return err
		}

		err = tx.Model(&PredictedReceiptAndPayment{}).Where("project_id = ?", d.ID).
			Updates(map[string]any{
				"deleted_at": time.Now(),
				"deleter":    d.Deleter,
			}).Error
		if err != nil {
			return err
		}

		err = tx.Model(&ProjectAndUser{}).Where("project_id = ?", d.ID).
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
