## Enhanced JVM Meomry Calculator

emc is enhanced [cloudfoundry/java-buildpack-memory-calculator](https://github.com/cloudfoundry/java-buildpack-memory-calculator).
emc counts classes in directory and classes in jar, and show memory options.

```
# emc counts class in directory recursively. also supports jar in directory.
emc --total-memory 1G <path/to/directory>
-XX:MaxDirectMemorySize=10M -XX:MaxMetaspaceSize=14447K -XX:ReservedCodeCacheSize=240M -Xmx266128K

# emc counts class in jar recursively. also supports UberJar.
emc --total-memory 1G <path/to/your.jar>
-XX:MaxDirectMemorySize=10M -XX:MaxMetaspaceSize=14447K -XX:ReservedCodeCacheSize=240M -Xmx266128K
```

## Install

```
# emc is not released yet. you can try with go get.
go get github.com/wreulicke/emc
```

## Usage

```
$ go run ./cmd/emc/ --help
usage: emc --total-memory=TOTAL-MEMORY [<flags>] [<jarOrDirectory>]

Enhanced java memory calculator

Flags:
      --help              Show context-sensitive help (also try --help-long and --help-man).
  -v, --verbose           verbose
      --total-memory=TOTAL-MEMORY  
                          Total memory
      --loaded-class-count=LOADED-CLASS-COUNT  
                          Loaded class count
      --thread-count=250  Loaded class count
      --java-options=""   JVM Options
      --head-room=0       percentage of total memory available which will be left unallocated to cover JVM overhead

Args:
  [<jarOrDirectory>]  jar or directory
```

## LICENSE

MIT LICENSE