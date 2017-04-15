# Universe

[![Build Status](https://travis-ci.org/qa-dev/universe.svg?branch=master)](https://travis-ci.org/qa-dev/universe)
[![Coverage Status](https://coveralls.io/repos/github/qa-dev/Universe/badge.svg?branch=master)](https://coveralls.io/github/qa-dev/Universe?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/qa-dev/universe)](https://goreportcard.com/report/github.com/qa-dev/universe)

## Requirements

- MongoDB
- RabbitMQ

## Setting up

- Run `make build`
- Edit `dist/config.json`
- Just run `dist/universe`

## How to use

There is only 3 endpoints and all of them are POST

1) Send event
2) Subscribe to event
3) Unsubscribe from event

Let's explain all 3 endpoints:

### *Send event*

**Path**: `/e/<event_name>`  
Where *event_name* is any url-allowed sequence of chars.

**Body**: put any information about event in request body.

### *Subscribe*

**Path**: `/subscribe/<plugin_name>`  
Today only log and webhook (in url named as `web`) plugins are available.

**Body**: example for webhook plugin
```
{
    "event_name": "myservice.job.done",
    "url": "http://example.com/webhook"
}
```

### *Unsubscribe*

**Path**: `/unsubscribe/<plugin_name>`  

**Body**: example for webhook plugin
```
{
    "event_name": "myservice.job.done",
    "url": "http://example.com/webhook"
}
```

# License

[MIT](LICENSE)
