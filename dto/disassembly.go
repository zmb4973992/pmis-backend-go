package dto

//以下dto均为输入

type DisassemblyTreeDTO struct {
	BaseDTO
}

// DisassemblyCreateOrUpdateDTO
// 除id外，所有字段都设置为必须绑定
type DisassemblyCreateOrUpdateDTO struct {
	BaseDTO
	Name       *string  `json:"name" binding:"required"`        //拆解项名称
	ProjectID  *int     `json:"project_id" binding:"required"`  //所属项目id
	Level      *int     `json:"level" binding:"required"`       //层级
	Weight     *float64 `json:"weight" binding:"required"`      //权重
	SuperiorID *int     `json:"superior_id" binding:"required"` //上级拆解项ID
}

// DisassemblyListDTO 是list查询的过滤器
// 在dto传递给sqlCondition时，空值会被忽略
type DisassemblyListDTO struct {
	ListDTO

	ProjectID  *int `form:"project_id"`
	SuperiorID *int `form:"superior_id"`
	Level      *int `form:"level"`
	LevelGte   *int `form:"level_gte"`
	LevelLte   *int `form:"level_lte"`
}

//以下dto均为输出

type DisassemblyOutputDTO struct {
	BaseDTO    `mapstructure:",squash"` //这里是嵌套结构体，mapstructure必须加squash，否则无法匹配
	Name       *string                  `json:"name" mapstructure:"name"`               //名称
	ProjectID  *int                     `json:"project_id" mapstructure:"project_id"`   //所属项目id
	Level      *int                     `json:"level" mapstructure:"level"`             //层级
	Weight     *float64                 `json:"weight" mapstructure:"weight"`           //权重
	SuperiorID *int                     `json:"superior_id" mapstructure:"superior_id"` //上级拆解项id
}

// DisassemblyOutputForTreeDTO 根据前端需求准备的数据结构
type DisassemblyOutputForTreeDTO struct {
	Name     *string                       `json:"title"`
	ID       int                           `json:"key"`
	Children []DisassemblyOutputForTreeDTO `json:"children" gorm:"-"`
}
