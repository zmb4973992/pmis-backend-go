package old_pmis

import (
	"fmt"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
)

type project struct {
	Id   int64
	Code string `gorm:"column:项目号"`
}

type disassembly struct {
	Id         int64
	Name       string  `gorm:"column:名称"`
	ProjectId  int64   `gorm:"column:项目id"`
	SuperiorId int64   `gorm:"column:上级id"`
	Level      int     `gorm:"column:层级"`
	Weight     float64 `gorm:"column:权重"`
}

func importDisassembly(userId int64) error {
	fmt.Println("★★★★★开始处理拆解情况记录......★★★★★")

	var oldDisassembliesOfLevel1 []disassembly
	err := global.DBForOldPmis.Table("拆解情况").
		Where("层级 = ?", 1).
		Find(&oldDisassembliesOfLevel1).Error
	if err != nil {
		return err
	}

	for i := range oldDisassembliesOfLevel1 {
		var oldProject project
		err = global.DBForOldPmis.Table("项目").
			Where("id = ?", oldDisassembliesOfLevel1[i].ProjectId).
			First(&oldProject).Error
		if err != nil {
			global.SugaredLogger.Errorln(err)
			continue
		}

		var count1 int64
		global.DB.Model(&model.Project{}).
			Where("code = ?", oldProject.Code).
			Count(&count1)
		if count1 == 0 {
			continue
		}

		var newProject model.Project
		err = global.DB.Where("code = ?", oldProject.Code).
			First(&newProject).Error
		if err != nil {
			global.SugaredLogger.Errorln(err)
			continue
		}

		var newDisassemblyOfLevel1 model.Disassembly
		err = global.DB.
			Where("project_id = ?", newProject.Id).
			Where("superior_id is null").
			Where("level = ?", 1).
			First(&newDisassemblyOfLevel1).Error
		if err != nil {
			global.SugaredLogger.Errorln(err)
			continue
		}

		var oldDisassembliesOfLevel2 []disassembly
		err = global.DBForOldPmis.Table("拆解情况").
			Where("上级id = ?", oldDisassembliesOfLevel1[i].Id).
			Find(&oldDisassembliesOfLevel2).Error
		if err != nil {
			global.SugaredLogger.Errorln(err)
			continue
		}

		for j := range oldDisassembliesOfLevel2 {
			var newDisassemblyOfLevel2 model.Disassembly
			newDisassemblyOfLevel2.Creator = &userId
			newDisassemblyOfLevel2.ProjectId = &newProject.Id
			newDisassemblyOfLevel2.SuperiorId = &newDisassemblyOfLevel1.Id
			newDisassemblyOfLevel2.Name = &oldDisassembliesOfLevel2[j].Name
			newDisassemblyOfLevel2.Weight = &oldDisassembliesOfLevel2[j].Weight
			newDisassemblyOfLevel2.Level = model.IntToPointer(2)
			newDisassemblyOfLevel2.ImportedIdFromOldPmis = &oldDisassembliesOfLevel2[j].Id

			err = global.DB.
				Where("project_id = ?", newProject.Id).
				Where("superior_id = ?", newDisassemblyOfLevel1.Id).
				Where("name = ?", oldDisassembliesOfLevel2[j].Name).
				Where("weight = ?", oldDisassembliesOfLevel2[j].Weight).
				Where("level = 2").
				Where("imported_id_from_old_pmis = ?", oldDisassembliesOfLevel2[j].Id).
				FirstOrCreate(&newDisassemblyOfLevel2).Error

			if err != nil {
				global.SugaredLogger.Errorln(err)
				continue
			}

			var oldDisassembliesOfLevel3 []disassembly
			err = global.DBForOldPmis.Table("拆解情况").
				Where("上级id = ?", oldDisassembliesOfLevel2[j].Id).
				Find(&oldDisassembliesOfLevel3).Error
			if err != nil {
				global.SugaredLogger.Errorln(err)
				continue
			}

			for k := range oldDisassembliesOfLevel3 {
				var newDisassemblyOfLevel3 model.Disassembly
				newDisassemblyOfLevel3.Creator = &userId
				newDisassemblyOfLevel3.ProjectId = &newProject.Id
				newDisassemblyOfLevel3.SuperiorId = &newDisassemblyOfLevel2.Id
				newDisassemblyOfLevel3.Level = model.IntToPointer(3)
				newDisassemblyOfLevel3.Weight = &oldDisassembliesOfLevel3[k].Weight
				newDisassemblyOfLevel3.Name = &oldDisassembliesOfLevel3[k].Name
				newDisassemblyOfLevel3.ImportedIdFromOldPmis = &oldDisassembliesOfLevel3[k].Id

				err = global.DB.Model(&model.Disassembly{}).
					Where("project_id = ?", newProject.Id).
					Where("superior_id = ?", newDisassemblyOfLevel2.Id).
					Where("name = ?", oldDisassembliesOfLevel3[k].Name).
					Where("weight = ?", oldDisassembliesOfLevel3[k].Weight).
					Where("level = ?", 3).
					Where("imported_id_from_old_pmis = ?", oldDisassembliesOfLevel3[k].Id).
					FirstOrCreate(&newDisassemblyOfLevel3).Error
				if err != nil {
					global.SugaredLogger.Errorln(err)
					continue
				}

				var oldDisassembliesOfLevel4 []disassembly
				err = global.DBForOldPmis.Table("拆解情况").
					Where("上级id = ?", oldDisassembliesOfLevel3[k].Id).
					Find(&oldDisassembliesOfLevel4).Error
				if err != nil {
					global.SugaredLogger.Errorln(err)
					continue
				}

				for l := range oldDisassembliesOfLevel4 {
					var newDisassemblyOfLevel4 model.Disassembly
					newDisassemblyOfLevel4.Creator = &userId
					newDisassemblyOfLevel4.ProjectId = &newProject.Id
					newDisassemblyOfLevel4.SuperiorId = &newDisassemblyOfLevel3.Id
					newDisassemblyOfLevel4.Level = model.IntToPointer(4)
					newDisassemblyOfLevel4.Weight = &oldDisassembliesOfLevel4[l].Weight
					newDisassemblyOfLevel4.Name = &oldDisassembliesOfLevel4[l].Name
					newDisassemblyOfLevel4.ImportedIdFromOldPmis = &oldDisassembliesOfLevel4[l].Id

					err = global.DB.Model(&model.Disassembly{}).
						Where("project_id = ?", newProject.Id).
						Where("superior_id = ?", newDisassemblyOfLevel3.Id).
						Where("name = ?", oldDisassembliesOfLevel4[l].Name).
						Where("weight = ?", oldDisassembliesOfLevel4[l].Weight).
						Where("level = ?", 4).
						Where("imported_id_from_old_pmis = ?", oldDisassembliesOfLevel4[l].Id).
						FirstOrCreate(&newDisassemblyOfLevel4).Error
					if err != nil {
						global.SugaredLogger.Errorln(err)
						continue
					}

					var oldDisassembliesOfLevel5 []disassembly
					err = global.DBForOldPmis.Table("拆解情况").
						Where("上级id = ?", oldDisassembliesOfLevel4[l].Id).
						Find(&oldDisassembliesOfLevel5).Error
					if err != nil {
						global.SugaredLogger.Errorln(err)
						continue
					}

					for m := range oldDisassembliesOfLevel5 {
						var newDisassemblyOfLevel5 model.Disassembly
						newDisassemblyOfLevel5.Creator = &userId
						newDisassemblyOfLevel5.ProjectId = &newProject.Id
						newDisassemblyOfLevel5.SuperiorId = &newDisassemblyOfLevel4.Id
						newDisassemblyOfLevel5.Level = model.IntToPointer(5)
						newDisassemblyOfLevel5.Name = &oldDisassembliesOfLevel5[m].Name
						newDisassemblyOfLevel5.Weight = &oldDisassembliesOfLevel5[m].Weight
						newDisassemblyOfLevel5.ImportedIdFromOldPmis = &oldDisassembliesOfLevel5[m].Id

						global.DB.
							Where("project_id = ?", newProject.Id).
							Where("superior_id = ?", newDisassemblyOfLevel4.Id).
							Where("name = ?", oldDisassembliesOfLevel5[m].Name).
							Where("weight = ?", oldDisassembliesOfLevel5[m].Weight).
							Where("level = ?", 5).
							Where("imported_id_from_old_pmis = ?", oldDisassembliesOfLevel5[m].Id).
							FirstOrCreate(&newDisassemblyOfLevel5)
						if err != nil {
							global.SugaredLogger.Errorln(err)
							continue
						}
					}
				}
			}
		}
	}
	fmt.Println("★★★★★拆解情况记录处理完成......★★★★★")

	return nil
}
