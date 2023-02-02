package dto

//以下为入参
//有些字段不用json tag，因为不从前端读取，而是在controller中处理

type DisassemblyTree struct {
	Creator      int
	LastModifier int
	ProjectID    int `json:"project_id" binding:"required"`
}

type DisassemblyCreate struct {
	Creator      int
	LastModifier int

	Name       string  `json:"name" binding:"required"`        //拆解项名称
	ProjectID  int     `json:"project_id" binding:"required"`  //所属项目id
	Level      int     `json:"level" binding:"required"`       //层级
	Weight     float64 `json:"weight" binding:"required"`      //权重
	SuperiorID int     `json:"superior_id" binding:"required"` //上级拆解项ID
}

//指针字段是为了区分入参为空或0与没有入参的情况，做到分别处理，通常用于update
//如果指针字段为空或0，那么数据库相应字段会改为null；
//如果指针字段没传，那么数据库不会修改该字段

type DisassemblyUpdate struct {
	LastModifier int
	ID           int

	Name       *string  `json:"name"`        //拆解项名称
	ProjectID  *int     `json:"project_id"`  //所属项目id
	Level      *int     `json:"level"`       //层级
	Weight     *float64 `json:"weight"`      //权重
	SuperiorID *int     `json:"superior_id"` //上级拆解项ID
}

type DisassemblyDelete struct {
	Deleter int
	ID      int
}

type DisassemblyList struct {
	ListInput
	NameInclude string `json:"name_include,omitempty"`

	ProjectID  int  `json:"project_id"`
	SuperiorID int  `json:"superior_id"`
	Level      int  `json:"level"`
	LevelGte   *int `json:"level_gte"`
	LevelLte   *int `json:"level_lte"`
}

//以下为出参

type DisassemblyOutput struct {
	Creator      *int `json:"creator" gorm:"creator"`
	LastModifier *int `json:"last_modifier" gorm:"last_modifier"`
	ID           int  `json:"id" gorm:"id"`

	Name       *string  `json:"name" gorm:"name"`               //名称
	ProjectID  *int     `json:"project_id" gorm:"project_id"`   //所属项目id
	Level      *int     `json:"level" gorm:"level"`             //层级
	Weight     *float64 `json:"weight" gorm:"weight"`           //权重
	SuperiorID *int     `json:"superior_id" gorm:"superior_id"` //上级拆解项id
}

// DisassemblyTreeOutput 根据前端需求准备的数据结构
type DisassemblyTreeOutput struct {
	Name     *string                 `json:"title" gorm:"name"`
	ID       int                     `json:"key" gorm:"id"`
	Level    int                     `json:"level" gorm:"level"`
	Children []DisassemblyTreeOutput `json:"children"`
}
