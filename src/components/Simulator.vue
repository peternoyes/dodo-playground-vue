<template>
  <div>    
    <canvas tabindex="0" 
      id="sim-canvas"
      v-on:keydown.a="keyADown"
      v-on:keyup.a="keyAUp"
      v-on:keydown.b="keyBDown"
      v-on:keyup.b="keyBUp"
      width="256" 
      height="128"></canvas>
  </div>
</template>

<script>
import dodo from './wasm/index.js'

export default {
  name: 'Simulator',
  props: ['firmware', 'binary'],
  data() {
    return {
      A: false,
      B: false
    }
  },
  computed: {
    keyState: function() {
      let keys = ""
      if (this.A) {
        keys += "A"
      }
      if (this.B) {
        keys += "B"
      }
      return keys
    }
  },
  watch: {
    keyState: async function(val) {
      await dodo.updateKeys(val)
    }
  },
  async mounted() {
    let canvas = document.getElementById('sim-canvas')
    let ctx = canvas.getContext('2d')

    var buffer
    switch (this.firmware) {
        case '1.0.1': 
            buffer = await import('../assets/firmware_1.0.1.bin')
            break
        case '1.1.0_beta2':
            buffer = await import('../assets/firmware_1.1.0_beta2.bin')
            break
    }

    var firmwareBytes = new Uint8Array(buffer.default)

    await dodo.init()
    await dodo.run(ctx, firmwareBytes, this.binary)
  },
  async unmounted() {
      console.log("Unmounting")
      await dodo.stop()
  },
  methods: {
    keyADown() {
      this.A = true
    },
    keyAUp() {
      this.A = false
    },
    keyBDown() {
      this.B = true
    },
    keyBUp() {
      this.B = false
    }
  }
}
</script>

<style scoped>
</style>
