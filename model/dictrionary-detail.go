package model

import (
	"pmis-backend-go/global"
)

type DictionaryDetail struct {
	BasicModel
	DictionaryTypeId int64   //字典类型的Id
	Name             string  //名称
	Sort             *int    //用于排序的值
	Status           *bool   //是否启用
	Remarks          *string //备注
}

// TableName 修改数据库的表名
func (d *DictionaryDetail) TableName() string {
	return "dictionary_detail"
}

type dictionaryDetailFormat struct {
	TypeName    string
	DetailNames []string
}

var initialDictionary = []dictionaryDetailFormat{
	{
		TypeName:    "省份",
		DetailNames: []string{"上海", "北京", "山东", "河南"},
	},
	{
		TypeName: "收付款方式",
		DetailNames: []string{"T/T(电汇)", "L/C(信用证)", "D/P(付款交单)", "D/A(承兑交单)",
			"D/D(票汇)", "Cash(现金)", "M/T(信汇)", "其他"},
	},
	{
		TypeName:    "进度类型",
		DetailNames: []string{"计划进度", "实际进度", "预测进度"},
	},
	{
		TypeName:    "币种",
		DetailNames: []string{"人民币", "美元", "欧元", "港币", "新加坡元", "马来西亚币"},
	},
	{
		TypeName: "合同类型",
		DetailNames: []string{
			"保险", "财务", "采购", "销售", "进口", "出口", "储运", "代理", "定金",
			"工程", "合作", "技服", "建安", "加工", "结算单", "贸易", "水泥批发采购",
			"水泥批发销售", "索赔", "总包", "租赁", "劳务", "其他",
		},
	},
	{
		TypeName:    "项目类型",
		DetailNames: []string{"工程EPC", "贸易", "服务", "投资", "其他"},
	},
	{
		TypeName:    "进度的数据来源",
		DetailNames: []string{"系统计算", "人工填写"},
	},
	{
		TypeName:    "数据范围",
		DetailNames: []string{"用户所在部门", "用户所在部门和子部门", "所有部门", "自定义部门"},
	},
	{
		TypeName:    "合同的资金方向",
		DetailNames: []string{"收款合同", "付款合同", "不涉及收付款"},
	},
	{
		TypeName:    "收付款的资金方向",
		DetailNames: []string{"收款", "付款"},
	},
	{
		TypeName:    "收付款的种类",
		DetailNames: []string{"计划", "实际", "预测"},
	},
	{
		TypeName: "国家",
		DetailNames: []string{"中国", "日本", "韩国", "朝鲜", "蒙古", "越南", "柬埔寨", "老挝",
			"泰国", "缅甸", "菲律宾", "文莱", "马来西亚", "新加坡", "印度尼西亚", "东帝汶", "尼泊尔",
			"不丹", "巴基斯坦", "印度", "孟加拉", "马尔代夫", "斯里兰卡", "哈萨克斯坦",
			"吉尔吉斯斯坦", "塔吉克斯坦", "乌兹别克斯坦", "土库曼斯坦", "阿富汗", "伊朗", "伊拉克",
			"叙利亚", "黎巴嫩", "以色列", "巴勒斯坦", "约旦", "沙特阿拉伯", "巴林", "卡塔尔",
			"科威特", "阿联酋", "阿曼", "也门", "格鲁吉亚", "亚美尼亚", "阿塞拜疆", "土耳其",
			"塞浦路斯", "冰岛", "丹麦", "挪威", "瑞典", "芬兰", "英国", "爱尔兰", "法国",
			"摩纳哥", "荷兰", "比利时", "卢森堡", "德国", "瑞士", "列支敦士登", "波兰", "捷克",
			"斯洛伐克", "奥地利", "匈牙利", "爱沙尼亚", "拉脱维亚", "立陶宛", "白俄罗斯",
			"乌克兰", "摩尔多瓦", "俄罗斯", "葡萄牙", "西班牙", "安道尔", "意大利", "圣马力诺",
			"梵蒂冈", "马耳他", "斯洛文尼亚", "克罗地亚", "波斯尼亚和黑塞哥维那", "黑山", "塞尔维亚",
			"阿尔巴尼亚", "北马其顿", "保加利亚", "希腊", "罗马尼亚", "塞浦路斯", "埃及", "利比亚",
			"突尼斯", "阿尔及利亚", "摩洛哥", "尼日尔", "布基纳法索", "马里", "毛里塔尼亚",
			"尼日利亚", "贝宁", "多哥", "加纳", "科特迪瓦", "利比里亚", "塞拉利昂", "几内亚",
			"几内亚比绍", "塞内加尔", "冈比亚", "佛得角", "乍得", "中非", "喀麦隆", "刚果民主共和国",
			"刚果共和国", "加蓬", "赤道几内亚", "圣多美和普林西比", "吉布提", "索马里", "厄立特里亚",
			"埃塞俄比亚", "苏丹", "南苏丹", "肯尼亚", "坦桑尼亚", "乌干达", "卢旺达", "布隆迪",
			"塞舌尔", "安哥拉", "赞比亚", "马拉维", "莫桑比克", "纳米比亚", "博茨瓦纳", "津巴布韦",
			"南非", "斯威士兰", "莱索托", "马达加斯加", "毛里求斯", "科摩罗", "澳大利亚", "新西兰",
			"帕劳", "密克罗尼西亚", "马绍尔群岛", "瑙鲁", "基里巴斯", "巴布亚新几内亚",
			"所罗门群岛", "瓦努阿图", "斐济", "图瓦卢", "萨摩亚", "汤加", "纽埃", "库克群岛",
			"加拿大", "美国", "墨西哥", "危地马拉", "伯利兹", "萨尔瓦多", "洪都拉斯", "尼加拉瓜",
			"哥斯达黎加", "巴拿马", "巴哈马", "古巴", "牙买加", "海地", "多米尼加",
			"圣基茨和尼维斯", "安提瓜和巴布达", "多米尼克", "圣卢西亚", "巴巴多斯",
			"圣文森特和格林纳丁斯", "格林纳达", "特立尼达和多巴哥", "哥伦比亚", "委内瑞拉", "圭亚那",
			"苏里南", "厄瓜多尔", "秘鲁", "玻利维亚", "巴西", "智利", "阿根廷", "乌拉圭", "巴拉圭"},
	},
	{
		TypeName:    "项目状态",
		DetailNames: []string{"未开始", "进行中", "已完成", "已中止"},
	},
	{
		TypeName: "我方签约主体",
		DetailNames: []string{"中国航空技术北京有限公司",
			"Avic International Engineering Holdings Pte. Ltd.",
			"北京凯玖科技发展有限责任公司", "北京凯祥恒业贸易有限公司",
			"江苏凯堡新材料科技有限公司",
			"浙江凯堡能源有限公司", "中航国际凯融有限公司"},
	},
	{
		TypeName: "款项类型",
		DetailNames: []string{"预付款", "定金", "进度款", "尾款", "质保款",
			"发货款", "港杂费", "调试款", "杂费", "租金", "服务费", "保证金",
			"保费", "其他"},
	},
	{
		TypeName:    "tabFukuan视图中不要导入的记录",
		DetailNames: []string{"CNC0155-001", "CNC0155-002", "CNS0419-001"},
	},
	{
		TypeName:    "操作类型",
		DetailNames: []string{"添加", "修改", "删除"},
	},
	{
		TypeName:    "收款的数据来源",
		DetailNames: []string{"收款", "收汇", "收票"},
	},
}

func generateDictionaryDetail() (err error) {
	var dictionaryDetails []DictionaryDetail
	for i := range initialDictionary {
		//先找到字典类型的记录
		var dictionaryTypeInfo DictionaryType
		err = global.DB.Where("name = ?", initialDictionary[i].TypeName).
			First(&dictionaryTypeInfo).Error
		if err != nil {
			return err
		}

		for j := range initialDictionary[i].DetailNames {
			dictionaryDetails = append(dictionaryDetails, DictionaryDetail{
				DictionaryTypeId: dictionaryTypeInfo.Id,
				Name:             initialDictionary[i].DetailNames[j],
				Remarks:          &initialDictionary[i].TypeName,
			})
		}
	}

	for _, dictionaryDetail := range dictionaryDetails {
		err = global.DB.
			Where("name = ?", dictionaryDetail.Name).
			Where("dictionary_type_id = ?", dictionaryDetail.DictionaryTypeId).
			Attrs(&DictionaryDetail{
				Status: BoolToPointer(true),
			}).
			FirstOrCreate(&dictionaryDetail).Error
		if err != nil {
			return err
		}
	}
	return nil
}
