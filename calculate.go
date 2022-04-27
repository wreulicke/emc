package emc

import (
	"github.com/wreulicke/emc/calculator"
)

func Calculate(totalMemory int64, loadedClassCount int64, threadCount int64, javaOptons string, headRoom int) ([]string, error) {
	options := calculator.JVMOptions{}
	options.Set(javaOptons)
	calc := calculator.Calculator{
		HeadRoom:         headRoom,
		JvmOptions:       &options,
		ThreadCount:      threadCount,
		TotalMemory:      calculator.Size(totalMemory),
		LoadedClassCount: loadedClassCount,
	}
	str, err := calc.Calculate()
	if err != nil {
		return nil, err
	}
	result := make([]string, 0, len(str))
	for _, v := range str {
		result = append(result, v.String())
	}
	return result, nil
}
