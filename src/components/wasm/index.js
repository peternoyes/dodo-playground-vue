import './wasm_exec.js'

const funcs = {
  async init() {    
    const go = new global.Go()
    const wasm = await import('./main.wasm')
    await WebAssembly.instantiate(wasm.default, go.importObject).then((result) => {
      go.run(result.instance)      
    })
  },
  async run(ctx, firmware, binary) {
    return global._dodo_simulator_run(ctx, firmware, binary)
  },
  async stop() {
    return global._dodo_simulator_stop()
  },
  async updateKeys(keys) {
    return global._dodo_simulator_update_keys(keys)
  }
}

export default funcs