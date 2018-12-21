#!/bin/bash

export JWT="Authorization:Bearer $(http post http://localhost:3000/login username=admin password=admin|jq -r .token)"

http :3000/api/purge "$JWT"
http :3000/api/import/1802 "$JWT"
http :3000/api/import/1803 "$JWT"
http :3000/api/import/1804 "$JWT"
http :3000/api/import/1805 "$JWT"
http :3000/api/import/1806 "$JWT"
http :3000/api/import/1807 "$JWT"
http :3000/api/import/1808 "$JWT"
http :3000/api/import/1809 "$JWT"
http :3000/api/import/1810 "$JWT"
http :3000/api/import/1811 "$JWT"
http :3000/api/import/1812 "$JWT"
