ECHO 'Installer made with GISG'

WHERE choco
IF %ERRORLEVEL% NEQ 0 ECHO choco is not installed. Please install it.

WHERE go
IF %ERRORLEVEL% NEQ 0 ECHO go is not installed. Please install it.

WHERE git
IF %ERRORLEVEL% NEQ 0 ECHO git is not installed. Please install it.

go get github.com/stretchr/testify/assert || exit /b

go build || exit 1
