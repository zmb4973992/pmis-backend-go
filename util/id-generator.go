package util

import "github.com/yitter/idgenerator-go/idgen"

func InitIDGenerator() {
	options := idgen.NewIdGeneratorOptions(123)
	idgen.SetIdGenerator(options)
}
