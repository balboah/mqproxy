# MQProxy

This is an http proxy using message queueing as the transport.

## Why?

Having a message queue brings several benefits:

- Requests are newer pushed to potentially broken servers, instead requests are pulled when we are ready to handle them.
- Auto scaling can monitor queue lengths to get an idea of current backp reassure.
- Logging can easily hook into the message queue without the need of additional implementations on the service itself.

## How?

There are two parts, the publisher (producer) and subscriber (consumer). Publisher is close to the client acting as a terminator for connections and new requests, subscriber is close to the service that should respond to requests.

## Example

Start message broker gnatsd:

	gnatsd &

Start a subscriber which is the one that replies to the request:

	./subscriber

Start a publisher which is the one open for http requests:

	./publisher

Make a request to the publisher, subscriber responds through the message queue:

	curl -q http://localhost:8080/nats
	Hello world!
