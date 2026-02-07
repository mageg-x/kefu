<!-- src/views/Home.vue -->
<template>
    <div class="flex h-screen bg-gray-50">
        <!-- 左侧导航栏 -->
        <aside class="w-64 bg-white border-r border-gray-200 flex flex-col shadow-sm">
            <!-- Logo区域 -->
            <div class="p-6 border-b border-gray-200 bg-gray-50">
                <div class="flex items-center gap-3">
                    <img src="@/assets/logo.png" alt="Logo" class="w-10 h-10 rounded-lg" />
                    <div>
                        <h1 class="text-xl font-bold text-gray-800">客服系统</h1>
                        <p class="text-xs text-gray-500">专业客服管理平台</p>
                    </div>
                </div>
            </div>

            <!-- 导航菜单 -->
            <nav class="flex-1 py-4 overflow-y-auto g-gray-b50">
                <ul class="space-y-1 px-3">
                    <li v-for="item in menuItems" :key="item.name">
                        <router-link :to="item.path"
                            class="flex items-center px-4 py-3 text-sm rounded-xl transition-all duration-200"
                            active-class="bg-gradient-to-r from-blue-600 to-indigo-700 text-white shadow-md"
                            :class="$route.path.startsWith(item.path) ? 'text-gray-600 hover:bg-gray-100' : 'text-gray-600 hover:bg-gray-100'">
                            <component :is="item.icon" class="w-5 h-5 mr-3" />
                            <span class="font-medium">{{ item.name }}</span>
                        </router-link>
                    </li>
                </ul>
            </nav>

            <!-- 用户信息区域 -->
            <div class="p-4 border-t border-gray-200 bg-gray-50">
                <div class="flex items-center gap-3 mb-4 px-2">
                    <el-avatar :size="40" :src="store.user?.avatar || ''">
                        {{ store.user?.name?.charAt(0) || 'U' }}
                    </el-avatar>
                    <div class="flex-1 min-w-0">
                        <p class="text-sm font-medium text-gray-800 truncate">{{ store.user?.name || '用户' }}</p>
                        <p class="text-xs text-gray-500">{{ store.isAdmin ? '管理员' : '客服专员' }}</p>
                    </div>
                </div>
                <el-button type="danger" plain @click="logout" class="w-full btn-logout">
                    <template #icon>
                        <el-icon>
                            <SwitchButton />
                        </el-icon>
                    </template>
                   退出登录
                </el-button>
            </div>
        </aside>

        <!-- 右侧内容区 -->
        <main class="flex-1 overflow-auto">
            <router-view />
        </main>
    </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useStore } from '@/script/store'
import api from '@/script/api'
import { House, Avatar,  Service,  Setting,   ChatDotRound,   Postcard, User} from '@element-plus/icons-vue'

const router = useRouter()
const route = useRoute()
const store = useStore()

const menuItems = computed(() => {
    if (store.isAdmin) {
        return [
            { name: '数据看板', path: '/home/dashboard', icon: House },
            { name: '访客管理', path: '/home/visitors', icon: Avatar },
            { name: '客服团队', path: '/home/staff', icon: Service },
            { name: '系统设置', path: '/home/settings', icon: Setting }
        ]
    } else {
        return [
            { name: '我的会话', path: '/home/inbox', icon: ChatDotRound },
            { name: '快捷回复', path: '/home/snippets', icon: Postcard },
            { name: '个人资料', path: '/home/profile', icon: User }
        ]
    }
})

const logout = async () => {
    try {
        await api.logout()
    } catch (error) {
        console.error('登出失败:', error)
    } finally {
        router.push('/login')
    }
}
</script>

<style scoped>
aside {
    background: linear-gradient(180deg, #ffffff 0%, #f8fafc 100%);
}

/* 滚动条样式 */
aside::-webkit-scrollbar {
    width: 4px;
}

aside::-webkit-scrollbar-track {
    background: transparent;
}

aside::-webkit-scrollbar-thumb {
    background: #cbd5e1;
    border-radius: 2px;
}

aside::-webkit-scrollbar-thumb:hover {
    background: #94a3b8;
}

/* 路由链接激活状态 */
.router-link-active {
    background: linear-gradient(135deg, #2563eb 0%, #4f46e5 100%) !important;
    color: white !important;
    box-shadow: 0 4px 6px -1px rgba(37, 99, 235, 0.2);
}

.router-link-active .el-icon {
    color: white !important;
}

.btn-logout {
    font-size: 1rem;
    font-weight: 700;
    height: 2rem;
}
</style>
