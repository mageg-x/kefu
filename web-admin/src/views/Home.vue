<!-- src/views/Home.vue -->
<template>
    <div class="flex h-screen bg-gray-50">
        <!-- 左侧导航栏 -->
        <aside :class="['bg-white border-r border-gray-200 flex flex-col shadow-sm', isCollapsed ? 'w-20' : 'w-64']">
            <!-- Logo区域 -->
            <div class="p-4 border-b border-gray-200 bg-gray-50 flex items-center justify-between">
                <div v-if="!isCollapsed" class="flex items-center">
                    <img  src="@/assets/logo.png" alt="Logo" class="rounded-lg w-10 h-10" />
                    <div  class="ml-2">
                        <h1 class="text-xl font-bold text-gray-800">客服系统</h1>
                        <p class="text-xs text-gray-500">专业客服管理平台</p>
                    </div>
                </div>

                <!-- 收起/展开按钮 -->
                <el-tooltip :content="isCollapsed ? '展开侧边栏' : '收起侧边栏'" placement="right">
                    <el-button 
                        @click="toggleCollapse" 
                        :icon="isCollapsed ? Expand : Fold"
                        circle
                        size="large"
                        class="text-xl!"
                    />
                </el-tooltip>
            </div>

            <!-- 导航菜单 -->
            <nav class="flex-1 py-4 overflow-y-auto g-gray-b50">
                <ul class="space-y-1 px-3">
                    <li v-for="item in menuItems" :key="item.name">
                        <el-tooltip :content="item.name" placement="right" :disabled="!isCollapsed">
                            <router-link :to="item.path"
                                class="flex items-center py-3 text-sm rounded-xl"
                                active-class="bg-gradient-to-r from-blue-600 to-indigo-700 text-white shadow-md"
                                :class="[$route.path.startsWith(item.path) ? 'text-gray-600 hover:bg-gray-100' : 'text-gray-600 hover:bg-gray-100', isCollapsed?'w-10 h-10 m-auto flex justify-center':'px-4']">
                                <component :is="item.icon" class="w-6 h-6" />
                                <span v-if="!isCollapsed" class="font-medium ml-3">{{ item.name }}</span>
                            </router-link>
                        </el-tooltip>
                    </li>
                </ul>
            </nav>

            <!-- 用户信息区域 -->
            <div class="p-4 border-t border-gray-200 bg-gray-50">
                <el-tooltip :content="store.user?.name || '用户'" placement="right" :disabled="!isCollapsed">
                    <div class="flex items-center gap-3 mb-4" :class="[isCollapsed ? 'px-0 aspect-square!' : 'px-2']">
                        <el-avatar size="large" :src="store.user?.avatar || ''" class="w-12! h-12!">
                            {{ store.user?.name?.charAt(0) || 'U' }}
                        </el-avatar>
                        <div v-if="!isCollapsed" class="flex-1 min-w-0">
                            <p class="text-sm font-medium text-gray-800 truncate">{{ store.user?.name || '用户' }}</p>
                            <p class="text-xs text-gray-500">{{ store.isAdmin ? '管理员' : '客服专员' }}</p>
                        </div>
                    </div>
                </el-tooltip>
                <el-tooltip content="退出登录" placement="right" :disabled="!isCollapsed">
                    <el-button type="danger" plain @click="logout" :class="['btn-logout', isCollapsed ? 'flex-1' : 'w-full']">
                        <template #icon>
                            <svg t="1770519176568" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="7394" width="200" height="200">
                                <path d="M136.2 513.8v366.7c0 26.4 23.4 47.8 52.4 47.8h261.8c28.9 0 52.4-21.4 52.4-47.8s-23.4-47.8-52.4-47.8H241V194.8h209.5c28.9 0 52.4-21.4 52.4-47.8s-23.4-47.8-52.4-47.8H188.6c-28.9 0-52.4 21.4-52.4 47.8v366.8z m757.5-34.6c10.1 8.8 16.4 21.5 16.4 35.6 0 14.2-6.3 26.8-16.4 35.6L728.3 693.8c-8.8 7.6-20.5 12.3-33.2 12.3-27.4 0-49.6-21.5-49.6-47.9 0-14.1 6.4-26.8 16.5-35.5l69.3-60.1h-301c-27.4 0-49.6-21.4-49.6-47.8s22.2-47.8 49.6-47.8h300.8l-69.3-60.1c-10.1-8.7-16.5-21.5-16.5-35.6 0-26.4 22.3-47.8 49.6-47.8 12.8 0 24.4 4.6 33.2 12.3l165.6 143.4z" fill="#160B0B" p-id="7395"></path>
                            </svg>
                        </template>
                        <span v-if="!isCollapsed">退出登录</span>
                    </el-button>
                </el-tooltip>
            </div>
        </aside>

        <!-- 右侧内容区 -->
        <main class="flex-1 overflow-auto">
            <router-view />
        </main>
    </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useStore } from '@/script/store'
import api from '@/script/api'
import { House, Avatar, Service, Setting,   ChatDotRound,   Postcard, User, Fold, Expand, Monitor} from '@element-plus/icons-vue'

const router = useRouter()
const route = useRoute()
const store = useStore()

// 侧边栏收起状态
const isCollapsed = ref(false)

// 初始化加载收起状态
onMounted(() => {
    const savedCollapsed = localStorage.getItem('sidebar_collapsed')
    if (savedCollapsed !== null) {
        isCollapsed.value = savedCollapsed === 'true'
    }
})

// 切换收起/展开
const toggleCollapse = () => {
    isCollapsed.value = !isCollapsed.value
    localStorage.setItem('sidebar_collapsed', isCollapsed.value.toString())
}

const menuItems = computed(() => {
    if (store.isAdmin) {
        return [
            { name: '数据看板', path: '/home/dashboard', icon: House },
            { name: '访客管理', path: '/home/visitors', icon: Avatar },
            { name: '客服团队', path: '/home/staff', icon: Service },
            { name: 'App管理', path: '/home/apps', icon: Monitor },
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
