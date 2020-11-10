/usr/bin/echo "Success" >> /tmp/scriptWorked
/usr/local/go/bin/go build ../main.go
/usr/sbin/service goweb restart
