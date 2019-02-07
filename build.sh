#!/bin/bash
docker build -t local/oauth2 .
docker build -t local/oauth2:dev -f Dockerfile-dev .
