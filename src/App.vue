<template>
  <div id="app">
    <b-container fluid class="d-flex flex-column flex-grow-1 vh-100 px-0 overflow-hidden">    
      <b-navbar toggleable="lg" type="dark" variant="dark" class="flex-shrink-0">
        <b-navbar-brand>Dodo Build</b-navbar-brand>
        <b-navbar-toggle target="nav-collapse"></b-navbar-toggle>
        <b-collapse id="nav-collapse" is-nav>
          <b-navbar-nav>
            <b-nav-item-dropdown text="C">
              <b-dropdown-item href="#">C</b-dropdown-item>
              <b-dropdown-item href="#">Assembly</b-dropdown-item>
            </b-nav-item-dropdown>
            <b-nav-item-dropdown text="1.0.1">
              <b-dropdown-item href="#">1.0.1</b-dropdown-item>
              <b-dropdown-item href="#">1.1.0_beta2</b-dropdown-item>
            </b-nav-item-dropdown>
            <b-nav-form>            
              <b-button pill variant="info" class="ml-2"><span><b-icon icon="play" /> Run</span></b-button>
              <b-button pill variant="success" class="ml-2"><span><b-icon icon="download" /> fram.bin</span></b-button>
              <b-button pill variant="warning"  class="ml-2"><span><b-icon icon="lightning" /> Flash</span></b-button>
            </b-nav-form>          
          </b-navbar-nav>
          <b-navbar-nav class="ml-auto">
            <b-nav-form>
              <b-button variant="dark"><span><b-icon icon="github" /> Sign with GitHub</span></b-button>
            </b-nav-form>
          </b-navbar-nav>
        </b-collapse>
      </b-navbar>
      <b-row class="flex-grow-1 overflow-hidden" no-gutters>
        <b-col cols="8" class="mh-100 overflow-auto py-2">          
            <editor v-model="content" @init="editorInit" lang="c_cpp" theme="chrome" width="100%" height="100%"></editor>           
        </b-col>
        <b-col class="mh-100 text-left px-2 mh-100 overflow-auto">
          <VueShowdown :markdown="apiContent" class="markdown-body"></VueShowdown>
        </b-col>        
      </b-row>         
      <b-row class="text-left flex-shrink-0 bg-success">
        <b-col cols="12" class="py-1">
          Test 2
        </b-col>
      </b-row>
    </b-container>
  </div>
</template>

<script>

import api from './assets/api.md'

export default {
  name: 'App',
  components: {
    editor: require('vue2-ace-editor')
  },
  data() {
    return {
      apiContent: '',
    }
  },
  mounted() {
    this.apiContent = api
  },
  methods: {
    editorInit: function () {
      require('brace/ext/language_tools') //language extension prerequsite...
      require('brace/mode/c_cpp')                
      require('brace/mode/less')
      require('brace/theme/chrome')
      require('brace/snippets/javascript') //snippet
    }
  },
}
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;  
}
</style>
