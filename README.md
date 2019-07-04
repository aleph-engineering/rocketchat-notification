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
        Wrap message in a code area (default false)
  -f string
        Configuration file (Optional)
  -m string
        Message to post
  -p string
        Rocket.Chat user's password
  -s string
        Rocket.Chat server (default "http://localhost:3000")
  -u string
        Rocket.Chat user
```
If the flag `-m` is not specified the program will read the message from standard input.

The following example will post the output of command ps to Rocket.Chat using code style.
Ex: `ps | ./rocket-notification -s http://meet.cu.aleph.engineering -u jenkins -p password -c general -code true`


## Using environment variables
You can also specify the basic configuration using environment variables

`ROCKET_CHAT_USER` : for user
`ROCKET_CHAT_PASSWORD`: for password
`ROCKET_CHAT_SERVER`: for server url
`ROCKET_CHAT_CHANNEL`: for channel 

## How to test

Install **github.com/stretchr/testify/assert** `go get -u -v github.com/stretchr/testify/assert`
and run `go test -coverprofile cp.out`
