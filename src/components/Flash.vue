<template>
  <div>
    <p>
      After choosing the appropriate COM port, ensure that Dodo is on the home screen before selecting 'Start' which will immediately flash Dodo. After successfully flashing Dodo, it will take a few additional seconds to load the game.
    </p>
    <b-alert v-model="showTestStatus" :variant="testVariant">{{ testStatus }}</b-alert>
    <b-button @click="getVersion">Test Connection</b-button>    
  </div>
</template>
<script>
export default {
  name: 'Flash',
  props: ['binary'],
  data() {
    return {
      showTestStatus: false,
      testStatus: '',
      testVariant: 'primary'
    }
  },
  methods: {    
    getVersion: async function() {
      this.testStatus = 'Requesting Port...'
      this.testVariant = 'primary'
      this.showTestStatus = true

      try
      {
        let port = await navigator.serial.requestPort();

        await port.open({ baudRate: 9600 })

        let decoder = new TextDecoderStream() // eslint-disable-line
        let inputDone = port.readable.pipeTo(decoder.writable)
        let inputStream = decoder.readable;
        let reader = inputStream.getReader()

        let encoder = new TextEncoderStream() // eslint-disable-line
        let outputDone = encoder.readable.pipeTo(port.writable)
        let outputStream = encoder.writable
        let writer = outputStream.getWriter()

        writer.write('V')
        writer.releaseLock()

        this.testStatus = 'Connecting...'

        // 2s timeout
        let timeout = setTimeout(async function() {           
          await reader.cancel()
        }, 2000)

        let text = ''      
        let failure = false
        while (true) { // eslint-disable-line
          let { value, done } = await reader.read()                
          if (value) {
            if (value.endsWith('\0')) {
              done = true
              value = value.slice(0, -1)
            }          
            text += value
          } else {          
            failure = true
            break
          }

          if (done) {
            await reader.cancel()                    
            break
          }
        }

        reader.releaseLock()          

        clearTimeout(timeout)

        if (failure) {
          this.testVariant = 'danger'
          this.testStatus = 'Dodo is not responding'
        } else {
          // Success!
          this.testVariant = 'success'
          this.testStatus = `Connected to Dodo with version: ${text}`
        }
        
        await inputDone.catch(() => {})

        await outputStream.getWriter().close()
        await outputDone

        await port.close()
      }
      catch (err) {
        this.testVariant = 'danger'
        this.testStatus = err.message
      }
    }
  }
}
</script>

<style scoped>
</style>