#!/usr/bin/env sh

bdjuno parse genesis-file --genesis-file-path /bdjuno/config/genesis.json

exec bdjuno start --home /bdjuno/config