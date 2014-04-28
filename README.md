gotest
======

Tool for renaming files to follow a particular order. First attempt at using go language.

Setting up
-----------

Configue environment variables

Path to go's binaries

    set GOROOT=c:\apps\go

Path to default go's workspace
  
    set GOPATH=c:\dev\go
  
The workspace should have following structure

    [c:\dev\go\bin]
    [c:\dev\go\pck]
    [c:\dev\go\src]
     c:\dev\go\src\github.com\deffer\gotest\.git

Althout it should be enough to set GOPATH, it may still complain about missing GOBIN.

    set GOBIN=%GOPATH%/bin
  
Running
--------
Compiling and installing as program.

    go install github.com/deffer/gotest/hello.go

