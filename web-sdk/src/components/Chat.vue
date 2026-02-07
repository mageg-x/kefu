<template>
  <div class="h-screen w-screen bg-gray-100 flex flex-col">
    <div class="bg-gradient-to-r from-blue-500 to-blue-600 px-4 py-3 flex items-center justify-between flex-shrink-0">
      <div class="flex items-center gap-3">
        <div class="w-10 h-10 rounded-full bg-white flex items-center justify-center">
          <img :src="logoUrl" alt="logo" class="w-10 h-10 object-contain" />
        </div>
        <div class="text-white">
          <h3 class="font-semibold text-base">唯一客服</h3>
          <div class="flex items-center gap-1 text-xs opacity-90 mt-0.5">
            <span class="w-2 h-2 rounded-full bg-green-400 animate-pulse"></span>
            <span>{{ t('online') }}</span>
          </div>
        </div>
      </div>
    </div>

    <div ref="chatMessages" class="flex-1 overflow-y-auto p-4 space-y-4">
      <div v-for="msg in messages" :key="msg.id">
        <div v-if="isTimestampMessage(msg)" class="text-center text-xs text-gray-500 mb-4">
          <span class="px-3 py-1 bg-gray-200 rounded-full">{{ msg.content }}</span>
        </div>

        <div v-else-if="msg.type === 'system'" class="flex justify-center">
          <div class="text-xs text-gray-500 bg-gray-100 px-4 py-2 rounded-lg max-w-xs text-center">
            {{ normalizeContent(msg.content) }}
          </div>
        </div>

        <div v-else-if="msg.type === 'robot'" class="flex gap-3">
          <div class="w-8 h-8 rounded-full bg-gray-200 flex-shrink-0 flex items-center justify-center">
            <span class="text-xs font-bold">AI</span>
          </div>
          <div class="flex flex-col gap-2 max-w-[70%]">
            <span class="text-xs text-gray-500">{{ msg.sender }}</span>
            <div class="bg-white px-4 py-1.5 rounded-lg rounded-tl-none shadow-sm border border-gray-100">
              <div class="text-sm" v-html="renderMarkdown(msg.content)"></div>
            </div>
          </div>
        </div>

        <div v-else-if="msg.type === 'service'" class="flex gap-3">
          <div class="w-8 h-8 rounded-full bg-blue-100 flex-shrink-0 flex items-center justify-center">
            <img :src="logoUrl" alt="客服" class="w-6 h-6 object-contain" />
          </div>
          <div class="flex flex-col gap-2 max-w-[70%]">
            <span class="text-xs text-gray-500">{{ msg.sender }}</span>
            <div class="bg-white px-4 py-1.5 rounded-lg rounded-tl-none shadow-sm border border-gray-100">
              <div v-if="hasImages(msg.content)" class="space-y-2">
                <img v-for="(src, index) in extractImages(msg.content)" :key="'img-' + msg.id + '-' + index" :src="src"
                  :alt="'图片' + index"
                  class="max-w-[200px] max-h-[200px] rounded cursor-pointer hover:opacity-90 transition-opacity" />
              </div>
              <div v-else-if="hasVoice(msg.content)" class="flex items-center gap-2 text-sm">
                <span v-html="VoiceIcon()"></span>
                {{ t('voiceMessage') }}
                <span class="text-xs text-gray-500">{{ formatDuration(msg.duration) }}</span>
              </div>
              <div v-else class="text-sm" v-html="renderMarkdown(msg.content)"></div>
            </div>
          </div>
        </div>

        <div v-else-if="msg.type === 'user'" class="flex gap-3 justify-end">
          <div class="flex flex-col max-w-[70%] items-end gap-2">
            <div class="bg-blue-100 text-black px-4 py-1.5 shadow-sm">
              <div v-if="hasImages(msg.content)" class="space-y-2">
                <img v-for="(src, index) in extractImages(msg.content)" :key="'img-' + msg.id + '-' + index" :src="src"
                  :alt="'图片' + index"
                  class="max-w-[100px] max-h-[100px] rounded cursor-pointer hover:opacity-90 transition-opacity" />
              </div>
              <div v-else-if="hasVoice(msg.content)" class="flex items-center gap-2 text-sm">
                <button @click="togglePlayAudio(msg)"
                  class="flex items-center justify-center text-blue-600 hover:text-blue-800">
                  <span v-html="PlayPauseIcon(msg.isPlaying)"></span>
                </button>
                <span>{{ t('voiceMessage') }}</span>
                <span class="text-xs text-gray-500">{{ formatDuration(msg.duration) }}</span>
                <audio ref="audioElements" :src="msg.audioUrl" @ended="onAudioEnded(msg)" class="hidden"></audio>
              </div>
              <div v-else class="text-sm break-all" v-html="renderMarkdown(msg.content)"></div>
            </div>
          </div>

          <span class="w-6 h-6 overflow-hidden" v-html="avatars[`head${avatarNumber}`] || avatars.head01"></span>

        </div>
      </div>
    </div>

    <div class="quick-replies px-2 py-1 bg-gray-50 border-y border-gray-200 flex-shrink-0">
      <div class="flex gap-1 overflow-x-auto">
        <button
          class="px-4 py-1.5 bg-blue-50 text-xs text-blue-700 rounded-full shadow-sm hover:bg-blue-100 transition whitespace-nowrap border border-blue-200">{{
            t('featureIntro') }}</button>
        <button
          class="px-4 py-1.5 bg-green-50 text-xs text-green-700 rounded-full shadow-sm hover:bg-green-100 transition whitespace-nowrap border border-green-200">{{
            t('humanService') }}</button>
        <button
          class="px-4 py-1.5 bg-purple-50 text-xs text-purple-700 rounded-full shadow-sm hover:bg-purple-100 transition whitespace-nowrap border border-purple-200">{{
            t('aiService') }}</button>
      </div>
    </div>

    <div class="bg-white border-t border-gray-200 p-3 flex-shrink-0">
      <div class="flex flex-col">
        <div class="flex items-center gap-2 pb-2 relative">
          <button class="text-gray-500 hover:text-blue-500 transition p-2 rounded-full hover:bg-gray-100 emoji-button"
            @click="toggleEmojiPicker" :title="t('emojiButton')">
            <span v-html="EmojiIcon()"></span>
          </button>

          <div v-if="showEmojiPicker"
            class="emoji-picker absolute bottom-full left-0 mb-2 bg-white border rounded-lg shadow-lg p-2 z-10 w-64 h-48 overflow-y-auto">
            <div class="grid grid-cols-8 gap-1">
              <button v-for="emoji in emojis" :key="emoji" class="text-xl hover:bg-gray-100 rounded p-1"
                @click="selectEmoji(emoji)">
                {{ emoji }}
              </button>
            </div>
          </div>

          <button class="text-gray-500 hover:text-blue-500 transition p-2 rounded-full hover:bg-gray-100"
            @click="selectImage" :title="t('imageButton')">
            <span v-html="ImageIcon()"></span>
          </button>

          <button class="text-gray-500 hover:text-blue-500 transition p-2 rounded-full hover:bg-gray-100 relative"
            :class="{ 'text-red-500': recording }" @click="toggleVoiceRecording"
            :title="recording ? t('releaseToStop') : t('holdToRecord')">
            <template v-if="recording">
              <div class="wave-animation">
                <div class="wave-bar"></div>
                <div class="wave-bar"></div>
                <div class="wave-bar"></div>
                <div class="wave-bar"></div>
                <div class="wave-bar"></div>
              </div>
            </template>
            <template v-else>
              <span v-html="MicrophoneIcon()"></span>
            </template>
          </button>

          <button
            class="text-gray-500 hover:text-blue-500 transition p-2 rounded-full hover:bg-gray-100 relative language-button"
            @click="toggleLanguageMenu" :title="t('languageButton')">
            <span v-html="LanguageIcon()"></span>

            <div v-if="showLanguageMenu"
              class="language-menu absolute bottom-full right-0 mb-2 bg-white border rounded-lg shadow-lg p-1 z-10 w-28"
              @click.stop>
              <button v-for="lang in languages" :key="lang.value"
                class="w-full text-left px-3 py-2 text-sm hover:bg-gray-100 rounded transition-colors"
                :class="{ 'bg-blue-100 text-blue-700': currentLanguage.value === lang.value }"
                @click="changeLanguage(lang.value)">
                {{ lang.label }}
              </button>
            </div>
          </button>
        </div>

        <div class="relative">
          <div ref="inputDiv" contenteditable="true" @input="handleInput" @paste="handlePaste"
            @keypress.enter.exact="sendMessage" @keypress.enter.shift="" placeholder="Please enter your message"
            class="w-full border border-gray-300 rounded-lg px-4 py-3 pr-16 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-y min-h-[80px] max-h-[150px] overflow-y-auto"
            style="white-space: pre-wrap; word-break: break-word;"></div>
          <button @click="sendMessage"
            class="absolute right-2 bottom-2 bg-blue-500 text-white px-2 py-1 mb-1 rounded-lg text-sm hover:bg-blue-600 transition shadow-sm disabled:bg-gray-300 disabled:cursor-not-allowed"
            :disabled="!inputText.trim()">
            {{ t('send') }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, watch, computed } from 'vue'
import { marked } from 'marked'
import { markedHighlight } from 'marked-highlight'
import hljs from 'highlight.js'
import 'highlight.js/styles/github.css'
import { VoiceIcon, PlayPauseIcon, MicrophoneIcon, LanguageIcon, EmojiIcon, ImageIcon, avatars } from './svg.js'

marked.use(markedHighlight({
  highlight: (code, lang) => {
    if (lang && hljs.getLanguage(lang)) {
      try {
        return hljs.highlight(code, { language: lang }).value
      } catch (err) {
        console.error('Highlight error:', err)
      }
    }
    return hljs.highlightAuto(code).value
  }
}))

// 计算 logo URL
const logoUrl = computed(() => {
    return './logo.png'
})

function calculateCRC(str) {
  let crc = 0
  for (let i = 0; i < str.length; i++) {
    crc = (crc * 31 + str.charCodeAt(i)) % 99 + 1
  }
  return crc.toString().padStart(2, '0')
}

function normalizeContent(content) {
  return typeof content === 'string' ? content : ''
}

function hasImages(content) {
  return typeof content === 'string' && content.includes('[img:')
}

function hasVoice(content) {
  if (typeof content !== 'string') return false
  return content.startsWith('[voice-message]') ||
    content.includes('[语音消息]') ||
    content.includes('[Voice Message]') ||
    content.includes('[वॉयस मैसेज]') ||
    content.includes('[Голосовое сообщение]') ||
    content.includes('[Sprachnachricht]') ||
    content.includes('[Message vocal]') ||
    content.includes('[音声メッセージ]')
}

function extractImages(content) {
  if (!hasImages(content)) return []
  const regex = /\[img:([^\]]+)\]/g
  const matches = []
  let match
  while ((match = regex.exec(content)) !== null) {
    matches.push(match[1])
  }
  return matches
}

function renderMarkdown(raw) {
  const text = normalizeContent(raw)
  if (!text.trim()) return ''

  const cleanText = text.replace(/\r\n/g, '\n').replace(/\r/g, '\n').replace(/\\n/g, '\n')

  const renderer = new marked.Renderer()

  renderer.heading = (token) => {
    const level = token.depth
    const cls = { 1: 'text-2xl', 2: 'text-xl', 3: 'text-lg', 4: 'text-base', 5: 'text-sm', 6: 'text-xs' }[level] || 'text-sm'
    return `<h${level} class="${cls} font-bold my-2">${token.text}</h${level}>`
  }

  renderer.listitem = (token) => {
    return `<li class="ml-6">${token.text}</li>`
  }

  renderer.code = (token) => {
    const code = token.text
    const lang = token.lang || ''
    return code.includes('\n')
      ? `<pre class="bg-gray-100 p-4 rounded text-sm overflow-x-auto"><code class="language-${lang}">${code}</code></pre>`
      : `<code class="bg-gray-100 px-1 py-0.5 rounded text-sm">${code}</code>`
  }

  renderer.link = (token) => {
    const href = token.href
    const txt = token.text || href
    return `<a href="${href}" target="_blank" rel="noopener noreferrer" class="text-blue-600 hover:underline">${txt}</a>`
  }

  try {
    const result = marked.parse(cleanText, { gfm: true, breaks: true, renderer })
    return typeof result === 'string' ? result : ''
  } catch (e) {
    console.error('Markdown error:', e)
    return `<div class="text-red-500">解析错误</div>`
  }
}

function getAppIdFromUrl() {
  const urlParams = new URLSearchParams(window.location.search)
  return urlParams.get('data-kefu-appid') || urlParams.get('appId') || 'default'
}

function getOrCreateUserId() {
  let userId

  if (typeof window !== 'undefined' && window.localStorage) {
    userId = localStorage.getItem('kefu_user_id')
  }

  if (!userId && typeof document !== 'undefined') {
    const cookieMatch = document.cookie.match(/kefu_user_id=([^;]+)/)
    if (cookieMatch) {
      userId = decodeURIComponent(cookieMatch[1])
    }
  }

  if (!userId) {
    const fingerprint = [
      navigator.userAgent,
      navigator.language,
      screen.width + 'x' + screen.height,
      (new Date()).getTimezoneOffset(),
      navigator.hardwareConcurrency,
      navigator.platform
    ].join('|')

    let hash = 0
    for (let i = 0; i < fingerprint.length; i++) {
      const char = fingerprint.charCodeAt(i)
      hash = ((hash << 5) - hash) + char
      hash = hash & hash
    }

    userId = Math.abs(hash).toString(16)

    if (typeof window !== 'undefined' && window.localStorage) {
      localStorage.setItem('kefu_user_id', userId)
    }

    if (typeof document !== 'undefined') {
      const expires = new Date()
      expires.setTime(expires.getTime() + (10 * 365 * 24 * 60 * 60 * 1000))
      document.cookie = `kefu_user_id=${encodeURIComponent(userId)};expires=${expires.toUTCString()};path=/`
    }
  }

  return userId
}

const messages = ref([
  { id: '1', type: 'system', content: '上午 10:30', timestamp: new Date().toISOString() },
  { id: '2', type: 'system', content: '唯一客服 正在为您服务!', timestamp: new Date().toISOString() },
  { id: '3', type: 'robot', sender: '扣子智能体助手', content: '您好！很高兴为您服务。请问有什么可以帮您？', timestamp: new Date().toISOString() },
  { id: '4', type: 'service', sender: '唯一客服', content: '好的，已收到您的♥号，稍等，我们一对一联系您，给您介绍。', timestamp: new Date().toISOString() },
  { id: '5', type: 'service', sender: '唯一客服', content: '没意向就没有聊的必要', timestamp: new Date().toISOString() }
])

const userId = ref(getOrCreateUserId())
const appId = ref(getAppIdFromUrl())
const avatarNumber = calculateCRC(userId.value)
const inputText = ref('')
const inputDiv = ref(null)
const chatMessages = ref(null)
const showEmojiPicker = ref(false)
const recording = ref(false)
const audioChunks = ref([])
const audioElements = ref([])
let mediaRecorder = null

const showLanguageMenu = ref(false)
const currentLanguage = ref(localStorage.getItem('chat_language') || 'zh')

const languages = [
  { label: '中文', value: 'zh' },
  { label: 'English', value: 'en' },
  { label: 'हिंदी', value: 'hi' },
  { label: 'Русский', value: 'ru' },
  { label: 'Deutsch', value: 'de' },
  { label: 'Français', value: 'fr' },
  { label: '日本語', value: 'ja' }
]

const translations = {
  zh: {
    holdToRecord: '开始录音',
    releaseToStop: '结束录音',
    voiceMessage: '[语音消息]',
    seconds: '秒',
    send: '发送',
    cancel: '取消',
    online: '在线',
    emojiButton: '表情',
    imageButton: '图片',
    languageButton: '语言',
    featureIntro: '功能介绍',
    humanService: '人工服务',
    aiService: 'AI客服'
  },
  en: {
    holdToRecord: 'Start recording',
    releaseToStop: 'Stop recording',
    voiceMessage: '[Voice Message]',
    seconds: 'sec',
    send: 'Send',
    cancel: 'Cancel',
    online: 'Online',
    emojiButton: 'Emoji',
    imageButton: 'Image',
    languageButton: 'Language',
    featureIntro: 'Features',
    humanService: 'Human Service',
    aiService: 'AI Service'
  },
  hi: {
    holdToRecord: 'रिकॉर्डिंग शुरू करें',
    releaseToStop: 'रिकॉर्डिंग रोकें',
    voiceMessage: '[वॉयस मैसेज]',
    seconds: 'सेकंड',
    send: 'भेजें',
    cancel: 'रद्द करें',
    online: 'ऑनलाइन',
    emojiButton: 'इमोजी',
    imageButton: 'छवि',
    languageButton: 'भाषा',
    featureIntro: 'विशेषताओं का परिचय',
    humanService: 'मानव सेवा',
    aiService: 'एआई सेवा'
  },
  ru: {
    holdToRecord: 'Начать запись',
    releaseToStop: 'Остановить запись',
    voiceMessage: '[Голосовое сообщение]',
    seconds: 'сек',
    send: 'Отправить',
    cancel: 'Отменить',
    online: 'Онлайн',
    emojiButton: 'Эмодзи',
    imageButton: 'Изображение',
    languageButton: 'Язык',
    featureIntro: 'Функции',
    humanService: 'Человеческая служба',
    aiService: 'AI-служба'
  },
  de: {
    holdToRecord: 'Aufnahme starten',
    releaseToStop: 'Aufnahme stoppen',
    voiceMessage: '[Sprachnachricht]',
    seconds: 'Sek',
    send: 'Senden',
    cancel: 'Abbrechen',
    online: 'Online',
    emojiButton: 'Emoji',
    imageButton: 'Bild',
    languageButton: 'Sprache',
    featureIntro: 'Funktionen',
    humanService: 'Menschlicher Service',
    aiService: 'KI-Service'
  },
  fr: {
    holdToRecord: 'Démarrer l\'enregistrement',
    releaseToStop: 'Arrêter l\'enregistrement',
    voiceMessage: '[Message vocal]',
    seconds: 'sec',
    send: 'Envoyer',
    cancel: 'Annuler',
    online: 'En ligne',
    emojiButton: 'Emoji',
    imageButton: 'Image',
    languageButton: 'Langue',
    featureIntro: 'Fonctionnalités',
    humanService: 'Service humain',
    aiService: 'Service IA'
  },
  ja: {
    holdToRecord: '録音開始',
    releaseToStop: '録音停止',
    voiceMessage: '[音声メッセージ]',
    seconds: '秒',
    send: '送信',
    cancel: 'キャンセル',
    online: 'オンライン',
    emojiButton: '絵文字',
    imageButton: '画像',
    languageButton: '言語',
    featureIntro: '機能紹介',
    humanService: '人間サービス',
    aiService: 'AIサービス'
  }
}

const t = (key) => {
  return translations[currentLanguage.value][key] || key
}

const toggleLanguageMenu = () => {
  showLanguageMenu.value = !showLanguageMenu.value
}

const changeLanguage = (lang) => {
  console.log('切换语言:', lang)
  currentLanguage.value = lang
  localStorage.setItem('chat_language', lang)
  showLanguageMenu.value = false
  scrollToBottom()
}

const emojis = [
  '👌', '💪', '✌️', '👍', '👎', '🌸', '💐', '🌹', '❤️', '💛', '💙', '💜', '🤝', '👏', '🎉', '🎊',
  '🙏', '😢', '😤', '😡', '😱', '😳', '🤯', '🤔', '😀', '😃', '😄', '😁', '😆', '😅', '🤣', '😂',
  '🙂', '🙃', '😉', '😊', '😇', '🥰', '😍', '🤩', '😘', '😗', '😚', '😙', '🥲', '😋', '😛', '😜',
  '🤪', '😝', '🤑', '🤗', '🤭', '🤫', '🤔', '🤐', '🤨', '😐', '😑', '😶', '😏', '😒', '🙄', '😬',
  '😮‍💨', '🤥', '😌', '😔', '😪', '🤤', '😴', '😷', '🤒', '🤕', '🤢', '🤮', '🤧', '🥵', '🥶', '🥴',
  '😵', '🤯', '🤠', '🥳', '😎', '🤓', '🧐', '😕'
]

function isTimestampMessage(msg) {
  return msg.type === 'system' && typeof msg.content === 'string' && msg.content.includes('上午')
}

function generateId() {
  return userId.value + '_' + Date.now().toString(36) + Math.random().toString(36).substr(2, 8)
}

function sendMessage() {
  const content = inputText.value.trim()
  if (!content) return

  messages.value.push({
    id: generateId(),
    type: 'user',
    content,
    timestamp: new Date().toISOString()
  })

  inputText.value = ''
  if (inputDiv.value) inputDiv.value.innerHTML = ''
  scrollToBottom()
}

function handleInput() {
  if (inputDiv.value) {
    inputText.value = inputDiv.value.textContent
  }
}

function handlePaste(event) {
  const hasImage = Array.from(event.clipboardData.items).some(item =>
    item.type.indexOf('image') !== -1
  );

  if (hasImage) {
    event.preventDefault()
    const items = event.clipboardData.items
    for (let i = 0; i < items.length; i++) {
      if (items[i].type.indexOf('image') !== -1) {
        const file = items[i].getAsFile()
        if (file) {
          const reader = new FileReader()
          reader.onload = (event) => {
            const imgSrc = event.target.result
            const imageMessage = `[img:${imgSrc}]`
            messages.value.push({ id: generateId(), type: 'user', content: imageMessage, timestamp: new Date().toISOString() })
            scrollToBottom()
            if (inputDiv.value) {
              inputDiv.value.innerHTML = ''
              inputText.value = ''
            }
          }
          reader.readAsDataURL(file)
        }
      }
    }
  } else {
    setTimeout(() => {
      if (inputDiv.value) {
        inputText.value = inputDiv.value.innerText || ''
      }
    }, 0)
  }
}

function toggleEmojiPicker() {
  showEmojiPicker.value = !showEmojiPicker.value
}

function selectEmoji(emoji) {
  inputText.value += emoji
  if (inputDiv.value) {
    inputDiv.value.focus()
    const range = window.getSelection().getRangeAt(0)
    range.deleteContents()
    range.insertNode(document.createTextNode(emoji))
    range.collapse(false)
    window.getSelection().removeAllRanges()
    window.getSelection().addRange(range)
  }
  showEmojiPicker.value = false
}

function selectImage() {
  const input = document.createElement('input')
  input.type = 'file'
  input.accept = 'image/*'
  input.onchange = (e) => {
    const file = e.target.files[0]
    if (file) {
      const reader = new FileReader()
      reader.onload = (event) => {
        const imgSrc = event.target.result
        const imageMessage = `[img:${imgSrc}]`
        messages.value.push({ id: generateId(), type: 'user', content: imageMessage, timestamp: new Date().toISOString() })
        scrollToBottom()
      }
      reader.readAsDataURL(file)
    }
  }
  input.click()
}

async function toggleVoiceRecording() {
  if (recording.value) {
    mediaRecorder.stop()
    recording.value = false
  } else {
    try {
      const stream = await navigator.mediaDevices.getUserMedia({ audio: true })
      mediaRecorder = new MediaRecorder(stream)
      audioChunks.value = []

      mediaRecorder.ondataavailable = (e) => e.data.size > 0 && audioChunks.value.push(e.data)
      mediaRecorder.onstop = () => {
        const audioBlob = new Blob(audioChunks.value, { type: 'audio/wav' })
        const audioUrl = URL.createObjectURL(audioBlob)

        const audio = new Audio(audioUrl)
        audio.onloadedmetadata = () => {
          const duration = audio.duration
          messages.value.push({ id: generateId(), type: 'user', content: '[voice-message]', audioUrl, duration, timestamp: new Date().toISOString() })
        }

        stream.getTracks().forEach(t => t.stop())
        scrollToBottom()
      }

      mediaRecorder.start()
      recording.value = true
    } catch (err) {
      console.error('录音失败:', err)
      alert('无法访问麦克风，请检查权限设置')
    }
  }
}

function togglePlayAudio(msg) {
  const audioEl = Array.from(audioElements.value).find(el => el.src === msg.audioUrl)
  if (!audioEl) return

  if (msg.isPlaying) {
    audioEl.pause()
    msg.isPlaying = false
  } else {
    messages.value.forEach(m => {
      if (m.isPlaying) {
        const el = Array.from(audioElements.value).find(e => e.src === m.audioUrl)
        if (el) el.pause()
        m.isPlaying = false
      }
    })
    audioEl.play()
    msg.isPlaying = true
  }
}

function onAudioEnded(msg) {
  msg.isPlaying = false
}

function formatDuration(duration) {
  if (!duration || isNaN(duration)) return '00:00'
  const minutes = Math.floor(duration / 60)
  const seconds = Math.floor(duration % 60)
  return `${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`
}

function scrollToBottom() {
  setTimeout(() => {
    if (chatMessages.value) {
      chatMessages.value.scrollTop = chatMessages.value.scrollHeight
    }
  }, 0)
}

function handleClickOutside(event) {
  const picker = document.querySelector('.emoji-picker')
  const button = document.querySelector('.emoji-button')
  if (picker && button && !picker.contains(event.target) && !button.contains(event.target)) {
    showEmojiPicker.value = false
  }

  const languageMenu = document.querySelector('.language-menu')
  const languageButton = event.target.closest('.language-button')
  if (languageMenu && !languageMenu.contains(event.target) && !languageButton) {
    showLanguageMenu.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
  scrollToBottom()
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})

watch(
  () => messages.value,
  () => scrollToBottom(),
  { deep: true, flush: 'post' }
)
</script>

<style scoped>
.no-scrollbar {
  scrollbar-width: none;
}

::-webkit-scrollbar {
  width: 4px;
}

::-webkit-scrollbar-track {
  background: transparent;
}

::-webkit-scrollbar-thumb {
  background: #cbd5e1;
  border-radius: 2px;
}

::-webkit-scrollbar-thumb:hover {
  background: #94a3b8;
}

.wave-animation {
  display: flex;
  align-items: center;
  gap: 2px;
  height: 20px;
}

.wave-bar {
  width: 3px;
  background: #ef4444;
  animation: wave 1s ease-in-out infinite;
}

.wave-bar:nth-child(1) {
  animation-delay: 0s;
}

.wave-bar:nth-child(2) {
  animation-delay: 0.1s;
}

.wave-bar:nth-child(3) {
  animation-delay: 0.2s;
}

.wave-bar:nth-child(4) {
  animation-delay: 0.3s;
}

.wave-bar:nth-child(5) {
  animation-delay: 0.4s;
}

@keyframes wave {

  0%,
  100% {
    height: 4px;
  }

  50% {
    height: 16px;
  }
}
</style>
