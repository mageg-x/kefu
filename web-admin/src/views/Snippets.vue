<template>
  <div class="p-8">
    <h1 class="text-2xl font-bold text-gray-800 mb-6">快捷回复</h1>
    
    <!-- 添加快捷回复 -->
    <div class="mb-6">
      <el-button type="primary" @click="showAddDialog = true">
        <template #icon>
          <el-icon><Plus /></el-icon>
        </template>
        添加快捷回复
      </el-button>
    </div>

    <!-- 快捷回复列表 -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      <el-card
        v-for="snippet in snippets"
        :key="snippet.id"
        class="cursor-pointer hover:shadow-lg transition-shadow"
        @click="useSnippet(snippet)"
      >
        <div class="flex items-start justify-between mb-3">
          <el-tag :type="snippet.category === '问候' ? 'success' : snippet.category === '常见问题' ? 'warning' : 'primary'" size="small">
            {{ snippet.category }}
          </el-tag>
          <div class="flex gap-1">
            <el-button type="primary" link size="small" @click.stop="editSnippet(snippet)">编辑</el-button>
            <el-button type="danger" link size="small" @click.stop="deleteSnippet(snippet.id)">删除</el-button>
          </div>
        </div>
        <p class="text-gray-800 font-medium mb-2">{{ snippet.title }}</p>
        <p class="text-gray-600 text-sm line-clamp-3">{{ snippet.content }}</p>
        <p class="text-xs text-gray-400 mt-3">使用 {{ snippet.usage }} 次</p>
      </el-card>
    </div>

    <!-- 添加/编辑对话框 -->
    <el-dialog v-model="showAddDialog" :title="editingSnippet ? '编辑快捷回复' : '添加快捷回复'" width="600px">
      <el-form :model="snippetForm" label-width="80px">
        <el-form-item label="标题">
          <el-input v-model="snippetForm.title" placeholder="请输入标题" />
        </el-form-item>
        <el-form-item label="分类">
          <el-select v-model="snippetForm.category">
            <el-option label="问候" value="问候" />
            <el-option label="常见问题" value="常见问题" />
            <el-option label="结束语" value="结束语" />
            <el-option label="其他" value="其他" />
          </el-select>
        </el-form-item>
        <el-form-item label="内容">
          <el-input v-model="snippetForm.content" type="textarea" :rows="6" placeholder="请输入回复内容" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showAddDialog = false">取消</el-button>
        <el-button type="primary" @click="handleSave">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

const showAddDialog = ref(false)
const editingSnippet = ref(null)
const snippetForm = ref({
  title: '',
  category: '问候',
  content: ''
})

const snippets = ref([
  { id: 1, title: '欢迎语', category: '问候', content: '您好，欢迎来到我们的客服中心，请问有什么可以帮您的吗？', usage: 128 },
  { id: 2, title: '感谢语', category: '结束语', content: '感谢您的咨询，祝您生活愉快！', usage: 96 },
  { id: 3, title: '产品价格', category: '常见问题', content: '关于产品价格，您可以查看我们的官网或咨询销售人员获取详细信息。', usage: 85 },
  { id: 4, title: '售后服务', category: '常见问题', content: '我们提供7天无理由退换货服务，具体请查看售后政策。', usage: 72 },
  { id: 5, title: '发货时间', category: '常见问题', content: '一般情况下，订单会在24小时内发货，偏远地区可能需要1-3天。', usage: 64 },
  { id: 6, title: '请稍等', category: '其他', content: '请稍等，我正在为您查询相关信息。', usage: 156 }
])

const useSnippet = (snippet) => {
  ElMessage.success(`已使用快捷回复：${snippet.title}`)
}

const editSnippet = (snippet) => {
  editingSnippet.value = snippet
  snippetForm.value = {
    title: snippet.title,
    category: snippet.category,
    content: snippet.content
  }
  showAddDialog.value = true
}

const deleteSnippet = (id) => {
  snippets.value = snippets.value.filter(s => s.id !== id)
  ElMessage.success('删除成功')
}

const handleSave = () => {
  if (editingSnippet.value) {
    const index = snippets.value.findIndex(s => s.id === editingSnippet.value.id)
    if (index !== -1) {
      snippets.value[index] = {
        ...snippets.value[index],
        ...snippetForm.value
      }
    }
    ElMessage.success('更新成功')
  } else {
    snippets.value.push({
      id: Date.now(),
      ...snippetForm.value,
      usage: 0
    })
    ElMessage.success('添加成功')
  }
  showAddDialog.value = false
  editingSnippet.value = null
  snippetForm.value = { title: '', category: '问候', content: '' }
}
</script>

<style scoped>
.line-clamp-3 {
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
