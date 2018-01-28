# gogo
Go to MIPS compiler implemented in Go. Made as a course project for CS335 (Compiler Design).

<p align="center">
  <img alt="Logo" src="gopher.svg">
</p>

*Gopher vector imported from [egonelbre/gophers
](https://github.com/egonelbre/gophers).*

- - -

## Setting up
Run `./scripts/setup.sh` from the root directory of the project to set up the pre-commit git hooks.

## Dependencies
* [gocc](https://github.com/goccmack/gocc)

## Build
The following should generate relevant binaries inside the directory `bin` -
```
make
```

Alternatively, individual components can be built via -
```
make lexer
make tac
```

The tokens and their corresponding lexemes can be produced via -
```
./bin/lexer test/test1.go
```

The three-address code (in-memory) data structure can be prepared via -
```
./bin/codegen test/test1.ir
```
