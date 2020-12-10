<template>
  <div class="home-container">
    <el-row>
      <el-card class="home-card">
        <h4>早安，管理员， 请开始您一天的工作吧！</h4>
        <p class="home-text">
          <i class="el-icon-sunny el-icon home-icon"></i>
          <span>今日晴，0℃ - 10℃，天气寒冷，注意添加衣物。</span>
        </p>
      </el-card>
    </el-row>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'

export default {
  name: 'Home',
  computed: {
    ...mapGetters([
      'name'
    ])
  }
}
</script>

<style lang="scss" scoped>
.home {
  &-container {
    margin: 30px;
  }
  &-card {
    width: 100%;
  }
  &-text {
     margin-top: 8px;
   }
  &-icon {
    margin-right: 8px;
  }
}
</style>
