<template>
  <div class="system">
    <el-form ref="form" :model="config" label-width="240px">
      <!--  System start  -->
      <el-collapse v-model="activeNames">
        <el-collapse-item title="系统配置" name="1">
          <el-form-item label="环境值">
            <!-- <el-input v-model="config.system.env" />-->
            <el-select v-model="config.system.env" style="width:100%">
              <el-option value="public" />
              <el-option value="develop" />
            </el-select>
          </el-form-item>
          <el-form-item label="端口值">
            <el-input v-model.number="config.system.addr" />
          </el-form-item>
          <el-form-item label="数据库类型">
            <el-select v-model="config.system['db-type']" style="width:100%">
              <el-option value="mysql" />
              <el-option value="pgsql" />
            </el-select>
          </el-form-item>
          <el-form-item label="多点登录拦截">
            <el-checkbox v-model="config.system['use-multipoint']">开启</el-checkbox>
          </el-form-item>
          <el-form-item label="开启redis">
            <el-checkbox v-model="config.system['use-redis']">开启</el-checkbox>
          </el-form-item>
          <el-form-item label="限流次数">
            <el-input-number v-model.number="config.system['iplimit-count']" />
          </el-form-item>
          <el-form-item label="限流时间">
            <el-input-number v-model.number="config.system['iplimit-time']" />
          </el-form-item>
          <el-tooltip
            content="请修改完成后，注意一并修改前端env环境下的VITE_BASE_PATH"
            placement="top-start"
          >
            <el-form-item label="全局路由前缀">
              <el-input v-model="config.system['router-prefix']" />
            </el-form-item>
          </el-tooltip>
        </el-collapse-item>
        <el-collapse-item title="jwt签名" name="2">
          <el-form-item label="jwt签名">
            <el-input v-model="config.jwt['signing-key']" />
          </el-form-item>
          <el-form-item label="有效期">
            <el-input v-model="config.jwt['expires-time']" />
          </el-form-item>
          <el-form-item label="缓冲期">
            <el-input v-model="config.jwt['buffer-time']" />
          </el-form-item>
          <el-form-item label="签发者">
            <el-input v-model="config.jwt.issuer" />
          </el-form-item>
        </el-collapse-item>
        <el-collapse-item title="Zap日志配置" name="3">
          <el-form-item label="级别">
            <el-input v-model.number="config.zap.level" />
          </el-form-item>
          <el-form-item label="输出">
            <el-input v-model="config.zap.format" />
          </el-form-item>
          <el-form-item label="日志前缀">
            <el-input v-model="config.zap.prefix" />
          </el-form-item>
          <el-form-item label="日志文件夹">
            <el-input v-model="config.zap.director" />
          </el-form-item>
          <el-form-item label="编码级">
            <el-input v-model="config.zap['encode-level']" />
          </el-form-item>
          <el-form-item label="栈名">
            <el-input v-model="config.zap['stacktrace-key']" />
          </el-form-item>
          <el-form-item label="日志留存时间(默认以天为单位)">
            <el-input v-model.number="config.zap['max-age']" />
          </el-form-item>
          <el-form-item label="显示行">
            <el-checkbox v-model="config.zap['show-line']" />
          </el-form-item>
          <el-form-item label="输出控制台">
            <el-checkbox v-model="config.zap['log-in-console']" />
          </el-form-item>
        </el-collapse-item>
        <el-collapse-item title="Redis admin数据库配置" name="4">
          <el-form-item label="库">
            <el-input v-model.number="config.redis.db" />
          </el-form-item>
          <el-form-item label="地址">
            <el-input v-model="config.redis.addr" />
          </el-form-item>
          <el-form-item label="密码">
            <el-input v-model="config.redis.password" />
          </el-form-item>
        </el-collapse-item>

        <el-collapse-item title="邮箱配置" name="5">
          <el-form-item label="接收者邮箱">
            <el-input v-model="config.email.to" placeholder="可多个，以逗号分隔" />
          </el-form-item>
          <el-form-item label="端口">
            <el-input v-model.number="config.email.port" />
          </el-form-item>
          <el-form-item label="发送者邮箱">
            <el-input v-model="config.email.from" />
          </el-form-item>
          <el-form-item label="host">
            <el-input v-model="config.email.host" />
          </el-form-item>
          <el-form-item label="是否为ssl">
            <el-checkbox v-model="config.email['is-ssl']" />
          </el-form-item>
          <el-form-item label="secret">
            <el-input v-model="config.email.secret" />
          </el-form-item>
          <el-form-item label="测试邮件">
            <el-button @click="email">测试邮件</el-button>
          </el-form-item>
        </el-collapse-item>
        <el-collapse-item title="验证码配置" name="6">
          <el-form-item label="字符长度">
            <el-input v-model.number="config.captcha['key-long']" />
          </el-form-item>
          <el-form-item label="平台宽度">
            <el-input v-model.number="config.captcha['img-width']" />
          </el-form-item>
          <el-form-item label="图片高度">
            <el-input v-model.number="config.captcha['img-height']" />
          </el-form-item>
        </el-collapse-item>
        <el-collapse-item title="数据库配置" name="7">
          <template v-if="config.system['db-type'] === 'mysql'">
            <el-form-item label="用户名">
              <el-input v-model="config.mysql.username" />
            </el-form-item>
            <el-form-item label="密码">
              <el-input v-model="config.mysql.password" />
            </el-form-item>
            <el-form-item label="地址">
              <el-input v-model="config.mysql.path" />
            </el-form-item>
            <el-form-item label="数据库">
              <el-input v-model="config.mysql['db-name']" />
            </el-form-item>
            <el-form-item label="前缀">
              <el-input v-model="config.mysql['refix']" />
            </el-form-item>
            <el-form-item label="复数表">
              <el-switch v-model="config.mysql['singular']" />
            </el-form-item>
            <el-form-item label="引擎">
              <el-input v-model="config.mysql['engine']" />
            </el-form-item>
            <el-form-item label="maxIdleConns">
              <el-input v-model.number="config.mysql['max-idle-conns']" />
            </el-form-item>
            <el-form-item label="maxOpenConns">
              <el-input v-model.number="config.mysql['max-open-conns']" />
            </el-form-item>
            <el-form-item label="写入日志">
              <el-checkbox v-model="config.mysql['log-zap']" />
            </el-form-item>
            <el-form-item label="日志模式">
              <el-input v-model="config.mysql['log-mode']" />
            </el-form-item>
          </template>
          <template v-if="config.system['db-type'] === 'pgsql'">
            <el-form-item label="用户名">
              <el-input v-model="config.pgsql.username" />
            </el-form-item>
            <el-form-item label="密码">
              <el-input v-model="config.pgsql.password" />
            </el-form-item>
            <el-form-item label="地址">
              <el-input v-model="config.pgsql.path" />
            </el-form-item>
            <el-form-item label="数据库">
              <el-input v-model="config.pgsql.dbname" />
            </el-form-item>
            <el-form-item label="前缀">
              <el-input v-model="config.pgsql['refix']" />
            </el-form-item>
            <el-form-item label="复数表">
              <el-switch v-model="config.pgsql['singular']" />
            </el-form-item>
            <el-form-item label="引擎">
              <el-input v-model="config.pgsql['engine']" />
            </el-form-item>
            <el-form-item label="maxIdleConns">
              <el-input v-model.number="config.pgsql['max-idle-conns']" />
            </el-form-item>
            <el-form-item label="maxOpenConns">
              <el-input v-model.number="config.pgsql['max-open-conns']" />
            </el-form-item>
            <el-form-item label="写入日志">
              <el-checkbox v-model="config.pgsql['log-zap']" />
            </el-form-item>
            <el-form-item label="日志模式">
              <el-input v-model="config.pgsql['log-mode']" />
            </el-form-item>
          </template>
        </el-collapse-item>

      
        <el-collapse-item title="Jenkins配置" name="8">
          <el-form-item label="url">
            <el-input v-model="config['jenkins-config'].url" />
          </el-form-item>
          <el-form-item label="username">
            <el-input v-model="config['jenkins-config'].username" />
          </el-form-item>
          <el-form-item label="password">
            <el-input v-model="config['jenkins-config'].password" />
          </el-form-item>
          <el-form-item label="timeout">
            <el-input v-model.number="config['jenkins-config'].timeout" />
          </el-form-item>
        </el-collapse-item>


        <el-collapse-item title="git凭证Id (jenkins)" name="9">
          <template v-for="(item,k) in config['git-credentials']">
            <div v-for="(_,k2) in item" :key="k2">
              <el-form-item :key="k+k2" :label="k2">
                <el-input v-model="item[k2]" />
              </el-form-item>
            </div>
            <el-divider></el-divider>
          </template>

        </el-collapse-item>


      </el-collapse>
    </el-form>
    <div class="gva-btn-list">
      <el-button type="primary" @click="update">立即更新</el-button>
      <el-button type="primary" @click="reload">重启服务（开发中）</el-button>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Config'
}
</script>
<script setup>
import { getSystemConfig, setSystemConfig } from '@/api/system'
import { emailTest } from '@/api/email'
import { ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'

const activeNames = reactive([])
const config = ref({
  system: {
    'iplimit-count': 0,
    'iplimit-time': 0
  },
  jwt: {},
  mysql: {},
  pgsql: {},
  excel: {},
  redis: {},
  captcha: {},
  zap: {},
  local: {},
  email: {},
  'jenkins-config': {},
  'super-admin-roles': {},
  'git-credentials': {},
})

const initForm = async() => {
  const res = await getSystemConfig()
  if (res.code === 0) {
    config.value = res.data.config
  }
}
initForm()
const reload = () => {}
const update = async() => {
  const res = await setSystemConfig({ config: config.value })
  if (res.code === 0) {
    ElMessage({
      type: 'success',
      message: '配置文件设置成功'
    })
    await initForm()
  }
}
const email = async() => {
  const res = await emailTest()
  if (res.code === 0) {
    ElMessage({
      type: 'success',
      message: '邮件发送成功'
    })
    await initForm()
  } else {
    ElMessage({
      type: 'error',
      message: '邮件发送失败'
    })
  }
}

</script>

<style lang="scss">
.system {
  background: #fff;
  padding:36px;
  border-radius: 2px;
  h2 {
    padding: 10px;
    margin: 10px 0;
    font-size: 16px;
    box-shadow: -4px 0px 0px 0px #e7e8e8;
  }
  ::v-deep(.el-input-number__increase){
    top:5px !important;
  }
  .gva-btn-list{
    margin-top:16px;
  }
}
</style>
