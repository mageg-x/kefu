<template>
  <div class="p-8">
    <h1 class="text-2xl font-bold text-gray-800 mb-6">系统设置</h1>
    
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- 基本设置 -->
      <el-card>
        <template #header>
          <span class="font-bold">基本设置</span>
        </template>
        <el-form label-width="120px">
          <el-form-item label="系统名称">
            <el-input v-model="settings.systemName" />
          </el-form-item>
          <el-form-item label="系统LOGO">
            <el-upload
              class="avatar-uploader"
              :show-file-list="false"
              :on-success="handleLogoSuccess"
              :before-upload="beforeUpload"
            >
              <img v-if="settings.logo" :src="settings.logo" class="logo" />
              <el-icon v-else class="avatar-uploader-icon"><Plus /></el-icon>
            </el-upload>
          </el-form-item>
          <el-form-item label="欢迎语">
            <el-input v-model="settings.welcomeMsg" type="textarea" :rows="3" />
          </el-form-item>
          <el-form-item>
            <el-button type="primary">保存</el-button>
          </el-form-item>
        </el-form>
      </el-card>

      <!-- 会话设置 -->
      <el-card>
        <template #header>
          <span class="font-bold">会话设置</span>
        </template>
        <el-form label-width="120px">
          <el-form-item label="最大会话数">
            <el-input-number v-model="settings.maxSessions" :min="1" :max="20" />
          </el-form-item>
          <el-form-item label="超时时间">
            <el-input-number v-model="settings.timeout" :min="30" :max="300" />
            <span class="ml-2 text-gray-500">秒</span>
          </el-form-item>
          <el-form-item label="自动分配">
            <el-switch v-model="settings.autoAssign" />
          </el-form-item>
          <el-form-item>
            <el-button type="primary">保存</el-button>
          </el-form-item>
        </el-form>
      </el-card>

      <!-- 通知设置 -->
      <el-card>
        <template #header>
          <span class="font-bold">通知设置</span>
        </template>
        <el-form label-width="120px">
          <el-form-item label="新会话通知">
            <el-switch v-model="settings.newSessionNotify" />
          </el-form-item>
          <el-form-item label="邮件通知">
            <el-switch v-model="settings.emailNotify" />
          </el-form-item>
          <el-form-item label="通知邮箱">
            <el-input v-model="settings.notifyEmail" />
          </el-form-item>
          <el-form-item>
            <el-button type="primary">保存</el-button>
          </el-form-item>
        </el-form>
      </el-card>

      <!-- 安全设置 -->
      <el-card>
        <template #header>
          <span class="font-bold">安全设置</span>
        </template>
        <el-form label-width="120px">
          <el-form-item label="会话加密">
            <el-switch v-model="settings.sessionEncrypt" />
          </el-form-item>
          <el-form-item label="IP限制">
            <el-switch v-model="settings.ipLimit" />
          </el-form-item>
          <el-form-item label="登录验证码">
            <el-switch v-model="settings.captcha" />
          </el-form-item>
          <el-form-item>
            <el-button type="primary">保存</el-button>
          </el-form-item>
        </el-form>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { Plus } from '@element-plus/icons-vue'

const settings = ref({
  systemName: '零点客服系统',
  logo: '',
  welcomeMsg: '您好，有什么可以帮您的吗？',
  maxSessions: 5,
  timeout: 180,
  autoAssign: true,
  newSessionNotify: true,
  emailNotify: false,
  notifyEmail: '',
  sessionEncrypt: true,
  ipLimit: false,
  captcha: false
})

const handleLogoSuccess = (response) => {
  settings.value.logo = URL.createObjectURL(response.raw)
}

const beforeUpload = (file) => {
  const isJPG = file.type === 'image/jpeg'
  const isPNG = file.type === 'image/png'
  const isLt2M = file.size / 1024 / 1024 < 2

  if (!isJPG && !isPNG) {
    ElMessage.error('只能上传JPG/PNG格式的图片！')
    return false
  }
  if (!isLt2M) {
    ElMessage.error('图片大小不能超过2MB！')
    return false
  }
  return true
}
</script>

<style scoped>
.logo {
  width: 100px;
  height: 100px;
  display: block;
}

.avatar-uploader .el-upload {
  border: 1px dashed #d9d9d9;
  border-radius: 6px;
  cursor: pointer;
  position: relative;
  overflow: hidden;
  transition: border-color 0.3s;
}

.avatar-uploader .el-upload:hover {
  border-color: #409eff;
}

.avatar-uploader-icon {
  font-size: 2rem;
  color: #8c939d;
  width: 100px;
  height: 100px;
  line-height: 100px;
  text-align: center;
}
</style>
