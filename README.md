# dodo-playground-vue

## Background

The Dodo Playground is a web application that allows users to write and run [Dodo](http://www.dodolabs.io) games. This repository is the second iteration of the playground. The [initial version](https://github.com/peternoyes/dodo-playground) of the playground was written entirely in Go using GopherJS and jquery. At this point GopherJS is trapped in version 1.12 of Go and is outdated as Go now officially supports WASM. 

This new version is a complete rewrite using VueJS with only the core simulator now in Go. The core simulator now uses WASM instead of GopherJS. Another signficant upgrade is the adoption of the WebSerial API for out of the box in browser flashing!  The previous app required a separate Chrome Application in order to flash games from the browser using a serial cable, but with the recent release of Chrome 89 we can finally use WebSerial!

This version also separates the [backend API](https://github.com/peternoyes/dodo-api) out into a separate repository. 

To see this version in action visit [next.dodolabs.io](https://next.dodolabs.io/?code=fc1811bd).

## WASM

Note that the actual simulator in the playground is written in Golang and requires compilation to WASM. At this point the WASM is not automatically compiled during the build. To compile execute the below command from within the /src/components/wasm folder.

```
GOOS=js GOARCH=wasm go build -o main.wasm
```

## Project setup
```
npm install
```

### Compiles and hot-reloads for development
```
npm run serve
```

### Compiles and minifies for production
```
npm run build
```

### Lints and fixes files
```
npm run lint
```

### Customize configuration
See [Configuration Reference](https://cli.vuejs.org/config/).
