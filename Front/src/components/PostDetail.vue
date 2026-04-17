<template>
  <div class="post-detail">
    <button class="back-btn" @click="goBack">返回</button>
    <div class="main-container">
      <!-- 左侧个人信息 -->
      <div class="left-sidebar">
        <div class="user-info-card">
          <div class="user-avatar">
            <img :src="userInfo.avatar && userInfo.avatar !== '' ? `http://localhost:8080${userInfo.avatar}` : defaultAvatar" alt="用户头像" />
          </div>
          <div class="user-basic-info">
            <h3 class="user-name">{{ userInfo.username }}</h3>
            <p class="user-id">用户: {{ userInfo.username }}</p>
            <div class="user-stats">
              <span class="stat">粉丝: {{ userInfo.followers }}</span>
              <span class="stat">关注: {{ userInfo.following }}</span>
            </div>
            <div class="user-actions">
              <button v-if="!isSelf" class="follow-btn" :class="{ following: isFollowing }" @click="toggleFollow">
                {{ isFollowing ? '已关注' : '关注' }}
              </button>
              <button class="profile-btn" @click="openUserProfile">个人主页</button>
            </div>
          </div>
          <div class="user-achievements">
            <h4>个人成就</h4>
            <div class="achievement-item">
              <span class="achievement-icon">👍</span>
              <span class="achievement-text">获得 {{ userInfo.likes }} 次点赞</span>
            </div>
            <div class="achievement-item">
              <span class="achievement-icon">💬</span>
              <span class="achievement-text">内容获得 {{ userInfo.comments }} 次评论</span>
            </div>
            <div class="achievement-item">
              <span class="achievement-icon">⭐</span>
              <span class="achievement-text">获得 {{ userInfo.collections }} 次收藏</span>
            </div>
          </div>
        </div>
      </div>
      
      <!-- 中间帖子内容 -->
      <div class="container">
        <div class="post-content">
          <h1 class="post-title">{{ post.title }}</h1>
          
          <div class="post-meta">
            <span class="post-time">{{ formatTime(post.create_time) }}</span>
            <span class="post-read">👁️ {{ post.read_count }}</span>
          </div>
          
          <div class="post-body" v-html="post.content"></div>
        </div>
      </div>
      
      <!-- 右侧工具栏 -->
      <div class="right-toolbar">
        <div class="toolbar-item">
          <button 
            class="toolbar-btn" 
            :class="{ active: isLiked }"
            @click="toggleLike"
            title="点赞"
          >
            <span class="toolbar-icon">👍</span>
            <span class="toolbar-count">{{ post.like_count }}</span>
          </button>
        </div>
        <div class="toolbar-item">
          <button 
            class="toolbar-btn" 
            :class="{ active: isCollected }"
            @click="toggleCollect"
            title="收藏"
          >
            <span class="toolbar-icon">⭐</span>
            <span class="toolbar-count">{{ post.collect_count }}</span>
          </button>
        </div>
        <div class="toolbar-item">
          <button 
            class="toolbar-btn"
            :class="{ active: showComments }"
            @click="showComments = !showComments"
            title="评论"
          >
            <span class="toolbar-icon">💬</span>
            <span class="toolbar-count">{{ comments.length }}</span>
          </button>
        </div>
      </div>
      
      <!-- 评论侧边栏 -->
      <div v-if="showComments" class="comments-sidebar">
        <div class="comments-header">
          <h3>评论 ({{ comments.length }})</h3>
          <button class="close-btn" @click="showComments = false">×</button>
        </div>
        
        <!-- 评论输入框 -->
        <div class="comment-input-area">
          <textarea 
            v-model="newComment"
            placeholder="写下你的评论..."
            maxlength="1000"
            class="comment-input"
          ></textarea>
          <div class="comment-input-footer">
            <span class="char-count">{{ 1000 - newComment.length }} 字符</span>
            <button class="comment-btn" @click="submitComment">发布评论</button>
          </div>
        </div>
        
        <!-- 评论列表 -->
        <div class="comments-list">
          <div v-for="comment in comments" :key="comment.id" class="comment-item">
            <div class="comment-header">
              <div class="comment-avatar">
                <img :src="comment.avatar && comment.avatar !== '' ? `http://localhost:8080${comment.avatar}` : defaultAvatar" alt="用户头像" class="avatar-img">
              </div>
              <div class="comment-info">
                <span class="comment-author">{{ comment.username || '未知用户' }}</span>
                <span class="comment-time">{{ formatTime(comment.create_time) }}</span>
              </div>
            </div>
            <div class="comment-content">{{ comment.content }}</div>
            <div class="comment-actions">
              <button class="comment-action-btn" @click="toggleCommentLike(comment.id)">
                <span class="like-icon" :class="{ 'liked': comment.is_liked }">{{ comment.is_liked ? '❤️' : '👍' }}</span>
                <span class="like-count">{{ comment.like_count || 0 }}</span>
              </button>
            </div>
          </div>
          
          <div v-if="comments.length === 0" class="no-comments">
            暂无评论，快来发表你的看法吧！
          </div>
        </div>
      </div>
    </div>
    
    <!-- 个人主页弹窗 -->
    <div v-if="showUserProfile" class="profile-modal">
      <div class="profile-modal-content">
        <div class="profile-modal-header">
          <h3>{{ userProfileInfo.username }}的个人主页</h3>
          <button class="close-btn" @click="closeUserProfile">×</button>
        </div>
        <div class="profile-modal-body">
          <div class="profile-left">
            <div class="user-info-section">
              <div class="user-avatar">
                <img :src="userProfileInfo.avatar && userProfileInfo.avatar !== '' ? `http://localhost:8080${userProfileInfo.avatar}` : defaultAvatar" alt="用户头像" />
              </div>
              <div class="user-basic-info">
                <h4>{{ userProfileInfo.username }}</h4>
                <div class="user-stats">
                  <span class="stat">粉丝: {{ userProfileInfo.followers }}</span>
                  <span class="stat">关注: {{ userProfileInfo.following }}</span>
                </div>
                <div class="user-actions">
                  <button v-if="showFollowButtonInProfile" class="follow-btn" :class="{ following: isFollowing }" @click="toggleFollow">
                    {{ isFollowing ? '已关注' : '关注' }}
                  </button>
                </div>
              </div>
            </div>
            
            <div class="user-achievements">
              <h4>个人成就</h4>
              <div class="achievement-item">
                <span class="achievement-icon">👍</span>
                <span class="achievement-text">获得 {{ userProfileInfo.likes }} 次点赞</span>
              </div>
              <div class="achievement-item">
                <span class="achievement-icon">💬</span>
                <span class="achievement-text">内容获得 {{ userProfileInfo.comments }} 次评论</span>
              </div>
              <div class="achievement-item">
                <span class="achievement-icon">⭐</span>
                <span class="achievement-text">获得 {{ userProfileInfo.collections }} 次收藏</span>
              </div>
            </div>
          </div>
          
          <div class="profile-right">
            <div class="profile-tabs">
              <button 
                class="tab-btn" 
                :class="{ active: activeProfileTab === 'posts' }"
                @click="activeProfileTab = 'posts'"
              >
                帖子
              </button>
              <button 
                class="tab-btn" 
                :class="{ active: activeProfileTab === 'collected' }"
                @click="activeProfileTab = 'collected'"
              >
                收藏
              </button>
            </div>
            
            <div class="profile-content">
              <div v-if="activeProfileTab === 'posts'" class="user-posts-section">
                <h4>{{ userProfileInfo.username }}的帖子</h4>
                <div v-if="userPosts.length === 0" class="empty-state">
                  暂无帖子
                </div>
                <div v-else class="posts-list">
                  <div class="post-item" v-for="post in userPosts" :key="post.post_id">
                    <h5 class="post-title">{{ post.title }}</h5>
                    <div class="post-meta">
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
              
              <div v-if="activeProfileTab === 'collected'" class="user-collected-section">
                <h4>{{ userProfileInfo.username }}收藏的帖子</h4>
                <div v-if="userCollectedPosts.length === 0" class="empty-state">
                  暂无收藏
                </div>
                <div v-else class="posts-list">
                  <div class="post-item" v-for="post in userCollectedPosts" :key="post.post_id">
                    <h5 class="post-title">{{ post.title }}</h5>
                    <div class="post-meta">
                      <span class="post-author">作者: {{ post.author }}</span>
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
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue';
import { useRoute } from 'vue-router';

const route = useRoute();
const postId = route.params.id;

const post = ref({
  post_id: '',
  user_id: '',
  title: '',
  content: '',
  like_count: 0,
  read_count: 0,
  collect_count: 0,
  create_time: ''
});

const comments = ref([]);
const showComments = ref(false);
const newComment = ref('');
const isLiked = ref(false);
const isCollected = ref(false);
const isFollowing = ref(false);

// 默认头像
import defaultAvatar from '../assets/login_1.jpg';

// 用户信息
const userInfo = ref({
  username: '用户',
  avatar: '',
  followers: 0,
  following: 0,
  likes: 0,
  comments: 0,
  collections: 0
});

// 当前登录用户信息
const currentUser = ref({
  user_id: 0,
  username: ''
});

// 计算属性：是否为自己
const isSelf = computed(() => {
  return currentUser && currentUser.value && currentUser.value.user_id && post && post.value && post.value.user_id && currentUser.value.user_id === post.value.user_id;
});

// 计算属性：是否显示关注按钮（个人主页弹窗中）
const showFollowButtonInProfile = computed(() => {
  return currentUser && currentUser.value && currentUser.value.user_id && post && post.value && post.value.user_id && currentUser.value.user_id !== post.value.user_id;
});

// 个人主页弹窗
const showUserProfile = ref(false);
const activeProfileTab = ref('posts');
const userProfileInfo = ref({
  username: '用户',
  avatar: '',
  followers: 0,
  following: 0,
  likes: 0,
  comments: 0,
  collections: 0
});
const userPosts = ref([]);
const userCollectedPosts = ref([]);

// 格式化时间
const formatTime = (time) => {
  if (!time) return '';
  const date = new Date(time);
  return date.toLocaleString('zh-CN');
};

// 获取用户信息
const getUserInfo = async (userId) => {
  try {
    // 检查userId是否有效
    if (!userId || userId === 0) {
      console.error('Invalid user ID:', userId);
      return;
    }
    
    const response = await fetch(`http://localhost:8080/api/v1/user/info/${userId}`, {
      credentials: 'include'
    });
    
    if (response.ok) {
      const data = await response.json();
      if (data.code === 0) {
        userInfo.value = {
          user_id: data.data.user_id || userId,
          username: data.data.username || '用户',
          avatar: data.data.avatar || '',
          followers: data.data.followers || 0,
          following: data.data.following || 0,
          likes: data.data.likes || 0,
          comments: data.data.comments || 0,
          collections: data.data.collections || 0
        };
      } else {
        console.error('Get user info failed:', data.message);
      }
    } else {
      console.error('Get user info failed: Network error');
    }
  } catch (error) {
    console.error('Get user info failed:', error);
  }
};

// 检查用户是否已点赞
const checkLikeStatus = async () => {
  try {
    const response = await fetch(`http://localhost:8080/api/v1/post/check-like/${postId}`, {
      credentials: 'include'
    });
    
    if (response.ok) {
      const data = await response.json();
      if (data.code === 0) {
        isLiked.value = data.data.liked;
      }
    }
  } catch (error) {
    console.error('Check like status failed:', error);
  }
};

// 获取帖子详情
const getPostDetail = async () => {
	try {
		const response = await fetch(`http://localhost:8080/api/v1/post/detail/${postId}`, {
			credentials: 'include'
		});
		
		if (response.ok) {
			const data = await response.json();
			if (data.code === 0) {
				post.value = data.data;
				// 获取用户信息
				getUserInfo(data.data.user_id);
				// 检查点赞状态
				checkLikeStatus();
				// 检查收藏状态
				checkCollectStatus();
				// 检查关注状态
				checkFollowStatus(data.data.user_id);
			} else {
				// 显示错误信息
				alert(data.message || '获取帖子详情失败');
			}
		} else {
			// 显示错误信息
			alert('网络错误，请稍后重试');
		}
	} catch (error) {
		console.error('Get post detail failed:', error);
		// 显示错误信息
		alert('网络错误，请稍后重试');
	}
};

// 获取评论列表
const getComments = async () => {
  try {
    const response = await fetch(`http://localhost:8080/api/v1/post/comments/${postId}`, {
      credentials: 'include'
    });
    
    if (response.ok) {
      const data = await response.json();
      if (data.code === 0) {
        // 确保data.data是数组
        if (Array.isArray(data.data)) {
          comments.value = data.data;
        } else {
          comments.value = [];
        }
      } else {
        comments.value = [];
      }
    } else {
      comments.value = [];
    }
  } catch (error) {
    console.error('Get comments failed:', error);
    comments.value = [];
  }
};

// 切换点赞状态
const toggleLike = async () => {
  try {
    const response = await fetch(`http://localhost:8080/api/v1/post/like/${postId}`, {
      method: 'POST',
      credentials: 'include'
    });
    
    if (response.ok) {
      const data = await response.json();
      if (data.code === 0) {
        isLiked.value = data.data.liked;
        post.value.like_count = data.data.like_count;
      }
    }
  } catch (error) {
    console.error('Toggle like failed:', error);
  }
};

// 检查用户是否已收藏
const checkCollectStatus = async () => {
	try {
		const response = await fetch(`http://localhost:8080/api/v1/post/check-collect/${postId}`, {
			credentials: 'include'
		});
		
		if (response.ok) {
			const data = await response.json();
			if (data.code === 0) {
				isCollected.value = data.data.collected;
			}
		}
	} catch (error) {
		console.error('Check collect status failed:', error);
	}
};

// 检查用户是否已关注
const checkFollowStatus = async (userId) => {
  try {
    if (!userId) {
      console.error('User ID not found');
      return;
    }
    
    const response = await fetch(`http://localhost:8080/api/v1/user/check-follow/${userId}`, {
      credentials: 'include'
    });
    
    if (response.ok) {
      const data = await response.json();
      if (data.code === 0) {
        isFollowing.value = data.data.following;
      }
    }
  } catch (error) {
    console.error('Check follow status failed:', error);
  }
};

// 评论点赞功能（无上限）
const toggleCommentLike = async (commentId) => {
  try {
    const response = await fetch(`http://localhost:8080/api/v1/post/comment/like/${commentId}`, {
      method: 'POST',
      credentials: 'include'
    });
    
    if (response.ok) {
      const data = await response.json();
      if (data.code === 0) {
        // 更新评论的点赞数量
        const comment = comments.value.find(c => c.id === commentId);
        if (comment) {
          comment.like_count = data.data.like_count;
        }
      }
    }
  } catch (error) {
    console.error('Comment like failed:', error);
  }
};

// 切换收藏状态
const toggleCollect = async () => {
	try {
		const response = await fetch(`http://localhost:8080/api/v1/post/collect/${postId}`, {
			method: 'POST',
			credentials: 'include'
		});
		
		if (response.ok) {
			const data = await response.json();
			if (data.code === 0) {
				isCollected.value = data.data.collected;
				post.value.collect_count = data.data.collect_count;
				// 重新获取用户信息，更新收藏总数
				if (post.value.user_id) {
					getUserInfo(post.value.user_id);
				}
			} else {
				// 显示错误信息
				alert(data.message || '操作失败');
			}
		}
	} catch (error) {
		console.error('Toggle collect failed:', error);
		alert('网络错误，请稍后重试');
	}
};

// 切换关注状态
const toggleFollow = async () => {
  try {
    if (!post.value) {
      console.error('Post not found');
      return;
    }
    const userId = post.value.user_id;
    if (!userId) {
      console.error('User ID not found');
      return;
    }
    
    const endpoint = isFollowing.value ? 
      `http://localhost:8080/api/v1/user/unfollow/${userId}` : 
      `http://localhost:8080/api/v1/user/follow/${userId}`;
    
    const response = await fetch(endpoint, {
      method: 'POST',
      credentials: 'include'
    });
    
    if (response.ok) {
      const data = await response.json();
      if (data.code === 0) {
        isFollowing.value = !isFollowing.value;
        // 重新获取用户信息，更新粉丝和关注数量
        if (userId) {
          getUserInfo(userId);
        }
      } else {
        console.error('Toggle follow failed:', data.message);
        // 恢复原来的状态
        isFollowing.value = !isFollowing.value;
      }
    } else {
      console.error('Toggle follow failed: Network error');
      // 恢复原来的状态
      isFollowing.value = !isFollowing.value;
    }
  } catch (error) {
    console.error('Toggle follow failed:', error);
    // 恢复原来的状态
    isFollowing.value = !isFollowing.value;
  }
};

// 提交评论
const submitComment = async () => {
  console.log('Submit comment called');
  if (!newComment.value.trim()) {
    alert('评论内容不能为空');
    return;
  }
  
  console.log('Comment content:', newComment.value);
  console.log('Post ID:', postId);
  
  // 确保postId是数字类型
  const numericPostId = parseInt(postId);
  console.log('Numeric Post ID:', numericPostId);
  
  if (isNaN(numericPostId)) {
    alert('帖子ID无效');
    return;
  }
  
  try {
    const response = await fetch('http://localhost:8080/api/v1/post/comment', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      credentials: 'include',
      body: JSON.stringify({
        post_id: numericPostId,
        content: newComment.value
      })
    });
    
    console.log('Response status:', response.status);
    
    if (response.ok) {
      const data = await response.json();
      console.log('Response data:', data);
      if (data.code === 0) {
        newComment.value = '';
        getComments();
      } else {
        console.error('API error:', data.message);
        alert(data.message || '发布评论失败');
      }
    } else {
      console.error('Network error:', response.status);
      alert('网络错误，请稍后重试');
    }
  } catch (error) {
    console.error('Submit comment failed:', error);
    alert('网络错误，请稍后重试');
  }
};

// 返回上一页
const goBack = () => {
  // 从URL参数中获取来源
  const urlParams = new URLSearchParams(window.location.search);
  const from = urlParams.get('from');
  
  if (from === 'user') {
    // 从"我"模块来的，返回到"我"模块
    window.location.href = '/home?tab=我&fromPost=true';
  } else {
    // 从首页来的，返回到首页
    window.location.href = '/home';
  }
};

// 打开个人主页
const openUserProfile = async () => {
  console.log('Open user profile clicked');
  console.log('post.value:', post.value);
  console.log('userInfo.value:', userInfo.value);
  
  // 获取用户信息
  let userId = null;
  if (post && post.value && post.value.user_id) {
    userId = post.value.user_id;
  } else if (userInfo && userInfo.value && userInfo.value.user_id) {
    userId = userInfo.value.user_id;
  }
  console.log('userId:', userId);
  
  if (userId) {
    try {
      await getUserProfileInfo(userId);
      await getUserPosts(userId);
      await getUserCollectedPosts(userId);
      showUserProfile.value = true;
      console.log('User profile opened successfully');
    } catch (error) {
      console.error('Open user profile failed:', error);
    }
  } else {
    console.error('User ID not found');
  }
};

// 关闭个人主页
const closeUserProfile = () => {
  showUserProfile.value = false;
};

// 查看帖子详情
const viewPost = (postId) => {
  window.location.href = `/post/${postId}?from=user`;
};

// 获取用户详细信息
const getUserProfileInfo = async (userId) => {
  try {
    const response = await fetch(`http://localhost:8080/api/v1/user/info/${userId}`, {
      credentials: 'include'
    });
    
    if (response.ok) {
      const data = await response.json();
      if (data.code === 0) {
        userProfileInfo.value = {
          username: data.data.username || '用户',
          avatar: data.data.avatar || '',
          followers: data.data.followers || 0,
          following: data.data.following || 0,
          likes: data.data.likes || 0,
          comments: data.data.comments || 0,
          collections: data.data.collections || 0
        };
      }
    }
  } catch (error) {
    console.error('Get user profile info failed:', error);
  }
};

// 获取用户帖子
const getUserPosts = async (userId) => {
  try {
    const response = await fetch(`http://localhost:8080/api/v1/post/user/${userId}?page=1`, {
      credentials: 'include'
    });
    
    if (response.ok) {
      const data = await response.json();
      if (data.code === 0) {
        userPosts.value = data.data || [];
      }
    }
  } catch (error) {
    console.error('Get user posts failed:', error);
  }
};

// 获取用户收藏的帖子
const getUserCollectedPosts = async (userId) => {
  try {
    const response = await fetch(`http://localhost:8080/api/v1/post/collected/${userId}?page=1`, {
      credentials: 'include'
    });
    
    if (response.ok) {
      const data = await response.json();
      if (data.code === 0) {
        userCollectedPosts.value = data.data || [];
      }
    }
  } catch (error) {
    console.error('Get user collected posts failed:', error);
  }
};

// 初始化
onMounted(async () => {
  await getCurrentUser();
  await getPostDetail();
  if (post.value.user_id) {
    await getUserInfo(post.value.user_id);
    await checkFollowStatus(post.value.user_id);
  }
  getComments();
});

// 获取当前登录用户信息
const getCurrentUser = async () => {
  try {
    const response = await fetch('http://localhost:8080/api/v1/user/info', {
      credentials: 'include'
    });
    
    if (response.ok) {
      const data = await response.json();
      if (data.code === 0) {
        currentUser.value = {
          user_id: data.data.user_id || 0,
          username: data.data.username || ''
        };
      }
    }
  } catch (error) {
    console.error('Get current user failed:', error);
  }
};
</script>

<style scoped>
.post-detail {
  min-height: 100vh;
  background-color: #f5f5f5;
  padding: 2rem 0;
}

.back-btn {
  position: fixed;
  top: 20px;
  left: 20px;
  background-color: #1890ff;
  color: white;
  border: none;
  border-radius: 4px;
  padding: 6px 12px;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.3s;
  z-index: 100;
}

.back-btn:hover {
  background-color: #40a9ff;
}

.main-container {
  position: relative;
  max-width: 1400px;
  margin: 0 auto;
  display: flex;
  gap: 2rem;
  padding: 0 1rem;
}

/* 左侧个人信息 */
.left-sidebar {
  width: 300px;
  flex-shrink: 0;
}

.user-info-card {
  position: sticky;
  top: 2rem;
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  padding: 2rem;
}

.user-avatar {
  text-align: center;
  margin-bottom: 1.5rem;
}

.user-avatar img {
  width: 120px;
  height: 120px;
  border-radius: 50%;
  object-fit: cover;
  border: 2px solid #e8e8e8;
}

.user-basic-info {
  text-align: center;
  margin-bottom: 2rem;
}

.user-name {
  font-size: 1.2rem;
  font-weight: bold;
  margin-bottom: 0.5rem;
  color: #333;
}

.user-id {
  color: #666;
  font-size: 0.9rem;
  margin-bottom: 1rem;
}

.user-stats {
  display: flex;
  justify-content: center;
  gap: 1.5rem;
  margin-bottom: 1.5rem;
}

.stat {
  color: #666;
  font-size: 0.9rem;
}

.follow-btn {
  background-color: #1890ff;
  color: white;
  border: none;
  border-radius: 20px;
  padding: 0.5rem 2rem;
  cursor: pointer;
  transition: all 0.3s;
  font-size: 0.9rem;
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

.user-achievements {
  border-top: 1px solid #e8e8e8;
  padding-top: 1.5rem;
}

.user-achievements h4 {
  font-size: 1rem;
  font-weight: bold;
  margin-bottom: 1rem;
  color: #333;
}

.achievement-item {
  display: flex;
  align-items: center;
  gap: 0.8rem;
  margin-bottom: 0.8rem;
  color: #666;
  font-size: 0.9rem;
}

.achievement-icon {
  font-size: 1.1rem;
}

/* 中间帖子内容 */
.container {
  flex: 1;
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  padding: 2rem;
  min-width: 0;
}

.post-title {
  font-size: 2rem;
  font-weight: bold;
  margin-bottom: 1rem;
  color: #333;
}

.post-meta {
  display: flex;
  gap: 1rem;
  margin-bottom: 1.5rem;
  color: #666;
  font-size: 0.9rem;
}

.post-body {
  line-height: 1.8;
  margin-bottom: 2rem;
  color: #333;
}

/* 右侧工具栏 */
.right-toolbar {
  position: sticky;
  top: 2rem;
  align-self: flex-start;
  display: flex;
  flex-direction: column;
  gap: 1rem;
  padding: 1rem;
  background-color: #e6f7ff;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  min-width: 80px;
}

.toolbar-item {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.toolbar-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.3rem;
  background: none;
  border: 1px solid #d9d9d9;
  border-radius: 8px;
  padding: 0.8rem;
  cursor: pointer;
  transition: all 0.3s;
  font-size: 0.9rem;
  width: 60px;
}

.toolbar-btn:hover {
  border-color: #1890ff;
  color: #1890ff;
  background-color: #f0f9ff;
}

.toolbar-btn.active {
  background-color: #e6f7ff;
  border-color: #1890ff;
  color: #1890ff;
}

.toolbar-icon {
  font-size: 1.2rem;
}

.toolbar-count {
  font-size: 0.8rem;
  font-weight: 500;
}

/* 评论侧边栏 */
.comments-sidebar {
  position: fixed;
  top: 0;
  right: 0;
  width: 400px;
  height: 100vh;
  background-color: white;
  box-shadow: -2px 0 10px rgba(0, 0, 0, 0.1);
  z-index: 1000;
  display: flex;
  flex-direction: column;
  animation: slideIn 0.3s ease-out;
}

@keyframes slideIn {
  from {
    transform: translateX(100%);
  }
  to {
    transform: translateX(0);
  }
}

.comments-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.5rem;
  border-bottom: 1px solid #e8e8e8;
}

.comments-header h3 {
  font-size: 1.2rem;
  margin: 0;
  color: #333;
}

.close-btn {
  background: none;
  border: none;
  font-size: 1.5rem;
  cursor: pointer;
  color: #666;
  transition: color 0.3s;
}

.close-btn:hover {
  color: #1890ff;
}

.comment-input-area {
  padding: 1.5rem;
  border-bottom: 1px solid #e8e8e8;
}

.comment-input {
  width: 100%;
  min-height: 100px;
  padding: 1rem;
  border: 1px solid #d9d9d9;
  border-radius: 4px;
  resize: vertical;
  font-size: 1rem;
  box-sizing: border-box;
}

.comment-input-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 0.5rem;
}

.char-count {
  color: #666;
  font-size: 0.9rem;
}

.comment-btn {
  background-color: #1890ff;
  color: white;
  border: none;
  border-radius: 4px;
  padding: 0.5rem 1rem;
  cursor: pointer;
  transition: background-color 0.3s;
}

.comment-btn:hover {
  background-color: #40a9ff;
}

.comments-list {
  flex: 1;
  overflow-y: auto;
  padding: 1.5rem;
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.comment-item {
  padding-bottom: 1.5rem;
  border-bottom: 1px solid #f0f0f0;
}

.comment-header {
  display: flex;
  align-items: flex-start;
  gap: 1rem;
  margin-bottom: 0.8rem;
}

.comment-avatar {
  flex-shrink: 0;
}

.comment-avatar img {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  object-fit: cover;
  border: 1px solid #e8e8e8;
}

.comment-info {
  flex: 1;
}

.comment-author {
  font-weight: bold;
  color: #333;
  display: block;
  margin-bottom: 0.25rem;
  font-size: 0.95rem;
}

.comment-time {
  color: #999;
  font-size: 0.8rem;
}

.comment-content {
  line-height: 1.6;
  margin-bottom: 0.8rem;
  color: #333;
  margin-left: 50px;
}

.comment-actions {
  display: flex;
  gap: 1rem;
  margin-left: 50px;
}

.comment-action-btn {
  background: none;
  border: none;
  cursor: pointer;
  font-size: 1rem;
  transition: color 0.3s;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.3rem 0;
}

.comment-action-btn:hover {
  color: #1890ff;
}

.like-icon {
  font-size: 1.1rem;
  transition: all 0.2s;
}

.like-icon.liked {
  color: #ff4d4f;
  transform: scale(1.1);
}

.like-count {
  font-size: 0.85rem;
  color: #666;
}

.no-comments {
  text-align: center;
  color: #999;
  padding: 2rem 0;
}

/* 响应式设计 */
@media (max-width: 1024px) {
  .main-container {
    flex-direction: column;
  }
  
  .left-sidebar {
    width: 100%;
  }
  
  .user-info-card {
    position: static;
  }
  
  .container {
    margin: 0;
    padding: 1.5rem;
  }
  
  .post-title {
    font-size: 1.8rem;
  }
  
  .right-toolbar {
    position: static;
    flex-direction: row;
    justify-content: center;
    min-width: auto;
  }
  
  .toolbar-btn {
    width: 80px;
  }
  
  .comments-sidebar {
    width: 100%;
  }
}

@media (max-width: 768px) {
  .post-title {
    font-size: 1.5rem;
  }
  
  .user-info-card {
    padding: 1.5rem;
  }
  
  .user-stats {
    gap: 1rem;
  }
  
  .achievement-item {
    font-size: 0.8rem;
  }
}

/* 个人主页弹窗样式 */
.profile-modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 2000;
}

.profile-modal-content {
  background-color: white;
  border-radius: 8px;
  width: 95%;
  max-width: 1400px;
  height: 90vh;
  overflow: hidden;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  display: flex;
  flex-direction: column;
}

.profile-modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.5rem 2rem;
  border-bottom: 1px solid #e8e8e8;
  flex-shrink: 0;
}

.profile-modal-header h3 {
  font-size: 1.2rem;
  font-weight: 500;
  margin: 0;
  color: #333;
}

.close-btn {
  background-color: transparent;
  border: none;
  font-size: 2rem;
  color: #999;
  cursor: pointer;
  padding: 0;
  width: 32px;
  height: 32px;
  line-height: 1;
  transition: all 0.3s;
}

.close-btn:hover {
  color: #333;
}

.profile-modal-body {
  padding: 0;
  display: flex;
  flex: 1;
  overflow: hidden;
}

.profile-left {
  width: 350px;
  padding: 2rem;
  border-right: 1px solid #e8e8e8;
  overflow-y: auto;
  flex-shrink: 0;
}

.profile-right {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.profile-tabs {
  display: flex;
  padding: 1rem 2rem;
  border-bottom: 1px solid #e8e8e8;
  gap: 1rem;
  flex-shrink: 0;
}

.tab-btn {
  padding: 0.5rem 1.5rem;
  background-color: transparent;
  color: #666;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 1rem;
  transition: all 0.3s;
}

.tab-btn:hover {
  background-color: #f0f0f0;
}

.tab-btn.active {
  background-color: #1890ff;
  color: white;
}

.profile-content {
  flex: 1;
  padding: 2rem;
  overflow-y: auto;
}

.user-info-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-bottom: 2rem;
  padding-bottom: 1.5rem;
  border-bottom: 1px solid #e8e8e8;
}

.user-actions {
  display: flex;
  gap: 10px;
  margin-top: 1rem;
  justify-content: center;
}

/* 确保两个按钮样式完全一致 */
.follow-btn,
.profile-btn {
  padding: 0.5rem 1.5rem;
  background-color: #1890ff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 1rem;
  transition: all 0.3s;
  min-width: 120px;
  text-align: center;
}

/* 调整个人主页弹窗的内部样式 */
.profile-modal-body {
  padding: 2rem;
}

.user-basic-info h4 {
  font-size: 1.2rem;
  margin-bottom: 0.5rem;
}

.user-achievements h4,
.user-posts-section h4 {
  font-size: 1.1rem;
  margin-bottom: 1rem;
}

.achievement-item {
  font-size: 1rem;
  margin-bottom: 0.5rem;
}

.posts-list .post-item {
  padding: 16px;
  margin-bottom: 12px;
}

.posts-list .post-title {
  font-size: 16px;
  margin-bottom: 8px;
}

.posts-list .post-meta {
  font-size: 12px;
  margin-bottom: 8px;
}

.posts-list .post-stats {
  font-size: 12px;
  margin-bottom: 10px;
}

.posts-list .view-btn {
  padding: 6px 12px;
  font-size: 12px;
}

.follow-btn:hover,
.profile-btn:hover {
  background-color: #40a9ff;
}

.follow-btn.following {
  background-color: #52c41a;
}

.follow-btn.following:hover {
  background-color: #73d13d;
}

.user-posts-section {
  margin-top: 2rem;
}

.user-posts-section h4 {
  font-size: 1rem;
  font-weight: 500;
  margin: 0 0 1rem;
  color: #333;
}

.posts-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.posts-list .post-item {
  padding: 12px;
  border: 1px solid #e8e8e8;
  border-radius: 6px;
  transition: all 0.3s;
}

.posts-list .post-item:hover {
  border-color: #1890ff;
  box-shadow: 0 2px 8px rgba(24, 144, 255, 0.1);
}

.posts-list .post-title {
  font-size: 14px;
  font-weight: 500;
  color: #333;
  margin: 0 0 6px;
}

.posts-list .post-meta {
  font-size: 11px;
  color: #999;
  margin-bottom: 6px;
  display: flex;
  gap: 12px;
}

.posts-list .post-stats {
  font-size: 11px;
  color: #666;
  margin-bottom: 8px;
  display: flex;
  gap: 12px;
}

.posts-list .view-btn {
  padding: 4px 8px;
  background-color: #1890ff;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 11px;
  cursor: pointer;
  transition: all 0.3s;
}

.posts-list .view-btn:hover {
  background-color: #40a9ff;
}

.empty-state {
  text-align: center;
  padding: 2rem 0;
  color: #999;
  font-size: 0.9rem;
}

/* 响应式设计 */
@media (max-width: 1200px) {
  .profile-modal-content {
    width: 95%;
    max-width: 1000px;
  }
  
  .profile-left {
    width: 300px;
  }
}

@media (max-width: 768px) {
  .profile-modal-content {
    width: 98%;
    height: 95vh;
  }
  
  .profile-modal-body {
    flex-direction: column;
  }
  
  .profile-left {
    width: 100%;
    border-right: none;
    border-bottom: 1px solid #e8e8e8;
    padding: 1.5rem;
    max-height: 40vh;
  }
  
  .profile-tabs {
    padding: 0.75rem 1.5rem;
  }
  
  .profile-content {
    padding: 1.5rem;
  }
  
  .user-actions {
    flex-direction: column;
    width: 100%;
  }
  
  .follow-btn,
  .profile-btn {
    width: 100%;
    text-align: center;
  }
}
</style>