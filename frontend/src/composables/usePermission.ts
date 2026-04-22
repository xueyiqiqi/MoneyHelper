import { computed } from 'vue'
import { useSpaceStore } from '@/stores/space'
import type { Role } from '@/types'

export function usePermission() {
  const spaceStore = useSpaceStore()

  const currentRole = computed<Role>(() => {
    return spaceStore.currentSpace?.role || 'member'
  })

  const canManageMembers = computed(() => {
    return ['creator', 'admin'].includes(currentRole.value)
  })

  const canInviteMember = computed(() => {
    return ['creator', 'admin'].includes(currentRole.value)
  })

  const canEditBill = computed(() => {
    return ['creator', 'admin', 'member'].includes(currentRole.value)
  })

  const canDeleteBill = computed(() => {
    return ['creator', 'admin'].includes(currentRole.value)
  })

  const canAddBill = computed(() => {
    return ['creator', 'admin', 'member'].includes(currentRole.value)
  })

  const canSetAdmin = computed(() => {
    return currentRole.value === 'creator'
  })

  const canRemoveMember = (targetRole: Role, targetUserId?: number, currentUserId?: number) => {
    if (targetRole === 'creator') return false
    if (currentUserId && targetUserId === currentUserId) return false

    if (currentRole.value === 'creator') {
      return ['admin', 'member', 'observer'].includes(targetRole)
    }

    if (currentRole.value === 'admin') {
      return ['member', 'observer'].includes(targetRole)
    }

    return false
  }

  const canEditMemberRole = (targetRole: Role, targetUserId?: number, currentUserId?: number) => {
    if (targetRole === 'creator') return false
    if (currentUserId && targetUserId === currentUserId) return false

    if (currentRole.value === 'creator') {
      return ['admin', 'member', 'observer'].includes(targetRole)
    }

    if (currentRole.value === 'admin') {
      return ['member', 'observer'].includes(targetRole)
    }

    return false
  }

  const assignableRoles = computed<Role[]>(() => {
    if (currentRole.value === 'creator') {
      return ['admin', 'member', 'observer']
    }
    if (currentRole.value === 'admin') {
      return ['member', 'observer']
    }
    return []
  })

  return {
    currentRole,
    canManageMembers,
    canInviteMember,
    canEditBill,
    canDeleteBill,
    canAddBill,
    canSetAdmin,
    canRemoveMember,
    canEditMemberRole,
    assignableRoles,
  }
}
