package dto

type Base struct {
	Creator      *int `json:"creator" mapstructure:"creator"`
	LastModifier *int `json:"last_modifier" mapstructure:"last_modifier"`
	//deleter不需要json tag，因为不从前端获取，而是从context读取
	Deleter *int
	ID      int `json:"id" mapstructure:"id"`
}
