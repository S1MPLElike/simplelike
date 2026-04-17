<template>
  <div class="user-profile">
    <!-- 个人信息区域 -->
    <div class="profile-header">
      <div class="avatar-container">
        <img :src="userInfo.avatar && userInfo.avatar !== '' ? `http://localhost:8080${userInfo.avatar}` : defaultAvatar" alt="头像" class="avatar" />
        <div class="avatar-upload">
          <input type="file" accept="image/*" @change="handleAvatarUpload" class="avatar-input" />
          <span class="upload-text">更换头像</span>
        </div>
      </div>
      <div class="user-info">
        <div class="user-header">
          <div class="user-name-section">
            <h2>{{ userInfo.username }}</h2>
            <p class="user-email" v-if="userInfo.email">{{ userInfo.email }}</p>
          </div>
          <button 
            class="follow-btn" 
            :class="{ following: isFollowing }" 
            @click="toggleFollow"
          >
            {{ isFollowing ? '已关注' : '关注' }}
          </button>
        </div>
        <div class="user-stats">
          <div class="stat-item">
            <span class="stat-value">{{ userInfo.likes || 0 }}</span>
            <span class="stat-label">获赞</span>
          </div>
          <div class="stat-item">
            <span class="stat-value">{{ userInfo.comments || 0 }}</span>
            <span class="stat-label">评论</span>
          </div>
          <div class="stat-item">
            <span class="stat-value">{{ userInfo.collections || 0 }}</span>
            <span class="stat-label">收藏</span>
          </div>
          <div class="stat-item">
            <span class="stat-value">{{ userInfo.followers || 0 }}</span>
            <span class="stat-label">粉丝</span>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 内容区域 -->
    <div class="profile-content">
      <div class="content-tabs">
        <button 
          class="tab-btn" 
          :class="{ active: activeTab === 'my-posts' }" 
          @click="switchTab('my-posts')"
        >
          我的帖子
        </button>
        <button 
          class="tab-btn" 
          :class="{ active: activeTab === 'collected-posts' }" 
          @click="switchTab('collected-posts')"
        >
          收藏的帖子
        </button>
      </div>
      
      <!-- 我的帖子 -->
      <div v-if="activeTab === 'my-posts'">
        <div class="loading" v-if="loading">加载中...</div>
        <div class="empty-state" v-else-if="myPosts.length === 0">
          <p>还没有发布帖子</p>
        </div>
        <div class="posts-container" v-else>
          <div class="post-item" v-for="post in myPosts" :key="post.post_id">
            <h3 class="post-title">{{ post.title }}</h3>
            <div class="post-meta">
              <span class="post-author">用户: {{ userInfo.username }}</span>
              <span class="post-time">{{ formatTime(post.create_time) }}</span>
            </div>
            <div class="post-stats">
              <span class="stat-item">👁️ {{ post.read_count }}</span>
              <span class="stat-item">❤️ {{ post.like_count }}</span>
              <span class="stat-item">⭐ {{ post.collect_count }}</span>
            </div>
            <button class="view-btn" @click="viewPost(post.post_id)">查看详情</button>
          </div>
        </div>
        <div class="pagination" v-if="currentPage > 1 || myPosts.length > 0">
            <button 
                class="page-btn" 
                :disabled="currentPage === 1"
                @click="goToPage(1, 'my-posts')"
            >
                首页
            </button>
            <button 
                class="page-btn" 
                :disabled="currentPage === 1"
                @click="goToPage(currentPage - 1, 'my-posts')"
            >
                上一页
            </button>
            <span class="page-info">
                第 {{ currentPage }} 页
            </span>
            <button 
                class="page-btn"
                @click="goToPage(currentPage + 1, 'my-posts')"
            >
                下一页
            </button>
        </div>
      </div>
      
      <!-- 收藏的帖子 -->
      <div v-else-if="activeTab === 'collected-posts'">
        <div class="loading" v-if="loading">加载中...</div>
        <div class="empty-state" v-else-if="collectedPosts.length === 0">
          <p>还没有收藏帖子</p>
        </div>
        <div class="posts-container" v-else>
          <div class="post-item" v-for="post in collectedPosts" :key="post.post_id">
            <h3 class="post-title">{{ post.title }}</h3>
            <div class="post-meta">
              <span class="post-author">用户: {{ getPostAuthor(post.user_id) }}</span>
              <span class="post-time">{{ formatTime(post.create_time) }}</span>
            </div>
            <div class="post-stats">
              <span class="stat-item">👁️ {{ post.read_count }}</span>
              <span class="stat-item">❤️ {{ post.like_count }}</span>
              <span class="stat-item">⭐ {{ post.collect_count }}</span>
            </div>
            <button class="view-btn" @click="viewPost(post.post_id)">查看详情</button>
          </div>
        </div>
        <div class="pagination" v-if="collectedPosts.length > 0">
          <button 
            class="page-btn" 
            :disabled="currentPage === 1"
            @click="goToPage(1, 'collected-posts')"
          >
            首页
          </button>
          <button 
            class="page-btn" 
            :disabled="currentPage === 1"
            @click="goToPage(currentPage - 1, 'collected-posts')"
          >
            上一页
          </button>
          <span class="page-info">
            第 {{ currentPage }} 页
          </span>
          <button 
            class="page-btn"
            @click="goToPage(currentPage + 1, 'collected-posts')"
          >
            下一页
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';

// 导入默认头像
import defaultAvatar from '../assets/login_1.jpg';
console.log('Default avatar:', defaultAvatar);

const userInfo = ref({
  user_id: '',
  username: '',
  email: '',
  avatar: '',
  likes: 0,
  comments: 0,
  collections: 0,
  followers: 0
});

const activeTab = ref('my-posts');
const myPosts = ref([]);
const collectedPosts = ref([]);
const isFollowing = ref(false);
const loading = ref(false);
const currentPage = ref(1);
const postAuthors = ref({});

const getUserProfile = async () => {
  try {
    // 首先获取当前登录用户的ID
    const loginResponse = await fetch('http://localhost:8080/api/v1/user/auto-login', {
      credentials: 'include'
    });
    
    if (loginResponse.ok) {
      const loginData = await loginResponse.json();
      if (loginData.code === 0) {
        const userID = loginData.data.user_id;
        
        // 然后获取用户详细信息
        const profileResponse = await fetch(`http://localhost:8080/api/v1/user/info/${userID}`, {
          credentials: 'include'
        });
        
        if (profileResponse.ok) {
          const profileData = await profileResponse.json();
          if (profileData.code === 0) {
            userInfo.value = profileData.data;
            console.log('User info:', userInfo.value);
            console.log('User avatar:', userInfo.value.avatar);
          }
        }
      }
    }
  } catch (error) {
    console.error('Get user profile failed:', error);
  }
};

const handleAvatarUpload = async (event) => {
  const file = event.target.files[0];
  if (!file) return;
  
  const formData = new FormData();
  formData.append('avatar', file);
  
  try {
    const response = await fetch('http://localhost:8080/api/v1/user/avatar', {
      method: 'POST',
      credentials: 'include',
      body: formData
    });
    
    if (response.ok) {
      const data = await response.json();
      if (data.code === 0) {
        // 更新本地用户信息
        userInfo.value.avatar = data.data.avatar;
        alert('头像上传成功');
      } else {
        alert('头像上传失败: ' + data.message);
      }
    } else {
      alert('网络错误，请稍后重试');
    }
  } catch (error) {
    console.error('Upload avatar failed:', error);
    alert('网络错误，请稍后重试');
  }
};

// 关注/取消关注功能
const toggleFollow = async () => {
  try {
    // 首先获取当前登录用户的ID
    const loginResponse = await fetch('http://localhost:8080/api/v1/user/auto-login', {
      credentials: 'include'
    });
    
    if (loginResponse.ok) {
      const loginData = await loginResponse.json();
      if (loginData.code === 0) {
        const currentUserID = loginData.data.user_id;
        
        // 检查是否关注自己
        if (currentUserID === userInfo.value.user_id) {
          alert('不能关注自己');
          return;
        }
        
        // 调用关注/取消关注API
        const url = isFollowing.value 
          ? `http://localhost:8080/api/v1/user/unfollow/${userInfo.value.user_id}` 
          : `http://localhost:8080/api/v1/user/follow/${userInfo.value.user_id}`;
        
        const response = await fetch(url, {
          method: 'POST',
          credentials: 'include'
        });
        
        if (response.ok) {
          const data = await response.json();
          if (data.code === 0) {
            isFollowing.value = !isFollowing.value;
            // 更新粉丝数
            if (isFollowing.value) {
              userInfo.value.followers++;
            } else {
              userInfo.value.followers--;
            }
            alert(isFollowing.value ? '关注成功' : '取消关注成功');
          } else {
            alert(data.message || '操作失败');
          }
        } else {
          alert('网络错误，请稍后重试');
        }
      }
    }
  } catch (error) {
    console.error('Toggle follow failed:', error);
    alert('网络错误，请稍后重试');
  }
};

// 切换标签页
const switchTab = (tab) => {
  activeTab.value = tab;
  if (tab === 'my-posts') {
    console.log('Switching to my posts, user ID:', userInfo.value.user_id);
    getMyPosts();
  } else if (tab === 'collected-posts') {
    console.log('Switching to collected posts, user ID:', userInfo.value.user_id);
    getCollectedPosts();
  }
};

// 获取我的帖子
const getMyPosts = async (page = 1) => {
  try {
    if (!userInfo.value.user_id) {
      console.log('User ID not available');
      return;
    }
    
    console.log('Fetching my posts for user:', userInfo.value.user_id, 'page:', page);
    loading.value = true;
    const response = await fetch(`http://localhost:8080/api/v1/post/user/${userInfo.value.user_id}?page=${page}`, {
      credentials: 'include'
    });
    
    console.log('Response status:', response.status);
    if (response.ok) {
      const data = await response.json();
      console.log('Response data:', data);
      if (data.code === 0) {
        myPosts.value = data.data || [];
        console.log('My posts:', myPosts.value);
      }
    }
  } catch (error) {
    console.error('Get my posts failed:', error);
  } finally {
    loading.value = false;
  }
};

// 获取收藏的帖子
const getCollectedPosts = async (page = 1) => {
  try {
    if (!userInfo.value.user_id) return;
    
    loading.value = true;
    const response = await fetch(`http://localhost:8080/api/v1/post/collected/${userInfo.value.user_id}?page=${page}`, {
      credentials: 'include'
    });
    
    if (response.ok) {
      const data = await response.json();
      if (data.code === 0) {
        collectedPosts.value = data.data || [];
        // 获取作者信息
        for (const post of collectedPosts.value) {
          if (post.user_id && !postAuthors.value[post.user_id]) {
            await getPostAuthorInfo(post.user_id);
          }
        }
      }
    }
  } catch (error) {
    console.error('Get collected posts failed:', error);
  } finally {
    loading.value = false;
  }
};

// 跳转到指定页面
const goToPage = (page, tab) => {
  if (page < 1) return;
  currentPage.value = page;
  if (tab === 'my-posts') {
    getMyPosts(page);
  } else if (tab === 'collected-posts') {
    getCollectedPosts(page);
  }
};

// 格式化时间
const formatTime = (timeString) => {
  const date = new Date(timeString);
  return date.toLocaleString();
};

// 查看帖子详情
const viewPost = (postId) => {
  window.location.href = `/post/${postId}?from=user`;
};

// 获取帖子作者信息
const getPostAuthorInfo = async (userId) => {
  try {
    const response = await fetch(`http://localhost:8080/api/v1/user/info/${userId}`, {
      credentials: 'include'
    });
    
    if (response.ok) {
      const data = await response.json();
      if (data.code === 0) {
        postAuthors.value[userId] = data.data.username;
      }
    }
  } catch (error) {
    console.error('Get post author info failed:', error);
  }
};

// 获取帖子作者
const getPostAuthor = (userId) => {
  return postAuthors.value[userId] || '未知用户';
};

onMounted(async () => {
  await getUserProfile();
  // 确保在用户信息加载后获取帖子列表
  setTimeout(() => {
    getMyPosts();
  }, 100);
});
</script>

<style scoped>
.user-profile {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.profile-header {
  flex: 1;
  min-height: 200px;
  background: linear-gradient(135deg, #e6f7ff 0%, #f0f9ff 100%);
  border-radius: 8px;
  padding: 30px;
  margin-bottom: 20px;
  display: flex;
  align-items: center;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
}

.avatar-container {
  position: relative;
  margin-right: 40px;
}

.avatar {
  width: 120px;
  height: 120px;
  border-radius: 50%;
  border: 4px solid white;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

.avatar-upload {
  position: absolute;
  bottom: 0;
  right: 0;
  background-color: rgba(0, 0, 0, 0.6);
  color: white;
  padding: 6px 12px;
  border-radius: 12px;
  cursor: pointer;
  font-size: 12px;
  display: flex;
  align-items: center;
  gap: 4px;
  transition: all 0.3s;
}

.avatar-upload:hover {
  background-color: rgba(0, 0, 0, 0.8);
}

.avatar-input {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  opacity: 0;
  cursor: pointer;
}

.user-info {
  flex: 1;
}

.user-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20px;
}

.user-header h2 {
  font-size: 24px;
  font-weight: 600;
  margin: 0;
  color: #333;
}

.user-email {
  font-size: 14px;
  font-weight: 400;
  margin: 4px 0 0 0;
  color: #666;
  line-height: 1.4;
}

.user-name-section {
  display: flex;
  flex-direction: column;
}

.follow-btn {
  padding: 8px 24px;
  background-color: #1890ff;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.3s;
}

.follow-btn:hover {
  background-color: #40a9ff;
}

.follow-btn.following {
  background-color: #f0f0f0;
  color: #666;
}

.follow-btn.following:hover {
  background-color: #e0e0e0;
}

.user-stats {
  display: flex;
  gap: 40px;
}

.stat-item {
  text-align: center;
}

.stat-value {
  display: block;
  font-size: 20px;
  font-weight: 600;
  color: #1890ff;
  margin-bottom: 4px;
}

.stat-label {
  display: block;
  font-size: 14px;
  color: #666;
}

.profile-content {
  flex: 2;
  background-color: white;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
}

.content-tabs {
  display: flex;
  margin-bottom: 20px;
  border-bottom: 1px solid #e8e8e8;
}

.tab-btn {
  padding: 10px 20px;
  background: none;
  border: none;
  font-size: 16px;
  font-weight: 500;
  color: #666;
  cursor: pointer;
  border-bottom: 2px solid transparent;
  transition: all 0.3s;
}

.tab-btn:hover {
  color: #1890ff;
}

.tab-btn.active {
  color: #1890ff;
  border-bottom-color: #1890ff;
}

.loading {
  text-align: center;
  padding: 40px;
  color: #999;
}

.posts-container {
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin-top: 20px;
}

.post-item {
  padding: 16px;
  border: 1px solid #e8e8e8;
  border-radius: 8px;
  transition: all 0.3s;
}

.post-item:hover {
  border-color: #1890ff;
  box-shadow: 0 2px 8px rgba(24, 144, 255, 0.1);
}

.post-title {
  font-size: 16px;
  font-weight: 500;
  color: #333;
  margin: 0 0 8px;
}

.post-meta {
  font-size: 12px;
  color: #999;
  margin-bottom: 8px;
  display: flex;
  gap: 16px;
}

.post-stats {
  font-size: 12px;
  color: #666;
  margin-bottom: 12px;
  display: flex;
  gap: 16px;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 4px;
}

.view-btn {
  padding: 6px 12px;
  background-color: #1890ff;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.3s;
}

.view-btn:hover {
  background-color: #40a9ff;
}

/* 分页样式 */
.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 10px;
  margin-top: 20px;
  padding: 10px;
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.09);
}

.page-btn {
  padding: 6px 12px;
  background-color: #f0f0f0;
  color: #333;
  border: 1px solid #d9d9d9;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.3s;
}

.page-btn:hover:not(:disabled) {
  background-color: #e6f7ff;
  border-color: #1890ff;
  color: #1890ff;
}

.page-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.page-info {
  font-size: 14px;
  color: #666;
  margin: 0 10px;
}

.empty-state {
  text-align: center;
  padding: 60px 0;
  color: #999;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .profile-header {
    flex-direction: column;
    text-align: center;
    padding: 30px;
  }
  
  .avatar-container {
    margin-right: 0;
    margin-bottom: 20px;
  }
  
  .user-stats {
    justify-content: center;
    gap: 30px;
  }
}
</style>