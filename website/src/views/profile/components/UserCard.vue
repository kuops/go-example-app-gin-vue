<template>
  <el-row :gutter="40">
    <el-col :xs="24" :sm="24" :lg="6">
      <el-card>
        <div slot="header" class="clearfix">
          <span>个人信息</span>
        </div>
        <div class="user-avatar">
          <el-avatar
            :size="240"
            :src="user.avatar"
          ></el-avatar>
        </div>
        <div class="user-personality">
          <p class="nickname">{{user.name}}</p>
          <p class="person-info">这个家伙很懒，什么都没有留下</p>
        </div>
        <div class="user-information">
          <ul>
            <li>
              <i class="el-icon-user"></i>资深前端工程师
            </li>
            <li>
              <i class="el-icon-data-analysis"></i>北京反转极光科技有限公司-技术部-前端事业群
            </li>
            <li>
              <i class="el-icon-video-camera-solid"></i>中国·北京市·朝阳区
            </li>
          </ul>
        </div>
      </el-card>
    </el-col>
    <el-col :xs="24" :sm="24" :lg="18">
      <el-card>
        <div slot="header" class="clearfix">
          <span>修改个人信息</span>
        </div>
        <el-form>
          <el-form-item label="Name">
            <el-input v-model.trim="user.name" />
          </el-form-item>
          <el-form-item label="Email">
            <el-input v-model.trim="user.email" />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="submit">Update</el-button>
          </el-form-item>
        </el-form>
      </el-card>
    </el-col>
  </el-row>
</template>

<script>
export default {
  props: {
    user: {
      type: Object,
      default: () => {
        return {
          name: '',
          email: '',
          avatar: '',
        }
      }
    }
  },
  methods: {
    submit() {
      this.$message({
        message: 'User information has been updated successfully',
        type: 'success',
      })
    }
  }
}
</script>

<style lang="scss" scoped>
.el-card {
  margin-bottom: 10px;
}
.user-avatar {
  text-align: center;
}
.user-personality {
  padding: 24px 0;
  text-align: center;

p {
  font-size: 16px;
}

.nickname {
  font-size: 26px;
}

.person-info {
  margin-top: 6px;
  font-size: 14px;
  color: #999
}

}
.user-information {
  width: 100%;
  height: 100%;
  text-align: left;
    ul {
      display: inline-block;
      height: 100%;
li {
i {
  margin-right: 8px;
}
padding: 20px 0;
font-size: 16px;
font-weight: 400;
color: #606266;
}
}
}
</style>
