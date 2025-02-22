<template>
  <header class="header">
    <div class="logo">
      <router-link to="/">Keep Learning</router-link>
    </div>
    
    <div class="nav">
      <router-link to="/" class="nav-item">首页</router-link>
      
      <!-- 未登录显示 -->
      <template v-if="!userStore.isLoggedIn">
        <router-link to="/login" class="nav-item">登录</router-link>
        <router-link to="/register" class="nav-item">注册</router-link>
      </template>
      
      <!-- 登录后显示 -->
      <template v-else>
        <router-link to="/create" class="nav-item">写文章</router-link>
        <el-dropdown @command="handleCommand">
          <span class="user-menu">
            {{ userStore.user.username }}
            <el-icon class="el-icon--right"><arrow-down /></el-icon>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="profile">个人中心</el-dropdown-item>
              <el-dropdown-item command="logout">退出登录</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </template>
    </div>
  </header>
</template>

<script setup>
import { ArrowDown } from '@element-plus/icons-vue'
import { useUserStore } from '../stores/user'
import { ElMessage } from 'element-plus'

const userStore = useUserStore()

const handleCommand = (command) => {
  if (command === 'logout') {
    userStore.logout()
    ElMessage.success('退出登录成功')
  } else if (command === 'profile') {
    // TODO: 跳转到个人中心
  }
}
</script>

<style scoped>
.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 0;
  border-bottom: 1px solid #eee;
  margin-bottom: 2rem;
}

.logo a {
  font-size: 1.5rem;
  font-weight: bold;
  color: #333;
  text-decoration: none;
}

.nav {
  display: flex;
  align-items: center;
  gap: 1.5rem;
}

.nav-item {
  color: #666;
  text-decoration: none;
}

.nav-item:hover {
  color: #409EFF;
}

.user-menu {
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 4px;
}
</style> 