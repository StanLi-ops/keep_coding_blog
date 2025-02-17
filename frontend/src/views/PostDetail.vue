<template>
  <div class="post-detail" v-loading="loading">
    <!-- 文章内容 -->
    <div v-if="post" class="post-content">
      <div class="post-header">
        <h1>{{ post.title }}</h1>
        <div class="post-meta">
          <span class="author">作者: {{ post.user.username }}</span>
          <span class="time">发布时间: {{ formatDate(post.created_at) }}</span>
          <div class="tags">
            <el-tag
              v-for="tag in post.tags"
              :key="tag"
              size="small"
              class="tag"
            >
              {{ tag.name }}
            </el-tag>
          </div>
        </div>
      </div>

      <div class="post-body" v-html="post.content"></div>

      <!-- 文章操作按钮 (仅作者可见) -->
      <div v-if="isAuthor" class="post-actions">
        <el-button type="primary" @click="handleEdit">编辑</el-button>
        <el-button type="danger" @click="handleDelete">删除</el-button>
      </div>
    </div>

    <!-- 评论区 -->
    <div class="comments-section">
      <h3>评论 ({{ comments.length }})</h3>
      
      <!-- 发表评论 -->
      <div v-if="userStore.isLoggedIn" class="comment-form">
        <el-input
          v-model="newComment"
          type="textarea"
          :rows="3"
          placeholder="写下你的评论..."
        />
        <el-button 
          type="primary"
          :loading="commentLoading"
          @click="submitComment"
          class="submit-comment"
        >
          发表评论
        </el-button>
      </div>
      <div v-else class="login-tip">
        <router-link to="/login">登录</router-link> 后发表评论
      </div>

      <!-- 评论列表 -->
      <div class="comments-list">
        <el-empty v-if="comments.length === 0" description="暂无评论" />
        
        <div v-for="comment in comments" :key="comment.id" class="comment-item">
          <div class="comment-header">
            <span class="comment-author">{{ comment.author }}</span>
            <span class="comment-time">{{ formatDate(comment.created_at) }}</span>
          </div>

          <!-- 评论内容：编辑模式/显示模式 -->
          <div v-if="editingCommentId === comment.id" class="comment-edit">
            <el-input
              v-model="editingContent"
              type="textarea"
              :rows="3"
              placeholder="编辑评论..."
            />
            <div class="edit-actions">
              <el-button type="primary" size="small" @click="handleUpdateComment(comment.id)">保存</el-button>
              <el-button size="small" @click="cancelEdit">取消</el-button>
            </div>
          </div>
          <div v-else class="comment-content">
            {{ comment.content }}
          </div>

          <!-- 评论操作按钮 (仅评论作者可见) -->
          <div v-if="comment.user_id === userStore.user.id" class="comment-actions">
            <el-button 
              type="text"
              class="edit-comment"
              @click="startEdit(comment)"
            >
              编辑
            </el-button>
            <el-button 
              type="text"
              class="delete-comment"
              @click="handleDeleteComment(comment.id)"
            >
              删除
            </el-button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useUserStore } from '../stores/user'
import { getPostById, deletePost, getComments, createComment, deleteComment as deleteCommentApi, updateComment } from '../api'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

// 状态
const loading = ref(false)
const commentLoading = ref(false)
const post = ref(null)
const comments = ref([])
const newComment = ref('')

// 评论编辑相关
const editingCommentId = ref(null)
const editingContent = ref('')

// 计算属性
const isAuthor = computed(() => {
  return post.value?.user_id === userStore.user.id
})

// 获取文章详情
const fetchPost = async () => {
  loading.value = true
  try {
    const response = await getPostById(route.params.id)
    post.value = response.data.post
  } catch (error) {
    ElMessage.error('获取文章失败')
  } finally {
    loading.value = false
  }
}

// 获取评论列表
const fetchComments = async () => {
  try {
    const response = await getComments({ post_id: route.params.id })
    comments.value = response.data.comments
  } catch (error) {
    ElMessage.error('获取评论失败')
  }
}

// 提交评论
const submitComment = async () => {
  if (!newComment.value.trim()) {
    ElMessage.warning('请输入评论内容')
    return
  }

  commentLoading.value = true
  try {
    await createComment({
      content: newComment.value,
      post_id: Number(route.params.id)
    })
    ElMessage.success('评论成功')
    newComment.value = ''
    await fetchComments()
  } catch (error) {
    ElMessage.error('评论失败')
  } finally {
    commentLoading.value = false
  }
}

// 开始编辑评论
const startEdit = (comment) => {
  editingCommentId.value = comment.id
  editingContent.value = comment.content
}

// 取消编辑
const cancelEdit = () => {
  editingCommentId.value = null
  editingContent.value = ''
}

// 更新评论
const handleUpdateComment = async (commentId) => {
  if (!editingContent.value.trim()) {
    ElMessage.warning('评论内容不能为空')
    return
  }

  try {
    await updateComment(commentId, {
      content: editingContent.value
    })
    ElMessage.success('更新成功')
    await fetchComments()  // 刷新评论列表
    cancelEdit()          // 退出编辑模式
  } catch (error) {
    ElMessage.error('更新失败')
  }
}

// 删除评论
const handleDeleteComment = async (commentId) => {
  try {
    await ElMessageBox.confirm('确定要删除这条评论吗？', '提示', {
      type: 'warning'
    })
    
    await deleteCommentApi(commentId)
    ElMessage.success('删除成功')
    await fetchComments()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

// 编辑文章
const handleEdit = () => {
  router.push(`/edit/${post.value.id}`)
}

// 删除文章
const handleDelete = async () => {
  try {
    await ElMessageBox.confirm('确定要删除这篇文章吗？', '提示', {
      type: 'warning'
    })
    
    await deletePost(post.value.id)
    ElMessage.success('删除成功')
    router.push('/')
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

// 格式化日期
const formatDate = (date) => {
  return new Date(date).toLocaleString()
}

// 生命周期
onMounted(() => {
  fetchPost()
  fetchComments()
})
</script>

<style scoped>
.post-detail {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}

.post-header {
  margin-bottom: 30px;
}

.post-meta {
  color: #666;
  margin: 10px 0;
  display: flex;
  align-items: center;
  gap: 20px;
}

.tags {
  display: inline-flex;
  gap: 8px;
}

.post-body {
  line-height: 1.8;
  font-size: 16px;
  margin-bottom: 40px;
}

/* 添加富文本内容的样式 */
:deep(.post-body) {
  h1, h2, h3, h4, h5, h6 {
    margin: 1em 0 0.5em;
    font-weight: 600;
  }
  
  p {
    margin: 1em 0;
  }
  
  ul, ol {
    padding-left: 2em;
    margin: 1em 0;
  }
  
  pre {
    background-color: #f6f8fa;
    padding: 1em;
    border-radius: 4px;
    overflow-x: auto;
  }
  
  code {
    background-color: #f6f8fa;
    padding: 0.2em 0.4em;
    border-radius: 3px;
  }
  
  blockquote {
    border-left: 4px solid #dfe2e5;
    padding-left: 1em;
    color: #666;
    margin: 1em 0;
  }
  
  img {
    max-width: 100%;
    height: auto;
  }
  
  table {
    border-collapse: collapse;
    width: 100%;
    margin: 1em 0;
  }
  
  th, td {
    border: 1px solid #dfe2e5;
    padding: 0.6em;
  }
}

.post-actions {
  margin: 20px 0;
  display: flex;
  gap: 10px;
}

.comments-section {
  margin-top: 40px;
  border-top: 1px solid #eee;
  padding-top: 20px;
}

.comment-form {
  margin: 20px 0;
}

.submit-comment {
  margin-top: 10px;
}

.login-tip {
  text-align: center;
  color: #666;
  padding: 20px 0;
}

.comment-item {
  padding: 15px 0;
  border-bottom: 1px solid #eee;
  position: relative;
}

.comment-header {
  margin-bottom: 10px;
}

.comment-author {
  font-weight: bold;
  margin-right: 10px;
}

.comment-time {
  color: #999;
  font-size: 0.9em;
}

.comment-content {
  line-height: 1.6;
}

.comment-edit {
  margin: 10px 0;
}

.edit-actions {
  margin-top: 10px;
  display: flex;
  gap: 10px;
}

.comment-actions {
  position: absolute;
  right: 0;
  top: 15px;
  display: flex;
  gap: 10px;
}

.edit-comment {
  color: #409EFF;
}

.edit-comment:hover {
  color: #79bbff;
}

.delete-comment {
  color: #F56C6C;
}

.delete-comment:hover {
  color: #f89898;
}
</style> 