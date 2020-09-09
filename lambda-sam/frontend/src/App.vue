<template>
  <div id="app">
    <div class="container">
      <div class="row">
        <AddNote />
      </div>
      <div class="row">
        <NotesGrid />
      </div>
    </div>
  </div>
</template>

<script>
import NotesGrid from './components/NotesGrid.vue'
import AddNote from './components/AddNote.vue'

export default {
  name: 'App',
  components: {
    NotesGrid,
    AddNote,
  },
}
</script>
