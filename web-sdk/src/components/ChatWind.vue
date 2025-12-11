<template>
  <div
    class="chat-window w-80 h-[500px] bg-white rounded-xl shadow-xl flex flex-col overflow-hidden border border-gray-200">
    <!-- 顶部标题栏 -->
    <div class="chat-header bg-gradient-to-r from-blue-500 to-blue-600 px-4 py-2 flex items-center justify-between">
      <div class="flex items-center gap-3">
        <div class="w-10 h-10 rounded-full bg-white flex items-center justify-center">
          <img src="../assets/logo.png" alt="logo" class="w-10 h-10 object-contain">
        </div>
        <div class="text-white">
          <h3 class="font-semibold text-base">唯一客服</h3>
          <div class="flex items-center gap-1 text-xs opacity-90 mt-0.5">
            <span class="w-2 h-2 rounded-full bg-green-400 animate-pulse"></span>
            <span>在线</span>
          </div>
        </div>
      </div>
      <div class="flex items-center gap-2">
        <button @click="$emit('close')" class="text-white hover:bg-blue-700 p-2 rounded-full transition">
          <svg class="w-6 h-6" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="none"
            viewBox="0 0 24 24">
            <path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 12h14" />
          </svg>
        </button>
      </div>
    </div>

    <!-- 消息区域 -->
    <div ref="chatMessages" class="chat-messages flex-1 p-4 bg-gray-50 overflow-y-auto space-y-4">
      <!-- 渲染所有消息 -->
      <div class="space-y-4">
        <div v-for="msg in messages" :key="msg.id">
          <!-- 时间戳消息 -->
          <div v-if="msg.type === 'system' && typeof msg.content === 'string' && msg.content.includes('上午')" class="text-center text-xs text-gray-500 mb-4">
            <span class="px-3 py-1 bg-gray-200 rounded-full">{{ msg.content }}</span>
          </div>

          <!-- 系统欢迎消息 -->
          <div v-else-if="msg.type === 'system'" class="flex justify-center">
            <div class="text-xs text-gray-500 bg-gray-100 px-4 py-2 rounded-lg max-w-xs text-center">
              {{ typeof msg.content === 'string' ? msg.content : '' }}
            </div>
          </div>

          <!-- 机器人消息 -->
          <div v-else-if="msg.type === 'robot'" class="flex gap-3">
            <div class="w-8 h-8 rounded-full bg-gray-200 flex-shrink-0 flex items-center justify-center">
              <span class="text-xs font-bold">AI</span>
            </div>
            <div class="flex flex-col gap-2">
              <span class="text-xs text-gray-500">{{ msg.sender }}</span>
              <div class="bg-white px-4 py-1.5 rounded-lg rounded-tl-none shadow-sm border border-gray-100 max-w-xs">
                <div class="text-sm" v-html="renderMarkdown(msg.content)"></div>
              </div>
            </div>
          </div>

          <!-- 客服消息 -->
          <div v-else-if="msg.type === 'service'" class="flex gap-3">
            <div class="w-8 h-8 rounded-full bg-blue-100 flex-shrink-0 flex items-center justify-center">
              <img src="../assets/logo.png" alt="客服" class="w-6 h-6 object-contain">
            </div>
            <div class="flex flex-col gap-2">
              <span class="text-xs text-gray-500">{{ msg.sender }}</span>
              <div class="bg-white px-4 py-1.5 rounded-lg rounded-tl-none shadow-sm border border-gray-100 max-w-xs">
                <!-- 检查消息内容是否包含图片 -->
                <div v-if="typeof msg.content === 'string' && msg.content.includes('[img:')" class="space-y-2">
                  <!-- 渲染图片 -->
                  <template v-for="(imgData, index) in extractImages(msg.content)" :key="'img-' + msg.id + '-' + index">
                    <img :src="imgData" :alt="'图片' + index" class="max-w-[200px] max-h-[200px] rounded cursor-pointer hover:opacity-90 transition-opacity" />
                  </template>
                </div>
                <!-- 检查消息内容是否包含语音消息 -->
                <div v-else-if="typeof msg.content === 'string' && msg.content.includes('[语音消息]')" class="flex items-center gap-2 text-sm">
                  <svg t="1765465928183" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="7054" xmlns:xlink="http://www.w3.org/1999/xlink" width="24" height="24"><path d="M487.648 240a16 16 0 0 1 16-16h16a16 16 0 0 1 16 16v546.784a16 16 0 0 1-16 16h-16a16 16 0 0 1-16-16V240z m155.84 89.04a16 16 0 0 1 16-16h16a16 16 0 0 1 16 16v346.432a16 16 0 0 1-16 16h-16a16 16 0 0 1-16-16V329.04z m155.824 144.704a16 16 0 0 1 16-16h16a16 16 0 0 1 16 16v123.824a16 16 0 0 1-16 16h-16a16 16 0 0 1-16-16v-123.84z m-467.488-144.704a16 16 0 0 1 16-16h16a16 16 0 0 1 16 16v346.432a16 16 0 0 1-16 16h-16a16 16 0 0 1-16-16V329.04zM176 473.76a16 16 0 0 1 16-16h16a16 16 0 0 1 16 16v112.688a16 16 0 0 1-16 16h-16a16 16 0 0 1-16-16V473.76z" fill="currentColor" p-id="7055"></path></svg>
                  {{ msg.content.replace('[语音消息]', '') || '语音消息' }}
                </div>
                <!-- 渲染Markdown内容 -->
                <div v-else class="text-sm" v-html="renderMarkdown(msg.content)"></div>
              </div>
            </div>
          </div>

          <!-- 用户消息 -->
          <div v-else-if="msg.type === 'user'" class="flex justify-end gap-3">
            <div class="flex flex-col max-w-3xs max-h-52 overflow-auto no-scrollbar gap-2 items-end">
              <div class="bg-blue-100 text-black px-4 py-1.5 rounded-lg rounded-tr-none shadow-sm  max-w-full">
                <!-- 检查消息内容是否包含图片 -->
                <div v-if="msg.content.includes('[img:')" class="space-y-2">
                  <!-- 渲染图片 -->
                  <template v-for="(imgData, index) in extractImages(msg.content)" :key="'img-' + msg.id + '-' + index">
                    <img :src="imgData" :alt="'图片' + index" class="max-w-[100px] max-h-[100px] rounded cursor-pointer hover:opacity-90 transition-opacity" />
                  </template>
                </div>
                <!-- 检查消息内容是否包含语音消息 -->
                <div v-else-if="msg.content.includes('[语音消息]')" class="flex items-center gap-2 text-sm">
                  <button 
                    @click="togglePlayAudio(msg)" 
                    class="flex items-center justify-center text-blue-600 hover:text-blue-800"
                  >
                    <svg 
                      t="1765465192033" 
                      class="icon" 
                      viewBox="0 0 1024 1024" 
                      version="1.1" 
                      xmlns="http://www.w3.org/2000/svg" 
                      p-id="6779" 
                      xmlns:xlink="http://www.w3.org/1999/xlink" 
                      width="24" 
                      height="24"
                    >
                      <path 
                        :d="msg.isPlaying ? 'M512 80c-238.656 0-432 193.344-432 432s193.344 432 432 432 432-193.344 432-432S750.656 80 512 80zM576 680c0 17.664-14.336 32-32 32H400c-17.664 0-32-14.336-32-32V344c0-17.664 14.336-32 32-32h144c17.664 0 32 14.336 32 32v336z' : 'M374.464 240.96l459.824 272.96a40 40 0 0 1 0 70.08L374.464 783.04a40 40 0 0 1-60.544-35.04V276c0-24.32 20.48-40.96 40.32-35.04zM352 400.96a16 16 0 0 1 16-16h16a16 16 0 0 1 16 16v224a16 16 0 0 1-16 16h-16a16 16 0 0 1-16-16v-224z'" 
                        fill="currentColor" 
                        p-id="6780"
                      ></path>
                    </svg>
                  </button>
                  <span>{{ typeof msg.content === 'string' ? (msg.content.replace('[语音消息]', '') || '语音消息') : '语音消息' }}</span>
                  <audio 
                    ref="audioElements" 
                    :src="msg.audioUrl" 
                    @ended="onAudioEnded(msg)"
                    class="hidden"
                  ></audio>
                </div>
                <!-- 渲染Markdown内容 -->
                <div v-else class="text-sm break-all" v-html="renderMarkdown(msg.content)"></div>
              </div>
            </div>
            <img 
              :src="`/src/assets/avatars/head${avatarNumber}.svg`" 
              class="w-8 h-8 rounded-full object-cover flex-shrink-0"
              alt="用户头像"
            >
          </div>
        </div>
      </div>


    </div>

    <!-- 快捷回复按钮 -->
    <div class="quick-replies px-2 py-1 bg-gray-50 border-y border-gray-200">
      <div class="flex gap-2 overflow-x-auto ">
        <button
          class="px-4 py-1.5 bg-blue-50 text-xs text-blue-700 rounded-full shadow-sm hover:bg-blue-100 transition whitespace-nowrap border border-blue-200">
          功能介绍
        </button>
        <button
          class="px-4 py-1.5 bg-green-50 text-xs text-green-700 rounded-full shadow-sm hover:bg-green-100 transition whitespace-nowrap border border-green-200">
          人工服务
        </button>
        <button
          class="px-4 py-1.5 bg-purple-50 text-xs text-purple-700 rounded-full shadow-sm hover:bg-purple-100 transition whitespace-nowrap border border-purple-200">
          AI客服
        </button>
      </div>
    </div>

    <!-- 输入区域 -->
    <div class="chat-input bg-white">
      <div class="flex flex-col  py-1">
        <!-- 功能按钮 -->
        <div class="flex items-center gap-2 pb-1 pl-2 relative">
          <!-- 输入表情 -->
          <button class="text-gray-500 hover:text-blue-500 transition p-2 rounded-full hover:bg-gray-100 emoji-button" @click="toggleEmojiPicker">
            <svg class="w-6 h-6" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="24" height="24"
              fill="none" viewBox="0 0 24 24">
              <path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M15 9h.01M8.99 9H9m12 3a9 9 0 1 1-18 0 9 9 0 0 1 18 0ZM6.6 13a5.5 5.5 0 0 0 10.81 0H6.6Z" />
            </svg>
          </button>
          <!-- 表情选择器 -->
          <div v-if="showEmojiPicker" class="emoji-picker absolute bottom-full left-0 mb-2 bg-white border rounded-lg shadow-lg p-2 z-10 w-64 h-48 overflow-y-auto">
            <div class="grid grid-cols-8 gap-1">
              <button 
                v-for="emoji in emojis" 
                :key="emoji"
                class="text-xl hover:bg-gray-100 rounded p-1"
                @click="selectEmoji(emoji)"
              >
                {{ emoji }}
              </button>
            </div>
          </div>
          <!-- 输入图片 -->
          <button class="text-gray-500 hover:text-blue-500 transition p-2 rounded-full hover:bg-gray-100" @click="selectImage">
            <svg class="w-6 h-6" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="24" height="24"
              fill="currentColor" viewBox="0 0 24 24">
              <path fill-rule="evenodd" d="M13 10a1 1 0 0 1 1-1h.01a1 1 0 1 1 0 2H14a1 1 0 0 1-1-1Z"
                clip-rule="evenodd" />
              <path fill-rule="evenodd"
                d="M2 6a2 2 0 0 1 2-2h16a2 2 0 0 1 2 2v12c0 .556-.227 1.06-.593 1.422A.999.999 0 0 1 20.5 20H4a2.002 2.002 0 0 1-2-2V6Zm6.892 12 3.833-5.356-3.99-4.322a1 1 0 0 0-1.549.097L4 12.879V6h16v9.95l-3.257-3.619a1 1 0 0 0-1.557.088L11.2 18H8.892Z"
                clip-rule="evenodd" />
            </svg>
          </button>
          <!-- 输入语音 -->
          <button class="text-gray-500 hover:text-blue-500 transition p-2 rounded-full hover:bg-gray-100" @click="toggleVoiceRecording" :class="{ 'text-red-500': recording }">
            <svg t="1765452216704" class="icon w-6 h-6" viewBox="0 0 1024 1024" version="1.1"
              xmlns="http://www.w3.org/2000/svg" p-id="4919" xmlns:xlink="http://www.w3.org/1999/xlink"
              fill="currentColor">
              <path
                d="M512 705.728c105.888 0 192-86.112 192-192L704 257.952c0-105.888-86.112-192-192-192s-192 86.112-192 192l0 255.776C320 619.584 406.112 705.728 512 705.728z"
                fill="currentColor" p-id="4920"></path>
              <path
                d="M864 479.776 864 352c0-17.664-14.304-32-32-32s-32 14.336-32 32l0 127.776c0 160.16-129.184 290.464-288 290.464-158.784 0-288-130.304-288-290.464L224 352c0-17.664-14.336-32-32-32s-32 14.336-32 32l0 127.776c0 184.608 140.864 336.48 320 352.832L480 896 288 896c-17.664 0-32 14.304-32 32s14.336 32 32 32l448 0c17.696 0 32-14.304 32-32s-14.304-32-32-32l-192 0 0-63.36C723.136 816.256 864 664.384 864 479.776z"
                fill="currentColor" p-id="4921"></path>
            </svg>
          </button>
          <!-- 界面多语言切换 -->
          <button class="text-gray-500 hover:text-blue-500 transition p-2 rounded-full hover:bg-gray-100">
            <svg class="w-6 h-6" t="1765451632654" viewBox="0 0 1024 1024" version="1.1"
              xmlns="http://www.w3.org/2000/svg" p-id="3101" xmlns:xlink="http://www.w3.org/1999/xlink"
              fill="currentColor">
              <path
                d="M848.805886 805.572222c70.998007-81.260745 109.779266-184.217628 109.779266-293.14448 0-119.204939-46.421262-231.277434-130.713041-315.569212C744.876861 113.862257 634.94103 67.61598 517.788843 66.213028c-1.924839-0.599657-10.290367-0.592494-12.227486 0.01535C388.878868 67.945485 279.434224 114.159016 196.73471 196.85853 113.863281 279.730982 67.630307 389.460106 66.095347 506.415818c-0.428765 1.64957-0.436952 8.601912-0.021489 10.226922 1.082658 117.628024 47.364751 228.058113 130.660852 311.354214 84.291778 84.291778 196.36325 130.713041 315.569212 130.713041 119.204939 0 231.277434-46.421262 315.569212-130.713041 6.139837-6.139837 12.054547-12.444427 17.789155-18.871813 0.50756-0.453325 1.001817-0.928139 1.471514-1.440815C847.750857 807.012014 848.295256 806.299793 848.805886 805.572222zM107.447151 532.043499l187.501418 0c1.322112 65.678862 9.253758 127.264499 22.505573 182.112688-61.690014 16.687054-100.819197 38.371936-121.076566 51.906184C144.30971 701.336206 111.676475 620.35687 107.447151 532.043499zM195.881272 259.408121c20.090571 13.556761 59.242266 35.461653 121.340579 52.260248-12.998035 54.127781-20.827351 114.778116-22.243607 179.432649L107.525945 491.101018C112.076588 403.731134 144.437623 323.612399 195.881272 259.408121zM917.081898 491.099994 729.628576 491.099994c-1.415232-64.630996-9.240455-125.260865-22.229281-179.37432 61.95505-16.693194 101.235682-38.444591 121.56673-52.020794C880.270505 323.860039 912.537396 403.866211 917.081898 491.099994zM688.677908 491.099994 532.167319 491.099994 532.167319 335.061149c52.209082-1.094938 97.103572-6.453992 135.272893-14.033621C680.000272 373.163955 687.286212 430.896844 688.677908 491.099994zM532.167319 294.115598 532.167319 109.918435c36.84107 10.398838 72.779583 49.205679 100.926644 110.015649 8.810666 19.035542 16.645099 39.641859 23.464411 61.521169C621.531626 288.227494 580.261687 293.062616 532.167319 294.115598zM491.223814 110.273523l0 183.805236c-47.504944-1.12666-88.378863-6.001691-123.120109-12.802584 6.807033-21.812795 14.623046-42.35976 23.409153-61.344137C419.351903 159.792333 454.809463 121.175827 491.223814 110.273523zM491.223814 335.040682l0 156.059312L335.928912 491.099994c1.391696-60.213383 8.679683-117.955482 21.243837-170.099073C395.008472 328.536548 439.487499 333.887416 491.223814 335.040682zM335.893096 532.043499l155.330718 0 0 158.667719c-51.609425 1.194198-96.019891 6.563486-133.821845 14.103206C344.576873 651.927913 337.193719 593.243349 335.893096 532.043499zM491.223814 731.672118l0 182.909843c-36.415374-10.902304-71.871911-49.51881-99.709933-109.659539-8.679683-18.752086-16.409738-39.034015-23.157419-60.551074C402.9964 737.645157 443.773106 732.820268 491.223814 731.672118zM532.167319 914.937049 532.167319 731.608673c47.904033 1.025353 89.103364 5.862521 124.116809 12.656251-6.755868 21.555945-14.497179 41.87369-23.190165 60.656475C604.946902 865.73137 569.008388 904.538211 532.167319 914.937049zM532.167319 690.660052 532.167319 532.043499l156.546406 0c-1.298576 61.096497-8.66024 119.68487-21.445428 172.502819C629.154233 697.013761 584.319096 691.710988 532.167319 690.660052zM729.659275 532.043499l187.501418 0c-4.221138 88.138386-36.732599 168.973436-88.620363 233.635131-20.469194-13.668301-59.635215-35.298947-121.30374-51.868321C720.43724 659.049101 728.33921 597.585237 729.659275 532.043499zM801.518906 228.742704c-18.329461 11.570523-52.309366 29.355585-104.858186 43.493583-19.295462-63.056128-46.110177-115.004267-78.06189-150.97655C689.00025 140.410913 751.833297 178.097234 801.518906 228.742704zM406.007991 121.259738c-31.905664 35.920094-58.690704 87.768973-77.979002 150.702304-52.40351-14.241352-86.370113-32.099069-104.581893-43.587728C273.076422 177.914062 335.777463 140.364865 406.007991 121.259738zM223.917816 796.963147c18.284435-11.535731 52.098565-29.230742 104.332207-43.335994 19.271926 62.60485 45.976124 114.186645 77.757968 149.968593C335.99952 884.550994 273.472442 847.181899 223.917816 796.963147zM618.59883 903.595746c31.801287-35.803437 58.517765-87.426165 77.792761-150.08218 51.984978 14.023388 85.972047 31.631418 104.533798 43.208081C751.3329 847.061149 688.718841 884.521319 618.59883 903.595746z"
                p-id="3102"></path>
            </svg>
          </button>

        </div>

        <!-- 多行输入框和发送按钮组合 -->
        <div class="relative  mx-2  w-[calc(100%-16px)]">
          <div
            ref="inputDiv"
            contenteditable="true"
            @input="handleInput"
            @keypress.enter.exact="sendMessage"
            @keypress.enter.shift=""
            placeholder="Please enter your message"
            class="w-full border border-gray-300 rounded-lg px-4 py-3 pr-16 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-y min-h-[80px] max-h-[150px] overflow-y-auto"
            style="white-space: pre-wrap; word-break: break-word;">
          </div>

          <button @click="sendMessage"
            class="absolute right-2 bottom-2 bg-blue-500 text-white px-2 py-1 mb-1 rounded-lg text-sm hover:bg-blue-600 transition shadow-sm disabled:bg-gray-300 disabled:cursor-not-allowed"
            :disabled="!inputText.trim()">
            发送
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, watch } from 'vue'
// 引入marked库
import { marked } from 'marked'

// 定义props
const props = defineProps({
  messages: {
    type: Array,
    default: () => []
  },
  userId: {
    type: String,
    default: ''
  }
})

// 定义emits
const emit = defineEmits(['close', 'send-message'])

// CRC计算函数，将字符串转换为01-99的数字
function calculateCRC(str) {
  let crc = 0
  for (let i = 0; i < str.length; i++) {
    crc = (crc * 31 + str.charCodeAt(i)) % 99 + 1 // 1-99
  }
  return crc.toString().padStart(2, '0') // 格式化为两位数，补前导零
}

function renderMarkdown(raw) {
  let text = raw == null ? '' : String(raw)
  if (!text.trim()) return ''

  text = text.replace(/\r\n/g, '\n').replace(/\r/g, '\n').replace(/\\n/g, '\n')

  const renderer = new marked.Renderer()

  // ✅ heading 接收 token 对象
  renderer.heading = (token) => {
    const level = token.depth
    const txt = token.text
    const cls = {1:'text-2xl',2:'text-xl',3:'text-lg',4:'text-base',5:'text-sm',6:'text-xs'}[level] || 'text-sm'
    return `<h${level} class="${cls} font-bold my-2">${txt}</h${level}>`
  }

  // ✅ listitem 也接收 token 对象
  renderer.listitem = (token) => {
    // 注意：token.text 可能不存在，需要用 marked.parseInline 或递归处理 tokens
    // 更安全的方式：让 marked 自己渲染子内容
    const content = marked.parser(token.tokens, { async: false, renderer })
    return `<li class="ml-6">${content}</li>`
  }

  renderer.code = (token) => {
    const code = token.text
    return code.includes('\n')
      ? `<pre class="bg-gray-100 p-4 rounded text-sm overflow-x-auto"><code>${code}</code></pre>`
      : `<code class="bg-gray-100 px-1 py-0.5 rounded text-sm">${code}</code>`
  }

  renderer.link = (token) => {
    const href = token.href
    const txt = token.text || href
    return `<a href="${href}" target="_blank" rel="noopener noreferrer" class="text-blue-600 hover:underline">${txt}</a>`
  }

  try {
    const result = marked.parse(text, { gfm: true, breaks: true, renderer })
    return typeof result === 'string' ? result : ''
  } catch (e) {
    console.error('Markdown error:', e)
    return `<div class="text-red-500">解析错误</div>`
  }
}

// 计算当前用户的头像编号
const avatarNumber = calculateCRC(props.userId)

// 音频元素引用
const audioElements = ref([])

// 生成唯一ID的函数
function generateId() {
  // 结合用户ID、时间戳和随机数，确保不同用户之间ID也不会重复
  return props.userId + '_' + Date.now().toString(36) + Math.random().toString(36).substr(2, 8)
}

// 响应式数据
const inputText = ref('')
const inputDiv = ref(null)
const chatMessages = ref(null)
const showEmojiPicker = ref(false)
const showVoiceRecorder = ref(false)
const recording = ref(false)
const audioChunks = ref([])
let mediaRecorder = null

// 模拟消息数据
const messages = ref([
  {
    id: '1',
    type: 'system',
    content: '上午 10:30',
    timestamp: new Date().toISOString()
  },
  {
    id: '2',
    type: 'system',
    content: '唯一客服 正在为您服务!',
    timestamp: new Date().toISOString()
  },
  {
    id: '3',
    type: 'robot',
    sender: '扣子智能体助手',
    content: '您好！很高兴为您服务。请问有什么可以帮您？',
    timestamp: new Date().toISOString()
  },
  {
    id: '4',
    type: 'service',
    sender: '唯一客服',
    content: '好的，已收到您的♥号，稍等，我们一对一联系您，给您介绍。',
    timestamp: new Date().toISOString()
  },
  {
    id: '5',
    type: 'service',
    sender: '唯一客服',
    content: '没意向就没有聊的必要',
    timestamp: new Date().toISOString()
  },
  // 用户消息将通过props传入
  ...props.messages
])

// 表情列表
const emojis = [
  '👌', '💪', '✌️', '👍', '👎', '🌸', '💐', '🌹',
  '❤️', '💛', '💙', '💜', '🤝', '👏', '🎉', '🎊',
  '🙏', '😢', '😤', '😡', '😱', '😳', '🤯', '🤔',
  '😀', '😃', '😄', '😁', '😆', '😅', '🤣', '😂',
  '🙂', '🙃', '😉', '😊', '😇', '🥰', '😍', '🤩',
  '😘', '😗', '😚', '😙', '🥲', '😋', '😛', '😜',
  '🤪', '😝', '🤑', '🤗', '🤭', '🤫', '🤔', '🤐',
  '🤨', '😐', '😑', '😶', '😏', '😒', '🙄', '😬',
  '😮‍💨', '🤥', '😌', '😔', '😪', '🤤', '😴', '😷',
  '🤒', '🤕', '🤢', '🤮', '🤧', '🥵', '🥶', '🥴',
  '😵', '🤯', '🤠', '🥳', '😎', '🤓', '🧐', '😕'
]

// 方法
// 提取消息中的图片
function extractImages(content) {
  // 检查输入类型，如果不是字符串或为空，直接返回空数组
  if (!content || typeof content !== 'string') return []
  
  const imgRegex = /\[img:([^\]]+)\]/g
  const images = []
  let match
  while ((match = imgRegex.exec(content)) !== null) {
    images.push(match[1])
  }
  return images
}

// 滚动到聊天窗口最底部// 滚动到底部
function scrollToBottom() {
  setTimeout(() => {
    if (chatMessages.value) {
      chatMessages.value.scrollTop = chatMessages.value.scrollHeight
    }
  }, 0)
}

// 播放/暂停音频
function togglePlayAudio(msg) {
  // 找到对应的音频元素
  const audioElement = Array.from(audioElements.value).find(el => el.src === msg.audioUrl)
  if (!audioElement) return
  
  // 如果当前音频正在播放，先暂停
  if (msg.isPlaying) {
    audioElement.pause()
    msg.isPlaying = false
  } else {
    // 暂停所有其他音频
    messages.value.forEach(m => {
      if (m.isPlaying) {
        const otherAudio = Array.from(audioElements.value).find(el => el.src === m.audioUrl)
        if (otherAudio) otherAudio.pause()
        m.isPlaying = false
      }
    })
    
    // 播放当前音频
    audioElement.play()
    msg.isPlaying = true
  }
}

// 音频播放结束处理
function onAudioEnded(msg) {
  msg.isPlaying = false
}

function sendMessage() {
  const content = String(inputText.value).trim()
  if (!content) return
  console.log(`input message : ${content}`)
  messages.value.push({
    id: generateId(),
    type: 'user',
    content: content, // ← 确保是字符串！
    timestamp: new Date().toISOString()
  })

  emit('send-message', content)
  inputText.value = ''
  if (inputDiv.value) {
    inputDiv.value.innerHTML = ''
  }
  scrollToBottom()
}

// 处理输入框内容变化
function handleInput() {
  if (inputDiv.value) {
    // 只处理文本内容，图片不再通过输入框发送
    inputText.value = inputDiv.value.textContent
  }
}

// 表情按钮功能
function toggleEmojiPicker() {
  showEmojiPicker.value = !showEmojiPicker.value
}

function selectEmoji(emoji) {
  // 更新 inputText
  inputText.value += emoji
  
  // 更新 contenteditable 元素的内容
  if (inputDiv.value) {
    // 确保输入框获得焦点
    inputDiv.value.focus()
    
    // 获取当前光标位置
    const selection = window.getSelection()
    
    // 检查是否有有效的选择范围
    if (selection.rangeCount > 0) {
      const range = selection.getRangeAt(0)
      
      // 创建表情文本节点
      const emojiNode = document.createTextNode(emoji)
      
      // 插入表情
      range.insertNode(emojiNode)
      
      // 将光标移动到表情后面
      range.setStartAfter(emojiNode)
      range.collapse(true)
      selection.removeAllRanges()
      selection.addRange(range)
    } else {
      // 如果没有有效的选择范围，直接追加到输入框末尾
      inputDiv.value.innerHTML += emoji
    }
  }
  
  showEmojiPicker.value = false
}

// 图片按钮功能
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
        // 直接发送图片消息
        const imageMessage = `[img:${imgSrc}]`
        
        // 添加新消息到消息列表
        messages.value.push({
          id: generateId(),
          type: 'user',
          content: imageMessage,
          timestamp: new Date().toISOString()
        })
        
        emit('send-message', imageMessage)
        
        // 滚动到底部
        scrollToBottom()
      }
      reader.readAsDataURL(file)
    }
  }
  input.click()
}

// 语音按钮功能
async function toggleVoiceRecording() {
  if (recording.value) {
    // 停止录音
    mediaRecorder.stop()
    recording.value = false
  } else {
    try {
      // 开始录音
      const stream = await navigator.mediaDevices.getUserMedia({ audio: true })
      mediaRecorder = new MediaRecorder(stream)
      audioChunks.value = []
      
      mediaRecorder.ondataavailable = (event) => {
        if (event.data.size > 0) {
          audioChunks.value.push(event.data)
        }
      }
      
      mediaRecorder.onstop = () => {
        const audioBlob = new Blob(audioChunks.value, { type: 'audio/wav' })
        // 创建音频URL
        const audioUrl = URL.createObjectURL(audioBlob)
        // 直接发送语音消息
        const voiceMessage = '[语音消息]'
        
        // 添加新消息到消息列表
        messages.value.push({
          id: generateId(),
          type: 'user',
          content: voiceMessage,
          audioUrl: audioUrl, // 存储音频URL
          timestamp: new Date().toISOString()
        })
        
        emit('send-message', voiceMessage)
        
        // 滚动到底部
        scrollToBottom()
        
        // 停止所有音频轨道
        stream.getTracks().forEach(track => track.stop())
      }
      
      mediaRecorder.start()
      recording.value = true
    } catch (error) {
      console.error('录音失败:', error)
      alert('无法访问麦克风，请检查权限设置')
    }
  }
}

// 点击外部关闭表情选择器
function handleClickOutside(event) {
  const emojiPicker = document.querySelector('.emoji-picker')
  const emojiButton = document.querySelector('.emoji-button')
  
  if (emojiPicker && emojiButton && !emojiPicker.contains(event.target) && !emojiButton.contains(event.target)) {
    showEmojiPicker.value = false
  }
}

// 生命周期钩子
onMounted(() => {
  document.addEventListener('click', handleClickOutside)
  // 组件挂载完成后滚动到底部
  scrollToBottom()
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})

// 监听消息变化，新消息时滚动到底部
watch(
  () => messages.value,
  () => {
    scrollToBottom()
  },
  { deep: true, flush: 'post' } // 使用flush: 'post'确保DOM已经更新
)
</script>

<style scoped>
.no-scrollbar {
    scrollbar-width: none;
}
/* 滚动条样式 */
.chat-messages::-webkit-scrollbar {
  width: 4px;
}

.chat-messages::-webkit-scrollbar-track {
  background: transparent;
}

.chat-messages::-webkit-scrollbar-thumb {
  background: rgba(0, 0, 0, 0.1);
  border-radius: 2px;
}

.chat-messages::-webkit-scrollbar-thumb:hover {
  background: rgba(0, 0, 0, 0.2);
}

/* 平滑滚动 */
.chat-messages {
  scroll-behavior: smooth;
}

/* 输入框聚焦效果 */
input:focus {
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}
</style>