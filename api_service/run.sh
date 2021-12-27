IP=`ifconfig | grep 192 | awk '{print $2}'`
go run main.go --listen=${IP}:10000
