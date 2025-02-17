<template>
  <div class="create-post">
    <h2>{{ isEdit ? '编辑文章' : '创建文章' }}</h2>
    
    <el-form
      ref="formRef"
      :model="postForm"
      :rules="rules"
      label-width="80px"
      class="post-form"
    >
      <!-- 标题 -->
      <el-form-item label="标题" prop="title">
        <el-input 
          v-model="postForm.title"
          placeholder="请输入文章标题"
        />
      </el-form-item>
      
      <!-- 标签 -->
      <el-form-item label="标签" prop="tags">
        <el-select
          v-model="postForm.tags"
          multiple
          filterable
          allow-create
          default-first-option
          placeholder="请选择或输入标签"
          class="tag-select"
        >
          <el-option
            v-for="tag in availableTags"
            :key="tag.id"
            :label="tag.name"
            :value="tag.name"
          />
        </el-select>
      </el-form-item>
      
      <!-- 内容（富文本编辑器） -->
      <el-form-item label="内容" prop="content">
        <div class="editor-container">
          <Toolbar
            :editor="editorRef"
            :defaultConfig="toolbarConfig"
            mode="default"
            style="border-bottom: 1px solid #ccc"
          />
          <Editor
            v-model="postForm.content"
            :defaultConfig="editorConfig"
            mode="default"
            style="height: 500px"
            @onCreated="handleCreated"
            @onChange="handleChange"
          />
        </div>
      </el-form-item>
      
      <!-- 按钮 -->
      <el-form-item>
        <el-button
          type="primary"
          :loading="loading"
          @click="handleSubmit"
        >
          {{ isEdit ? '保存修改' : '发布文章' }}
        </el-button>
        <el-button @click="handleCancel">取消</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, onBeforeUnmount, shallowRef } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { createPost, updatePost, getPostById, getTags } from '../api'
import '@wangeditor/editor/dist/css/style.css'
import { Editor, Toolbar } from '@wangeditor/editor-for-vue'

const route = useRoute()
const router = useRouter()
const formRef = ref(null)
const loading = ref(false)
const availableTags = ref([])

// 编辑器实例，必须用 shallowRef
const editorRef = shallowRef()

// 编辑器配置
const toolbarConfig = {
  excludeKeys: [
    'uploadImage',
    'uploadVideo',
    // 可以继续排除不需要的功能按钮
  ]
}

const editorConfig = {
  placeholder: '请输入内容...',
  autoFocus: false
}

// 组件销毁时，也及时销毁编辑器
onBeforeUnmount(() => {
  const editor = editorRef.value
  if (editor == null) return
  editor.destroy()
})

// 编辑器回调函数
const handleCreated = (editor) => {
  editorRef.value = editor
}

const handleChange = (editor) => {
  postForm.content = editor.getHtml()
}

// 判断是否为编辑模式
const isEdit = computed(() => !!route.params.id)

// 表单数据
const postForm = reactive({
  title: '',
  content: '',
  tags: []
})

// 表单验证规则
const rules = {
  title: [
    { required: true, message: '请输入文章标题', trigger: 'blur' },
    { min: 2, max: 100, message: '标题长度在 2 到 100 个字符', trigger: 'blur' }
  ],
  content: [
    { required: true, message: '请输入文章内容', trigger: 'blur' },
    { min: 10, message: '内容至少 10 个字符', trigger: 'blur' }
  ],
  tags: [
    { required: true, message: '请至少选择一个标签', trigger: 'change' },
    { 
      validator: (rule, value, callback) => {
        if (value.length > 5) {
          callback(new Error('最多选择 5 个标签'))
        } else {
          callback()
        }
      },
      trigger: 'change'
    }
  ]
}

// 获取所有标签
const fetchTags = async () => {
  try {
    const response = await getTags()
    availableTags.value = response.data.tags
  } catch (error) {
    ElMessage.error('获取标签列表失败')
  }
}

// 获取文章详情
const fetchPost = async () => {
  if (!route.params.id) return
  
  try {
    const response = await getPostById(route.params.id)
    const post = response.data.post
    
    postForm.title = post.title
    postForm.content = post.content // 富文本内容
    postForm.tags = post.tags.map(tag => tag.name)
  } catch (error) {
    ElMessage.error('获取文章失败')
    router.push('/')
  }
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        if (isEdit.value) {
          // 编辑模式：调用更新接口
          await updatePost(route.params.id, postForm)
          ElMessage.success('更新成功')
        } else {
          // 创建模式：调用创建接口
          await createPost(postForm)
          ElMessage.success('发布成功')
        }
        router.push('/')
      } catch (error) {
        ElMessage.error(isEdit ? '更新失败' : '发布失败')
      } finally {
        loading.value = false
      }
    }
  })
}

// 取消
const handleCancel = () => {
  ElMessage.info('已取消')
  router.back()
}

// 生命周期
onMounted(async () => {
  await fetchTags()
  if (isEdit) {
    await fetchPost()
  }
})
</script>

<style scoped>
.create-post {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}

.post-form {
  margin-top: 20px;
}

.tag-select {
  width: 100%;
}

.editor-container {
  border: 1px solid #ccc;
  z-index: 100; /* 确保编辑器工具栏在其他元素之上 */
}

:deep(.w-e-text-container) {
  min-height: 500px !important;
}
</style> 