## Enhanced JVM Meomry Calculator

emc is enhanced [cloudfoundry/java-buildpack-memory-calculator](https://github.com/cloudfoundry/java-buildpack-memory-calculator) and inspired by [making/memory-calculator-cnb](https://github.com/making/memory-calculator-cnb).
emc counts classes in directory and classes in jar, and show memory options.



```bash
# emc counts class in directory recursively. also supports jar in directory.
emc --total-memory 1G <path/to/directory>
-XX:MaxDirectMemorySize=10M -XX:MaxMetaspaceSize=14447K -XX:ReservedCodeCacheSize=240M -Xmx266128K

# emc counts class in jar recursively. also supports UberJar.
emc --total-memory 1G <path/to/your.jar>
-XX:MaxDirectMemorySize=10M -XX:MaxMetaspaceSize=14447K -XX:ReservedCodeCacheSize=240M -Xmx266128K
```

## Install

```bash
# emc is not released yet. you can try with go get.
go get github.com/wreulicke/emc
```

## Usage

```bash
$ go run ./cmd/emc/ --help
usage: emc [<flags>] [<jarOrDirectory>]

Enhanced java memory calculator

Flags:
      --help              Show context-sensitive help (also try --help-long and --help-man).
  -v, --verbose           Verbose
      --total-memory=TOTAL-MEMORY  
                          Total memory. Required if is not limited by cgroup
      --loaded-class-count=LOADED-CLASS-COUNT  
                          Loaded class count
      --thread-count=250  Loaded class count
      --java-options=""   JVM Options
      --head-room=0       Percentage of total memory available which will be left unallocated to cover JVM overhead

Args:
  [<jarOrDirectory>]  jar or directory
```

## LICENSE

MIT LICENSE