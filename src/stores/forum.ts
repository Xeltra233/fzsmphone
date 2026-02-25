import { defineStore } from 'pinia'
import { ref } from 'vue'
import { postApi } from '@/api/services'
import type { Post, Comment } from '@/api/types'

export const useForumStore = defineStore('forum', () => {
  const posts = ref<Post[]>([])
  const currentPost = ref<Post | null>(null)
  const comments = ref<Comment[]>([])
  const loading = ref(false)

  async function fetchPosts(board?: string, page?: number) {
    loading.value = true
    try {
      const res = await postApi.list({ board, page })
      posts.value = res.data || []
    } catch (e) {
      console.error('获取帖子失败:', e)
    } finally {
      loading.value = false
    }
  }

  async function fetchPost(id: number) {
    loading.value = true
    try {
      currentPost.value = await postApi.get(id)
    } catch (e) {
      console.error('获取帖子详情失败:', e)
    } finally {
      loading.value = false
    }
  }

  async function createPost(data: { board: string; title: string; content: string; images?: string[] }) {
    try {
      const res = await postApi.create(data)
      await fetchPosts(data.board)
      return res.id
    } catch (e) {
      console.error('创建帖子失败:', e)
      return null
    }
  }

  async function deletePost(id: number) {
    try {
      await postApi.delete(id)
      posts.value = posts.value.filter(p => p.id !== id)
    } catch (e) {
      console.error('删除帖子失败:', e)
    }
  }

  async function fetchComments(postId: number) {
    try {
      const res = await postApi.listComments(postId)
      comments.value = res.data || []
    } catch (e) {
      console.error('获取评论失败:', e)
    }
  }

  async function addComment(postId: number, content: string, parentId?: number) {
    try {
      await postApi.createComment(postId, { content, parent_id: parentId })
      await fetchComments(postId)
      // 更新帖子评论数
      const post = posts.value.find(p => p.id === postId)
      if (post) post.comments++
    } catch (e) {
      console.error('发表评论失败:', e)
    }
  }

  return {
    posts,
    currentPost,
    comments,
    loading,
    fetchPosts,
    fetchPost,
    createPost,
    deletePost,
    fetchComments,
    addComment,
  }
})
