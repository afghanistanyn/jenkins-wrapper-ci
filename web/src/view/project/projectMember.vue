<template>
    <div>
        <div class="manager" v-if="isSuperAdmin() || isOps()">
            <el-tag class="project_member_trans" size="large" type='danger'>设置项目管理员</el-tag>
            <div class="project_member_trans">
                <el-transfer  style="text-align: left; display: inline-block" v-model="projectManagers" :data="allUsers4mangers" filterable :titles="['用户列表', '项目管理员列表']" >
                    <template #right-footer>
                    <el-button class="transfer-footer" size="small" :disabled="disableSubmitButton && !(isSuperAdmin() || isOps()) " @click="updateProjectManagers" style="float: right;margin-top: 8px;margin-right: 8px;">提交</el-button>
                </template>
                </el-transfer>
            </div>
        </div>

        <el-divider></el-divider>
        <div class="members">
            <el-tag class="project_member_trans" size="medium" type='danger'>设置项目普通成员</el-tag>
            <div class="project_member_trans">
                <el-transfer  style="text-align: left; display: inline-block" v-model="projectMembers" :data="allUsers4members" filterable :titles="['用户列表', '成员列表']" >
                    <template #right-footer>
                    <el-button class="transfer-footer" size="small" :disabled="disableSubmitButton && !(isSuperAdmin() || isOps() || isProjectManager)" @click="updateProjectMembers" style="float: right;margin-top: 8px;margin-right: 8px;">提交</el-button>
                </template>
                </el-transfer>
            </div>
        </div>
        <el-divider></el-divider>

    </div>
</template>
  
<script>
export default {
    name: 'projectMember'
}
</script>
  
<script setup>
import {
    findProject,
    setProjectMembers
} from '@/api/project'


import {
    getUserList
} from '@/api/user'

// 全量引入格式化工具 请按需保留
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive, onActivated, onMounted} from 'vue'
import { useUserStore } from '@/pinia/modules/user'
import {  useRoute } from 'vue-router'

//  router
const route = useRoute()

// userInfo store
const userStore = useUserStore()
const currentUserId =  userStore.userInfo.ID
const currentUserAuthorityId = userStore.userInfo.authority.authorityId

// console.log("currentUserId", currentUserId)
// console.log("currentUserAuthorityId", currentUserAuthorityId)

const disableSubmitButton = ref(true)
const project = ref()
const isProjectManager = ref(false)
const allUsers4mangers = ref([])
const allUsers4members = ref([])
const projectManagers = ref([])
const projectMembers = ref([])


onMounted(() => {
    initProjectMembers(route.params.projectId)
})

onActivated(() => {
  initProjectMembers(route.params.projectId)
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


const initProjectMembers = async (projectId) => {
    if (project.value) {
        project.value.managers.map(user => {
            projectManagers.value.push(user.ID)
            if (currentUserId === user.ID) {
               isProjectManager.value = true
            }
        })

        project.value.members.map(user => {
            projectMembers.value.push(user.ID)
        })

        // get without pageinfo for all users 
        const users_resp = await getUserList({})
        if (users_resp.code ===  0) {
            for (let user of users_resp.data.list) {
                if (user.enable === 1 ){
                    let userOption = {
                    key: user.ID,
                    label: user.nickName + " (" + user.userName + ")",
                    disabled: false
                }
                allUsers4mangers.value.push(userOption)
                allUsers4members.value.push(userOption)
                }

            }
        }

        // 数据加载完成
        disableSubmitButton.value = false
    } else {
        const project_resp = await findProject({ "ID": projectId })
        if (project_resp.code === 0) {
            project.value = project_resp.data.project
            project.value.managers.map(user => {
                projectManagers.value.push(user.ID)
                if (currentUserId === user.ID) {
                    isProjectManager.value = true
                }
            })

            project.value.members.map(user => {
                projectMembers.value.push(user.ID)
            })

            // get without pageinfo for all users 
            const users_resp = await getUserList({})
            if (users_resp.code ===  0) {
                for (let user of users_resp.data.list) {
                    if (user.enable === 1 ){
                        let userOption = {
                        key: user.ID,
                        label: user.nickName + " (" + user.userName + ")",
                        disabled: false
                    }
                    allUsers4mangers.value.push(userOption)
                    allUsers4members.value.push(userOption)
                    }
                }
            }

            // 数据加载完成
            disableSubmitButton.value = false
        }
    }
}


// 更新项目管理员
const updateProjectManagers = async (row) => {
    let data = {
        projectId: project.value.ID,
        memberType: "managers",
        memberIds: projectManagers.value
    }
    const res = await setProjectMembers(data)
    if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: "设置项目管理员成功"
        })
    }
}

// 更新项目管理员
const updateProjectMembers = async (row) => {
    let data = {
        projectId: project.value.ID,
        memberType: "members",
        memberIds: projectMembers.value
    }
    const res = await setProjectMembers(data)
    if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: "设置项目成员成功"
        })
    }
}

</script>
  
<style scoped>
.project_member_trans {
    text-align: left;
    margin-top: 20px;
    margin-left: 15%;
    --el-tag-font-size: 15px

}

.project_member_trans >>> .el-transfer-panel {
    width: 300px;
}

.transfer-footer {
    margin-left: 20px;
    padding: 6px 5px;
  }

</style>
  