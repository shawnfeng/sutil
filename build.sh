#!/bin/bash

go build -v $(go list ./... | grep -v vendor/)
