#!/bin/bash

export JWT="Authorization:Bearer $(http post http://localhost:3000/login username=admin password=admin|jq -r .token)"

http :3000/api/compact/etablissement "$JWT"
http :3000/api/compact/entreprise "$JWT"
http :3000/api/reduce/algo2/1812 "$JWT"
