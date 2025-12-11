<template>
    <div class="fixed bottom-5 md:bottom-20 right-5 z-50 font-sans">
        <!-- 悬浮按钮 -->
        <button v-if="!isOpen" @click="isOpen = true"
            class="flex flex-col items-center p-1 rounded-2xl bg-white text-blue-600 shadow-lg hover:shadow-xl transition-all duration-300 hover:scale-105 border border-blue-200">
            <div class="w-12 h-12 rounded-full bg-transparent shadow-md flex items-center justify-center">
                <img src="../assets/logo.png" alt="logo" class="w-12 h-12 object-contain">
            </div>
            <div></div>
            <span class="font-medium text-xs p-2 whitespace-nowrap">{{ t('onlineConsultation') }}</span>
        </button>

        <!-- 聊天窗口 - 使用新的 ChatWind 组件 -->
        <ChatWind 
            v-else 
            :messages="messages"
            :user-id="userId"
            @close="isOpen = false"
            @send-message="handleSendMessage"
            @language-changed="handleLanguageChanged"
        />
    </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import ChatWind from './ChatWind.vue'

// 定义props
const props = defineProps({
    appId: { type: String, default: 'default' },
    backendUrl: { type: String, default: 'http://localhost:9000' }
})

// 生成唯一机器ID的函数
function generateMachineId() {
    // 创建一个基于设备信息的唯一ID
    const fingerprint = [
        navigator.userAgent,
        navigator.language,
        screen.width + 'x' + screen.height,
        (new Date()).getTimezoneOffset(),
        navigator.hardwareConcurrency,
        navigator.platform
    ].join('|');
    
    // 使用简单的哈希函数生成ID
    let hash = 0;
    for (let i = 0; i < fingerprint.length; i++) {
        const char = fingerprint.charCodeAt(i);
        hash = ((hash << 5) - hash) + char;
        hash = hash & hash; // 转换为32位整数
    }
    
    // 转换为16进制并添加前缀
    return Math.abs(hash).toString(16);
}

// 获取或生成唯一用户ID
function getOrCreateUserId() {
    let userId;
    
    // 尝试从localStorage获取
    if (typeof window !== 'undefined' && window.localStorage) {
        userId = localStorage.getItem('kefu_user_id');
    }
    
    // 如果localStorage没有，尝试从cookie获取
    if (!userId && typeof document !== 'undefined') {
        const cookieMatch = document.cookie.match(/kefu_user_id=([^;]+)/);
        if (cookieMatch) {
            userId = decodeURIComponent(cookieMatch[1]);
        }
    }
    
    // 如果都没有，生成新的ID
    if (!userId) {
        userId = generateMachineId();
        
        // 存储到localStorage
        if (typeof window !== 'undefined' && window.localStorage) {
            localStorage.setItem('kefu_user_id', userId);
        }
        
        // 存储到cookie，过期时间设置为10年
        if (typeof document !== 'undefined') {
            const expires = new Date();
            expires.setTime(expires.getTime() + (10 * 365 * 24 * 60 * 60 * 1000));
            document.cookie = `kefu_user_id=${encodeURIComponent(userId)};expires=${expires.toUTCString()};path=/`;
        }
    }
    
    return userId;
}

// 语言相关状态
// 从localStorage加载语言设置，如果没有则使用默认中文
const currentLanguage = ref(localStorage.getItem('chat_language') || 'zh')

// 国际化翻译
const translations = {
  zh: {
    onlineConsultation: '在线咨询'
  },
  en: {
    onlineConsultation: 'Online Consultation'
  },
  hi: {
    onlineConsultation: 'ऑनलाइन परामर्श'
  },
  ru: {
    onlineConsultation: 'Онлайн-консультация'
  },
  de: {
    onlineConsultation: 'Online-Beratung'
  },
  fr: {
    onlineConsultation: 'Consultation en ligne'
  },
  ja: {
    onlineConsultation: 'オンライン相談'
  }
}

// 获取翻译
const t = (key) => {
  return translations[currentLanguage.value][key] || key
}

// 响应式数据
const isOpen = ref(false)
const messages = ref([])
const userId = ref(getOrCreateUserId())

// 处理语言变化
const handleLanguageChanged = (lang) => {
  currentLanguage.value = lang
  // 保存语言设置到localStorage
  localStorage.setItem('chat_language', lang)
}

// 方法
async function handleSendMessage(text) {
    messages.value.push({ from: 'me', text })
    // TODO: 调用后端或直接连悟空IM
    console.log('发送消息:', text)
}
</script>