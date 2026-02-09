<template>
    <div class="p-6">
        <div class="mb-6">
            <h2 class="text-2xl font-bold text-gray-800">App 管理</h2>
            <p class="text-gray-600 mt-1">管理接入客服系统的应用配置</p>
        </div>

        <div class="bg-white rounded-lg shadow">
            <div class="p-4 border-b border-gray-200 flex justify-between items-center">
                <div class="flex items-center gap-4">
                    <el-input v-model="searchKeyword" placeholder="搜索应用名称或 AppID" clearable style="width: 300px"
                        @clear="loadApps" @keyup.enter="loadApps">
                        <template #prefix>
                            <el-icon> <Search /> </el-icon>
                        </template>
                    </el-input>
                    <el-button type="primary" @click="loadApps">搜索</el-button>
                    <el-select v-model="statusFilter" placeholder="状态筛选" clearable style="width: 150px"
                        @change="loadApps">
                        <el-option label="全部" value="" />
                        <el-option label="启用" :value="1" />
                        <el-option label="禁用" :value="0" />
                    </el-select>
                </div>
                <el-button type="primary" @click="openDialog()">
                    <el-icon class="mr-1"> <Plus /> </el-icon>
                    新增应用
                </el-button>
            </div>

            <el-table :data="apps" v-loading="loading" stripe>
                <el-table-column prop="name" label="应用名称" width="100" align="center" />
                <el-table-column prop="app_id" label="AppID" width="180" align="center" />
                <el-table-column prop="logo" label="Logo" width="100" align="center">
                    <template #default="{ row }">
                        <el-image v-if="row.logo" :src="row.logo" fit="cover"
                            style="width: 40px; height: 40px; border-radius: 4px" />
                        <div v-else
                            class="w-10 h-10 bg-gray-200 rounded flex items-center justify-center text-gray-400 text-xs">
                            暂无
                        </div>
                    </template>
                </el-table-column>
                <el-table-column prop="allow_domain" label="允许域名" min-width="120"  align="center" show-overflow-tooltip />
                <el-table-column prop="contact" label="联系人" width="120" align="center" />
                <el-table-column prop="status" label="状态" width="50" align="center">
                    <template #default="{ row }">
                        <el-tag :type="row.status === 1 ? 'success' : 'danger'">
                            {{ row.status === 1 ? '启用' : '禁用' }}
                        </el-tag>
                    </template>
                </el-table-column>
                <el-table-column prop="created_at" label="创建时间" width="120" align="center">
                    <template #default="{ row }">
                        {{ formatDate(row.CreatedAt) }}
                    </template>
                </el-table-column>
                <el-table-column label="操作" width="120" align="center">
                    <template #default="{ row }">
                        <el-button type="primary" link size="small" @click="openDialog(row)">编辑</el-button>
                        <el-button :type="row.status === 1 ? 'warning' : 'success'" link size="small"
                            @click="toggleStatus(row)">
                            {{ row.status === 1 ? '禁用' : '启用' }}
                        </el-button>
                        <el-button type="danger" link size="small" @click="deleteApp(row)">删除</el-button>
                    </template>
                </el-table-column>
            </el-table>

            <div class="p-4 border-t border-gray-200 flex justify-end">
                <el-pagination v-model:current-page="currentPage" v-model:page-size="pageSize" :total="total"
                    :page-sizes="[10, 20, 50, 100]" layout="total, sizes, prev, pager, next, jumper"
                    @size-change="loadApps" @current-change="loadApps" />
            </div>
        </div>

        <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px" @close="resetForm">
            <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
                <el-form-item label="应用名称" prop="name" class="mr-8">
                    <el-input v-model="form.name" placeholder="请输入应用名称" />
                </el-form-item>
                <el-form-item label="AppID" prop="app_id" class="mr-8">
                    <el-input v-model="form.app_id" placeholder="请输入 AppID" :disabled="!!form.id" />
                    <div v-if="!form.id" class="text-xs text-gray-500 mt-1">留空则自动生成</div>
                </el-form-item>
                <el-form-item label="Logo" prop="logo" class="mr-8">
                    <el-input v-model="form.logo" placeholder="请输入 Logo URL" />
                </el-form-item>
                <el-form-item label="允许域名" prop="allow_domain" class="mr-8">
                    <el-input v-model="form.allow_domain" placeholder="请输入允许的域名，多个用逗号分隔" />
                </el-form-item>
                <el-form-item label="欢迎语" prop="welcome_msg" class="mr-8">
                    <el-input v-model="form.welcome_msg" type="textarea" :rows="3" placeholder="请输入欢迎消息" />
                </el-form-item>
                <el-form-item label="联系人" prop="contact" class="mr-8">
                    <el-input v-model="form.contact" placeholder="请输入联系人信息" />
                </el-form-item>
                <el-form-item label="状态" prop="status">
                    <el-radio-group v-model="form.status">
                        <el-radio :value="1">启用</el-radio>
                        <el-radio :value="0">禁用</el-radio>
                    </el-radio-group>
                </el-form-item>
            </el-form>
            <template #footer>
                <el-button @click="dialogVisible = false">取消</el-button>
                <el-button class="mr-8!" type="primary" @click="submitForm" :loading="submitting">确定</el-button>
            </template>
        </el-dialog>
    </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Plus } from '@element-plus/icons-vue'
import api from '@/script/api'

const loading = ref(false)
const submitting = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref('新增应用')
const searchKeyword = ref('')
const statusFilter = ref('')
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const apps = ref([])
const formRef = ref(null)

const form = ref({
    id: null,
    name: '',
    app_id: '',
    logo: '',
    allow_domain: '',
    welcome_msg: '',
    contact: '',
    status: 1
})

const rules = {
    name: [{ required: true, message: '请输入应用名称', trigger: 'blur' }],
    app_id: [{ required: true, message: '请输入 AppID', trigger: 'blur' }],
    status: [{ required: true, message: '请选择状态', trigger: 'change' }]
}

const loadApps = async () => {
    loading.value = true
    try {
        const params = {
            page: currentPage.value,
            page_size: pageSize.value
        }
        if (searchKeyword.value) {
            params.keyword = searchKeyword.value
        }
        if (statusFilter.value !== '') {
            params.status = statusFilter.value
        }
        const response = await api.listApps(params)
        apps.value = response.data?.data?.data || []
        total.value = response.data?.data?.total || 0
    } catch (error) {
        ElMessage.error('加载应用列表失败')
        console.error(error)
    } finally {
        loading.value = false
    }
}

const openDialog = (row = null) => {
    if (row) {
        dialogTitle.value = '编辑应用'
        form.value = { ...row }
    } else {
        dialogTitle.value = '新增应用'
        resetForm()
    }
    dialogVisible.value = true
}

const resetForm = () => {
    form.value = {
        id: null,
        name: '',
        app_id: '',
        logo: '',
        allow_domain: '',
        welcome_msg: '',
        contact: '',
        status: 1
    }
    formRef.value?.clearValidate()
}

const submitForm = async () => {
    if (!formRef.value) return
    await formRef.value.validate()
    submitting.value = true
    try {
        if (form.value.app_id) {
            await api.updateApp(form.value.app_id, form.value)
            ElMessage.success('更新成功')
        } else {
            await api.createApp(form.value)
            ElMessage.success('创建成功')
        }
        dialogVisible.value = false
        loadApps()
    } catch (error) {
        ElMessage.error(error.message || '操作失败')
        console.error(error)
    } finally {
        submitting.value = false
    }
}

const toggleStatus = async (row) => {
    const newStatus = row.status === 1 ? 0 : 1
    try {
        await ElMessageBox.confirm(
            `确定要${newStatus === 1 ? '启用' : '禁用'}该应用吗？`,
            '提示',
            { type: 'warning' }
        )
        await api.updateApp(row.app_id, { status: newStatus })
        ElMessage.success(`${newStatus === 1 ? '启用' : '禁用'}成功`)
        loadApps()
    } catch (error) {
        if (error !== 'cancel') {
            ElMessage.error('操作失败')
            console.error(error)
        }
    }
}

const deleteApp = async (row) => {
    try {
        await ElMessageBox.confirm(
            `确定要删除应用"${row.name}"吗？此操作不可恢复！`,
            '警告',
            { type: 'warning', confirmButtonText: '确定删除', confirmButtonClass: 'el-button--danger' }
        )
        await api.deleteApp(row.app_id)
        ElMessage.success('删除成功')
        loadApps()
    } catch (error) {
        if (error !== 'cancel') {
            ElMessage.error('删除失败')
            console.error(error)
        }
    }
}

const formatDate = (dateString) => {
    if (!dateString) return '-'
    const date = new Date(dateString)
    return date.toLocaleString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit'
    })
}

onMounted(() => {
    loadApps()
})
</script>

<style scoped>
:deep(.el-table) {
    font-size: 14px;
}

:deep(.el-table .cell) {
    padding: 8px 0;
}
</style>
