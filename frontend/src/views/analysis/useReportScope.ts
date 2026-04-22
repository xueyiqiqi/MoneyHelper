import { computed, ref, watch, type Ref } from 'vue'
import type { Space } from '@/types'

export const PERSONAL_SCOPE = 'personal' as const
export type ReportScopeValue = typeof PERSONAL_SCOPE | number

interface CreateReportScopeControllerOptions {
  spaces: Ref<Space[]>
  currentSpace: Ref<Space | null>
  setCurrentSpace: (space: Space | null) => void
}

export function createReportScopeController({
  spaces,
  currentSpace,
  setCurrentSpace,
}: CreateReportScopeControllerOptions) {
  const scope = ref<ReportScopeValue>(PERSONAL_SCOPE)

  const syncFromStore = () => {
    if (currentSpace.value && spaces.value.some((space) => space.id === currentSpace.value?.id)) {
      scope.value = currentSpace.value.id
      return
    }

    if (currentSpace.value) {
      setCurrentSpace(null)
    }

    scope.value = PERSONAL_SCOPE
  }

  const setScope = (value: ReportScopeValue) => {
    if (value === PERSONAL_SCOPE) {
      scope.value = PERSONAL_SCOPE
      setCurrentSpace(null)
      return
    }

    const matchedSpace = spaces.value.find((space) => space.id === value) ?? null
    if (!matchedSpace) {
      scope.value = PERSONAL_SCOPE
      setCurrentSpace(null)
      return
    }

    scope.value = matchedSpace.id
    setCurrentSpace(matchedSpace)
  }

  const targetLabel = computed(() => {
    if (scope.value === PERSONAL_SCOPE) {
      return '个人账本'
    }

    return spaces.value.find((space) => space.id === scope.value)?.name ?? '个人账本'
  })

  const requestParams = computed(() => {
    if (scope.value === PERSONAL_SCOPE) {
      return {
        isPersonal: true,
        spaceId: undefined,
      }
    }

    return {
      isPersonal: false,
      spaceId: scope.value,
    }
  })

  watch([currentSpace, spaces], syncFromStore, { immediate: true, deep: true })

  return {
    scope,
    targetLabel,
    requestParams,
    setScope,
  }
}
