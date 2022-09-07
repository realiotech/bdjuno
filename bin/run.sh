#!/usr/bin/env sh

bdjuno init

bdjuno parse genesis-file --genesis-file-path /bdjuno/config/genesis.json

exec bdjuno start --home /bdjuno/config