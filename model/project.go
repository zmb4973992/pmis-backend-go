package model

import (
	"gorm.io/gorm"
	"time"
)

type Project struct {
	BaseModel
	ProjectCode          *string
	ProjectFullName      *string
	ProjectShortName     *string
	Country              *string
	Province             *string
	ProjectType          *string
	SegmentedProjectType *string //细分的项目类型
	Amount               *float64
	Currency             *string
	ExchangeRate         *float64
	ProjectStatus        *string
	OurSignatory         *string    //我方签约主体
	ConstructionPeriod   *int       //工期，天
	SigningDate          *time.Time `gorm:"type:date"` //签约日期
	EffectiveDate        *time.Time `gorm:"type:date"` //生效日期
	CommissioningDate    *time.Time `gorm:"type:date"` //调试日期
	DepartmentID         *int
	RelatedPartyID       *int
	Task                 *string //工作内容
	//外键
	//Disassemblies               []Disassembly                `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	//ActualReceiptAndPayments    []ActualReceiptAndPayment    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	//PlannedReceiptAndPayments   []PlannedReceiptAndPayment   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	//PredictedReceiptAndPayments []PredictedReceiptAndPayment `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	//ProjectAndUsers             []ProjectAndUser             `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
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
