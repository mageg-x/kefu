import { createApp } from 'vue'
import './style.css'
import KefuWidget from './components/KefuWidget.vue'

export function mountKefuWidget(props) {
  const div = document.createElement('div')
  div.id = 'kefu-widget-root'
  // 设置临时内联样式，确保在CSS加载完成前元素不会显示在错误位置
  div.style.position = 'fixed'
  div.style.bottom = '5px'
  div.style.right = '5px'
  div.style.zIndex = '50'
  div.style.opacity = '0' // 初始隐藏
  div.style.transition = 'opacity 0.3s ease'
  document.body.appendChild(div)

  // 👇 关键：这里实际使用了 KefuWidget，Vite 才会打包它
  const app = createApp(KefuWidget, props)
  app.mount('#kefu-widget-root')
  
  // 应用挂载完成后显示widget，实现平滑过渡
  setTimeout(() => {
    div.style.opacity = '1'
  }, 0)
  
  return app
}