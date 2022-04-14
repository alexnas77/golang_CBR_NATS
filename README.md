# golang_CBR_NATS
Obtain currency rates from www.cbr.ru and publish it through NATS message broker

First install and run NATS message broker

for example

`go install github.com/nats-io/nats.go/@latest`

`nats-server -m 8222 -DV`

Then build and run 

`go run receive.go`

Finally, build and run

`go run cbr.go`
