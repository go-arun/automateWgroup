/usr/bin/echo "Success" >> /tmp/scriptWorked
/usr/local/go/bin/go build ../main.go
/usr/bin/mv -f /var/www/html/WordPress/Scripts/main /home/ubuntu/automateWgroup
/usr/sbin/service goweb restart
