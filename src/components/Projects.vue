<template>
  <div>
    <b-container>
      <b-row>
        <b-col>
          <div class="mt-2 mb-4">
            <b-img center :src="user.gravatar" rounded="circle" width="96" height="96"></b-img>
            <p class="text-center"><small class="text-muted">{{ user.email }}</small></p>
          </div>
        </b-col>
      </b-row>
      <b-row>        
        <b-col>                    
          <div class="mb-3">
            <b-input-group>
              <b-form-input v-model="newTitle" placeholder="New Project"></b-form-input>
              <b-input-group-append>
                <b-button :disabled="!canCreate" @click="createProject">Create</b-button>
              </b-input-group-append>              
            </b-input-group>            
          </div>          
          <b-list-group>
            <b-list-group-item href="#" v-for="project in projects" 
              :key="project.title" 
              :active="project.title==selected.title"
              @click="select(project.title)">
              {{ project.title }}
            </b-list-group-item>
          </b-list-group>
          <div class="mt-2 mb-3">
            <b-button variant="danger" :disabled="!hasSelection" @click="deleteProject">Delete Project</b-button>
          </div>
        </b-col>
      </b-row>
    </b-container>
  </div>
</template>

<script>
export default {
  name: 'Projects',  
  props: ['user'],
  data() {
    return {      
      selected: {},
      projects: [],
      newTitle: ''
    }
  },
  computed: {
    hasSelection: function() {
      return Object.keys(this.selected).length !== 0
    },
    canCreate: function() {
      return this.newTitle && this.newTitle.length !== 0
    }
  }, 
  mounted() {   
    this.axios.get(`${process.env.VUE_APP_API_URL}/projects`)
      .then((response) => {
        this.projects = response.data
      })
      .catch((reason) => {          
        console.log(reason)          
      }) 
  },
  methods: {
    select(title) {
      this.axios.get(`${process.env.VUE_APP_API_URL}/projects/${title}`)
        .then((response) => {
          this.selected = response.data          
          this.$emit('project-loaded', response.data)
        })
        .catch((reason) => {          
          console.log(reason)          
        }) 
    },
    createProject() {
      let title = this.newTitle
      let source = `#include <stdio.h>
#include <stdlib.h>
#include "api.h"

int main() {
	// Initialize the API
	api_init();

	// Clear the graphics in video memory
	CLEAR();

	for (;;) {
		// Game Logic
		
		SET_CURSOR(0, 0);
		DRAW_STRING("Hello, World!");

		// Push video memory to the OLED
		DISPLAY();

		// Wait for next interrupt
		WAIT();
	}

	return 0;
}`
      
      let config = {
        headers: {
          'Content-Type': 'application/text',
          'X-Language': 'c',
          'X-Version': '1.0.1'
        }
      }

      this.axios.post(`${process.env.VUE_APP_API_URL}/projects/${title}`, source, config)
        .then(() => {
          // Activate Project      
          let newProjects = [...this.projects]
          newProjects.push({title: title}) // Just the title, not the full project as below
          newProjects.sort((a, b) => (a.title > b.title) ? 1 : -1)          

          let project = {}
          project.title = title
          project.source = source
          project.language = 'c'
          project.firmware = '1.0.1'
          
          this.projects = newProjects
          this.selected = project

          this.$emit('project-loaded', project)
        })
        .catch((reason) => {
          console.log(reason)
          // TODO: show message
        })
    },
    deleteProject() {
      let title = this.selected.title
      this.axios.delete(`${process.env.VUE_APP_API_URL}/projects/${title}`)
        .then(() => {          
          this.selected = {}
          
          // Deselect project
          let project = {}
          project.title = ''
          project.content = ''
          project.language = 'c'
          project.firmware = '1.0.1'
          this.$emit('project-loaded', project)

          // Remove from Project List
          this.projects = this.projects.filter(p => p.title !== title)
        })
        .catch((reason) => {
          console.log(reason)
          // TODO: show message
        })
    }
  }
}
</script>

<style scoped>
</style>
