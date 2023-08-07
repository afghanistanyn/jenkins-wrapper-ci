<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline" :rules="searchRule" @keyup.enter="onSubmit">
        <el-form-item label="项目"> 
          <el-input v-model="searchInfo.projectInfo" placeholder="" clearable size="small" style="width: 240px" ></el-input>
        </el-form-item>
        <el-form-item label="应用"> 
          <el-input v-model="searchInfo.appInfo" placeholder="" clearable size="small" style="width: 240px" ></el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
        <div class="gva-btn-list">
            <el-button v-if="isSuperAdmin() || isOps()" type="primary" icon="plus" @click="openDialog">新增项目</el-button>
        </div>
        <el-table
        ref="multipleTable"
        style="width: 100%"
        tooltip-effect="dark"
        :data="tableData"
        row-key="ID"
        border
        @selection-change="handleSelectionChange"
        >
        <!-- <el-table-column type="selection" width="55" /> -->

        <el-table-column align="left" label="英文名" prop="name" width="150" />
        <el-table-column align="left" label="中文名" prop="nameCn" width="150" />
        <el-table-column align="left" label="描述" prop="description" width="200" />

        <el-table-column align="left" label="应用列表" prop="apps">
          <template #default="scope">
            <div v-for="app in scope.row.apps">
                <el-tag class="big-tag" effect="plain" :hit=false>{{ app.name }}</el-tag>
            </div>
          </template>
        </el-table-column>
        <el-table-column align="left" label="管理人员" prop="managers" >
          <template v-slot="scope">
            <div class="avatar-group" v-for="user in scope.row.managers">
              <div v-if="user.enable === 1">
                <el-avatar :size="40" :src="user.headerImg" fit="cover" :alt="user.nickName" >  </el-avatar>  
                <div class="avatar-with-name">{{ user.nickName }} </div>
               </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column align="left" label="普通成员" prop="members" >
          <template v-slot="scope">
            <div class="avatar-group" v-for="user in scope.row.members">
              <div v-if="user.enable === 1">
              <el-avatar :size="40" :src="user.headerImg" fit="cover" :alt="user.nickName" >  </el-avatar>  
              <div class="avatar-with-name">{{ user.nickName }} </div>
            </div>
          </div>
          </template>
        </el-table-column>

        <el-table-column align="left" label="创建日期" width="110">
            <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>

        <el-table-column align="left" label="操作" class-name="small-padding fixed-width" width="130">
            <template #default="scope">
            <div class="table-button-group">
            <el-button :disabled="!(hasProjectManagerPermission(scope.row))" size="" type="text"  icon="user" class="table-button" @click="projectMemberManageFunc(scope.row.ID)">成员管理</el-button>
            <el-button :disabled="!(hasProjectManagerPermission(scope.row))" size="" type="text" icon="film" class="table-button" @click="projectAppManageFunc(scope.row.ID)">应用管理</el-button>
            <el-button :disabled="!(hasProjectManagerPermission(scope.row))" size="" type="text" icon="edit" class="table-button" @click="updateProjectFunc(scope.row)">修改</el-button>
            <el-button :disabled="!(isOps() || isSuperAdmin())" size="" type="text" icon="delete" class="table-button" @click="deleteRow(scope.row)">删除</el-button>
            </div>
            </template>
        </el-table-column>
        </el-table>
        <div class="gva-pagination">
            <el-pagination
            layout="total, sizes, prev, pager, next, jumper"
            :current-page="page"
            :page-size="pageSize"
            :page-sizes="[10, 20, 30, 50]"
            :total="total"
            @current-change="handleCurrentChange"
            @size-change="handleSizeChange"
            />
        </div>
    </div>
    <el-dialog v-model="dialogFormVisible" :before-close="closeDialog" :title="type==='create'?'添加项目':'修改项目'" destroy-on-close>
      <el-form :model="formData" label-position="right" ref="elFormRef" :rules="rule" label-width="150px">
        <el-form-item label="英文名:"  prop="name" >
          <el-input v-model="formData.name" :clearable="true"  :disabled="type==='update'" placeholder="请输入" />
        </el-form-item>
        <el-form-item label="中文名:"  prop="nameCn" >
          <el-input v-model="formData.nameCn" :clearable="true"  placeholder="请输入" />
        </el-form-item>
        <el-form-item label="描述:"  prop="description" >
          <el-input v-model="formData.description" :clearable="true"  placeholder="请输入" />
        </el-form-item>
        <el-form-item label="企业微信webhook:"  prop="weworkWebhook" >
          <el-input v-model="formData.weworkWebhook" :clearable="true"  placeholder="通知机器人" />
        </el-form-item>
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
  name: 'Project'
}
</script>

<script setup>
import {
  createProject,
  deleteProject,
  deleteProjectByIds,
  updateProject,
  findProject,
  getProjectList
} from '@/api/project'

// 全量引入格式化工具 请按需保留
import { getDictFunc, formatDate, formatBoolean, filterDict } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive, onActivated} from 'vue'
import { useUserStore } from '@/pinia/modules/user'
import { useRouter } from 'vue-router'

//  router
const router = useRouter()

// userInfo store
const userStore = useUserStore()
const currentUserId =  userStore.userInfo.ID
const currentUserAuthorityId = userStore.userInfo.authority.authorityId

// console.log("currentUserId", currentUserId)
// console.log("currentUserAuthorityId", currentUserAuthorityId)


// 自动化生成的字典（可能为空）以及字段
const formData = ref({
        name: '',
        nameCn: '',
        description: '',
        weworkWebhook: undefined,
        })

// 验证规则
const rule = reactive({
               name : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               },
               {
                   whitespace: true,
                   message: '不能只输入空格',
                   trigger: ['input', 'blur'],
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
  mine: true,
  projectInfo: undefined,
  appInfo: undefined
})


const isSuperAdmin = () => {
  // todo get superAdminRoles from backend
  const superAdminRoles = [888, 999]
  for(let authority of superAdminRoles) {
    if(currentUserAuthorityId === authority) {
      return true
    }
  }
  return false
}

const isOps = () => {
    const opsRoleId = 999
    if(currentUserAuthorityId === opsRoleId) {
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

const hasProjectManagerPermission = (project) => {
  let hasPermission = isSuperAdmin() || isOps() || isProjectManager(project)
  // console.log("hasPermission", hasPermission)
  return hasPermission
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
  elSearchFormRef.value?.validate(async(valid) => {
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

onActivated(() => {
  getTableData()
})


// 查询
const getTableData = async() => {


  const table = await getProjectList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
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
const setOptions = async () =>{
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
    ElMessageBox.confirm('注意: 项目下所有应用也将被删除! 确定要删除吗?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
    }).then(() => {
            deleteProjectFunc(row)
        })
    }



// 批量删除控制标记
const deleteVisible = ref(false)

// 多选删除
const onDelete = async() => {
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
      const res = await deleteProjectByIds({ ids })
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
const updateProjectFunc = async(row) => {
    const res = await findProject({ ID: row.ID })
    type.value = 'update'
    if (res.code === 0) {
        formData.value = res.data.project
        dialogFormVisible.value = true
    }
}


// 删除行
const deleteProjectFunc = async (row) => {
    const res = await deleteProject({ ID: row.ID })
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
}

// 关闭弹窗
const closeDialog = () => {
    dialogFormVisible.value = false
    formData.value = {
        name: '',
        nameCn: '',
        }
}
// 弹窗确定
const enterDialog = async () => {
     elFormRef.value?.validate( async (valid) => {
             if (!valid) return
              let res
              switch (type.value) {
                case 'create':
                  res = await createProject(formData.value)
                  break
                case 'update':
                  res = await updateProject(formData.value)
                  break
                default:
                  res = await createProject(formData.value)
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

// 应用管理
const projectAppManageFunc = (projectId) => {
  router.push({
      name: 'projectApp',
      params: {'projectId': projectId},
  })  
}

// 成员管理
const projectMemberManageFunc = (projectId) => {
  router.push({
      name: 'projectMember',
      params: {'projectId': projectId},
  })  
}


</script>

<style scoped>

.headerAvatar{
    display: flex;
    justify-content: center;
    align-items: center;
    margin-right: 8px;
}

.avatar-group {
  display: flex;
  justify-content: left;
}


.avatar-with-name {
  display: inline-block;
  vertical-align: middle;
  margin-left: 10px;
  margin-left: 10px;
  margin-bottom: 5px;
}

.table-button-group {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
}

.table-button {
    margin-top: 2px;
    margin-left: 10px;
}

.big-tag {
  --el-tag-font-size: 13px
}

</style>
