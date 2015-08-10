# httpgo practice code
  This is my simple ip tracer, or it may call Dynamic DNS.
```bash
# runserver
go install
httpsrv -password="yourpassword" -port="12345"

# client register ip with username
curl -k -u username:yourpassword http://ip:12345/updateIp
# and you could put this command to crontab

# other machine get the ip by username
curl -k -u username:yourpassword http://ip:1234/getIp
```
