<template>
  <div>
    <h2>Notes</h2>
     <table class="u-full-width">
       <thead>
          <tr>
            <th>Title</th>
            <th>Description</th>
            <th>Actions</th>
          </tr>
       </thead>
       <tbody>
          <tr v-for="note in notes" v-bind:key="note.id">
            <td>{{ note.title }}</td>
            <td>{{ note.description }}</td>
            <td>
              <button class="danger" @click="removeNote(note.id)">Delete</button>
            </td>
          </tr>
       </tbody>
    </table> 
  </div>
</template>

<style lang="css" scoped>
.danger {
  background: #ff5e5e;
  color: white;
}

</style>

<script>

export default {
  name: 'NotesGrid',

  data() {
    return {
      intervalPtr: undefined,
    }
  },

  mounted() {
    this.intervalPtr = setInterval(() => {
      this.$store.dispatch('fetchNotes')
    }, 2000)
  },

  beforeDestroy() {
    if (this.intervalPtr !== undefined) {
      clearInterval(this.intervalPtr)
    }
  },

  computed: {
    notes() {
      return this.$store.state.notes
    }
  },

  methods: {
    removeNote(noteID) {
      this.$store.dispatch('removeNote', noteID).then(() => {
        alert('removed')
      })
    }
  }
}
</script>