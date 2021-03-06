# mockingjay server

[![Build Status](https://travis-ci.org/quii/mockingjay-server.svg?branch=master)](https://travis-ci.org/quii/mockingjay-server)[![Coverage Status](https://coveralls.io/repos/quii/mockingjay-server/badge.svg?branch=master)](https://coveralls.io/r/quii/mockingjay-server?branch=master)[![GoDoc](https://godoc.org/github.com/quii/mockingjay-server?status.svg)](https://godoc.org/github.com/quii/mockingjay-server)

Mockingjay lets you define the contract between a consumer and producer and with just a configuration file you get:

- A fast to launch fake server for your integration tests
 - Configurable to simulate the eratic nature of calling other services
- [Consumer driven contracts (CDCs)](http://martinfowler.com/articles/consumerDrivenContracts.html) to run against your real downstream services.

**Mockingjay makes it really easy to check your integration points**. It's fast, requires no coding and is better than other solutions because it will ensure your mock servers and real integration points are consistent

## Rationale

In the hip exciting world of SOA/microservices with heavy investment in PaaS/IaaS you want to be able to quickly iterate over small services and deploy to live quickly and without fear of breaking things.

If you are using this kind of architecture you will be faced with the challenge of ensuring that your huge numbers of services can actually talk to each other.

You will probably employ things like versioning to help but you might also be spending time writing consumer driven contracts (CDCs) to ensure your integration points are working.

In addition you might be writing integration tests against fakes/stubs to check your code can send the correct requests and be able to parse responses.

![alt tag](http://i.imgur.com/oC6BjGn.png)

If you squint hard enough, you can imagine that the requirements for both CDCs and fake servers are the same. *Given a particular request, I expect a particular kind of response*. Yet with this set up you are duplicating this information _with code_ in two different files which obviously isn't ideal.

What mockingjay enables you to do is to capture these requirements in one configuration file.

````yaml
---
 - name: My very important integration point
   request:
     uri: /hello
     method: POST
     body: "Chris" # * matches any body
   response:
     code: 200
     body: '{"message": "hello, Chris"}'   # * matches any body
     headers:
       content-type: application/json

# define as many as you need...
````

From this you can create a fake server to write integration tests with and also check the service you are dependant on is consistent with what you expect.

#### Main advantages

- No coding whatsoever, so no naughtiness in fake servers overcomplicating things. Even non developers can add new scenarios to test with.
- The contract is defined once, rather than dispersed across different scripts which you have to keep in sync.
- Entirely language agnostic. If you speak HTTP you can use mockingjay.
- Checks the structure of the data (currently JSON is the only type checked) rather than the contents, which will reduce flakiness of your builds.
- Both the fake server and CDCs are really fast to run, to help keep your builds fast.

#### Drawbacks/constraints

- You can only express your consumer-producer interaction in terms of isolated request/responses. Sometimes you might need to test a number of requests which are dependant on each other.

## Installation

     $ go get github.com/quii/mockingjay-server

If you don't have Go installed you can find a binary for your platform in the [releases tab](https://github.com/quii/mockingjay-server/releases)

## Running a fake server

````bash
$ mockingjay-server -config=example.yaml -port=1234 &
2015/04/13 14:27:54 Serving 3 endpoints defined from example.yaml on port 1234
$ curl http://localhost:1234/hello
{"message": "hello, world"}
````

## Check configuration is compatible with a real server

````bash
$ mockingjay-server -config=example.yaml -realURL=http://some-real-api.com
2015/04/13 21:06:06 Test endpoint (GET /hello) is incompatible with http://some-real-api - Couldn't reach real server
2015/04/13 21:06:06 Test endpoint 2 (DELETE /world) is incompatible with http://some-real-api - Couldn't reach real server
2015/04/13 21:06:06 Failing endpoint (POST /card) is incompatible with http://some-real-api - Couldn't reach real server
2015/04/13 21:06:06 At least one endpoint was incompatible with the real URL supplied
````

### Inspect what requests mockingjay has received

     http://{mockingjayhost}:{port}/requests

Calling this will return you a JSON list of requests

## Make your fake server flaky

Mockingjay has an annoying friend, a monkey. Given a monkey configuration you can make your fake service misbehave. This can be useful for performance tests where you want to simulate a more realistic scenario (i.e all integration points are painful).

````yaml
---
# Writes a different body 50% of the time
- body: "This is wrong :( "
  frequency: 0.5

# Delays initial writing of response by a second 20% of the time
- delay: 1000
  frequency: 0.2

# Returns a 404 30% of the time
- status: 404
  frequency: 0.3

# Write 10,000,000 garbage bytes 9% of the time
- garbage: 10000000
  frequency: 0.09
````

````bash
$ mockingjay-server -config=examples/example.yaml -monkeyConfig=examples/monkey-business.yaml
2015/04/17 14:19:53 Serving 3 endpoints defined from examples/example.yaml on port 9090
2015/04/17 14:19:53 Monkey config loaded
2015/04/17 14:19:53 50% of the time | Body: This is wrong :(
2015/04/17 14:19:53 20% of the time | Delay: 1s
2015/04/17 14:19:53 30% of the time | Status: 404
2015/04/17 14:19:53  9% of the time | Garbage bytes: 10000000
````

## Building

### Requirements

- Go 1.3+ installed ($GOPATH set, et al)
- godep https://github.com/tools/godep
- golint https://github.com/golang/lint

````bash
$ go get https://github.com/quii/mockingjay-server.git
$ cd $GOPATH/src/github.com/quii/mockingjay-server
$ ./build.sh
````
