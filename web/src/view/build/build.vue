<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline" :rules="searchRule"
        @keyup.enter="onSubmit">
        <el-form-item label="项目组">
          <el-select v-model="searchInfo.projectId" filterable :clearable="true">
            <el-option v-for="proj in projects" :key="proj.ID" :label="proj.name" :value="proj.ID" />
          </el-select>
        </el-form-item>
        <el-form-item label="应用">
          <el-input v-model="searchInfo.appName" />
        </el-form-item>
        <el-form-item label="发布状态">
          <el-input v-model="searchInfo.result" />
        </el-form-item>

        <el-form-item label="待审核" v-if="tobeApproveTotal > 0">
          <el-badge :value="tobeApproveTotal" class="item">
            <el-switch v-model="searchInfo.approveStatus" active-color="#409EFF" inactive-color="#C0CCDA" active-value=1
              inactive-value=0>
            </el-switch>
          </el-badge>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="openDialog">发布应用</el-button>
      </div>
      <el-table ref="multipleTable" style="width: 100%" tooltip-effect="dark" :data="tableData" row-key="ID" border
        @selection-change="handleSelectionChange">

        <el-table-column align="left" label="项目" prop="projectName" width="120" />
        <el-table-column align="left" label="应用" width="170" >
          <template #default="scope">
            <div v-if="scope.row.buildNumber !== 0">
            {{ scope.row.appName }} #{{ scope.row.buildNumber }}
            </div>
            <div v-else>
              {{ scope.row.appName }}
            </div>
          </template>
        </el-table-column>
        <el-table-column align="left" label="发布参数" prop="buildParamValues" width="200">
          <template #default="scope">
            <div v-for="(paramValue, param) in scope.row.buildParamValues" :key="param" style="display: flex; align-items: center;">
              {{ param }}: <div style="flex:1;display: flex;justify-content: flex-end;"><el-tag  type="info" style="width:70px;margin-right: 10px;">{{ paramValue }}</el-tag></div>
            </div>
          </template>
        </el-table-column>
        <!-- <el-table-column align="left" label="buildNumber" width="110">
          <template #default="scope">
            <div v-if="scope.row.buildNumber !== 0">
              <el-tag type="info" style="margin-left: 30px">{{ scope.row.buildNumber }}</el-tag>
            </div>
          </template>
        </el-table-column> -->
        <el-table-column align="left" label="发布者" width="90">
          <template #default="scope"> {{ formatUser(scope.row.createdBy) }}</template>
        </el-table-column>
        <el-table-column align="left" label="发布内容" prop="buildInfo" />


        <el-table-column align="left" label="审核人员" width="90">
          <template #default="scope"> {{ formatUser(scope.row.approveBy) }}</template>
        </el-table-column>
        <!-- <el-table-column align="left" label="git变更" prop="changes" /> -->
        <!-- <el-table-column align="left" label="gitCommit" prop="gitCommit" width="120" /> -->
        <el-table-column align="left" label="审核状态" prop="approveStatus" width="120">
          <template #default="scope">
            <el-tag v-if="scope.row.approveStatus === 1" type="" effect="dark">待审批</el-tag>
            <el-tag v-if="scope.row.approveStatus === 2" type="success" effect="dark">已批准</el-tag>
            <el-tag v-if="scope.row.approveStatus === 3" type="danger" effect="dark">已拒绝</el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="发布结果" width="120">
          <template #default="scope">
            <el-tag v-if="scope.row.result === 1" type="success" effect="dark">成功</el-tag>
            <el-tag v-if="scope.row.result === 2" type="danger" effect="dark">失败</el-tag>
            <el-tag v-if="scope.row.result === 3" type="" effect="dark">发布中</el-tag>
            <el-tag v-if="scope.row.result === 4" type="info" effect="dark">未知</el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="创建时间" width="160">
          <template #default="scope">{{ formatDateTime(scope.row.CreatedAt) }}</template>
        </el-table-column>

        <el-table-column align="left" label="操作" width="130">
          <template #default="scope">
            <div class="table-button-group">
              <el-button v-if="scope.row.approveStatus === 1 && hasProjectManagerPermission(scope.row.projectId)" size=""
                type="text" icon="unlock" class="table-button" @click="handleApproveBuild(scope.row)">批准发布</el-button>
              <el-button v-if="scope.row.approveStatus === 1 && hasProjectManagerPermission(scope.row.projectId)" size=""
                type="text" icon="close" class="table-button" @click="handleRejectBuild(scope.row)">拒绝发布</el-button>
              <el-button size="" :disabled="!(scope.row.approveStatus === 2)" type="text" icon="tickets"
                class="table-button" @click="handleViewBuildLog(scope.row)">查看日志</el-button>
              <el-button size="" type="text" icon="document-copy" class="table-button"
                @click="handleReCreateBuild(scope.row)">再次发布</el-button>
              <el-button
                v-if="scope.row.approveStatus === 1 && (hasBuildOwnerPermission(scope.row.createBy) || hasProjectManagerPermission(scope.row.projectId))"
                size="" type="text" icon="delete" @click="handleDeleteBuild(scope.row)">删除发布</el-button>
              <!-- <el-button v-if="scope.row.approveStatus === 2 && isSuperAdmin()" size="" type="text" icon="unlock"
                class="table-button" @click="viewBuildDetail(scope.row.ID)">查看详情</el-button> -->
            </div>
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination layout="total, sizes, prev, pager, next, jumper" :current-page="page" :page-size="pageSize"
          :page-sizes="[10, 30, 50, 100]" :total="total" @current-change="handleCurrentChange"
          @size-change="handleSizeChange" />
      </div>
    </div>
    <el-dialog v-model="dialogFormVisible" :before-close="closeDialog" :title="type === 'create' ? '发布应用' : '修改发布'"
      destroy-on-close :close-on-click-modal="false">
      <el-form :model="formData" label-position="right" ref="elFormRef" :rules="rule" label-width="100px">
        <el-form-item label="项目组:" prop="projectId">
          <el-select v-model="formData.projectId" filterable :clearable="true" @change="onProjectChanged">
            <el-option v-for=" proj  in  projects " :key="proj.ID" :label="proj.name" :value="proj.ID" />
          </el-select>
        </el-form-item>
        <!-- <el-form-item label="应用名" prop="appName">
          <el-input v-model="formData.appName" disabled />
        </el-form-item> -->

        <el-form-item label="应用:" prop="appId" v-show="formData.projectId !== undefined && formData.projectId !== ''">
          <el-select v-model="formData.appId" filterable :clearable="true" @change="onAppChanged">
            <el-option v-for=" app  in  project.apps " :key="app.ID" :label="app.name" :value="app.ID" />
          </el-select>
        </el-form-item>

        <template v-for="( paramValue, param ) in  formData.buildParamValues " :key="param">
          <el-form-item :label="param" v-if:="formData.appId !== undefined && formData.appId !== ''"
            :prop="`buildParamValues.${param}`" :rules="{ required: true, message: '发布参数不能为空', trigger: 'blur' }">
            <el-input v-model="formData.buildParamValues[param]" :placeholder="buildParamLabel" />
          </el-form-item>
        </template>

        <el-form-item label="发布内容:" prop="BuildInfo">
          <el-input v-model="formData.BuildInfo" type="textarea" :autosize="{ minRows: 4 }" clearable
            placeholder="请填写发布内容" />
        </el-form-item>

        <el-form-item label="审核状态:" prop="status">
          <el-select v-model="formData.approveStatus" placeholder="请选择" style="width:100%" :clearable="true" disabled>
            <el-option v-for=" item  in  approveStatusOptions " :key="item.key" :label="item.label" :value="item.key" />
          </el-select>
        </el-form-item>

      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closeDialog">取 消</el-button>
          <el-button type="primary" @click="enterDialog">确 定</el-button>
        </div>
      </template>
    </el-dialog>



    <el-dialog class="logDialog" v-model="viewBuildLogEnabled" :close="closeViewBuildLogDialog" :close-on-click-modal="false" top="2%" width="75%">
      <template #header="{ titleId, titleClass }">
      <div class="el-dialog-header">
        <h4 :id="titleId" :class="titleClass">{{ viewLogTitle }}</h4>
        <span class="header-sub-title">{{viewLogSubTitle}}</span>
      </div>
    </template>
      <div style="max-height: 700px; overflow-y: auto; overflow-x: hidden;">
        <Codemirror id="codemirrorLog" v-model:value="buildLog" :options="cmOptions" border ref="cmRef" />
      </div>
    </el-dialog>

  </div>
</template>

<script>
export default {
  name: 'Build'
}
</script>

<script setup>
import {
  createBuild,
  deleteBuild,
  deleteBuildByIds,
  updateBuild,
  findBuild,
  getBuildList,
  approveBuild,
  getBuildInfo,
} from '@/api/build'


import {
  findProject,
  getProjectList
} from '@/api/project'

import {
  getUserList
} from '@/api/user'

// 全量引入格式化工具 请按需保留
import { getDictFunc, formatDateTime, formatBoolean, filterDict } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive, onMounted, onActivated, computed, nextTick } from 'vue'


// codemirror log
import Codemirror from "codemirror-editor-vue3"
import 'codemirror/theme/idea.css'


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
  projectId: undefined,
  projectName: undefined,
  appId: undefined,
  appName: undefined,
  buildParamValues: {},
  approveStatus: 1,
})
const buildParamLabel = "发布参数"

var users = ref()
var projects = ref()
var project = ref({
  apps: [],
}
)
var app = ref({})

const approveStatusOptions = [
  {
    key: 1,
    label: "待审批"
  },
  {
    key: 2,
    label: "已批准"
  },
  {
    key: 3,
    label: "已拒绝"
  },
]

// 验证规则
const rule = reactive({
  projectId: [
    {
      required: true,
      message: '请选择项目',
      trigger: ['blur'],
    }],
  appId: [{
    required: true,
    message: '请选择应用',
    trigger: ['blur'],
  }],
  BuildInfo: [{
    required: true,
    message: '请填写发布内容',
    trigger: ['blur'],
  }],
  approveStatus: [{
    required: true,
    message: '请选择审核状态',
    trigger: ['blur'],
  }],
})

const searchRule = reactive({
})

const elFormRef = ref()
const elSearchFormRef = ref()

// =========== 表格控制部分 ===========
const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({
  projectId: undefined,
  projectName: undefined,
  appName: undefined,
  result: undefined,
  mine: true,
})

const viewBuildLogEnabled = ref(false)
const viewLogTitle = ref()
const viewLogSubTitle = ref()
const buildLog = ref()
const buildDuration = ref()
const timer = ref()
const cmRef = ref()
const cmOptions = {
  mode: "fclog",
  theme: "idea",
  tabSize: 4,
  indentUnit: 4,
  lineNumber: true,
  readOnly: true,
}


// 重置
const onReset = () => {
  searchInfo.value = {
    mine: true
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


//  user format
const formatUser = (userId) => {
  if (userId === 0 || userId == undefined) {
    return
  }
  for (var u of users.value) {
    if (u.ID === userId) {
      return u.nickName
    }
  }

}

// 查询
const getTableData = async () => {
  const table = await getBuildList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize

  }
}



// my projects
const getProjects = async () => {
  const resp = await getProjectList({ mine: true })
  if (resp.code === 0) {
    projects.value = resp.data.list
  }
}


//  all users
const getUsers = async () => {
  // get without pageinfo for all users 
  const resp = await getUserList({})
  if (resp.code === 0) {
    users.value = resp.data.list
  }
}


const getProjectConsole = async (projectId) => {
  const resp = await findProject({ ID: row.ID })
  if (resp.code === 0) {
    projects.value = resp.data.list
  }
}

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

const isProjectManager = (project) => {
  for (let user of project.managers) {
    if (currentUserId === user.ID) {
      return true
    }
  }
  return false
}

const getProjectById = (projectId) => {
  for (let proj in projects.value) {
    if (proj.Id === projectId) {
      return proj
    }
  }
}

const hasProjectManagerPermission = (projectId) => {
  //  find project
  let project = getProjectById(projectId)

  let hasPermission = isSuperAdmin() || isOps() || isProjectManager(project)
  // console.log("hasPermission", hasPermission)
  return hasPermission
}

const hasBuildOwnerPermission = (ownerId) => {
  if (ownerId === currentUserId) {
    return true
  }
  return false
}

const initBuildDatas = async () => {
  // get without pageinfo for all users 
  const user_resp = await getUserList({})
  if (user_resp.code === 0) {
    users.value = user_resp.data.list
  }

  const proj_resp = await getProjectList({ mine: true })
  if (proj_resp.code === 0) {
    projects.value = proj_resp.data.list
  }
  getTableData()
}

onActivated(() => {
  getTableData()
})

const tobeApproveTotal = computed(() => {
  var t = 0
  for (var build of tableData.value) {
    if (build.approveStatus === 1) {
      // 
      if (hasProjectManagerPermission(build.projectId)) {
        t += 1
      }
    }
  }
  return t
})


initBuildDatas()

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
const handleDeleteBuild = (row) => {
  ElMessageBox.confirm('确定要删除本次发布吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    deleteBuildFunc(row)
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
  const res = await deleteBuildByIds({ ids })
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
const updateBuildFunc = async (row) => {
  const res = await findBuild({ ID: row.ID })
  type.value = 'update'
  if (res.code === 0) {
    formData.value = res.data.rebuild
    dialogFormVisible.value = true
  }
}


// 删除行
const deleteBuildFunc = async (row) => {
  const res = await deleteBuild({ ID: row.ID })
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
  dialogFormVisible.value = true
  // console.log("projectId", formData.value.projectId)
}

// 关闭弹窗
const closeDialog = () => {
  dialogFormVisible.value = false
  formData.value = {
    projectId: undefined,
    projectName: undefined,
    appId: undefined,
    appName: undefined,
    approveStatus: 1,
  }
}
// 弹窗确定
const enterDialog = async () => {
  elFormRef.value?.validate(async (valid) => {
    if (!valid) return
    let res
    switch (type.value) {
      case 'create':
        res = await createBuild(formData.value)
        break
      case 'update':
        res = await updateBuild(formData.value)
        break
      default:
        res = await createBuild(formData.value)
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

const onProjectChanged = (value) => {
  console.log("projectChanged", value)

  if (value === '' || value === undefined) {
    project = ref({})
    formData.value.projectName = undefined
    return
  }

  projects.value.map(proj => {
    if (proj.ID === value) {
      project.value = proj
      formData.value.projectName = project.value.name
    }
  })

}

const onAppChanged = (value) => {
  console.log("appChanged", value)

  if (value === '') {
    formData.value.appName = undefined
    return
  }
  // reset params
  formData.value.buildParamValues = {}
  project.value.apps.map(app_ => {
    if (app_.ID === value) {
      app.value = app_
      formData.value.appName = app.value.name
      // formData.value.buildParamValues
      if (app.value.buildParams !== null) {
        app.value.buildParams.map(param => {
          formData.value.buildParamValues[param] = undefined
        })
      }
    }
  })

  // console.log("current project", project.value)
  // console.log('current app', app.value)
  // console.log('formData', formData.value)
}

const handleApproveBuild = (row) => {
  ElMessageBox.confirm('确定允许本次发布?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    var approveBuildData = {
      id: row.ID,
      approveStatus: 2
    }
    const resp = await approveBuild(approveBuildData)
    if (resp.code === 0) {
      ElMessage({
        type: 'success',
        message: '审批成功'
      })
      getTableData()
    }

  })
}

const handleRejectBuild = (row) => {
  ElMessageBox.confirm('确定拒绝本次发布?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    var refuseBuildData = {
      id: row.ID,
      approveStatus: 3
    }
    const resp = await approveBuild(refuseBuildData)
    if (resp.code === 0) {
      ElMessage({
        type: 'success',
        message: '审批完成'
      })
      getTableData()
    }

  })
}

const handleViewBuildLog = async (row) => {
  const resp = await findBuild({ ID: row.ID })
  if (resp.code === 0) {
    // console.log("findBuild:", resp)
    viewBuildLogEnabled.value = true
    await nextTick()
    viewLogTitle.value = "查看日志 [ " + resp.data.build.appName + "#" + resp.data.build.buildNumber + " ]"
    buildLog.value = resp.data.build.log
    buildDuration.value = (resp.data.build.duration / 1000).toFixed(2)
    viewLogSubTitle.value =  "本次发布共消耗" + buildDuration.value + "s"

    if (resp.data.build.result === 3) {
        timer.value = setInterval(async () => {
        const resp = await findBuild({ ID: row.ID })
        if(!viewBuildLogEnabled.value || (resp.code === 0 && resp.data.build.result !== 3)) {
          buildLog.value = resp.data.build.log
          buildDuration.value = (resp.data.build.duration / 1000).toFixed(2)
          viewLogSubTitle.value = "本次发布共消耗" + buildDuration.value + "s"
          clearInterval(timer.value)
          timer.value = null
        } else if(resp.code === 0) {
          buildLog.value = resp.data.build.log
          buildDuration.value = (resp.data.build.duration / 1000).toFixed(2)
          viewLogSubTitle.value =  "本次发布共消耗" + buildDuration.value + "s"
          nextTick(() => {
            const scrollDom = document.getElementById('codemirrorLog')
            if (scrollDom) {
              const parent = scrollDom.parentElement
              // console.log(scrollDom.scrollHeight)
              parent.scrollTo(0, scrollDom.scrollHeight)
            }
          })
        }
      }, 2000)
    }
  }

}

const closeViewBuildLogDialog = () => {
  viewBuildLogEnabled.value = false
  if(timer.value) {
    clearInterval(timer.value)
    timer.value = null
  }
  buildLog.value = ""
  buildDuration.value = 0
}

const handleReCreateBuild = (row) => {
  formData.value = {
    projectId: row.projectId,
    projectName: row.projectName,
    appId: row.appId,
    appName: row.appName,
    buildParamValues: row.buildParamValues,
    approveStatus: 1,
  }
  projects.value.map(proj => {
    if (proj.ID === row.projectId) {
      project.value = proj
    }
  })
  project.value.apps.map(app_ => {
    if (app_.ID === row.appId) {
      app.value = app_
    }
  })

  type.value = 'create'
  dialogFormVisible.value = true
}

const viewBuildDetail = async (buildId) => {
  let data = {
    id: buildId
  }
  const resp = await getBuildInfo(data)
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

<style>

</style>