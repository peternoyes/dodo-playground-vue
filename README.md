# dodo-playground-vue

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
