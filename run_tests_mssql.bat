@echo off
rem TODO: Make these deletes a little less noisy
del "template_*.go"
del "gen_*.go"
cd tmpl_client
del "template_*.go"
cd ..
del "gosora.exe"

echo Generating the dynamic code
go generate
if %errorlevel% neq 0 (
	pause
	exit /b %errorlevel%
)

echo Building the router generator
go build ./router_gen
if %errorlevel% neq 0 (
	pause
	exit /b %errorlevel%
)
echo Running the router generator
router_gen.exe
if %errorlevel% neq 0 (
	pause
	exit /b %errorlevel%
)

echo Building the query generator
go build "./cmd/query_gen"
if %errorlevel% neq 0 (
	pause
	exit /b %errorlevel%
)
echo Running the query generator
query_gen.exe
if %errorlevel% neq 0 (
	pause
	exit /b %errorlevel%
)

echo Generating the JSON handlers
easyjson -pkg common

echo Building the executable
go test -tags mssql
if %errorlevel% neq 0 (
	pause
	exit /b %errorlevel%
)
pause