<template>
  <div class="app-container">
    <div class="filter-container">
      <el-form :inline="true">
        <el-form-item label="命名空间">
          <el-select clearable v-model="selectedNS" placeholder="请选择" @change="getNamaspaceDeploymentsList(selectedNS)">
            <el-option
              v-for="item in deploymentsList"
              :key="item.metadata.namespace"
              :label="item.metadata.namespace"
              :value="item.metadata.namespace"
            >
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="createNS">创建示例 deployment </el-button>
        </el-form-item>
      </el-form>
    </div>

    <el-table
      v-loading="loading"
      :data="deploymentsList"
      border>
      <el-table-column align="center" label="ID" width="95">
        <template slot-scope="scope">
          {{ scope.$index }}
        </template>
      </el-table-column>
      <el-table-column label="NAMESPACE">
        <template slot-scope="scope">
          {{ scope.row.metadata.namespace }}
        </template>
      </el-table-column>
      <el-table-column label="NAME">
        <template slot-scope="scope">
          {{ scope.row.metadata.name }}
        </template>
      </el-table-column>
      <el-table-column label="Ready">
        <template slot-scope="scope">
          {{ scope.row.status.updatedReplicas || 0 }}/{{ scope.row.spec.replicas }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="200">
        <template slot-scope="scope">
          <el-button type="primary" @click="showEditDialog(scope.row)">编辑</el-button>
          <el-button type="danger" @click="showDeleteDialog(scope.row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog
      title="提示"
      :visible.sync="editDialogVisible"
      width="60%">
      <div class="editor-container">
        <yaml-editor  ref="yamlEditor" v-model="editDeployment" />
      </div>
      <span slot="footer" class="dialog-footer">
        <el-button @click="editDialogVisible = false">取 消</el-button>
        <el-button type="primary" @click="pathDeployment()">确 定</el-button>
      </span>
    </el-dialog>

    <el-dialog
      title="提示"
      :visible.sync="deleteDialogVisible"
      width="30%">
      <span v-if="deleteDeployment !== ''">确认删除 {{ deleteDeployment.metadata.name }} ?</span>
      <span slot="footer" class="dialog-footer">
        <el-button @click="deleteDialogVisible = false">取 消</el-button>
        <el-button type="primary" @click="deleteDeploy(deleteDeployment)">确 定</el-button>
      </span>
    </el-dialog>

    <pagination v-show="total>0" :total="total" :page.sync="listQuery.page" :limit.sync="listQuery.pageSize" @pagination="getDeploys()" />
  </div>
</template>

<script>
import {getDeploymentsList,pathNamespaceDeployments,deleteNamespaceDeployment,createExampleDeploy} from '@/api/app/cluster/deployment';
import YamlEditor from '@/components/YamlEditor'
import Pagination from '@/components/Pagination' // secondary package based on el-pagination

export default {
  name: 'Deployment',
  components: { YamlEditor, Pagination },
  data() {
    return {
      total: 0,
      namespacesList: null,
      selectedNS: '',
      deploymentsList: null,
      editDialogVisible: false,
      deleteDialogVisible: false,
      deleteDeployment: '',
      editDeployment: '',
      listQuery: {
        page: 1,
        pageSize: 5,
        key: '',
        namespace: ''
      }
    }
  },
  methods: {
    createNS() {
      createExampleDeploy().then(response => {
        if (response.code === 20000) {
          this.$message.success("创建 deployment 成功")
        }
      })
      this.$router.go(0)
    },
    getDeploys() {
      this.loading = true
      getDeploymentsList(this.listQuery).then(response => {
        this.deploymentsList = response.data.list.items
        this.loading = false
      })
    },
    getNamaspaceDeploymentsList(namespace) {
      this.listQuery.namespace = namespace
      getDeploymentsList(this.listQuery).then(response => {
        this.deploymentsList = response.data.list.items
        this.loading = false
      })
    },
    showEditDialog(data) {
      this.editDeployment = jsyaml.safeDump(data)
      this.editDialogVisible = true
    },
    pathDeployment() {
      const data = jsyaml.safeLoad(this.editDeployment)
      console.log(data)
      pathNamespaceDeployments(data).then(response => {
        if (response.code === 20000) {
          this.$message.success("修改 deployment 成功")
        }
        this.getDeploys()
      })
      this.editDialogVisible = false
    },
    showDeleteDialog(data) {
      this.deleteDeployment = data
      this.deleteDialogVisible = true
    },
    deleteDeploy(deploy) {
      const data = { name: deploy.metadata.name,namespace:  deploy.metadata.namespace}
      deleteNamespaceDeployment(data).then(response => {
        if (response.code === 20000) {
          this.$message.success("删除 deployment 成功")
        }
        this.getDeploys()
      })
      this.deleteDialogVisible = false
    }
  },
  created() {
    this.getDeploys()
  }
}
</script>

<style lang="scss" scoped>
.el-button {
  padding: 10px 15px;
}
.delete-dialog {
  text-align: center;
}
</style>
