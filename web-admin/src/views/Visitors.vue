<template>
  <div class="p-8">
    <h1 class="text-2xl font-bold text-gray-800 mb-6">访客管理</h1>
    
    <!-- 搜索和筛选 -->
    <el-card class="mb-6">
      <div class="flex gap-4">
        <el-input v-model="searchQuery" placeholder="搜索访客" clearable style="width: 300px">
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
        <el-select v-model="statusFilter" placeholder="状态筛选" clearable style="width: 150px">
          <el-option label="全部" value="" />
          <el-option label="在线" value="online" />
          <el-option label="离线" value="offline" />
        </el-select>
      </div>
    </el-card>

    <!-- 访客列表 -->
    <el-card>
      <el-table :data="visitors" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="访客名称">
          <template #default="{ row }">
            <div class="flex items-center gap-2">
              <el-avatar :size="32" :src="row.avatar">
                {{ row.name?.charAt(0) }}
              </el-avatar>
              <span>{{ row.name }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="email" label="邮箱" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'online' ? 'success' : 'info'" size="small">
              {{ row.status === 'online' ? '在线' : '离线' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="sessions" label="会话数" width="100" />
        <el-table-column prop="lastVisit" label="最后访问" />
        <el-table-column label="操作" width="150">
          <template #default>
            <el-button type="primary" link size="small">查看</el-button>
            <el-button type="danger" link size="small">封禁</el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <!-- 分页 -->
      <div class="flex justify-end mt-4">
        <el-pagination v-model:current-page="currentPage" :page-size="pageSize" :total="total" layout="prev, pager, next" />
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { Search } from '@element-plus/icons-vue'

const searchQuery = ref('')
const statusFilter = ref('')
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(50)

const visitors = ref([
  { id: 1, name: '访客A', email: 'visitor1@example.com', status: 'online', sessions: 5, lastVisit: '10:30', avatar: '' },
  { id: 2, name: '访客B', email: 'visitor2@example.com', status: 'offline', sessions: 3, lastVisit: '09:45', avatar: '' },
  { id: 3, name: '访客C', email: 'visitor3@example.com', status: 'online', sessions: 8, lastVisit: '10:25', avatar: '' },
  { id: 4, name: '访客D', email: 'visitor4@example.com', status: 'offline', sessions: 2, lastVisit: '08:30', avatar: '' },
  { id: 5, name: '访客E', email: 'visitor5@example.com', status: 'online', sessions: 6, lastVisit: '10:15', avatar: '' }
])
</script>
