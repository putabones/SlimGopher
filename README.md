# SlimGopher

SlimGopher is a wrapper written in Golang.  It compresses, encrypts, and embeds the executable.

### Shrink

First we build the compression and encryption executable.  There will will be 3 files generated that will be used from the output of Shrink to build the final executable.

```
cd shrink
go build -o shrink doit.go
```

### Shrink help

```
./shrink

usage: shrink [-h|--help] -p|--path "<value>" -e|--embed "<value>"

              Shrink compresses and encrypts your embedded, and also creates
              the variables need for your payload

Arguments:

  -h  --help   Print help information
  -p  --path   Path where you want the variable files saved (REQUIRED)
  -e  --embed  Executable you want to embed (REQUIRED)

```

### Then we build the final executable

You will need the 3 files from Shrink which need to be built with the final executable.

```
cd ..
go build -o final *.go
```

### Acknowledgments

[dtoebe](https://github.com/dtoebe/embed-binary)

[amenzhinsky](https://github.com/amenzhinsky/go-memexec)

[whit3rabbit](https://github.com/whit3rabbit)
