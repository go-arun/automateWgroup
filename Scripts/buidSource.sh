/usr/bin/echo "Success" >> /tmp/scriptWorked
sudo /usr/local/go/bin/go build ../main.go
sudo /usr/bin/mv -f main /home/ubuntu/automateWgroup
sudo /usr/sbin/service goweb restart
