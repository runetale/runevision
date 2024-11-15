# thor
automated red team tools

# how to start
`make up`
`go run main.go`

# arch
web <-> application api server <-> thor local backend server (daemon service) <-> thor engine
           |                        |
           |                        |
           |                        |
     thor database                black hat module or white hat module <-> ml api server <-> python ml model