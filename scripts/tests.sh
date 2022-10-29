#!/bin/bash -e
R='\033[0;31m'
G='\033[0;32m'
B='\033[0;34m'
NC='\033[0m'

printf "\n${R}UNIT TESTS...${NC}\n"
go test -shuffle=on $@ --tags=unit ./... | grep -v "?"
printf "${R}UNIT TESTS...done${NC}\n"

printf "\n${B}INTEGRATION TESTS...${NC}\n"
go test -count=1 $@ --tags=integration ./... | grep -v "?"
printf "${B}INTEGRATION TESTS...done${NC}\n"
