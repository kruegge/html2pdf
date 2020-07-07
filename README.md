# simple HTML2PDF converter in golang

## requirements
* need a running headless chrome browser in a container or local on your device
* get dependencies
* run the code
* have a look at the result :)

### run a headless chrome
> `google-chrome --headless --disable-gpu --remote-debugging-port=9222`

### get all dependencies
> `go mod vendor`


### run programm
> `go run main.go`

## results

when everything was going right, then you have two new files in your folder named `screenshot.png` and `page.pdf`. Please have a look inside and try to modify the code or capture some things on other sites.