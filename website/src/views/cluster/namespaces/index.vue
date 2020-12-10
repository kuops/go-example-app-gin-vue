<template>
  <div class="app-container">
    <div class="create-namespace">
          <el-input placeholder="请输入要创建的命名空间" v-model="namespaceObj.name"></el-input>
          <el-button type="primary" @click="createNS">创建命名空间</el-button>
    </div>
    <el-table
      v-loading="loading"
      :data="namespacesList"
      border>
      <el-table-column align="center" label="ID" width="95">
        <template slot-scope="scope">
          {{ scope.$index }}
        </template>
      </el-table-column>
      <el-table-column label="NAME">
        <template slot-scope="scope">
          {{ scope.row.metadata.name }}
        </template>
      </el-table-column>
      <el-table-column label="STATUS">
        <template slot-scope="scope">
          {{ scope.row.status.phase }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="200">
        <template slot-scope="scope">
          <el-button type="danger" @click="deleteNS(scope.row.metadata.name)">删除</el-button>
        </template>

      </el-table-column>
    </el-table>
  </div>
</template>

<script>
import { getNamespaces,createNamespace,deleteNamespace } from '@/api/cluster'

export default {
  data() {
    return {
      loading: true,
      namespaceObj: {},
      namespacesList: null,
    }
  },
  created() {
    this.fetchData()
  },
  methods: {
    fetchData() {
      this.loading = true
      getNamespaces().then(response => {
        this.namespacesList = response.data.items
        this.loading = false
      })
    },
    createNS() {
      createNamespace(this.namespaceObj).then(response => {
        console.log(response.code)
        if (response.code === 20000) {
          this.fetchData()
          this.$message.success("创建命名空间成功")
        } else  {
          this.$message.error("创建命名空间失败")
        }
        })
    },
    deleteNS(name) {
      const deleteObj = Object.assign({},{name:name})
      deleteNamespace(deleteObj).then(response => {
        this.fetchData()
        }
      )
    }
  }
}
</script>

<style lang="scss" scoped>
.el-button {
  padding: 10px 15px;
}
.create-namespace {
  margin-bottom: 15px;
}
.el-input {
  width: 200px;
  margin-right: 10px;
}
</style>
