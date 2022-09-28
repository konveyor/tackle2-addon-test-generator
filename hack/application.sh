curl -X POST ${host}/applications -d \
'{
    "name":"test",
    "description": "Cat application.",
    "repository": {
      "kind": "git",
      "url": "https://github.com/shashirajraja/onlinebookstore",
    },
}' | jq -M .
