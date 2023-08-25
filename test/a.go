package test

import (
	"github.com/yitter/idgenerator-go/idgen"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"time"
)

var chanForTest = make(chan model.Test, 100)

func generateData() {
	time.Sleep(2 * time.Second)
	number := idgen.NextId()

	var record model.Test
	record.A = number

	chanForTest <- record

	println(number, "已进入管道")
}

func saveData(record model.Test) {
	time.Sleep(1 * time.Second)
	global.DB.Create(&record)
	println(record.A, "已保存到数据库")
}

func Test() {
	for i := 0; i < 300; i++ {
		i := i
		go func() {
			println(i)
			generateData()
		}()
	}

	go func() {
		for {
			select {
			case record := <-chanForTest:
				saveData(record)
			}
		}
	}()
}
