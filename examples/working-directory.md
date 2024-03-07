# Title of the example

## Code

```yaml
verbose: true
tasks:
  - alias: myls
    cmd: ls
    args:
      - "-al"
      - "-h"
    directory: /var/log
  - alias: mypwd
    cmd: pwd
    directory: /var/log
```

## Description

## Execution

```bash
❯ xec --config examples/working-directgory.xec.yaml myls                                                     ☸ docker-desktop in xec on  shell-support [!+] via 🐹 v1.22.0 with unknown env 
2024-03-08T00:42:51+01:00 | SUCCESS | Loaded config file, [examples/working-directgory.xec.yaml]
2024-03-08T00:42:51+01:00 | SUCCESS | Loaded config file, [$HOME/.xec.yaml]
2024-03-08T00:42:51+01:00 | SUCCESS | Loaded config file, [./.xec.yaml]
2024-03-08T00:42:51+01:00 | INFO | Task myls is starting.
2024-03-08T00:42:51+01:00 | INFO | Task myls is not logged.
total 43760
drwxr-xr-x  46 root            wheel   1.4K Mar  8 00:30 .
drwxr-xr-x  36 root            wheel   1.1K Feb 23 21:15 ..
-rw-r--r--   1 root            wheel     0B Dec 14 18:46 alf.log
drwxr-xr-x   2 root            wheel    64B Feb  2 18:19 apache2
...
```


```bash
❯ xec --config examples/working-directgory.xec.yaml mypwd                                                    ☸ docker-desktop in xec on  shell-support [!+] via 🐹 v1.22.0 with unknown env 
2024-03-08T00:45:59+01:00 | SUCCESS | Loaded config file, [examples/working-directgory.xec.yaml]
2024-03-08T00:45:59+01:00 | SUCCESS | Loaded config file, [$HOME/.xec.yaml]
2024-03-08T00:45:59+01:00 | SUCCESS | Loaded config file, [./.xec.yaml]
2024-03-08T00:45:59+01:00 | INFO | Task mypwd is starting.
2024-03-08T00:45:59+01:00 | INFO | Task mypwd is not logged.
/var/log
2024-03-08T00:45:59+01:00 | INFO | Task mypwd finished in 2.155875ms.
2024-03-08T00:45:59+01:00 | SUCCESS | Task mypwd completed successfully in 2.155875ms.

```