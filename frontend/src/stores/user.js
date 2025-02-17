import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { login as loginApi, register as registerApi } from '../api'

export const useUserStore = defineStore('user', () => {
  // 状态
  const token = ref(localStorage.getItem('token') || '')
  const user = ref(JSON.parse(localStorage.getItem('user') || '{}'))

  // 计算属性
  const isLoggedIn = computed(() => !!token.value)

  // 方法
  const setToken = (newToken) => {
    token.value = newToken
    localStorage.setItem('token', newToken)
  }

  const setUser = (newUser) => {
    user.value = newUser
    localStorage.setItem('user', JSON.stringify(newUser))
  }

  const login = async (username, password) => {
    try {
      const response = await loginApi({ username, password })
      setToken(response.data.token)
      setUser(response.data.user)
      return Promise.resolve(response)
    } catch (error) {
      return Promise.reject(error)
    }
  }

  const register = async (userData) => {
    try {
      const response = await registerApi(userData)
      return Promise.resolve(response)
    } catch (error) {
      return Promise.reject(error)
    }
  }

  const logout = () => {
    token.value = ''
    user.value = {}
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  }

  return {
    token,
    user,
    isLoggedIn,
    login,
    register,
    logout
  }
}) 