import { ref, computed } from 'vue'
import { PAGE_SIZE } from '@/constants'

/** 分页组合式函数 */
export function usePagination(initialPageSize = PAGE_SIZE) {
  const page = ref(1)
  const pageSize = ref(initialPageSize)
  const total = ref(0)

  const totalPages = computed(() => Math.ceil(total.value / pageSize.value) || 1)

  function setPage(value: number) {
    page.value = value
  }

  function setTotal(value: number) {
    total.value = value
  }

  function reset() {
    page.value = 1
    total.value = 0
  }

  return {
    page,
    pageSize,
    total,
    totalPages,
    setPage,
    setTotal,
    reset,
  }
}
