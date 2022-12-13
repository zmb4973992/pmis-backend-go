package dto

type DictionaryGetDTO struct {
	ProjectType  []string `json:"project_type"`
	Province     []string `json:"province"`
	ContractType []string `json:"contract_type"`
}

type DictionaryListDTO struct {
	//ListDTO
	//SuperiorID *int    `form:"superior_id"`
	//Level      *string `form:"level"`
	//
	//Name     *string `form:"name"`
	//NameLike *string `form:"name_like"`
}
