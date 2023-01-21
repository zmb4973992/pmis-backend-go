package snowflake

import (
	"github.com/sony/sonyflake"
	"pmis-backend-go/global"
)

var snowFlakeInstance *sonyflake.Sonyflake

func Init() {
	settings := sonyflake.Settings{}
	snowFlakeInstance = sonyflake.NewSonyflake(settings)
	if snowFlakeInstance == nil {
		global.SugaredLogger.Panicln("生成snowflake实例失败，请重试")
	}
}

func GenerateID() (id uint64, err error) {
	id, err = snowFlakeInstance.NextID()
	return
}
