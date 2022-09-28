#!/bin/bash

host="${HOST:-localhost:8080}"

curl -X POST ${host}/tasks -d \
'{
	"name": "Test Generator",
	"state": "Ready",
	"locator": "testgen",
	"addon": "testgen",
	"application": {
		"id": 4
	},
	"data": {
    "name": "TKLTEST_CONFIG_FILE",
    "general": {
      "app_name": "onlinebookstore",
      "monolith_app_path": [
      "OnlineBookStore/servlets"
      ],
      "app_classpath_file": "",
      "java_jdk_home": "",
      "offline_instrumentation": true,
      "build_type": "maven"
    },
    "generate": {
      "time_limit": 2,
      "add_assertions": false,
      "app_build_files": [
      "pom.xml"
      ],
      "target_class_list": [],
      "ctd_amplified": {
        "base_test_generator": "combined",
        "interaction_level": 1,
        "no_ctd_coverage": false,
        "num_seq_executions": 2
      }
    }
	}
}'
