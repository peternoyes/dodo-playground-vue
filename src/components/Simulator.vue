<template>
  <div 
    tabindex="0"
    v-on:keydown.a="keyADown"
    v-on:keyup.a="keyAUp"
    v-on:keydown.b="keyBDown"
    v-on:keyup.b="keyBUp">
    <center>
      <img src="../assets/dodo_large.png" width=764 usemap="#game_pad">
      <map name="game_pad">
        <area shape="rect" coords="509,374,574,439" href="#" @mousedown="keyADown" @mouseup="keyAUp" />
      </map>
      <canvas
        id="sim-canvas"
        style="position:absolute; top:104px; left:271px; z-index:1"
        width="256" 
        height="128"></canvas>
    </center>
    <p>Frames Per Cycle: {{ cycles }}</p>
  </div>
</template>

<script>
import dodo from './wasm/index.js'

export default {
  name: 'Simulator',
  props: ['firmware', 'binary', 'stopped'],
  data() {
    return {
      cycles: 0,
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
    },
    stopped: async function(val) {
      if (val) {
        console.log('Stopping')
        await dodo.stop()
      }
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
    await dodo.run(ctx, firmwareBytes, this.binary, (cycles) => {
      this.cycles = cycles
    })
  },
  async unmounted() {
      console.log('Unmounting')
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
