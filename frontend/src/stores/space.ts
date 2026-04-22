import { defineStore } from 'pinia'
import { ref } from 'vue'
import { spaceApi } from '@/api/space'
import type { Space } from '@/types'

export const useSpaceStore = defineStore('space', () => {
  const spaces = ref<Space[]>([])
  const currentSpace = ref<Space | null>(null)

  const fetchSpaces = async () => {
    try {
      const { data } = await spaceApi.getList()
      // 后端直接返回数组
      spaces.value = Array.isArray(data) ? data : []
    } catch (error) {
      spaces.value = []
    }
  }

  const setCurrentSpace = (space: Space | null) => {
    currentSpace.value = space
  }

  return { spaces, currentSpace, fetchSpaces, setCurrentSpace }
})
