package service

import (
	"pmis-backend-go/global"
	"pmis-backend-go/util"
	"strconv"
)

//以下为入参
//有些字段不用json tag，因为不从前端读取，而是在controller中处理

//指针字段是为了区分入参为空或0与没有入参的情况，做到分别处理，通常用于update
//如果指针字段为空或0，那么数据库相应字段会改为null；
//如果指针字段没传，那么数据库不会修改该字段

type RBACUpdate struct {
	LastModifier int
	RoleIDs      []int `json:"role_ids,omitempty"`
	MenuSnowIDs  []int `json:"menu_snow_ids,omitempty"`
	APISnowIDs   []int `json:"api_snow_ids,omitempty"`
	//RBACInfos    []RBACInfo
	//连接关联表的id

	//连接dictionary_item表的id

	//日期

	//允许为0的数字

	//允许为null的字符串

}

type RBACInfo struct {
	Path   string `json:"path"`
	Method string `json:"method"`
}

//type ContractDelete struct {
//	ID int
//}
//
//type ContractGetList struct {
//	dto.ListInput
//	dto.DataScopeInput
//	ProjectID     int    `json:"project_id,omitempty"`
//	FundDirection int    `json:"fund_direction,omitempty"`
//	NameInclude   string `json:"name_include,omitempty"`
//}

//以下为出参

//type ContractOutput struct {
//	Creator      *int `json:"creator"`
//	LastModifier *int `json:"last_modifier"`
//	ID           int  `json:"id"`
//	//连接关联表的id，只用来给gorm查询，不在json中显示
//	ProjectID      *int `json:"-"`
//	RoleID *int `json:"-"`
//	RelatedPartyID *int `json:"-"`
//	//连接dictionary_item表的id，只用来给gorm查询，不在json中显示
//	FundDirection *int `json:"-"`
//	OurSignatory  *int `json:"-"`
//	Currency      *int `json:"-"`
//	Type          *int `json:"-"`
//	//关联表的详情，不需要gorm查询，需要在json中显示
//	ProjectExternal      *ProjectOutput      `json:"project" gorm:"-"`
//	OrganizationExternal *OrganizationOutput `json:"organization" gorm:"-"`
//	RelatedPartyExternal *RelatedPartyOutput `json:"related_party" gorm:"-"`
//	//dictionary_item表的详情，不需要gorm查询，需要在json中显示
//	FundDirectionExternal *DictionaryDetailOutput `json:"fund_direction" gorm:"-"`
//	OurSignatoryExternal  *DictionaryDetailOutput `json:"our_signatory" gorm:"-"`
//	CurrencyExternal      *DictionaryDetailOutput `json:"currency" gorm:"-"`
//	TypeExternal          *DictionaryDetailOutput `json:"type" gorm:"-"`
//	//其他属性
//	SigningDate       *string `json:"signing_date"`
//	EffectiveDate     *string `json:"effective_date"`
//	CommissioningDate *string `json:"commissioning_date"`
//	CompletionDate    *string `json:"completion_date"`
//
//	Amount             *float64 `json:"amount"`
//	ExchangeRate       *float64 `json:"exchange_rate"`
//	ConstructionPeriod *int     `json:"construction_period"`
//
//	Name        *string `json:"name"`
//	Code        *string `json:"code"`
//	Content     *string `json:"content"`
//	Deliverable *string `json:"deliverable"`
//	PenaltyRule *string `json:"penalty_rule"`
//	Attachment  *string `json:"attachment"`
//	Operator    *string `json:"operator"`
//}

func (r *RBACUpdate) Update() error {
	cachedEnforcer, err := util.NewCachedEnforcer()
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return err
	}

	for i := range r.RoleIDs {
		subject := strconv.Itoa(r.RoleIDs[i])
		var ok bool
		ok, err = cachedEnforcer.RemoveFilteredNamedPolicy("p", 0, subject)
		if err != nil || !ok {
			global.SugaredLogger.Errorln(err)
			return err
		}

		//var RBACRules [][]string
		//for _, v := range r.RBACInfos {
		//	RBACRules = append(RBACRules, []string{subject, v.Path, v.Method})
		//}

		//ok, err = cachedEnforcer.AddPolicies(RBACRules)
		//if err != nil || !ok {
		//	global.SugaredLogger.Errorln(err)
		//	return err
		//}

		//修改了policy以后，因为用的是cachedEnforcer，所以要清除缓存
		err = cachedEnforcer.InvalidateCache()
		if err != nil {
			global.SugaredLogger.Errorln(err)
			return err
		}
	}

	return nil
}
