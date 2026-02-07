<template>
  <div class="p-8">
    <h1 class="text-2xl font-bold text-gray-800 mb-6">客服团队管理</h1>
    
    <!-- 添加客服按钮 -->
    <div class="mb-6">
      <el-button type="primary" @click="showAddDialog = true">
        <template #icon>
          <el-icon><Plus /></el-icon>
        </template>
        添加客服
      </el-button>
    </div>

    <!-- 客服列表 -->
    <el-card>
      <el-table :data="staffList" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="姓名">
          <template #default="{ row }">
            <div class="flex items-center gap-2">
              <el-avatar :size="32" :src="row.avatar">
                {{ row.name?.charAt(0) }}
              </el-avatar>
              <span>{{ row.name }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="username" label="用户名" />
        <el-table-column prop="role" label="角色">
          <template #default="{ row }">
            <el-tag :type="row.role === 'admin' ? 'danger' : 'primary'" size="small">
              {{ row.role === 'admin' ? '管理员' : '客服' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'online' ? 'success' : 'info'" size="small">
              {{ row.status === 'online' ? '在线' : '离线' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="sessions" label="今日会话" width="100" />
        <el-table-column prop="rating" label="评分" width="100">
          <template #default="{ row }">
            <el-rate v-model="row.rating" disabled show-score text-color="#ff9900" />
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200">
          <template #default>
            <el-button type="primary" link size="small">编辑</el-button>
            <el-button type="warning" link size="small">重置密码</el-button>
            <el-button type="danger" link size="small">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 添加客服对话框 -->
    <el-dialog v-model="showAddDialog" title="添加客服" width="500px">
      <el-form :model="newStaff" label-width="80px">
        <el-form-item label="姓名">
          <el-input v-model="newStaff.name" />
        </el-form-item>
        <el-form-item label="用户名">
          <el-input v-model="newStaff.username" />
        </el-form-item>
        <el-form-item label="密码">
          <el-input v-model="newStaff.password" type="password" />
        </el-form-item>
        <el-form-item label="角色">
          <el-select v-model="newStaff.role">
            <el-option label="客服" value="agent" />
            <el-option label="管理员" value="admin" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showAddDialog = false">取消</el-button>
        <el-button type="primary" @click="handleAdd">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { Plus } from '@element-plus/icons-vue'

const showAddDialog = ref(false)
const newStaff = ref({
  name: '',
  username: '',
  password: '',
  role: 'agent'
})

const staffList = ref([
  { id: 1, name: '张三', username: 'agent1', role: 'agent', status: 'online', sessions: 12, rating: 4.5, avatar: '' },
  { id: 2, name: '李四', username: 'agent2', role: 'agent', status: 'online', sessions: 8, rating: 4.8, avatar: '' },
  { id: 3, name: '王五', username: 'agent3', role: 'agent', status: 'offline', sessions: 5, rating: 4.2, avatar: '' },
  { id: 4, name: '赵六', username: 'admin2', role: 'admin', status: 'online', sessions: 15, rating: 4.9, avatar: '' }
])

const handleAdd = () => {
  console.log('添加客服:', newStaff.value)
  showAddDialog.value = false
}
</script>
