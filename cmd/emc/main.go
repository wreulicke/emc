package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/alecthomas/units"
	"github.com/wreulicke/emc"
	"github.com/wreulicke/emc/calculator"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	DefaultMemoryLimitPathV1 = "/sys/fs/cgroup/memory/memory.limit_in_bytes"
	DefaultMemoryLimitPathV2 = "/sys/fs/cgroup/memory.max"
	UnsetTotalMemory         = calculator.Size(9_223_372_036_854_771_712)
)

func getMemory() (calculator.Size, error) {
	v, err := getMemoryLimit(DefaultMemoryLimitPathV1)
	if err != nil {
		return -1, err
	}
	if v == UnsetTotalMemory {
		v, err = getMemoryLimit(DefaultMemoryLimitPathV2)
		if err != nil {
			return -1, err
		}
	}
	// TODO parse /proc/meminfo

	if v == UnsetTotalMemory {
		return calculator.Size(1 * calculator.Gibi), nil
	}
	return v, nil
}

func getMemoryLimit(path string) (calculator.Size, error) {
	bs, err := os.ReadFile(path)
	if !os.IsNotExist(err) {
		return UnsetTotalMemory, nil
	} else if err == nil {
		s, err := calculator.ParseSize(strings.TrimSpace(string(bs)))
		if err != nil {
			return UnsetTotalMemory, nil
		}
		return s, nil
	}
	return UnsetTotalMemory, nil
}

func main() {
	app := kingpin.New("emc", "Enhanced java memory calculator")
	verbose := app.Flag("verbose", "Verbose").Default("false").Short('v').Bool()
	// TODO: currently total memory is not same interface for jvm...
	totalMemory := app.Flag("total-memory", "Total memory. Required if is not limited by cgroup").Bytes()
	loadedClassCount := app.Flag("loaded-class-count", "Loaded class count").Int64()
	threadCount := app.Flag("thread-count", "thread count").Default("250").Int64()
	javaOpts := app.Flag("java-options", "JVM Options").Envar("JAVA_OPTS").Default("").String()
	headRoom := app.Flag("head-room", "Percentage of total memory available which will be left unallocated to cover JVM overhead").Default("0").Int()
	findLambda := app.Flag("find-lambda", "find lambda").Default("false").Hidden().Bool() // experimental
	javaVersion := app.Flag("java-version", "Java version").Default("11").Int()
	jarOrDirectory := app.Arg("jarOrDirectory", "jar or directory").File()
	app.Action(func(c *kingpin.ParseContext) error {
		if *jarOrDirectory == nil && *loadedClassCount == 0 {
			return fmt.Errorf("jarOrDirectory or loaded-class-count is not specified")
		}
		if *jarOrDirectory != nil && *loadedClassCount > 0 {
			return fmt.Errorf("please specify either jarOrDirectory or loaded-class-count")
		}
		if int64(*totalMemory) == 0 {
			t, err := getMemory()
			if err != nil {
				return err
			}
			*totalMemory = units.Base2Bytes(t)
		}

		if *verbose {
			log.SetOutput(os.Stderr)
		} else {
			log.SetOutput(ioutil.Discard)
		}
		if j := *jarOrDirectory; j != nil {
			fi, err := j.Stat()
			if err != nil {
				return err
			}
			actualClassCount, err := emc.CountClassFile(j, fi, *findLambda)
			if err != nil {
				return err
			}
			stdLibClassCount := emc.CountClassInStandardLibrary(*javaVersion)
			*loadedClassCount = int64(0.35 * float64(actualClassCount+stdLibClassCount))
		}
		r, err := emc.Calculate(int64(*totalMemory), *loadedClassCount, *threadCount, *javaOpts, *headRoom)
		if err != nil {
			return fmt.Errorf("cannot calculate memory options. err=%v", err)
		}
		fmt.Println(strings.Join(r, " "))
		return nil
	})
	kingpin.MustParse(app.Parse(os.Args[1:]))
}
