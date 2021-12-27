while true
do
    for process in client data_server api_service
    do
	num=`ps aux | grep $process | grep -v grep | wc -l`
	if [ $num -eq 1 ];then
	    PID=`ps aux | grep $process | grep -v grep | awk '{print $2}'`
	    Count=`lsof -p $PID | wc -l`
	    echo ">>>>> $PID $process $Count"
	fi
    done
    echo -e "\n"
    sleep 1
done
