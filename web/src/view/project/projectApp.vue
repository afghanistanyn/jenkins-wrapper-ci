<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline" :rules="searchRule"
        @keyup.enter="onSubmit">

        <el-form-item label="应用名">
          <el-input v-model="searchInfo.name" placeholder="" clearable size="small" style="width: 240px"></el-input>
        </el-form-item>
        <el-form-item label="git仓库">
          <el-input v-model="searchInfo.gitRepo" placeholder="" clearable size="small" style="width: 240px"></el-input>
        </el-form-item>
        <el-form-item label="镜像地址">
          <el-input v-model="searchInfo.image" placeholder="" clearable size="small" style="width: 240px"></el-input>
        </el-form-item>
        <el-form-item label="发布参数">
          <el-input v-model="searchInfo.buildParam" placeholder="" clearable size="small" style="width: 240px"></el-input>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button v-if="isSuperAdmin() || isOps() || isProjectManager" type="primary" icon="plus"
          @click="openDialog">新增应用</el-button>
      </div>
      <el-table ref="multipleTable" style="width: 100%" tooltip-effect="dark" :data="tableData" row-key="ID" border
        @selection-change="handleSelectionChange">

        <!-- <el-table-column type="selection" width="55" /> -->
        <el-table-column align="left" label="项目名称" width="150">
          {{ project.nameCn }}
        </el-table-column>
        <el-table-column align="left" label="应用名称" prop="name" width="150" />
        <el-table-column align="left" label="描述" prop="description" width="150" />
        <el-table-column align="left" label="git仓库" prop="gitRepo" width="350" />
        <el-table-column align="left" label="发布参数" prop="buildParams">
          <template #default="scope">
            {{ formatArray(scope.row.buildParams) }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="镜像地址" prop="image" />
        <el-table-column align="left" label="创建日期" width="110" >
          <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>
        <el-table-column align="left" label="操作" width="130">
          <template #default="scope">
            <div class="table-button-group">
              <el-button v-if="isSuperAdmin() || isOps() || isProjectManager" type="primary" link icon="edit"
                class="table-button" @click="updateAppFunc(scope.row)">修改</el-button>
              <el-button v-if="isSuperAdmin() || isOps() || isProjectManager" type="primary" link icon="delete"
                class="table-button" @click="deleteRow(scope.row)">删除</el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination layout="total, sizes, prev, pager, next, jumper" :current-page="page" :page-size="pageSize"
          :page-sizes="[10, 20, 30, 50]" :total="total" @current-change="handleCurrentChange"
          @size-change="handleSizeChange" />
      </div>
    </div>
    <el-dialog v-model="dialogFormVisible" :before-close="closeDialog" :title="type === 'create' ? '添加应用' : '修改应用'"
      destroy-on-close>
      <el-form :model="formData" label-position="right" ref="elFormRef" :rules="rule" label-width="100px">
        <el-form-item label="所属项目id:" prop="projectId">
          <el-input v-model="formData.projectId" disabled />
        </el-form-item>
        <el-form-item label="所属项目:" prop="projectName">
          <el-input v-model="formData.projectName" disabled />
        </el-form-item>
        <el-form-item label="应用名称:" prop="name">
          <el-input v-model="formData.name" :clearable="true" :disabled="type === 'update'" placeholder="请输入" />
        </el-form-item>
        <el-form-item label="描述:" prop="description">
          <el-input v-model="formData.description" :clearable="true" :disabled="formData.customConfig !== ''" placeholder="请输入项目描述" />
        </el-form-item>
        <el-form-item label="发布参数:" prop="buildParams">
          <el-input v-model="formData.buildParams" :clearable="true" :disabled="formData.customConfig !== ''" placeholder="请输入发布参数" />
        </el-form-item>

        <el-form-item label="git仓库:" prop="gitRepo">
          <el-input v-model="formData.gitRepo" :clearable="true" :disabled="formData.customConfig !== ''" placeholder="请输入git仓库地址" />
        </el-form-item>

        <el-form-item label="Jenkinsfile:" prop="jenkinsFile">
          <el-input v-model="formData.jenkinsFile" :clearable="true" :disabled="formData.customConfig !== ''" placeholder="请输入Jenkinsfile路径" />
        </el-form-item>

        <el-form-item label="镜像地址:" prop="image">
          <el-input v-model="formData.image" :clearable="true" :disabled="formData.customConfig !== ''" placeholder="请输入镜像仓库地址" />
        </el-form-item>

        <el-collapse v-model="activeCollapses">
          <el-collapse-item title="高级设置" name="0">
          <el-form-item label="自定义jenkinsConfig" prop="customConfig">
            <el-input v-model="formData.customConfig" @input="customConfigChange" type="textarea" :autosize=" { minRows: 10 } " />
          </el-form-item>
        </el-collapse-item>
        </el-collapse>


      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closeDialog">取 消</el-button>
          <el-button type="primary" @click="enterDialog">确 定</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script>
export default {
  name: 'App'
}
</script>

<script setup>
import {
  createApp,
  deleteApp,
  deleteAppByIds,
  updateApp,
  findApp,
  getAppList
} from '@/api/app'


import {
  findProject,
} from '@/api/project'

// 全量引入格式化工具 请按需保留
import { getDictFunc, formatDateTime, formatDate, formatBoolean, filterDict, formatArray } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive, onActivated, onMounted } from 'vue'
import { useUserStore } from '@/pinia/modules/user'
import { useRouter, useRoute } from 'vue-router'

//  router
const router = useRouter()
const route = useRoute()

// userInfo store
const userStore = useUserStore()
const currentUserId = userStore.userInfo.ID
const currentUserAuthorityId = userStore.userInfo.authority.authorityId

// 自动化生成的字典（可能为空）以及字段
const formData = ref({
  projectId: '',
  projectName: '',
  name: '',
  description: '',
  gitRepo: '',
  jenkinsFile: 'Jenkinsfile',
  buildParams: '',
  image: '',
  customConfig: '',
})

// 验证规则
const rule = reactive({
  name: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
  {
    whitespace: true,
    message: '不能只输入空格',
    trigger: ['input', 'blur'],
  }],
})

const searchRule = reactive({
  createdAt: [
    {
      validator: (rule, value, callback) => {
        if (searchInfo.value.startCreatedAt && !searchInfo.value.endCreatedAt) {
          callback(new Error('请填写结束日期'))
        } else if (!searchInfo.value.startCreatedAt && searchInfo.value.endCreatedAt) {
          callback(new Error('请填写开始日期'))
        } else if (searchInfo.value.startCreatedAt && searchInfo.value.endCreatedAt && (searchInfo.value.startCreatedAt.getTime() === searchInfo.value.endCreatedAt.getTime() || searchInfo.value.startCreatedAt.getTime() > searchInfo.value.endCreatedAt.getTime())) {
          callback(new Error('开始日期应当早于结束日期'))
        } else {
          callback()
        }
      }, trigger: 'change'
    }
  ],
})

const elFormRef = ref()
const elSearchFormRef = ref()

// =========== 表格控制部分 ===========
const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({
  'projectId': route.params.projectId,
  'name': undefined,
  'buildParam': undefined,
  'gitRepo': undefined,
  'image': undefined,
})
const project = ref({
  'ID': undefined,
  'name': undefined,
  'nameCn': undefined,
})
const isProjectManager = ref(false)
const defaultBuildParams = ['branch', 'build_env']
const activeCollapses = ref([])

const isSuperAdmin = () => {
  // todo get superAdminRoles from backend
  const superAdminRoles = [888, 999]
  for (let authority of superAdminRoles) {
    if (currentUserAuthorityId === authority) {
      return true
    }
  }
  return false
}

const isOps = () => {
  const opsRoleId = 999
  if (currentUserAuthorityId === opsRoleId) {
    return true
  }
  return false
}

// 重置
const onReset = () => {
  searchInfo.value = {
    'projectId': route.params.projectId,
  }
  getTableData()
}

// 搜索
const onSubmit = () => {
  elSearchFormRef.value?.validate(async (valid) => {
    if (!valid) return
    page.value = 1
    pageSize.value = 10
    getTableData()
  })
}

// 分页
const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

// 修改页面容量
const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

onMounted(() => {
  initProject(route.params.projectId)
})

onActivated(() => {
  initProject(route.params.projectId)
  getTableData()
})


const initProject = async (projectId) => {
  const project_resp = await findProject({ "ID": projectId })
  if (project_resp.code === 0) {
    project.value = project_resp.data.project
    project.value.managers.map(user => {
      if (currentUserId === user.ID) {
        isProjectManager.value = true
      }
    })
  }
}

// 查询
const getTableData = async () => {
  const table = await getAppList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}

getTableData()

// ============== 表格控制部分结束 ===============

// 获取需要的字典 可能为空 按需保留
const setOptions = async () => {
}

// 获取需要的字典 可能为空 按需保留
setOptions()


// 多选数据
const multipleSelection = ref([])
// 多选
const handleSelectionChange = (val) => {
  multipleSelection.value = val
}

// 删除行
const deleteRow = (row) => {
  ElMessageBox.confirm('确定要删除吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    deleteAppFunc(row)
  })
}


// 批量删除控制标记
const deleteVisible = ref(false)

// 多选删除
const onDelete = async () => {
  const ids = []
  if (multipleSelection.value.length === 0) {
    ElMessage({
      type: 'warning',
      message: '请选择要删除的数据'
    })
    return
  }
  multipleSelection.value &&
    multipleSelection.value.map(item => {
      ids.push(item.ID)
    })
  const res = await deleteAppByIds({ ids })
  if (res.code === 0) {
    ElMessage({
      type: 'success',
      message: '删除成功'
    })
    if (tableData.value.length === ids.length && page.value > 1) {
      page.value--
    }
    deleteVisible.value = false
    getTableData()
  }
}

// 行为控制标记（弹窗内部需要增还是改）
const type = ref('')

// 更新行
const updateAppFunc = async (row) => {
  const res = await findApp({ ID: row.ID })
  type.value = 'update'
  if (res.code === 0) {
    formData.value = res.data.app
    formData.value.projectName = project.value.nameCn
    if(formData.value.customConfig !== '') {
      formData.value.gitRepo = ''
      formData.value.buildParams = ''
      formData.value.jenkinsFile = ''
      activeCollapses.value  = ['0']
    }

    dialogFormVisible.value = true
  }
}


// 删除行
const deleteAppFunc = async (row) => {
  const res = await deleteApp({ ID: row.ID })
  if (res.code === 0) {
    ElMessage({
      type: 'success',
      message: '删除成功'
    })
    if (tableData.value.length === 1 && page.value > 1) {
      page.value--
    }
    getTableData()
  }
}

// 弹窗控制标记
const dialogFormVisible = ref(false)

// 打开弹窗
const openDialog = () => {
  type.value = 'create'
  formData.value.buildParams = defaultBuildParams
  formData.value.projectId = parseInt(route.params.projectId)
  formData.value.projectName = project.value.nameCn
  dialogFormVisible.value = true
}

// 关闭弹窗
const closeDialog = () => {
  dialogFormVisible.value = false
  formData.value = {
    projectId: '',
    name: '',
    description: '',
    gitRepo: '',
    image: '',
    buildParams: '',
    jenkinsFile: 'Jenkinsfile',
    customConfig: '',
  }
  activeCollapses.value = []
}

const customConfigChange = (val) => {
    formData.value.customConfig = val
    if(formData.value.customConfig !== '') {
      formData.value.gitRepo = undefined
      formData.value.buildParams = undefined
      formData.value.jenkinsFile = undefined
      formData.value.image = undefined
    }
}


// 弹窗确定
const enterDialog = async () => {
  elFormRef.value?.validate(async (valid) => {
    if (!valid) return
    let res

    if(formData.value.customConfig !== '') {
      formData.value.gitRepo = undefined
      formData.value.buildParams = undefined
      formData.value.jenkinsFile = undefined
      formData.value.image = undefined
    }

    //  reformat formData.buildParams
    if (typeof formData.value.buildParams === 'string') {
      formData.value.buildParams = formData.value.buildParams.split(',')
    }
    switch (type.value) {
      case 'create':
        res = await createApp(formData.value)
        break
      case 'update':
        res = await updateApp(formData.value)
        break
      default:
        res = await createApp(formData.value)
        break
    }
    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: '创建/更改成功'
      })
      closeDialog()
      getTableData()
    }
  })
}


</script>

<style scoped>
.table-button-group {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
}

.table-button {
  margin-top: 2px;
  margin-left: 10px;
}
</style>
