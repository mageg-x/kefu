<template>
  <div class="p-8">
    <h1 class="text-2xl font-bold text-gray-800 mb-6">个人资料</h1>
    
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- 个人信息卡片 -->
      <el-card class="lg:col-span-1">
        <div class="text-center">
          <el-avatar :size="100" :src="userInfo.avatar">
            {{ userInfo.name?.charAt(0) || 'U' }}
          </el-avatar>
          <h2 class="text-xl font-bold text-gray-800 mt-4">{{ userInfo.name }}</h2>
          <p class="text-gray-500">{{ userInfo.role === 'admin' ? '管理员' : '客服专员' }}</p>
          <div class="flex justify-center gap-4 mt-4">
            <div class="text-center">
              <p class="text-2xl font-bold text-blue-600">{{ userInfo.sessions }}</p>
              <p class="text-sm text-gray-500">今日会话</p>
            </div>
            <div class="text-center">
              <p class="text-2xl font-bold text-green-600">{{ userInfo.rating }}</p>
              <p class="text-sm text-gray-500">评分</p>
            </div>
          </div>
        </div>
      </el-card>

      <!-- 编辑资料表单 -->
      <el-card class="lg:col-span-2">
        <template #header>
          <span class="font-bold">编辑资料</span>
        </template>
        <el-form :model="form" label-width="100px">
          <el-form-item label="用户名">
            <el-input v-model="form.username" disabled />
          </el-form-item>
          <el-form-item label="姓名">
            <el-input v-model="form.name" />
          </el-form-item>
          <el-form-item label="邮箱">
            <el-input v-model="form.email" />
          </el-form-item>
          <el-form-item label="手机号">
            <el-input v-model="form.phone" />
          </el-form-item>
          <el-form-item label="头像">
            <el-upload
              class="avatar-uploader"
              :show-file-list="false"
              :on-success="handleAvatarSuccess"
              :before-upload="beforeAvatarUpload"
            >
              <img v-if="form.avatar" :src="form.avatar" class="avatar" />
              <el-icon v-else class="avatar-uploader-icon"><Plus /></el-icon>
            </el-upload>
          </el-form-item>
          <el-form-item label="个性签名">
            <el-input v-model="form.bio" type="textarea" :rows="3" />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="handleSave">保存修改</el-button>
          </el-form-item>
        </el-form>
      </el-card>
    </div>

    <!-- 修改密码 -->
    <el-card class="mt-6">
      <template #header>
        <span class="font-bold">修改密码</span>
      </template>
      <el-form :model="passwordForm" label-width="100px">
        <el-form-item label="当前密码">
          <el-input v-model="passwordForm.currentPassword" type="password" />
        </el-form-item>
        <el-form-item label="新密码">
          <el-input v-model="passwordForm.newPassword" type="password" />
        </el-form-item>
        <el-form-item label="确认密码">
          <el-input v-model="passwordForm.confirmPassword" type="password" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handlePasswordChange">修改密码</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { useStore } from '@/script/store'

const store = useStore()

const userInfo = ref({
  name: '客服专员',
  username: 'agent1',
  email: 'agent@example.com',
  phone: '13800138000',
  avatar: '',
  bio: '用心服务每一位客户',
  role: 'agent',
  sessions: 12,
  rating: 4.8
})

const form = ref({
  username: '',
  name: '',
  email: '',
  phone: '',
  avatar: '',
  bio: ''
})

const passwordForm = ref({
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
})

onMounted(() => {
  form.value = { ...userInfo.value }
})

const handleAvatarSuccess = (response) => {
  form.value.avatar = URL.createObjectURL(response.raw)
  ElMessage.success('头像上传成功')
}

const beforeAvatarUpload = (file) => {
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

const handleSave = () => {
  userInfo.value = { ...form.value }
  ElMessage.success('保存成功')
}

const handlePasswordChange = () => {
  if (!passwordForm.value.currentPassword) {
    ElMessage.error('请输入当前密码')
    return
  }
  if (!passwordForm.value.newPassword) {
    ElMessage.error('请输入新密码')
    return
  }
  if (passwordForm.value.newPassword !== passwordForm.value.confirmPassword) {
    ElMessage.error('两次输入的密码不一致')
    return
  }
  ElMessage.success('密码修改成功')
  passwordForm.value = {
    currentPassword: '',
    newPassword: '',
    confirmPassword: ''
  }
}
</script>

<style scoped>
.avatar {
  width: 100px;
  height: 100px;
  display: block;
  border-radius: 50%;
}

.avatar-uploader .el-upload {
  border: 1px dashed #d9d9d9;
  border-radius: 50%;
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
