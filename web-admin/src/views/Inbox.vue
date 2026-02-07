<template>
  <div class="flex h-screen">
    <!-- 左侧会话列表 -->
    <div class="w-80 border-r border-gray-200 bg-white flex flex-col">
      <div class="p-4 border-b border-gray-200">
        <h2 class="font-bold text-gray-800">会话列表</h2>
        <p class="text-sm text-gray-500">{{ activeSessions.length }} 个进行中</p>
      </div>
      
      <div class="flex-1 overflow-y-auto">
        <div
          v-for="session in sessions"
          :key="session.id"
          @click="selectSession(session)"
          class="p-4 border-b border-gray-100 cursor-pointer transition-colors"
          :class="selectedSession?.id === session.id ? 'bg-blue-50 border-l-4 border-l-blue-600' : 'hover:bg-gray-50'"
        >
          <div class="flex items-start gap-3">
            <el-avatar :size="40" :src="session.avatar">
              {{ session.name?.charAt(0) }}
            </el-avatar>
            <div class="flex-1 min-w-0">
              <div class="flex items-center justify-between">
                <span class="font-medium text-gray-800">{{ session.name }}</span>
                <span class="text-xs text-gray-500">{{ session.time }}</span>
              </div>
              <p class="text-sm text-gray-600 truncate">{{ session.lastMessage }}</p>
              <div class="flex items-center gap-2 mt-1">
                <el-tag :type="session.status === '进行中' ? 'success' : 'info'" size="small">
                  {{ session.status }}
                </el-tag>
                <span v-if="session.unread > 0" class="text-xs text-red-500 font-medium">
                  {{ session.unread }} 条未读
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 右侧聊天区域 -->
    <div class="flex-1 flex flex-col bg-gray-50">
      <div v-if="selectedSession" class="flex-1 flex flex-col">
        <!-- 聊天头部 -->
        <div class="p-4 bg-white border-b border-gray-200 flex items-center justify-between">
          <div class="flex items-center gap-3">
            <el-avatar :size="40" :src="selectedSession.avatar">
              {{ selectedSession.name?.charAt(0) }}
            </el-avatar>
            <div>
              <p class="font-medium text-gray-800">{{ selectedSession.name }}</p>
              <p class="text-xs text-gray-500">{{ selectedSession.email }}</p>
            </div>
          </div>
          <div class="flex gap-2">
            <el-button type="primary" size="small">转接</el-button>
            <el-button type="danger" size="small">结束会话</el-button>
          </div>
        </div>

        <!-- 消息列表 -->
        <div class="flex-1 overflow-y-auto p-4 space-y-4">
          <div
            v-for="message in messages"
            :key="message.id"
            class="flex"
            :class="message.isSelf ? 'justify-end' : 'justify-start'"
          >
            <div
              class="max-w-xs px-4 py-2 rounded-lg"
              :class="message.isSelf ? 'bg-blue-600 text-white' : 'bg-white text-gray-800'"
            >
              <p>{{ message.content }}</p>
              <p class="text-xs mt-1 opacity-70">{{ message.time }}</p>
            </div>
          </div>
        </div>

        <!-- 输入区域 -->
        <div class="p-4 bg-white border-t border-gray-200">
          <div class="flex gap-3">
            <el-input
              v-model="inputMessage"
              type="textarea"
              :rows="3"
              placeholder="输入消息..."
              @keyup.ctrl.enter="sendMessage"
            />
            <el-button type="primary" @click="sendMessage">发送</el-button>
          </div>
          <div class="flex items-center gap-4 mt-3">
            <el-button text size="small">
              <el-icon><Picture /></el-icon>
              图片
            </el-button>
            <el-button text size="small">
              <el-icon><Document /></el-icon>
              文件
            </el-button>
            <el-button text size="small">
              <el-icon><Folder /></el-icon>
              快捷回复
            </el-button>
          </div>
        </div>
      </div>

      <!-- 空状态 -->
      <div v-else class="flex-1 flex items-center justify-center text-gray-400">
        <div class="text-center">
          <el-icon :size="64" class="mb-4"><ChatDotRound /></el-icon>
          <p>选择一个会话开始聊天</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { ChatDotRound, Picture, Document, Folder } from '@element-plus/icons-vue'

const selectedSession = ref(null)
const inputMessage = ref('')

const sessions = ref([
  { id: 1, name: '访客A', email: 'visitor1@example.com', avatar: '', status: '进行中', time: '10:30', lastMessage: '你好，请问有什么可以帮您的吗？', unread: 0 },
  { id: 2, name: '访客B', email: 'visitor2@example.com', avatar: '', status: '进行中', time: '10:25', lastMessage: '我想咨询一下产品价格', unread: 2 },
  { id: 3, name: '访客C', email: 'visitor3@example.com', avatar: '', status: '进行中', time: '10:20', lastMessage: '好的，谢谢', unread: 0 },
  { id: 4, name: '访客D', email: 'visitor4@example.com', avatar: '', status: '已结束', time: '10:15', lastMessage: '再见', unread: 0 },
  { id: 5, name: '访客E', email: 'visitor5@example.com', avatar: '', status: '进行中', time: '10:10', lastMessage: '请问有优惠活动吗？', unread: 1 }
])

const activeSessions = computed(() => sessions.value.filter(s => s.status === '进行中'))

const messages = ref([])

const selectSession = (session) => {
  selectedSession.value = session
  messages.value = [
    { id: 1, content: '你好，请问有什么可以帮您的吗？', time: '10:30', isSelf: false },
    { id: 2, content: '我想咨询一下产品价格', time: '10:31', isSelf: true },
    { id: 3, content: '好的，请问您想了解哪个产品？', time: '10:32', isSelf: false }
  ]
}

const sendMessage = () => {
  if (!inputMessage.value.trim()) return
  
  messages.value.push({
    id: messages.value.length + 1,
    content: inputMessage.value,
    time: new Date().toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' }),
    isSelf: true
  })
  
  inputMessage.value = ''
}
</script>
