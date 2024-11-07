package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/qri-io/jsonschema"
)

var schemaData = []byte(`{
  "$schema": "https://json-schema.org/draft/2019-09/schema",
  "type": "object",
  "properties": {
    "opentofu": { "$ref": "#/$defs/tool" },
    "terraform": { "$ref": "#/$defs/tool" }
  },
  "required": ["opentofu", "terraform"],
  "$defs": {
    "tool": {
      "type": "object",
      "properties": {
        "license": {
          "type": "string"
        },
        "registry": {
          "type": "string",
          "format": "uri"
        },
        "features": {
          "type": "array",
          "items": {
            "type": "object",
            "properties": {
              "name": {
                "type": "string"
              },
              "url": {
                "type": "string",
                "format": "uri"
              },
              "version": {
                "type": "string"
              }
            },
						"required": ["name", "version"]
          }
        }
      },
      "required": ["license", "registry", "features"]
    }
  }
}`)

var tools = []byte(`{
  "opentofu": {
    "license": "MPL 2.0",
    "registry": "https://search.opentofu.org/",
    "features": [
      {
        "name": "Test",
        "version": "1.6.0"
      },
      {
        "name": "State encryption",
        "url": "https://opentofu.org/docs/language/state/encryption/",
        "version": "1.7.0"
      },
      {
        "name": "Removed block",
        "version": "1.7.0"
      },
      {
        "name": "Provider-defined functions",
        "version": "1.7.0"
      },
      {
        "name": "Configured provider-defined functions",
        "version": "1.7.0"
      },
      {
        "name": "Loopable import blocks",
        "version": "1.7.0"
      },
      {
        "name": "Backend configuration using locals and variables",
        "version": "1.8.0"
      },
      {
        "name": ".tofu extension",
        "version": "1.8.0"
      },
      {
        "name": "Provider mocking",
        "version": "1.8.0"
      },
      {
        "name": "override_resource, override_data, override_module",
        "version": "1.8.0"
      },
      {
        "name": "templatefile() and templatestring() recursion",
        "version": "1.7.0"
      }
    ]
  },
  "terraform": {
    "license": "BSL",
    "registry": "https://registry.terraform.io/",
    "features": [
      {
        "name": "Test",
        "version": "1.6.0"
      },
      {
        "name": "Removed block",
        "version": "1.7.0"
      },
      {
        "name": "Provider-defined functions",
        "version": "1.8.0"
      },
      {
        "name": "Provider mocking",
        "version": "1.7.0"
      },
      {
        "name": "override_resource, override_data, override_module",
        "version": "1.7.0"
      }
    ]
  }
}
`)

func main() {
	ctx := context.Background()

	rs := &jsonschema.Schema{}
	if err := json.Unmarshal(schemaData, rs); err != nil {
		panic("unmarshal schema: " + err.Error())
	}

	errs, err := rs.ValidateBytes(ctx, tools)
	if err != nil {
		panic(err)
	}
	if len(errs) > 0 {
		for _, e := range errs {
			fmt.Println(e.Error())
		}
	}

	// File for Hugo to template the table
	f, err := os.Create("data/tools.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.Write(tools)

	// allow access to https://cani.tf/tools.json
	static, err := os.Create("data/tools.json")
	if err != nil {
		panic(err)
	}
	defer static.Close()
	static.Write(tools)

}