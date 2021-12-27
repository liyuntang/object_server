IP=`ifconfig | grep 192 | awk '{print $2}'`
go run main.go  --storage=/Users/liyuntang/data --listen=${IP}:9200
