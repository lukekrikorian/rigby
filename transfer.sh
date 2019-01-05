echo "Building go project..."
go build
echo "Built"
echo "Killing website process..."
ssh -t root@rigby.space 'fuser -k 443/tcp'
echo "Process killed"
echo "Transfering site file..."
scp site root@rigby.space:~/rigby
echo "File transfered"
echo "Pulling changes from origin on remote machine..."
ssh -t root@rigby.space 'cd ~/rigby && git pull'
echo "Changes pulled"
echo "Restarting server..."
ssh root@rigby.space
echo "Testing site"
sleep 2s
curl rigby.space