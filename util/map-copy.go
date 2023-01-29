package util

func MapCopy(originalMap map[string]any, keysToBeOmitted ...string) (targetMap map[string]any) {
	targetMap = make(map[string]any)

	//如果有需要忽略的字段
	if len(keysToBeOmitted) > 0 {
		for k1, v1 := range originalMap {
			for _, v2 := range keysToBeOmitted {
				if k1 != v2 {
					targetMap[k1] = v1
				}
			}
		}
		return targetMap
	}

	//没有需要忽略的字段，即所有字段都保留
	for k1, v1 := range originalMap {
		targetMap[k1] = v1
	}
	return targetMap
}
