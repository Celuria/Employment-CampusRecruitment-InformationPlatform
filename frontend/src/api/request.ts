import axios, { type AxiosInstance, type AxiosRequestConfig, type AxiosResponse } from 'axios'
import { ElMessage } from 'element-plus'
import { API_PREFIX } from '@/constants'
import { storage } from '@/utils/storage'
import type { ApiResponse } from '@/types'

const request: AxiosInstance = axios.create({
  baseURL: API_PREFIX,
  timeout: 15000,
  headers: {
    'Content-Type': 'application/json',
  },
})

request.interceptors.request.use((config) => {
  const token = storage.getToken()
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

function getErrorMessage(error: unknown): string {
  if (axios.isAxiosError(error)) {
    const data = error.response?.data as ApiResponse | undefined
    if (data?.message) return data.message
  }
  return '网络异常，请稍后重试'
}

function isAuthEndpoint(url?: string): boolean {
  if (!url) return false
  return url.includes('/auth/login') || url.includes('/auth/register')
}

request.interceptors.response.use(
  (response: AxiosResponse<ApiResponse>) => {
    const { code, message, data } = response.data
    if (code === 0) {
      return data as never
    }
    ElMessage.error(message || '请求失败')
    return Promise.reject(new Error(message || '请求失败'))
  },
  (error) => {
    const status = error.response?.status
    const message = getErrorMessage(error)
    const url = error.config?.url as string | undefined

    if (status === 401 && !isAuthEndpoint(url)) {
      storage.removeToken()
      ElMessage.warning('登录已过期，请重新登录')
      window.location.href = `/login?redirect=${encodeURIComponent(window.location.pathname)}`
    } else if (status === 423) {
      ElMessage.error(message || '账号已锁定，请稍后再试')
    } else {
      ElMessage.error(message)
    }
    return Promise.reject(error)
  },
)

export function get<T>(url: string, config?: AxiosRequestConfig): Promise<T> {
  return request.get(url, config)
}

export function post<T>(url: string, data?: unknown, config?: AxiosRequestConfig): Promise<T> {
  return request.post(url, data, config)
}

export function put<T>(url: string, data?: unknown, config?: AxiosRequestConfig): Promise<T> {
  return request.put(url, data, config)
}

export function del<T>(url: string, config?: AxiosRequestConfig): Promise<T> {
  return request.delete(url, config)
}

export function patch<T>(url: string, data?: unknown, config?: AxiosRequestConfig): Promise<T> {
  return request.patch(url, data, config)
}

export default request
