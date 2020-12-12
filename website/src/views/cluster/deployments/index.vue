<template>
  <div class="app-container">
    <el-form>
      <el-form-item label="命名空间">
        <el-select clearable v-model="selectedNS" placeholder="请选择" @change="getNamespaceDeploys(selectedNS)">
          <el-option
            v-for="item in namespacesList"
            :key="item.metadata.name"
            :label="item.metadata.name"
            :value="item.metadata.name"
            >
          </el-option>
        </el-select>
      </el-form-item>
    </el-form>
    <el-table
      v-loading="loading"
      :data="deploymentsList"
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
  </div>
</template>

<script>
import {getNamespaces, getDeployments,getNamespaceDeployments,pathNamespaceDeployments,deleteNamespaceDeployment} from "@/api/cluster";
import YamlEditor from '@/components/YamlEditor'

export default {
  components: { YamlEditor },
  data() {
    return {
      namespacesList: null,
      selectedNS: '',
      deploymentsList: null,
      editDialogVisible: false,
      deleteDialogVisible: false,
      deleteDeployment: '',
      editDeployment: '',
    }
  },
  methods: {
    getNS() {
      this.loading = true
      getNamespaces().then(response => {
        this.namespacesList = response.data.items
        this.loading = false
      })
    },
    getDeploys() {
      this.loading = true
      getDeployments().then(response => {
        this.deploymentsList = response.data.items
        this.loading = false
      })
    },
    getNamespaceDeploys(name) {
      this.loading = true
      if (name === "") {
        return getDeployments().then(response => {
          this.deploymentsList = response.data.items
          this.loading = false
        })
      }
        getNamespaceDeployments(name).then(response => {
          this.deploymentsList = response.data.items
          this.loading = false
        })
      console.log(this.deploymentsList)
    },
    showEditDialog(data) {
      this.editDeployment = jsyaml.safeDump(data)
      this.editDialogVisible = true
    },
    pathDeployment() {
      const data = jsyaml.safeLoad(this.editDeployment)
      console.log(data)
      pathNamespaceDeployments(data.metadata.namespace,data.metadata.name,data).then(response => {
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
    deleteDeploy(data) {
      deleteNamespaceDeployment(data.metadata.namespace,data.metadata.name).then(response => {
        if (response.code === 20000) {
          this.$message.success("删除 deployment 成功")
        }
        this.getDeploys()
      })
      this.deleteDialogVisible = false
    }
  },
  created() {
    this.getNS()
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
