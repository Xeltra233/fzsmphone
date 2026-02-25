import { defineStore } from 'pinia'
import { ref } from 'vue'
import { weiboApi } from '@/api/services'
import type { WeiboPost } from '@/api/types'

export const useWeiboStore = defineStore('weibo', () => {
  const posts = ref<WeiboPost[]>([])
  const loading = ref(false)
  const currentTab = ref('latest')

  async function fetchPosts(tab?: string) {
    loading.value = true
    if (tab) currentTab.value = tab
    try {
      const res = await weiboApi.list(currentTab.value)
      posts.value = res.data || []
    } catch (e) {
      console.error('获取微博失败:', e)
    } finally {
      loading.value = false
    }
  }

  async function createPost(content: string, images?: string[]) {
    try {
      await weiboApi.create({ content, images })
      await fetchPosts()
    } catch (e) {
      console.error('发布微博失败:', e)
    }
  }

  async function deletePost(id: number) {
    try {
      await weiboApi.delete(id)
      posts.value = posts.value.filter(p => p.id !== id)
    } catch (e) {
      console.error('删除微博失败:', e)
    }
  }

  return {
    posts,
    loading,
    currentTab,
    fetchPosts,
    createPost,
    deletePost,
  }
})
