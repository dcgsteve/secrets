#!/bin/bash

SV=$(cat version)
SV=$((${SV}+1))
echo $SV > version