package dto

//以下为入参

type DictionaryList struct {
	//GetList
	//SuperiorID *int    `form:"superior_id"`
	//LevelName      *string `form:"level"`
	//
	//Name     *string `form:"name"`
	//NameLike *string `form:"name_like"`
}

//以下为出参

type DictionaryOutput struct {
	ProjectType  []string `json:"project_type"`
	Province     []string `json:"province"`
	ContractType []string `json:"contract_type"`
}
