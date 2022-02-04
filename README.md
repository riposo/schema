# Schema Plugin

Schema rules plugin for [Riposo](https://github.com/riposo/riposo).

This plugin allows to configure [JSON Schema](https://json-schema.org/)
validations rules ensuring user-submitted records adhere to a pre-defined data
schema.

## Configuration

The configure, add `type: schema` rules to your configuration. Schema rules must
include a `schema` URL and a `path` array.

```yaml
rules:
  - type: schema
    schema: "https://example.com/person.schema.json"
    path:
      - "/buckets/*/collections/people/records/*"
  - type: schema
    schema: "file:///tmp/geographical-location.schema.json"
    path:
      - "/buckets/{one,two}/collections/coordinates/records/*"
      - "!/**/records/special"
```

## License

Copyright 2022 Black Square Media Ltd

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this material except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
