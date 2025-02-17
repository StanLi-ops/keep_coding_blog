<template>
  <div class="home">
    <!-- 文章列表 -->
    <div class="post-list">
      <el-empty v-if="posts.length === 0" description="暂无文章" />
      
      <div v-for="post in posts" :key="post.id" class="post-item">
        <h2 class="post-title">
          <router-link :to="`/post/${post.id}`">{{ post.title }}</router-link>
        </h2>
        
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
        
        <p class="post-content" v-html="truncateContent(post.content)"></p>
        
        <div class="post-footer">
          <router-link :to="`/post/${post.id}`" class="read-more">
            阅读全文
          </router-link>
        </div>
      </div>
    </div>

    <!-- 分页 -->
    <div class="pagination">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :total="total"
        :page-sizes="[10, 20, 30, 50]"
        layout="total, sizes, prev, pager, next"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getPosts } from '../api'
import { ElMessage } from 'element-plus'

// 状态
const posts = ref([])
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

// 获取文章列表
const fetchPosts = async () => {
  try {
    const response = await getPosts({
      page: currentPage.value,
      page_size: pageSize.value
    })
    posts.value = response.data.posts
    total.value = response.data.total
  } catch (error) {
    ElMessage.error('获取文章列表失败')
  }
}

// 分页处理
const handleSizeChange = (val) => {
  pageSize.value = val
  fetchPosts()
}

const handleCurrentChange = (val) => {
  currentPage.value = val
  fetchPosts()
}

// 修改截断函数，处理 HTML 内容
const truncateContent = (content) => {
  // 创建临时 div 去除 HTML 标签
  const div = document.createElement('div')
  div.innerHTML = content
  const text = div.textContent || div.innerText
  return text.length > 200 ? text.slice(0, 200) + '...' : text
}

const formatDate = (date) => {
  return new Date(date).toLocaleDateString()
}

// 生命周期
onMounted(() => {
  fetchPosts()
})
</script>

<style scoped>
.home {
  padding: 20px 0;
}

.post-list {
  margin-bottom: 20px;
}

.post-item {
  padding: 20px;
  margin-bottom: 20px;
  border-bottom: 1px solid #eee;
}

.post-title {
  margin: 0 0 10px 0;
}

.post-title a {
  color: #333;
  text-decoration: none;
  font-size: 1.5rem;
}

.post-title a:hover {
  color: #409EFF;
}

.post-meta {
  font-size: 0.9rem;
  color: #666;
  margin-bottom: 10px;
  display: flex;
  align-items: center;
  gap: 20px;
}

.tags {
  display: inline-flex;
  gap: 8px;
}

.post-content {
  color: #666;
  line-height: 1.6;
  margin: 15px 0;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
}

.post-footer {
  margin-top: 15px;
}

.read-more {
  color: #409EFF;
  text-decoration: none;
}

.read-more:hover {
  text-decoration: underline;
}

.pagination {
  display: flex;
  justify-content: center;
  margin-top: 30px;
}
</style> 