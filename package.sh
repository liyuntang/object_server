export GOOS=linux

cd client && go build main.go && mv main ../process/client && cd ../data_server/ && go build main.go && mv main ../process/data_server && cd ../api_service/ && go build main.go && mv main ../process/api_service && cd ../ && rsync -avzP process bmon11::soft/
