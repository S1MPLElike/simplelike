<template>
  <div class="home-container">
    <div class="home-header">
      <h1>欢迎回来，{{ username }}</h1>
      <div class="header-buttons">
        <button class="logout-btn" @click="handleLogout">退出登录</button>
      </div>
    </div>
    <div class="home-content">
      <div class="sidebar">
        <ul>
          <li :class="{ active: activeTab === '首页' }" @click="activeTab = '首页'">首页</li>
          <li :class="{ active: activeTab === '发布' }" @click="activeTab = '发布'">发布</li>
          <li :class="{ active: activeTab === '好友' }" @click="activeTab = '好友'">好友</li>
          <li :class="{ active: activeTab === '设置' }" @click="activeTab = '设置'">设置</li>
          <li :class="{ active: activeTab === '我' }" @click="activeTab = '我'">我</li>
        </ul>
      </div>
      <div class="main-content">
        <PostList v-if="activeTab === '首页'" />
        <PostEditor v-if="activeTab === '发布'" />
        <FriendChat v-if="activeTab === '好友'" />
        <Settings v-if="activeTab === '设置'" />
        <UserProfile v-if="activeTab === '我'" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import PostList from './PostList.vue';
import PostEditor from './PostEditor.vue';
import FriendChat from './FriendChat.vue';
import UserProfile from './UserProfile.vue';
import Settings from './Settings.vue';

const username = ref('');
const friends = ref([]);
const activeTab = ref('首页');

// 从URL参数中获取tab值
const getTabFromUrl = () => {
  const urlParams = new URLSearchParams(window.location.search);
  const tab = urlParams.get('tab');
  const fromPost = urlParams.get('fromPost');
  if (tab && fromPost) {
    activeTab.value = tab;
    // 移除URL参数，避免刷新页面时再次触发
    const newUrl = new URL(window.location.href);
    newUrl.searchParams.delete('tab');
    newUrl.searchParams.delete('fromPost');
    window.history.replaceState({}, '', newUrl);
  }
};

const handleLogout = async () => {
  try {
    const response = await fetch('http://localhost:8080/api/v1/user/logout', {
      method: 'POST',
      credentials: 'include'
    });
    
    if (response.ok) {
      // 跳转到登录页
      window.location.href = '/';
    }
  } catch (error) {
    console.error('Logout failed:', error);
  }
};

const getUserInfo = async () => {
  try {
    const response = await fetch('http://localhost:8080/api/v1/user/auto-login', {
      method: 'GET',
      credentials: 'include'
    });
    
    if (response.ok) {
      const data = await response.json();
      if (data.code === 0) {
        username.value = data.data.username;
      }
    }
  } catch (error) {
    console.error('Get user info failed:', error);
  }
};

onMounted(() => {
  getUserInfo();
  getTabFromUrl();
});
</script>

<style scoped>
.home-container {
  height: 100vh;
  width: 100vw;
  background-color: #f5f5f5;
  display: flex;
  flex-direction: column;
}

.home-header {
  background-color: #1890ff;
  color: white;
  padding: 0 20px;
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.header-buttons {
  display: flex;
  gap: 10px;
}

.home-header h1 {
  font-size: 18px;
  font-weight: 500;
}

.logout-btn {
  background-color: transparent;
  border: 1px solid white;
  color: white;
  padding: 6px 12px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.3s;
}

.logout-btn:hover {
  background-color: rgba(255, 255, 255, 0.1);
}

.home-content {
  flex: 1;
  display: flex;
  overflow: hidden;
}

.sidebar {
  width: 200px;
  background-color: white;
  border-right: 1px solid #e8e8e8;
  padding: 20px 0;
}

.sidebar h2 {
  font-size: 14px;
  color: #999;
  padding: 0 20px 10px;
  margin: 0;
}

.sidebar ul {
  list-style: none;
  padding: 0;
  margin: 0;
}

.sidebar li {
  padding: 10px 20px;
  cursor: pointer;
  transition: all 0.3s;
  font-size: 14px;
}

.sidebar li:hover {
  background-color: #f0f8ff;
  color: #1890ff;
}

.sidebar li.active {
  background-color: #e6f7ff;
  color: #1890ff;
  border-right: 3px solid #1890ff;
}

.main-content {
  flex: 1;
  padding: 20px;
  overflow-y: auto;
}

/* 模态框样式 */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background-color: white;
  border-radius: 8px;
  width: 90%;
  max-width: 500px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.modal-header {
  padding: 16px 20px;
  border-bottom: 1px solid #e8e8e8;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.modal-header h2 {
  font-size: 16px;
  font-weight: 500;
  color: #333;
  margin: 0;
}

.close-btn {
  background: none;
  border: none;
  font-size: 20px;
  color: #999;
  cursor: pointer;
  padding: 0;
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.close-btn:hover {
  color: #333;
}

.modal-body {
  padding: 20px;
}

.form-group {
  margin-bottom: 16px;
}

.form-group label {
  display: block;
  font-size: 14px;
  font-weight: 500;
  color: #333;
  margin-bottom: 8px;
}

.form-group input,
.form-group textarea {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #d9d9d9;
  border-radius: 4px;
  font-size: 14px;
  transition: all 0.3s;
  box-sizing: border-box;
}

.form-group input:focus,
.form-group textarea:focus {
  outline: none;
  border-color: #1890ff;
  box-shadow: 0 0 0 2px rgba(24, 144, 255, 0.2);
}

.form-group textarea {
  resize: vertical;
  min-height: 120px;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 20px;
}

.cancel-btn {
  padding: 8px 16px;
  background-color: white;
  border: 1px solid #d9d9d9;
  color: #333;
  border-radius: 4px;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.3s;
}

.cancel-btn:hover {
  border-color: #1890ff;
  color: #1890ff;
}

.submit-btn {
  padding: 8px 16px;
  background-color: #1890ff;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.3s;
}

.submit-btn:hover {
  background-color: #40a9ff;
}
</style>