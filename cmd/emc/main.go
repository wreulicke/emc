package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/alecthomas/units"
	"github.com/wreulicke/emc"
	"gopkg.in/alecthomas/kingpin.v2"
)

func getMemory() (units.Base2Bytes, error) {
	bs, err := ioutil.ReadFile("/sys/fs/cgroup/memory/memory.limit_in_bytes")
	if err != nil {
		return -1, fmt.Errorf("cannot detect memory limit automatically: %v", err)
	}
	v := strings.TrimSpace(string(bs))
	if v == "9223372036854771712" {
		return 1 * units.GiB, nil
	}
	i, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return -1, fmt.Errorf("cannot parse memory limit from /sys/fs/cgroup/memory/memory.limit_in_bytes: %v", err)
	}
	return units.Base2Bytes(i), nil

}

func main() {
	app := kingpin.New("emc", "Enhanced java memory calculator")
	verbose := app.Flag("verbose", "Verbose").Default("false").Short('v').Bool()
	totalMemory := app.Flag("total-memory", "Total memory. Required if is not limited by cgroup").Bytes()
	loadedClassCount := app.Flag("loaded-class-count", "Loaded class count").Int64()
	threadCount := app.Flag("thread-count", "thread count").Default("250").Int64()
	javaOpts := app.Flag("java-options", "JVM Options").Envar("JAVA_OPTS").Default("").String()
	headRoom := app.Flag("head-room", "Percentage of total memory available which will be left unallocated to cover JVM overhead").Default("0").Int()
	findLambda := app.Flag("find-lambda", "find lambda").Default("false").Hidden().Bool() // experimental
	javaVersion := app.Flag("java-version", "Java version").Default("11").Int()
	jarOrDirectory := app.Arg("jarOrDirectory", "jar or directory or http/https schema").String()
	app.Action(func(c *kingpin.ParseContext) error {
		if jarOrDirectory == nil && *loadedClassCount == 0 {
			return fmt.Errorf("jarOrDirectory or loaded-class-count is not specified")
		}
		if jarOrDirectory != nil && *loadedClassCount > 0 {
			return fmt.Errorf("please specify either jarOrDirectory or loaded-class-count")
		}
		if int64(*totalMemory) == 0 {
			t, err := getMemory()
			if err != nil {
				return err
			}
			*totalMemory = t
		}

		if *verbose {
			log.SetOutput(os.Stderr)
		} else {
			log.SetOutput(ioutil.Discard)
		}
		if jarOrDirectory != nil {
			actualClassCount, err := emc.CountClassFileWithPath(*jarOrDirectory, *findLambda)
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
