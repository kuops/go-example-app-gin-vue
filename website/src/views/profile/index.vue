<template>
  <div class="app-container">
    <div v-if="user">
      <user-card :user="user" />
    </div>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import UserCard from './components/UserCard'

export default {
  name: 'Profile',
  components: { UserCard },
  computed: {
    ...mapGetters([
      'name',
      'avatar',
    ])
  },
  data() {
    return {
      user: {},
    }
  },
  created() {
    this.getUser()
  },
  methods: {
    getUser() {
      this.user = {
        name: this.name,
        email: 'admin@test.com',
        avatar: this.avatar
      }
    }
  }
}
</script>

<style lang="scss" scoped>
</style>
