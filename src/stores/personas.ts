import { defineStore } from 'pinia'
import { ref } from 'vue'
import { personaApi } from '@/api/services'
import type { Persona, PersonaInput } from '@/api/types'

export const usePersonasStore = defineStore('personas', () => {
  const personas = ref<Persona[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function fetchPersonas() {
    loading.value = true
    error.value = null
    try {
      const res = await personaApi.list()
      personas.value = res.data || []
    } catch (e: any) {
      error.value = e.message || '加载人设失败'
    } finally {
      loading.value = false
    }
  }

  async function createPersona(data: PersonaInput) {
    const res = await personaApi.create(data)
    await fetchPersonas()
    return res
  }

  async function updatePersona(id: number, data: PersonaInput) {
    await personaApi.update(id, data)
    await fetchPersonas()
  }

  async function deletePersona(id: number) {
    await personaApi.delete(id)
    personas.value = personas.value.filter(p => p.id !== id)
  }

  const defaultPersona = () => personas.value.find(p => p.is_default)

  return {
    personas,
    loading,
    error,
    fetchPersonas,
    createPersona,
    updatePersona,
    deletePersona,
    defaultPersona,
  }
})
