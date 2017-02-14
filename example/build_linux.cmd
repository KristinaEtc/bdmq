set GOOS=linux
rem call govvv build -o stomp-server stomp-server.go 
call go build example.go
copy stomp-server o:\work\go-stomp\
