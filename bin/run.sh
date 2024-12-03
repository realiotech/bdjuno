#!/usr/bin/env sh

callisto init

callisto parse genesis-file --genesis-file-path /callisto/config/genesis.json

exec callisto start --home /callisto/config