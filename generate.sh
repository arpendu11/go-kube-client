#!/bin/bash

protoc kubepb/*.proto --go_out=plugins=grpc:kubepb/