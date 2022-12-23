<!-- vale Microsoft.HeadingAcronyms = NO -->
# Xec

## Features

- Reads .env by default.
- Can import other `.xec.yaml` files

Looks for ./.xec.(yaml|yml) if not found it traverses up to /

## Examples

### 

```yaml
verbose: true
debug: true
logFile: "xec.log"
environment:
  values:
    - key=value
    - environment=dev
  passOn: true
  acceptFilterRegex:
    - "XEC_*"
  rejectFilterRegex:
    - "*SECRET*"
tasks:
  - alias: printenv
    description: execute printenv, usuallu for debugging.
    cmd: printenv
    wait: true # false will cause xec spawn another process
    retry: 3
```