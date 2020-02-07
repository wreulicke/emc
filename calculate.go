package emc

import (
	"github.com/cloudfoundry/java-buildpack-memory-calculator/calculator"
	"github.com/cloudfoundry/java-buildpack-memory-calculator/flags"
)

func Calculate(totalMemory int64, loadedClassCount int64, threadCount int64, javaOptons string, headRoom int) ([]string, error) {
	h := flags.HeadRoom(headRoom)
	options := flags.JVMOptions{}
	options.Set(javaOptons)
	tc := flags.ThreadCount(threadCount)
	tm := flags.TotalMemory(totalMemory)
	lcc := flags.LoadedClassCount(loadedClassCount)
	calc := calculator.Calculator{
		HeadRoom:         &h,
		JvmOptions:       &options,
		ThreadCount:      &tc,
		TotalMemory:      &tm,
		LoadedClassCount: &lcc,
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
