echo "Building go project..."
go build
echo "Built"
echo "Killing website process..."
ssh -t root@rigby.space 'fuser -k 443/tcp'
echo "Process killed"
echo "Transfering site file..."
scp site root@rigby.space:~/rigby
echo "File transfered"