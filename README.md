# StickyD

_Status: ✅✅ (i.e. the second stage of done)_

I wanted to pick a small project and have a look at what's going on under the hood, so I picked [etsy/statsd](https://github.com/etsy/statsd).

## Installation

	$ go install github.com/jpatel531/stickyd

## Running

	$ $GOPATH/bin/sticky path/to/config.json

## Configuration

The configuration file is similar to [the original](https://github.com/etsy/statsd/blob/master/exampleConfig.js). The difference is that this port doesn't support the pickle protocol, and there is no dynamic importing of backend or frontend interfaces.

## Caveat

This wasn't built to be used, just as a nice educational kata.
