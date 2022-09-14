#!/bin/bash

wget $(curl -s https://api.github.com/repos/dcgsteve/secrets/releases/latest | grep 'browser_' | grep 'secrets-linux-amd64' | cut -d\" -f4)
chmod +x secrets-linux-amd64
sudo mv secrets-linux-amd64 /usr/bin/secrets