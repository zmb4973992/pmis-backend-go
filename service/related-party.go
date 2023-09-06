package service

import (
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/serializer/list"
	"pmis-backend-go/util"
	"strconv"
	"strings"
)

//以下为入参
//有些字段不用json tag，因为不从前端读取，而是在controller中处理

type RelatedPartyGet struct {
	ID int64
}

type RelatedPartyCreate struct {
	UserID       int64
	LastModifier int64

	Name                    string `json:"name,omitempty"`
	EnglishName             string `json:"english_name,omitempty"`
	Address                 string `json:"address,omitempty"`
	UniformSocialCreditCode string `json:"uniform_social_credit_code,omitempty"` //统一社会信用代码
	Telephone               string `json:"telephone,omitempty"`
	Remarks                 string `json:"remarks,omitempty"`
	FileIDs                 string `json:"file_ids,omitempty"`
	ImportedOriginalName    string `json:"imported_original_name,omitempty"`
}

//指针字段是为了区分入参为空或0与没有入参的情况，做到分别处理，通常用于update
//如果指针字段为空或0，那么数据库相应字段会改为null；
//如果指针字段没传，那么数据库不会修改该字段

type RelatedPartyUpdate struct {
	UserID int64
	ID     int64

	Name                    *string `json:"name"`
	EnglishName             *string `json:"english_name"`
	Address                 *string `json:"address"`
	UniformSocialCreditCode *string `json:"uniform_social_credit_code"` //统一社会信用代码
	Telephone               *string `json:"telephone"`
	Remarks                 *string `json:"remarks"`
	FileIDs                 *string `json:"file_ids"`
}

type RelatedPartyDelete struct {
	ID int64
}

type RelatedPartyGetList struct {
	list.Input

	NameInclude        string `json:"name_include,omitempty"`
	EnglishNameInclude string `json:"english_name_include,omitempty"`
}

//以下为出参

type RelatedPartyOutput struct {
	Creator      *int64 `json:"creator"`
	LastModifier *int64 `json:"last_modifier"`
	ID           int64  `json:"id"`

	Name                    *string      `json:"name"`
	EnglishName             *string      `json:"english_name"`
	Address                 *string      `json:"address"`
	UniformSocialCreditCode *string      `json:"uniform_social_credit_code"` //统一社会信用代码
	Telephone               *string      `json:"telephone"`
	Remarks                 *string      `json:"remarks"`
	FileIDs                 *string      `json:"-"`
	FilesExternal           []FileOutput `json:"files" gorm:"-"`
}

func (r *RelatedPartyGet) Get() (output *RelatedPartyOutput, errCode int) {
	err := global.DB.Model(&model.RelatedParty{}).
		Where("id = ?", r.ID).
		First(&output).Error
	if err != nil {
		return nil, util.ErrorRecordNotFound
	}

	//查文件信息
	if output.FileIDs != nil {
		tempFileIDs := strings.Split(*output.FileIDs, ",")
		var fileIDs []int64
		for i := range tempFileIDs {
			fileID, err1 := strconv.ParseInt(tempFileIDs[i], 10, 64)
			if err1 != nil {
				continue
			}
			fileIDs = append(fileIDs, fileID)
		}

		var records []FileOutput
		global.DB.Model(&model.File{}).
			Where("id in ?", fileIDs).
			Find(&records)
		output.FilesExternal = records
	}

	return output, util.Success
}

func (r *RelatedPartyCreate) Create() (errCode int) {
	var paramOut model.RelatedParty
	if r.UserID > 0 {
		paramOut.Creator = &r.UserID
	}

	if r.LastModifier > 0 {
		paramOut.LastModifier = &r.LastModifier
	}

	if r.Name != "" {
		paramOut.Name = &r.Name
	}

	if r.EnglishName != "" {
		paramOut.EnglishName = &r.EnglishName
	}

	if r.Address != "" {
		paramOut.Address = &r.Address
	}

	if r.UniformSocialCreditCode != "" {
		paramOut.UniformSocialCreditCode = &r.UniformSocialCreditCode
	}

	if r.Telephone != "" {
		paramOut.Telephone = &r.Telephone
	}

	if r.ImportedOriginalName != "" {
		paramOut.ImportedOriginalName = &r.ImportedOriginalName
	}

	if r.Remarks != "" {
		paramOut.Remarks = &r.Remarks
	}

	if r.FileIDs != "" {
		paramOut.FileIDs = &r.FileIDs
	}

	err := global.DB.Create(&paramOut).Error
	if err != nil {
		return util.ErrorFailToCreateRecord
	}
	return util.Success
}

func (r *RelatedPartyUpdate) Update() (errCode int) {
	paramOut := make(map[string]any)

	if r.UserID > 0 {
		paramOut["last_modifier"] = r.UserID
	}

	if r.Name != nil {
		if *r.Name != "" {
			paramOut["name"] = r.Name
		} else {
			paramOut["name"] = nil
		}
	}

	if r.EnglishName != nil {
		if *r.EnglishName != "" {
			paramOut["english_name"] = r.EnglishName
		} else {
			paramOut["english_name"] = nil
		}
	}

	if r.Address != nil {
		if *r.Address != "" {
			paramOut["address"] = r.Address
		} else {
			paramOut["address"] = nil
		}
	}

	if r.UniformSocialCreditCode != nil {
		if *r.UniformSocialCreditCode != "" {
			paramOut["uniform_social_credit_code"] = r.UniformSocialCreditCode
		} else {
			paramOut["uniform_social_credit_code"] = nil
		}
	}

	if r.Telephone != nil {
		if *r.Telephone != "" {
			paramOut["telephone"] = r.Telephone
		} else {
			paramOut["telephone"] = nil
		}
	}

	if r.Remarks != nil {
		if *r.Remarks != "" {
			paramOut["remarks"] = r.Remarks
		} else {
			paramOut["remarks"] = nil
		}
	}

	//查文件信息
	if r.FileIDs != nil {
		if *r.FileIDs != "" {
			paramOut["file_ids"] = r.FileIDs
		} else {
			paramOut["file_ids"] = nil
		}
	}

	err := global.DB.Model(&model.RelatedParty{}).
		Where("id = ?", r.ID).
		Updates(paramOut).Error
	if err != nil {
		return util.ErrorFailToUpdateRecord
	}

	return util.Success
}

func (r *RelatedPartyDelete) Delete() (errCode int) {
	//先找到记录，然后把deleter赋值给记录方便传给钩子函数，再删除记录，详见：
	var record model.RelatedParty
	err := global.DB.Where("id = ?", r.ID).
		Find(&record).
		Delete(&record).Error

	if err != nil {
		return util.ErrorFailToDeleteRecord
	}

	return util.Success
}

func (r *RelatedPartyGetList) GetList() (outputs []RelatedPartyOutput,
	errCode int, paging *list.PagingOutput) {
	db := global.DB.Model(&model.RelatedParty{})
	// 顺序：where -> count -> Order -> limit -> offset -> outputs

	//where
	if r.NameInclude != "" {
		db = db.Where("name like ?", "%"+r.NameInclude+"%")
	}

	if r.EnglishNameInclude != "" {
		db = db.Where("english_name like ?", "%"+r.EnglishNameInclude+"%")
	}

	// count
	var count int64
	db.Count(&count)

	//Order
	orderBy := r.SortingInput.OrderBy
	desc := r.SortingInput.Desc
	//如果排序字段为空
	if orderBy == "" {
		//如果要求降序排列
		if desc == true {
			db = db.Order("id desc")
		}
	} else { //如果有排序字段
		//先看排序字段是否存在于表中
		exists := util.FieldIsInModel(&model.RelatedParty{}, orderBy)
		if !exists {
			return nil, util.ErrorSortingFieldDoesNotExist, nil
		}
		//如果要求降序排列
		if desc == true {
			db = db.Order(orderBy + " desc")
		} else { //如果没有要求排序方式
			db = db.Order(orderBy)
		}
	}

	//limit
	page := 1
	if r.PagingInput.Page > 0 {
		page = r.PagingInput.Page
	}
	pageSize := global.Config.DefaultPageSize
	if r.PagingInput.PageSize != nil && *r.PagingInput.PageSize >= 0 &&
		*r.PagingInput.PageSize <= global.Config.MaxPageSize {
		pageSize = *r.PagingInput.PageSize
	}
	if pageSize > 0 {
		db = db.Limit(pageSize)
	}

	//offset
	offset := (page - 1) * pageSize
	db = db.Offset(offset)

	//outputs
	db.Model(&model.RelatedParty{}).Find(&outputs)

	if len(outputs) == 0 {
		return nil, util.ErrorRecordNotFound, nil
	}

	numberOfRecords := int(count)
	numberOfPages := util.GetNumberOfPages(numberOfRecords, pageSize)

	return outputs,
		util.Success,
		&list.PagingOutput{
			Page:            page,
			PageSize:        pageSize,
			NumberOfPages:   numberOfPages,
			NumberOfRecords: numberOfRecords,
		}
}
