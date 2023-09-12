package util

import "github.com/yitter/idgenerator-go/idgen"

func InitIdGenerator() {
	//workerId范围：0-63
	options := idgen.NewIdGeneratorOptions(0)
	idgen.SetIdGenerator(options)
}
