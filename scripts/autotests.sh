#!/bin/env bash

go build -buildvcs=false -o ./cmd/server/server ./cmd/server && metricstest -test.v -test.run=^TestIteration1$ -binary-path=cmd/server/server
