import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { characterApi } from '@/api/services'
import type { Character, CharacterInput } from '@/api/types'

const CURRENT_CHARACTER_KEY = 'currentUserCharId'

export const useCharactersStore = defineStore('characters', () => {
  const items = ref<Character[]>([])
  const loading = ref(false)
  const storage = ref({ quota_bytes: 0, used_bytes: 0 })

  const charItems = computed(() => items.value.filter(item => (item.extra?.type || 'char') === 'char'))

  async function fetchCharacters() {
    loading.value = true
    try {
      const res = await characterApi.listWithStorage()
      items.value = Array.isArray(res.data) ? res.data : []
      if (res.storage) {
        storage.value = res.storage
      }
      reconcileSelections()
    } finally {
      loading.value = false
    }
  }

  function getCharacterById(id: string | number) {
    return items.value.find(item => String(item.id) === String(id)) || null
  }

  function reconcileSelections() {
    const currentCharacterId = localStorage.getItem(CURRENT_CHARACTER_KEY)
    if (currentCharacterId && !getCharacterById(currentCharacterId)) {
      const fallback = charItems.value[0]
      if (fallback) localStorage.setItem(CURRENT_CHARACTER_KEY, String(fallback.id))
      else localStorage.removeItem(CURRENT_CHARACTER_KEY)
    }
  }

  async function createCharacter(data: CharacterInput) {
    await characterApi.create(data)
    await fetchCharacters()
  }

  async function updateCharacter(id: number, data: CharacterInput) {
    await characterApi.update(id, data)
    await fetchCharacters()
  }

  async function deleteCharacter(id: number) {
    await characterApi.delete(id)
    items.value = items.value.filter(item => item.id !== id)
    reconcileSelections()
  }

  return {
    items,
    loading,
    storage,
    charItems,
    fetchCharacters,
    getCharacterById,
    reconcileSelections,
    createCharacter,
    updateCharacter,
    deleteCharacter,
  }
})
