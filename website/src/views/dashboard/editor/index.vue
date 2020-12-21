<template>
  <div class="dashboard-editor-container">
  </div>
</template>

<script>
import { mapGetters } from 'vuex'

export default {
  name: 'DashboardEditor',
  data() {
    return {
    }
  },
  computed: {
    ...mapGetters([
      'name',
      'avatar',
      'roles'
    ])
  }
}
</script>

<style lang="scss" scoped>
.dashboard-editor {
  &-container {
    position: relative;
    padding: 32px;
    background: #f0f2f5;
  }
}
</style>
