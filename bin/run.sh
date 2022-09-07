#!/usr/bin/env bash

bdjuno parse genesis-file --genesis-file-path /bdjuno/config/genesis.json

exec bdjuno start --home /bdjuno/config