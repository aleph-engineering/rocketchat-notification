# rocket-notification
Push messages to Rocket.Chat

## How to build

Using GO >= 1.9 run
`go build .`


## How to use
```
Usage of ./rocket-notification:
  -c string
        Channel used to post the message (default "general")
  -code
        Wrap message in a code area
  -m string
        Message to post
  -p string
        Rocket.Chat user's password (default "password")
  -s string
        Rocket.Chat server (default "http://localhost:3000")
  -u string
        Rocket.Chat user (default "user")
```
If the flag `-m` is not specified the program will read the message from standard input.


## How to test

Install **github.com/stretchr/testify/assert** `go get -u -v github.com/stretchr/testify/assert`
and run `go test`