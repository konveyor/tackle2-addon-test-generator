#!/bin/bash

host="${HOST:-localhost:8080}"

curl -X POST ${host}/tasks -d \
'{
    "name":"Test Generator",
    "state": "Ready",
    "locator": "testgen",
    "addon": "testgen",
    "application": {"id": 1},
    "data": {
        "foo": {
            "bar": "",
    }
}' | jq -M .
