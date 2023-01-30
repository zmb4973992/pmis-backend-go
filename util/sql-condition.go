package util

import (
	"gorm.io/gorm"
	"pmis-backend-go/dto"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
)

//这里的思路来自zendea
//https://github.com/zendea/zendea

type SqlCondition struct {
	SelectedColumns []string //暂时弃用，因为比较麻烦，要考虑dto和model转换的问题
	ParamPairs      []ParamPair
	Sorting         dto.SortingInput
	Paging          dto.PagingInput
}

type ParamPair struct {
	Key   string //查询参数的名称，如 age>=, name include, id=
	Value any    //查询参数的值
}

// NewSqlCondition 生成自定义的查询条件,参数可不填
// 必须为指针，因为下面的方法要用到指针进行修改入参
func NewSqlCondition() *SqlCondition {
	pageSize := global.Config.PagingConfig.DefaultPageSize
	return &SqlCondition{
		Paging: dto.PagingInput{
			Page:     1,
			PageSize: pageSize,
		},
	}
}

// Where 建议不要直接用，如果是“等于赋值”，可以用equal
// 给SqlCondition自定义where方法，将参数保存到ParameterPair中
func (s *SqlCondition) Where(key string, value any) *SqlCondition {
	s.ParamPairs = append(s.ParamPairs, ParamPair{
		Key:   key,
		Value: value,
	})
	return s
}

func (s *SqlCondition) Equal(paramKey string, paramValue any) *SqlCondition {
	s.Where(paramKey+" = ?", paramValue)
	return s
}

func (s *SqlCondition) NotEqual(paramKey string, paramValue any) *SqlCondition {
	s.Where(paramKey+" <> ?", paramValue)
	return s
}

func (s *SqlCondition) Gt(paramKey string, paramValue any) *SqlCondition {
	s.Where(paramKey+" > ?", paramValue)
	return s
}

func (s *SqlCondition) Gte(paramKey string, paramValue any) *SqlCondition {
	s.Where(paramKey+" >= ?", paramValue)
	return s
}

func (s *SqlCondition) Lt(paramKey string, paramValue any) *SqlCondition {
	s.Where(paramKey+" < ?", paramValue)
	return s
}

func (s *SqlCondition) Lte(paramKey string, paramValue any) *SqlCondition {
	s.Where(paramKey+" <= ?", paramValue)
	return s
}

//func (s *SqlCondition) Include(paramKey string, paramValue string) *SqlCondition {
//	s.Where(paramKey+" LIKE ?", "%"+paramValue+"%")
//	return s
//}

func (s *SqlCondition) Like(paramKey string, paramValue string) *SqlCondition {
	s.Where(paramKey+" LIKE ?", "%"+paramValue+"%")
	return s
}

func (s *SqlCondition) StartWith(paramKey string, paramValue string) *SqlCondition {
	s.Where(paramKey+" LIKE ?", paramValue+"%")
	return s
}

func (s *SqlCondition) EndWith(paramKey string, paramValue string) *SqlCondition {
	s.Where(paramKey+" LIKE ?", "%"+paramValue)
	return s
}

func (s *SqlCondition) In(paramKey string, paramValue any) *SqlCondition {
	s.Where(paramKey+" IN ?", paramValue)
	return s
}

func (s *SqlCondition) Build(db *gorm.DB) *gorm.DB {
	//处理顺序：select → Where → order → limit → offset
	//select

	//选择要显示哪些字段。如果不填，就显示全部字段
	//selectedColumns暂时弃用，因为比较麻烦，涉及到dto、model的转换
	if len(s.SelectedColumns) > 0 {
		db = db.Select(s.SelectedColumns)
	}

	//定义绝对不传给前端的字段，比如密码等
	OmittedColumns := global.Config.DBConfig.OmittedColumns
	db = db.Omit(OmittedColumns...)

	//Where
	if len(s.ParamPairs) > 0 {
		for _, parameterPair := range s.ParamPairs {
			db = db.Where(parameterPair.Key, parameterPair.Value)
		}
	}

	//orderBy
	orderBy := s.Sorting.OrderBy
	if orderBy == "" { //如果排序字段为空
		if s.Sorting.Desc == true { //如果要求降序排列
			db = db.Order("id desc")
		}
	} else { //如果有排序字段
		if s.Sorting.Desc == true { //如果要求降序排列
			db = db.Order(s.Sorting.OrderBy + " desc")
		} else { //如果没有要求排序方式
			db = db.Order(s.Sorting.OrderBy)
		}
	}

	//limit
	db = db.Limit(s.Paging.PageSize)

	//原offset方法，数据量超过50万后会明显变慢。好处是不用考虑id的缺失
	offset := (s.Paging.Page - 1) * s.Paging.PageSize
	db = db.Offset(offset)

	//新offset方法，数据量哪怕达到几千万也不会产生查询瓶颈，已测试过
	//任何数据库的 offset 1000000 都比 Where id > 1000000 要慢很多
	//问题在于如果id不连续，会导致偏移出现错误
	//offset := (s.Paging.Page - 1) * s.Paging.PageSize
	//if offset > 0 {
	//	db = db.Where("id > ?", offset)
	//}

	return db
}

// Count 第二个参数应为model struct，如：model.User{}
// 不理解的话可以看该方法的源码，因为使用了gorm的db.model()方法
func (s *SqlCondition) Count(db *gorm.DB, modelName model.IModel) int {
	// Where
	if len(s.ParamPairs) > 0 {
		for _, parameterPair := range s.ParamPairs {
			db = db.Where(parameterPair.Key, parameterPair.Value)
		}
	}
	var totalRecords int64
	err := db.Debug().Model(&modelName).Count(&totalRecords).Error
	if err != nil {
		return 0
	}
	return int(totalRecords)
}

func (s *SqlCondition) Find(tempDb *gorm.DB, modelName model.IModel) (list []map[string]any) {
	//根据sqlCondition处理db
	tempDb = s.Build(tempDb)

	//出结果
	err := tempDb.Debug().Model(&modelName).Find(&list).Error
	if err != nil {
		return nil
	}
	return
}

// FieldIsInModel 验证提交的单个字段是否存在于表中（即数据表是否有相应的字段）
func (s *SqlCondition) FieldIsInModel(model model.IModel, field string) bool {
	//获取自定义的数据库表名
	tableName := model.TableName()
	//自行拼接的sql，找出对应表名的所有字段名
	//sqlStatement server的标准写法为：Select Name FROM SysColumns Where id = Object_Id('[某某表]')
	//给 某某表 加上中括号，是因为当表名中含有特殊字符时，直接使用单引号，会出现表名不被识别的问题
	var existedFields []string
	//这里goland编译器莫名报错，函数可以正常运行，可忽略
	sqlStatement := "Select Name FROM SysColumns Where id = OBJECT_ID('[" + tableName + "]')"
	global.DB.Raw(sqlStatement).Find(&existedFields)
	//如果表中字段数量>0且该字段在表的这些字段中
	if len(existedFields) > 0 && IsInSlice(field, existedFields) {
		return true
	}
	return false
}

// FieldIsInModel 验证提交的单个字段是否存在于表中（即数据表是否有相应的字段）
func FieldIsInModel(model model.IModel, field string) bool {
	//获取自定义的数据库表名
	tableName := model.TableName()
	var existedFields []string
	//自行拼接的sql，找出对应表名的所有字段名
	//sqlStatement server的标准写法为：Select Name FROM SysColumns Where id = Object_Id('[某某表]')
	//给 某某表 加上中括号，是因为当表名中含有特殊字符时，直接使用单引号，会出现表名不被识别的问题
	//这里goland编译器莫名报错，函数可以正常运行，可忽略
	sqlStatement := "Select Name FROM SysColumns Where id = OBJECT_ID('[" + tableName + "]')"
	global.DB.Raw(sqlStatement).Find(&existedFields)
	//如果表中字段数量>0且该字段在表的这些字段中
	if len(existedFields) > 0 && IsInSlice(field, existedFields) {
		return true
	}
	return false
}

// FieldsAreInModel 验证提交的多个字段是否存在于表中（即数据表是否有相应的字段）
func (s *SqlCondition) FieldsAreInModel(model model.IModel, fields ...string) bool {
	for _, field := range fields {
		res := s.FieldIsInModel(model, field)
		//如果有任何一个字段不符合要求,则直接返回false
		if res == false {
			return false
		}
	}
	return true
}

// FieldsAreInModel 验证提交的多个字段是否存在于表中（即数据表是否有相应的字段）
func FieldsAreInModel(model model.IModel, fields ...string) bool {
	for _, field := range fields {
		res := FieldIsInModel(model, field)
		//如果有任何一个字段不符合要求,则直接返回false
		if res == false {
			return false
		}
	}
	return true
}

func GetTotalNumberOfPages(numberOfRecords int, pageSize int) (numberOfPages int) {
	if numberOfRecords <= 0 || pageSize <= 0 {
		return 0
	}
	numberOfPages = numberOfRecords / pageSize
	if numberOfRecords%pageSize != 0 {
		numberOfPages++
	}
	return numberOfPages
}
