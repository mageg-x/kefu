import { mountKefuWidget } from './main.js'
// 👇 动态加载 CSS（关键！）
function loadCSS() {
  return new Promise((resolve, reject) => {
    const existingLink = document.getElementById('kefu-sdk-styles');
    if (existingLink) {
      resolve();
      return;
    }

    const link = document.createElement('link');
    link.id = 'kefu-sdk-styles';
    link.rel = 'stylesheet';
    // 注意：这里假设 kefu.min.js 和 web-sdk.css 在同一个目录
    link.href = new URL('./web-sdk.css', import.meta.url).href;
    
    // 监听样式加载完成事件
    link.onload = () => resolve();
    link.onerror = () => {
      console.error('Failed to load KefuChat stylesheet');
      reject(new Error('Failed to load KefuChat stylesheet'));
    };
    
    document.head.appendChild(link);
  });
}

window.KefuChat = {
  init: async (options = {}) => {
    console.log('KefuChat initializing...');
    try {
      // 等待CSS加载完成后再初始化widget
      await loadCSS();
      console.log('KefuChat initialized with:', options);
      return mountKefuWidget(options);
    } catch (error) {
      console.error('KefuChat initialization failed:', error);
      // 样式加载失败，不显示聊天窗口
      return null;
    }
  }
};