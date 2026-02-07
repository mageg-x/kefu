// API 封装

import axios from 'axios'
import { useStore } from './store'

const API_BASE_URL = 'http://localhost:5300/api/v1'

class ApiService {
  constructor() {
    this.baseURL = API_BASE_URL
    
    // 创建axios实例
    this.api = axios.create({
      baseURL: this.baseURL,
      timeout: 10000,
      headers: {
        'Content-Type': 'application/json'
      }
    })
    
    // 请求拦截器
    this.api.interceptors.request.use(
      config => {
        const store = useStore()
        if (store.token) {
          config.headers.Authorization = `Bearer ${store.token}`
        }
        return config
      },
      error => {
        return Promise.reject(error)
      }
    )
    
    // 响应拦截器
    this.api.interceptors.response.use(
      response => {
        return response.data
      },
      error => {
        const errorData = error.response?.data || {}
        const errorMsg = errorData.msg || error.message || '请求失败'
        return Promise.reject(new Error(errorMsg))
      }
    )
  }

  setToken(token) {
    const store = useStore()
    store.token = token
  }

  removeToken() {
    const store = useStore()
    store.clearUser()
  }

  // SHA256 hash密码
  async hashPassword(password) {
    const encoder = new TextEncoder()
    const data = encoder.encode(password)
    const hashBuffer = await crypto.subtle.digest('SHA-256', data)
    const hashArray = Array.from(new Uint8Array(hashBuffer))
    const hashHex = hashArray.map(b => b.toString(16).padStart(2, '0')).join('')
    return hashHex
  }

  // 生成签名（使用密码的1次hash作为密钥）
  async generateSignature(username, passwordHash1, passwordHash2, timestamp, nonce) {
    const crypto = window.crypto || window.msCrypto
    const data = username + passwordHash2 + timestamp + nonce
    const secretKey = passwordHash1 // 使用密码的1次hash作为密钥
    
    return new Promise((resolve) => {
      crypto.subtle.importKey(
        'raw',
        new TextEncoder().encode(secretKey),
        { name: 'HMAC', hash: { name: 'SHA-256' } },
        false,
        ['sign']
      ).then(key => {
        crypto.subtle.sign(
          'HMAC',
          key,
          new TextEncoder().encode(data)
        ).then(signature => {
          const hexSignature = Array.from(new Uint8Array(signature))
            .map(b => b.toString(16).padStart(2, '0'))
            .join('')
          resolve(hexSignature)
        })
      })
    })
  }

  // 生成随机nonce
  generateNonce() {
    return Math.random().toString(36).substring(2, 15) + Math.random().toString(36).substring(2, 15)
  }

  // 登录
  async login(username, password) {
    // 对密码进行1次hash
    const passwordHash1 = await this.hashPassword(password)
    // 对密码的1次hash再进行1次hash，得到密码的2次hash
    const passwordHash2 = await this.hashPassword(passwordHash1)
    const timestamp = Date.now().toString()
    const nonce = this.generateNonce()
    // 使用密码的1次hash作为签名密钥
    const signature = await this.generateSignature(username, passwordHash1, passwordHash2, timestamp, nonce)
    
    const data = await this.api.post('/login', {
      username,
      password: passwordHash2, // 传递2次hash后的密码
      timestamp,
      nonce,
      signature
    })
    
    this.setToken(data.data.token)
    return data
  }

  // 获取用户信息
  async getUserInfo() {
    return this.api.get('/user/info')
  }

  // 登出
  async logout() {
    await this.api.post('/logout')
    this.removeToken()
  }
}

export default new ApiService()
