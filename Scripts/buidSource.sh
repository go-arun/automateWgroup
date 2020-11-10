#!/bin/bash
/usr/bin/echo "Success-" >> /tmp/scriptWorked
/usr/local/go/bin/go build -o /home/ubuntu/automateWgroup/main /var/www/html/WordPress/main.go
sleep 10
/usr/sbin/service goweb restart
