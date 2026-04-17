<template>
  <div class="post-container">
    <div class="search-and-sort">
      <div class="search-module">
        <input 
          type="text" 
          v-model="searchKeyword" 
          placeholder="搜索帖子标题" 
          class="search-input"
          @keyup.enter="searchPosts"
        />
        <button class="search-btn" @click="searchPosts">搜索</button>
      </div>
      <div class="sort-buttons">
        <button 
          class="sort-btn" 
          :class="{ active: sortBy === 'time' }"
          @click="setSortBy('time')"
        >
          最新发布
        </button>
        <button 
          class="sort-btn" 
          :class="{ active: sortBy === 'likes' }"
          @click="setSortBy('likes')"
        >
          最多点赞
        </button>
        <button 
          class="sort-btn" 
          :class="{ active: sortBy === 'views' }"
          @click="setSortBy('views')"
        >
          最多访问
        </button>
      </div>
    </div>
    <div class="posts-module">
      <div v-if="loading" class="loading">加载中...</div>
      <div v-else-if="posts.length === 0" class="empty">暂无帖子</div>
      <div v-else class="posts-container">
        <div class="post-item" v-for="post in posts" :key="post.post_id">
          <h3 class="post-title">{{ post.title }}</h3>
          <div class="post-meta">
            <span class="post-author">用户: {{ post.username }}</span>
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
    </div>
    <div class="pagination" v-if="posts.length > 0">
      <button 
        class="page-btn" 
        :disabled="currentPage === 1"
        @click="goToPage(1)"
      >
        首页
      </button>
      <button 
        class="page-btn" 
        :disabled="currentPage === 1"
        @click="goToPage(currentPage - 1)"
      >
        上一页
      </button>
      <span class="page-info">
        第 {{ currentPage }} 页
      </span>
      <button 
        class="page-btn"
        @click="goToPage(currentPage + 1)"
      >
        下一页
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';

const router = useRouter();

const posts = ref([]);
const loading = ref(true);
const searchKeyword = ref('');
const sortBy = ref('time'); // time, likes, views
const currentPage = ref(1);

const fetchPosts = async (keyword = '') => {
  try {
    let url = 'http://localhost:8080/api/v1/post/hot';
    if (keyword) {
      url = `http://localhost:8080/api/v1/post/search?keyword=${encodeURIComponent(keyword)}&sort=${sortBy.value}&page=${currentPage.value}`;
    } else {
      url = `http://localhost:8080/api/v1/post/hot?sort=${sortBy.value}&page=${currentPage.value}`;
    }
    
    const response = await fetch(url, {
      credentials: 'include'
    });
    
    if (response.ok) {
      const data = await response.json();
      if (data.code === 0) {
        posts.value = data.data;
      }
    }
  } catch (error) {
    console.error('Fetch posts failed:', error);
  } finally {
    loading.value = false;
  }
};

const goToPage = (page) => {
  if (page < 1) return;
  currentPage.value = page;
  loading.value = true;
  fetchPosts(searchKeyword.value);
};

const searchPosts = () => {
  loading.value = true;
  fetchPosts(searchKeyword.value);
};

const setSortBy = (sort) => {
  if (sortBy.value !== sort) {
    sortBy.value = sort;
    currentPage.value = 1; // 切换排序方式时重置到第一页
    loading.value = true;
    fetchPosts(searchKeyword.value);
  }
};

const formatTime = (timeString) => {
  const date = new Date(timeString);
  return date.toLocaleString();
};

const viewPost = (postId) => {
  router.push(`/post/${postId}?from=home`);
};

onMounted(() => {
  fetchPosts();
});
</script>

<style scoped>
.post-container {
  display: flex;
  flex-direction: column;
  gap: 16px;
  width: 100%;
  height: 100%;
  min-height: 600px;
}

.search-and-sort {
  display: flex;
  gap: 16px;
  align-items: center;
}

.search-module {
  background-color: white;
  border-radius: 8px;
  padding: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.09);
  display: flex;
  gap: 10px;
  flex: 1;
  max-width: 60%;
}

.search-input {
  flex: 1;
  padding: 10px 12px;
  border: 1px solid #d9d9d9;
  border-radius: 4px;
  font-size: 14px;
  transition: all 0.3s;
}

.search-input:focus {
  outline: none;
  border-color: #1890ff;
  box-shadow: 0 0 0 2px rgba(24, 144, 255, 0.2);
}

.search-btn {
  padding: 0 16px;
  background-color: #1890ff;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.3s;
}

.search-btn:hover {
  background-color: #40a9ff;
}

.sort-buttons {
  display: flex;
  gap: 8px;
}

.sort-btn {
  padding: 8px 12px;
  background-color: white;
  border: 1px solid #d9d9d9;
  border-radius: 4px;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.3s;
}

.sort-btn:hover {
  border-color: #1890ff;
  color: #1890ff;
}

.sort-btn.active {
  background-color: #1890ff;
  color: white;
  border-color: #1890ff;
}

.posts-module {
  background-color: white;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.09);
  flex: 1;
  min-height: 400px;
  overflow-y: auto;
}

.loading, .empty {
  text-align: center;
  padding: 40px;
  color: #999;
}

.posts-container {
  display: flex;
  flex-direction: column;
  gap: 16px;
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

/* 滚动条样式 */
.posts-module::-webkit-scrollbar {
  width: 8px;
}

.posts-module::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 4px;
}

.posts-module::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 4px;
}

.posts-module::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8;
}
</style>