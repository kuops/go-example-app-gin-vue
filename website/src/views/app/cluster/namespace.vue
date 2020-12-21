<template>
  <div class="app-container">
    <div class="filter-container">
      <el-form :inline="true">
        <el-form-item>
          <el-input placeholder="请输入要创建的命名空间" v-model="namespaceObj.name"></el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="createNS">创建命名空间</el-button>
        </el-form-item>
      </el-form>
    </div>

    <el-table
      v-loading="listLoading"
      :data="list"
      border
    >
      <el-table-column type="selection" width="55" />
      <el-table-column align="center" label="ID" width="95">
        <template slot-scope="scope">
          {{ (listQuery.page -1) * listQuery.pageSize + scope.$index + 1 }}
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

    <pagination v-show="total>0" :total="total" :page.sync="listQuery.page" :limit.sync="listQuery.pageSize" @pagination="getList" />

  </div>
</template>

<script>
import { getNamespace, deleteNamespace,createNamespace } from '@/api/app/cluster/namespace'
import waves from '@/directive/waves' // waves directive
import Pagination from '@/components/Pagination' // secondary package based on el-pagination

export default {
  name: 'Namespace',
  components: { Pagination },
  directives: { waves },
  data() {
    return {
      namespaceObj: {
        name: '',
      },
      tableKey: 0,
      list: null,
      total: 0,
      listLoading: true,
      listQuery: {
        page: 1,
        pageSize: 5,
        key: ''
      }
    }
  },
  created() {
    this.getList()
  },
  methods: {
    getList() {
      this.listLoading = true
      getNamespace(this.listQuery).then(response => {
        console.log(response.data)
        this.list = response.data.list.items
        this.total = response.data.total
        // Just to simulate the time of the request
        setTimeout(() => {
          this.listLoading = false
        }, 1.5 * 1000)
      })
    },
    deleteNS(name) {
      const data = { name: name }
      deleteNamespace(data).then(response => {
        this.getList()
      }
      )
    },
    handleFilter() {
      this.listQuery.page = 1
      this.getList()
    },
    handleCreate() {
      this.resetTemp()
      this.dialogStatus = 'create'
      this.dialogFormVisible = true
      this.$nextTick(() => {
        this.$refs['dataForm'].clearValidate()
      })
    },
    handleUpdate(row) {
      this.temp = Object.assign({}, row) // copy obj
      this.temp.timestamp = new Date(this.temp.timestamp)
      this.dialogStatus = 'update'
      this.dialogFormVisible = true
      this.$nextTick(() => {
        this.$refs['dataForm'].clearValidate()
      })
    },
    createNS() {
      createNamespace(this.namespaceObj).then(response => {
        if (response.code === 20000) {
          this.fetchData()
          this.$message.success("创建命名空间成功")
        } else  {
          this.$message.error("创建命名空间失败")
        }
        this.namespaceObj.name = ''
      })
    },
  }
}
</script>

<style lang="scss" scoped>
.filter-item {
  margin-left: 12px;
}
</style>
