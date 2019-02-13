#!/bin/bash
docker build -t local/oauth2:ropc .
docker build -t local/oauth2:ropc-dev -f Dockerfile-dev .
