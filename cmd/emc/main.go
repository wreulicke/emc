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
	totalMemory := app.Flag("total-memory", "Total memory").Bytes()
	loadedClassCount := app.Flag("loaded-class-count", "Loaded class count").Int64()
	threadCount := app.Flag("thread-count", "Loaded class count").Default("250").Int64()
	javaOpts := app.Flag("java-options", "JVM Options").Envar("JAVA_OPTS").Default("").String()
	headRoom := app.Flag("head-room", "Percentage of total memory available which will be left unallocated to cover JVM overhead").Default("0").Int()
	jarOrDirectory := app.Arg("jarOrDirectory", "jar or directory").File()
	app.Action(func(c *kingpin.ParseContext) error {
		if *jarOrDirectory == nil && *loadedClassCount == 0 {
			return fmt.Errorf("jarOrDirectory or loaded-class-count is not specified")
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
		if j := *jarOrDirectory; j != nil {
			fi, err := j.Stat()
			if err != nil {
				return err
			}
			*loadedClassCount, err = emc.CountClassFile(j, fi)
			if err != nil {
				return err
			}
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
