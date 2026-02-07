<template>
  <div class="p-8">
    <h1 class="text-2xl font-bold text-gray-800 mb-6">数据看板</h1>
    
    <!-- 统计卡片 -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
      <el-card class="bg-gradient-to-br from-blue-500 to-blue-600 text-white">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-blue-100 text-sm">今日会话</p>
            <p class="text-3xl font-bold mt-2">128</p>
          </div>
          <el-icon :size="40" class="opacity-80"><ChatDotRound /></el-icon>
        </div>
      </el-card>
      
      <el-card class="bg-gradient-to-br from-green-500 to-green-600 text-white">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-green-100 text-sm">已解决</p>
            <p class="text-3xl font-bold mt-2">96</p>
          </div>
          <el-icon :size="40" class="opacity-80"><CircleCheck /></el-icon>
        </div>
      </el-card>
      
      <el-card class="bg-gradient-to-br from-orange-500 to-orange-600 text-white">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-orange-100 text-sm">待处理</p>
            <p class="text-3xl font-bold mt-2">32</p>
          </div>
          <el-icon :size="40" class="opacity-80"><Clock /></el-icon>
        </div>
      </el-card>
      
      <el-card class="bg-gradient-to-br from-purple-500 to-purple-600 text-white">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-purple-100 text-sm">在线客服</p>
            <p class="text-3xl font-bold mt-2">8</p>
          </div>
          <el-icon :size="40" class="opacity-80"><User /></el-icon>
        </div>
      </el-card>
    </div>

    <!-- 图表区域 -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
      <el-card>
        <template #header>
          <span class="font-bold">会话趋势</span>
        </template>
        <div class="h-64 flex items-center justify-center text-gray-400">
          <p>图表区域</p>
        </div>
      </el-card>
      
      <el-card>
        <template #header>
          <span class="font-bold">客服绩效</span>
        </template>
        <div class="h-64 flex items-center justify-center text-gray-400">
          <p>图表区域</p>
        </div>
      </el-card>
    </div>

    <!-- 最近会话 -->
    <el-card>
      <template #header>
        <span class="font-bold">最近会话</span>
      </template>
      <el-table :data="recentSessions" stripe>
        <el-table-column prop="id" label="会话ID" width="100" />
        <el-table-column prop="visitor" label="访客" />
        <el-table-column prop="agent" label="客服" />
        <el-table-column prop="status" label="状态">
          <template #default="{ row }">
            <el-tag :type="row.status === '已解决' ? 'success' : 'warning'" size="small">
              {{ row.status }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="time" label="时间" />
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { ChatDotRound, CircleCheck, Clock, User } from '@element-plus/icons-vue'

const recentSessions = ref([
  { id: '1001', visitor: '访客A', agent: '客服1', status: '进行中', time: '10:30' },
  { id: '1002', visitor: '访客B', agent: '客服2', status: '已解决', time: '10:25' },
  { id: '1003', visitor: '访客C', agent: '客服1', status: '已解决', time: '10:20' },
  { id: '1004', visitor: '访客D', agent: '客服3', status: '进行中', time: '10:15' },
  { id: '1005', visitor: '访客E', agent: '客服2', status: '已解决', time: '10:10' }
])
</script>
