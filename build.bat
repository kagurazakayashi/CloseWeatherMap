@ECHO OFF
SET name=xlsweather
SET version=1.2.0
RD /S /Q bin
MKDIR bin
SET CGO_ENABLED=0

SET GOOS=windows
ECHO Compiling Windows x86
SET GOARCH=386
go generate
go build -o bin\%name%_%version%_windows-x86.exe
START MAKECAB /D compressiontype=lzx /D compressionmemory=21 bin\%name%_%version%_windows-x86.exe bin\%name%_%version%_windows-x86.cab
DEL /Q *.syso
ECHO Compiling Windows x64
SET GOARCH=amd64
go generate
go build -o bin\%name%_%version%_windows-x64.exe
START MAKECAB /D compressiontype=lzx /D compressionmemory=21 bin\%name%_%version%_windows-x64.exe bin\%name%_%version%_windows-x64.cab
DEL /Q *.syso
ECHO Compiling Windows ARM32
SET GOARCH=arm
go generate
go build -o bin\%name%_%version%_windows-arm32.exe
START MAKECAB /D compressiontype=lzx /D compressionmemory=21 bin\%name%_%version%_windows-arm32.exe bin\%name%_%version%_windows-arm32.cab
DEL /Q *.syso
ECHO Compiling Windows ARM64
SET GOARCH=arm64
go generate
go build -o bin\%name%_%version%_windows-arm64.exe
START /WAIT MAKECAB /D compressiontype=lzx /D compressionmemory=21 bin\%name%_%version%_windows-arm64.exe bin\%name%_%version%_windows-arm64.cab
DEL /Q *.syso
DEL bin\*.exe

SET GOOS=darwin
ECHO Compiling macOS x64
SET GOARCH=amd64
go build -o bin\%name%_%version%_macos-x64
START xz -z -e -9 -v bin\%name%_%version%_macos-x64
ECHO Compiling macOS ARM64
SET GOARCH=arm64
go build -o bin\%name%_%version%_macos-arm64
START xz -z -e -9 -v bin\%name%_%version%_macos-arm64

SET GOOS=linux
ECHO Compiling Linux x86
SET GOARCH=386
go build -o bin\%name%_%version%_linux-x86
START xz -z -e -9 -v bin\%name%_%version%_linux-x86
ECHO Compiling Linux x64
SET GOARCH=amd64
go build -o bin\%name%_%version%_linux-x64
START xz -z -e -9 -v bin\%name%_%version%_linux-x64
ECHO Compiling Linux ARM32
SET GOARCH=arm
go build -o bin\%name%_%version%_linux-arm32
START xz -z -e -9 -v bin\%name%_%version%_linux-arm32
ECHO Compiling Linux ARM64
SET GOARCH=arm64
go build -o bin\%name%_%version%_linux-arm64
START xz -z -e -9 -v bin\%name%_%version%_linux-arm64

SET name=
SET version=
SET CGO_ENABLED=
SET GOOS=
SET GOARCH=
