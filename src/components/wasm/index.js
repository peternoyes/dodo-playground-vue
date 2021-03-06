import './wasm_exec.js'

const funcs = {
  async init() {    
    const go = new global.Go()
    const wasm = await import('./main.wasm')
    await WebAssembly.instantiate(wasm.default, go.importObject).then((result) => {
      go.run(result.instance)      
    })
  },
  async run(ctx, firmware, binary, callback) {
    return global._dodo_simulator_run(ctx, firmware, binary, callback)
  },
  async stop() {
    return global._dodo_simulator_stop()
  },
  async updateKeys(keys) {
    return global._dodo_simulator_update_keys(keys)
  },
  async setFramFlusher(callback) {
    return global._dodo_simulator_set_fram_flusher(callback)
  }
}

export default funcs