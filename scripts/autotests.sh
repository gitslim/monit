#!/bin/env bash

go build -buildvcs=false -o ./cmd/server/server ./cmd/server && \
go build -buildvcs=false -o ./cmd/agent/agent ./cmd/agent && \
          metricstest -test.v -test.run=^TestIteration1$ \
            -binary-path=cmd/server/server && \
          metricstest -test.v -test.run=^TestIteration2[AB]*$ \
            -source-path=. \
            -agent-binary-path=cmd/agent/agent && \
          metricstest -test.v -test.run=^TestIteration3[AB]*$ \
            -source-path=. \
            -agent-binary-path=cmd/agent/agent \
            -binary-path=cmd/server/server
