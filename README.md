WARNING
=======
This is mega super alpha. It has lots of bugs and is probably very insecure. Use at your own risk.

Building
========
server
------
6g server.go model.go common.go
6l -o server server.6

client
------
6g lunch.go common.go 
6l -o lunch lunch.6

migrate
-------
mv lunch_config ~/lunch_config
