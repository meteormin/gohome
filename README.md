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
$ git clone https://github.com/meteormin/gohome.git
# CLI version. using gocron package
$ go run cmd/gohome-cli/main.go
# GUI version. using gocv package
$ go run cmd/gohome-gui/main.go


```
