/usr/bin/echo "Success" >> /tmp/scriptWorked
/usr/local/go/bin/go build -o /home/ubuntu/automateWgroup/main /var/www/html/WordPress/main.go
/usr/sbin/service goweb restart
