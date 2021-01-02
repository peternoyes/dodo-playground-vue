<template>
  <div id="app">
    <b-container fluid class="d-flex flex-column flex-grow-1 vh-100 px-0 overflow-hidden">    
      <b-navbar toggleable="lg" type="dark" variant="dark" class="flex-shrink-0">
        <b-navbar-brand>Dodo Playground</b-navbar-brand>
        <b-navbar-toggle target="nav-collapse"></b-navbar-toggle>
        <b-collapse id="nav-collapse" is-nav>
          <b-navbar-nav>
            <b-nav-item-dropdown :text="displayLanguage">
              <b-dropdown-item v-for="option in languageOptions"
                :key="option.value"
                :value="option.value"
                @click="language = option.value">{{ option.text }}</b-dropdown-item>
            </b-nav-item-dropdown>            
            <b-nav-item-dropdown :text="firmware">
              <b-dropdown-item v-for="version in firmwareVerisons"
                :key="version"
                :value="version"
                @click="firmware = version">{{ version }}</b-dropdown-item>
            </b-nav-item-dropdown>
            <b-nav-form>            
              <b-button pill variant="info" class="ml-2" @click="run"><span><b-icon icon="play" /> Run</span></b-button>
              <b-button pill variant="success" class="ml-2" @click="download"><span><b-icon icon="download" /> fram.bin</span></b-button>
              <b-button pill variant="warning"  class="ml-2"><span><b-icon icon="lightning" /> Flash</span></b-button>
            </b-nav-form>          
          </b-navbar-nav>
          <b-navbar-nav class="ml-auto">
            <b-nav-form>
              <b-button variant="dark"><span><b-icon icon="github" /> Sign in with GitHub</span></b-button>
            </b-nav-form>
          </b-navbar-nav>
        </b-collapse>
      </b-navbar>
      <b-row class="flex-grow-1 overflow-hidden" no-gutters>
        <b-col cols="8" class="mh-100 overflow-auto py-2">          
            <editor v-model="content" @init="editorInit" lang="c_cpp" theme="chrome" width="100%" height="100%"></editor>           
        </b-col>
        <b-col class="mh-100 px-2 mh-100 overflow-auto">
          <Documentation />
        </b-col>        
      </b-row>         
      <b-row :class="`text-white px-2 flex-shrink-0 ${status.background}`">
        <b-col cols="12" class="py-1">
          {{ status.message }}
        </b-col>
      </b-row>
    </b-container>
    <b-modal id="bv-modal-simulator" 
      header-bg-variant="dark"
      header-text-variant="light"
      centered hide-footer title="Simulator">
      <Simulator />
    </b-modal>
  </div>
</template>

<script>

import Documentation from '../components/Documentation.vue'
import Simulator from '../components/Simulator.vue'
import DefaultCode from '../assets/default.c'
import helpers from '../helpers.js'

export default {
  name: 'App',
  components: {    
    Documentation,
    Simulator,
    editor: require('vue2-ace-editor')    
  },
  data() {
    return {
      status: {
        message: "Ready. Click 'Run' to try your game in the simulator.",
        background: 'bg-success'
      },      
      content: '',
      language: 'c',
      languageOptions: [
        { value: 'c', text: 'C'},
        { value: 'assembly', text: 'Assembly'}      
      ],
      firmware: '1.0.1',
      firmwareVerisons: [
        '1.0.1',
        '1.1.0_beta2'
      ]
    }
  },
  computed: {
    displayLanguage: function() {
      return this.languageOptions.find(o => o.value == this.language).text
    }
  },
  async mounted() {
    let code = this.$route.query.code

    if (code) {
      this.axios.get(`${process.env.VUE_APP_API_URL}/code/${code}`)
        .then((response) => {
          this.content = response.data
          this.language = response.headers['x-language']        
        })
        .catch((reason) => {          
          // todo change background color
          this.setStatus(reason.message, 'bg-error')          
        })

    } else {    
      this.content = DefaultCode
      this.language = 'c'
    }
  },  
  methods: {
    editorInit: function () {
      require('brace/ext/language_tools') //language extension prerequsite...
      require('brace/mode/c_cpp')                
      require('brace/mode/less')
      require('brace/theme/chrome')
      require('brace/snippets/javascript') //snippet
    },
    compileCode: async function() {      
      let config = {
        headers: {
          'Content-Type': 'application/text',
          'X-Language': this.language,
          'X-Version': this.firmware          
        }
      }

      return await this.axios.post(`${process.env.VUE_APP_API_URL}/build`, this.content, config)
    },
    setStatus: function(message, background) {
      this.status.message = message
      this.status.background = background
    },
    run: async function() {
      // TODO Save      
      this.setStatus('Compiling...', 'bg-info')
     
      try {
      let response = await this.compileCode()
      
      // Start Simulator
      this.$bvModal.show('bv-modal-simulator')

        console.log(response.data.binary)
      } catch (error) {
        this.setStatus(error.response.data, 'bg-danger')
      }
    },
    download: async function() {
      this.setStatus('Compiling...', 'bg-info')
      
      try {
        let response = await this.compileCode()              
        
        let blob = helpers.b64toBlob(response.data.binary, "application/octet-string")

        let link = document.createElement("a")        
        link.href = window.URL.createObjectURL(blob)        
        link.download = "fram.bin"
        link.click()

        this.setStatus('Successfully created fram.bin', 'bg-success')
      } catch (error) {
        this.setStatus(error.response.data, 'bg-danger')
      }
    }
  },
}
</script>

<style>
#app {
  font-family: "Helvetica Neue", Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: left;
}
</style>
