# Go Home.

Detect motion and person in your home.
so, right now go home.

## Install

### Requirements

- go 1.22
- opencv 4.10.x

### Set env

```bash
# OpenCV Config for MacOS installed by homebrew
PKG_CONFIG_PATH=/opt/homebrew/opt/opencv/lib/pkgconfig
```

### Build and Run

```bash
 GO HOME Deployments Makefile

 Usage:

        make {command}

 Commands:

        help                                                        helps (default)
        clean                                                       clean build directory
        run cmd={gohome-cli,gohome-gui,detector} flags={go flags}   run $(cmd)
        build cmd={gohome-cli,gohome-gui,detector}                  build $(cmd)
```
