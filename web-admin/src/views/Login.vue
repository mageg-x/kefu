<template>
    <div class="min-h-screen bg-white overflow-hidden relative">
        <!-- 左侧蓝色渐变背景 -->
        <div class="absolute left-0 top-0 w-full h-full bg-gradient-to-br from-blue-600 via-blue-700 to-indigo-900 z-0">
            <!-- 装饰三角形 -->
            <div class="absolute top-[20%] left-[15%] w-40 h-40 bg-purple-700/30 rounded-2xl transform rotate-12"></div>
            <div class="absolute bottom-[25%] left-[30%] w-64 h-64 bg-indigo-800/20 rounded-2xl transform -rotate-12">
            </div>
            <div class="absolute top-[50%] left-[20%] w-32 h-32 bg-blue-600/25 rounded-2xl transform rotate-6"></div>

            <!-- 左侧文字内容 -->
            <div class="absolute left-0 top-0 w-3/5 h-full flex flex-col justify-center px-20 text-white z-10">
                <div class="mb-8 flex items-center animate-pulse">
                    <img src="@/assets/logo.png" alt="零点客服系统" class="w-24 h-24 mr-3 drop-shadow-lg animate-glow">
                    <span class="text-4xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-white to-blue-200 drop-shadow-md">零点</span>
                    <span class="text-4xl font-bold ml-3 bg-clip-text text-transparent bg-gradient-to-r from-blue-300 to-cyan-200 drop-shadow-md">客服系统</span>
                </div>
                <h1 class="text-5xl font-bold mb-6 leading-tight bg-clip-text text-transparent bg-gradient-to-r from-white via-blue-100 to-cyan-200 drop-shadow-lg animate-pulse">让客户服务如此简单</h1>
                <p class="text-blue-100 opacity-95 max-w-md text-lg drop-shadow-md animate-pulse">
                    <span class="text-cyan-300 font-medium">『魔力空间软件技术工作室』</span>官方推出的专业客服系统<br>
                    助力企业打造一流的服务体验
                </p>
            </div>
        </div>

        <!-- 右侧白色背景（带弧度） -->
        <div
            class="absolute right-0 top-0 w-1/2 h-full bg-gradient-to-bl from-white to-blue-50 z-20 transform -skew-x-6 origin-top-right">
        </div>

        <!-- 登录卡片 -->
        <div
            class="absolute right-[10%] top-1/2 transform -translate-y-1/2 bg-white/95 backdrop-blur-sm rounded-2xl shadow-2xl p-8 border border-blue-200 w-96 z-30 hover:shadow-blue-500/30 transition-all duration-300 animate-float animate-border-glow">
            <!-- 标题和二维码 -->
            <div class="flex justify-between items-center mb-8">
                <h2 class="text-2xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-blue-800 to-indigo-600">客服登录</h2>
            </div>

            <!-- 登录表单 -->
            <el-form :model="loginForm" @submit.prevent="handleLogin">
                <!-- 账号输入 -->
                <el-form-item class="mb-4">
                    <el-input v-model="loginForm.username" placeholder="请输入账号" :prefix-icon="UserFilled" size="large"
                        class="rounded-xl border-blue-200 hover:border-blue-400 focus:border-blue-500 focus:ring-2 focus:ring-blue-300/50 transition-all duration-300" :disabled="loading" />
                </el-form-item>

                <!-- 密码输入 -->
                <el-form-item class="mb-6">
                    <el-input v-model="loginForm.password" type="password" placeholder="请输入密码" :prefix-icon="Lock"
                        size="large" class="rounded-xl border-blue-200 hover:border-blue-400 focus:border-blue-500 focus:ring-2 focus:ring-blue-300/50 transition-all duration-300" :disabled="loading" />
                </el-form-item>

                <!-- 记住密码 -->
                <el-form-item class="mb-6">
                    <el-checkbox v-model="rememberPassword" size="large" class="text-blue-700 hover:text-blue-500 transition-colors duration-300">记住密码</el-checkbox>
                </el-form-item>

                <!-- 错误提示 -->
                <el-alert v-if="errorMessage" :title="errorMessage" type="error" show-icon class="mb-4 rounded-lg"
                    :closable="false" />

                <!-- 登录按钮 -->
                <el-form-item>
                    <el-button type="primary" native-type="submit" size="large" :loading="loading"
                        class="w-full py-3 rounded-xl bg-gradient-to-r from-blue-600 to-indigo-600 hover:from-blue-700 hover:to-indigo-700 text-white font-medium shadow-lg hover:shadow-blue-500/50 transition-all duration-300 transform hover:scale-[1.02] active:scale-[0.98]">
                        登录
                    </el-button>
                </el-form-item>
            </el-form>
        </div>

        <!-- 底部版权信息 -->
        <div class="absolute bottom-4 left-0 right-0 text-center text-xs z-40 drop-shadow-md animate-pulse">
            魔力空间软件技术工作室 | 商业版 客服系统
            <p class="mt-1.5 bg-clip-text text-transparent bg-gradient-to-r from-gray-100 to-gray-900 drop-shadow-sm">&copy; 2026 ZeroVpn. All rights reserved.</p>
        </div>
    </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { UserFilled, Lock } from '@element-plus/icons-vue'
import api from '@/script/api'
import { useStore } from '@/script/store'

const store = useStore()
// 路由
const router = useRouter()

// 表单数据
const loginForm = ref({
    username: '',
    password: ''
})

// 记住密码
const rememberPassword = ref(false)

// 状态
const loading = ref(false)
const errorMessage = ref('')

// 错误提示映射
const errorMessages = {
    'invalid': '用户或密码不正确',
    'Invalid': '用户或密码不正确',
    'credentials': '用户或密码不正确',
    'password': '用户或密码不正确',
    'username': '用户或密码不正确',
    '参数': '用户或密码不正确',
    'parameter': '用户或密码不正确',
    'network': '网络连接失败，请检查网络设置',
    'Network': '网络连接失败，请检查网络设置',
    'connection': '网络连接失败，请检查网络设置',
    'Connection': '网络连接失败，请检查网络设置',
    'timeout': '网络连接失败，请检查网络设置',
    'Timeout': '网络连接失败，请检查网络设置',
    'server': '服务器错误，请稍后重试',
    'Server': '服务器错误，请稍后重试',
    '500': '服务器错误，请稍后重试',
    '502': '服务器错误，请稍后重试',
    '503': '服务器错误，请稍后重试',
    '504': '服务器错误，请稍后重试'
}


// 加载保存的登录信息
const loadSavedLogin = () => {
    try {
        const savedLogin = localStorage.getItem('login_credentials')
        if (savedLogin) {
            const { username, password, remember } = JSON.parse(savedLogin)
            loginForm.value.username = username
            loginForm.value.password = password
            rememberPassword.value = remember
        }
    } catch (error) {
        console.error('加载保存的登录信息失败:', error)
    }
}

// 保存登录信息
const saveLogin = (username, password, remember) => {
    try {
        if (remember) {
            localStorage.setItem('login_credentials', JSON.stringify({
                username,
                password,
                remember
            }))
        } else {
            localStorage.removeItem('login_credentials')
        }
    } catch (error) {
        console.error('保存登录信息失败:', error)
    }
}

// 初始化加载
loadSavedLogin()


// 登录处理
const handleLogin = async () => {
    try {
        // 重置错误信息
        errorMessage.value = ''
        // 设置加载状态
        loading.value = true

        // 验证表单
        if (!loginForm.value.username || !loginForm.value.password) {
            errorMessage.value = '请输入账号和密码'
            return
        }

        // 调用登录 API
        const response = await api.login(loginForm.value.username, loginForm.value.password)
        console.log('登录成功:', response)

        // 保存登录信息
        saveLogin(loginForm.value.username, loginForm.value.password, rememberPassword.value)

        store.setUser(response?.data?.token, {
            id: response?.data?.user?.ID,
            name: response?.data?.user?.username,
            role: response?.data?.user?.role,
            avatar: response?.data?.user?.avatar,
        })

        // 登录成功后，路由守卫会自动根据角色跳转到对应页面
        router.push('/home')
    } catch (error) {
        console.error('登录失败:', error)
        const errorMsg = error.message || ''
        const matchedKey = Object.keys(errorMessages).find(key => errorMsg.includes(key))
        errorMessage.value = matchedKey ? errorMessages[matchedKey] : '登录失败，请检查账号和密码'
    } finally {
        // 清除加载状态
        loading.value = false
    }
}
</script>

<style scoped>
/* 自定义样式 */
@media (max-width: 768px) {
    .absolute.left-0.top-0.w-3\/5 {
        width: 100%;
        height: 400px;
    }

    .absolute.right-0.top-0.w-1\/2 {
        width: 100%;
        top: 400px;
        height: calc(100% - 400px);
        transform: none;
    }

    .absolute.right-[10%].top-1\/2 {
        right: 5%;
        left: 5%;
        top: 300px;
        transform: none;
        width: 90%;
        max-width: 400px;
    }
}

/* 科技感样式 */
@keyframes pulse {
    0%, 100% {
        opacity: 1;
    }
    50% {
        opacity: 0.8;
    }
}

@keyframes glow {
    0%, 100% {
        filter: drop-shadow(0 0 8px rgba(59, 130, 246, 0.5));
    }
    50% {
        filter: drop-shadow(0 0 16px rgba(59, 130, 246, 0.8));
    }
}

@keyframes float {
    0%, 100% {
        transform: translateY(0) rotate(0deg);
    }
    50% {
        transform: translateY(-10px) rotate(2deg);
    }
}

@keyframes border-glow {
    0%, 100% {
        border-color: rgba(59, 130, 246, 0.5);
        box-shadow: 0 0 8px rgba(59, 130, 246, 0.3);
    }
    50% {
        border-color: rgba(59, 130, 246, 0.8);
        box-shadow: 0 0 16px rgba(59, 130, 246, 0.6);
    }
}

.animate-pulse {
    animation: pulse 3s ease-in-out infinite;
}

.animate-glow {
    animation: glow 2s ease-in-out infinite;
}

.animate-float {
    animation: float 6s ease-in-out infinite;
}

.animate-border-glow {
    animation: border-glow 2s ease-in-out infinite;
}

/* 自定义 Element Plus 样式 */
:deep(.el-input__wrapper) {
    border-radius: 0.5rem;
}

:deep(.el-button--primary) {
    border-radius: 0.5rem;
    font-weight: 500;
}

:deep(.el-form-item) {
    margin-bottom: 1rem;
}

:deep(.el-alert) {
    border-radius: 0.5rem;
}
</style>
