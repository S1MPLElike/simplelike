<template>
  <div class="chat-container">
    <!-- 左侧好友列表 -->
    <div class="friend-list">
      <div class="search-header">
        <div class="search-box">
          <input type="text" placeholder="搜索" v-model="searchKeyword" />
          <button class="add-btn">+</button>
        </div>
      </div>
      <div class="friend-items">
        <div 
          v-for="friend in filteredFriends" 
          :key="friend.id"
          class="friend-item"
          :class="{ active: selectedFriend?.id === friend.id }"
          @click="selectFriend(friend)"
        >
          <img :src="friend.avatar && friend.avatar !== '' ? `http://localhost:8080${friend.avatar}` : defaultAvatar" class="avatar" alt="头像" />
          <div class="friend-info">
            <div class="friend-name">{{ friend.name }}</div>
            <div class="last-message" v-if="friend.lastMessage">{{ friend.lastMessage }}</div>
          </div>
          <div class="friend-meta">
            <div class="message-time">{{ friend.time || '' }}</div>
            <div v-if="friend.unread > 0" class="unread-badge">{{ friend.unread }}</div>
          </div>
        </div>
      </div>
    </div>

    <!-- 右侧聊天窗口 -->
    <div class="chat-window" v-if="selectedFriend">
      <div class="chat-header">
        <div class="chat-title">{{ selectedFriend.name }}</div>
        <div class="chat-actions">
          <button class="action-btn" @click="openFriendProfile">👤</button>
          <button class="action-btn">⋯</button>
        </div>
      </div>

      <div class="chat-messages" ref="messageContainer">
        <div 
          v-for="msg in messages" 
          :key="msg.id"
          class="message-item"
          :class="{ self: msg.isSelf }"
        >
          <img 
            :src="msg.isSelf ? userAvatar : (selectedFriend && selectedFriend.avatar && selectedFriend.avatar !== '' ? `http://localhost:8080${selectedFriend.avatar}` : defaultAvatar)" 
            class="msg-avatar" 
            alt="头像" 
          />
          <div class="message-content">
            <div class="message-bubble" v-if="msg.type === 'text' || msg.type === 'emoji'">{{ msg.content }}</div>
            <div class="message-bubble" v-else-if="msg.type === 'image'">
              <img :src="`http://localhost:8080${msg.content}`" class="message-image" alt="图片" />
            </div>
            <div class="message-time">{{ formatTime(msg.create_time) }}</div>
          </div>
        </div>
      </div>

      <div class="chat-input-area">
        <div class="input-toolbar">
          <button class="toolbar-btn" @click="toggleEmojiPicker">😊</button>
          <button class="toolbar-btn" @click="triggerImageUpload">📎</button>
        </div>
        <textarea 
          v-model="inputMessage" 
          class="message-input" 
          placeholder="请输入消息"
          @keyup.enter.prevent="sendMessage('text', inputMessage)"
        ></textarea>
        <input 
          type="file" 
          ref="imageInput" 
          accept="image/*" 
          style="display: none" 
          @change="handleImageUpload"
        />
        <div class="input-actions">
          <button class="send-btn" @click="sendMessage('text', inputMessage)" :disabled="!inputMessage.trim()">发送</button>
        </div>
      </div>

      <!-- Emoji选择器 -->
      <div class="emoji-picker" v-if="showEmojiPicker">
        <div class="emoji-list">
          <span 
            v-for="emoji in emojis" 
            :key="emoji"
            class="emoji-item"
            @click="selectEmoji(emoji)"
          >{{ emoji }}</span>
        </div>
      </div>
    </div>

    <!-- 未选择好友时的默认显示 -->
    <div class="chat-window empty" v-else>
      <div class="empty-tip">选择一个好友开始聊天</div>
    </div>

    <!-- 好友个人主页弹窗 -->
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
import { ref, computed, onMounted, nextTick, watch, onUnmounted } from 'vue';
import defaultAvatar from '../assets/login_1.jpg';

const searchKeyword = ref('');
const selectedFriend = ref(null);
const inputMessage = ref('');
const messageContainer = ref(null);
const friends = ref([]);
const messages = ref([]);
const showEmojiPicker = ref(false);
const imageInput = ref(null);
const userID = ref('');
const userAvatar = ref('');
const ws = ref(null);
const wsConnected = ref(false);

// 好友个人主页相关
const showUserProfile = ref(false);
const userProfileInfo = ref({});
const activeProfileTab = ref('posts');
const userPosts = ref([]);
const userCollectedPosts = ref([]);
const isFollowing = ref(false);
const showFollowButtonInProfile = ref(true);

// 常用emoji
const emojis = [
  '😊', '😂', '😍', '🤔', '😢', '😡', '👍', '👎', '❤️', '🎉',
  '🔥', '🤣', '😭', '😱', '🤗', '🤩', '😎', '🤬', '🤭', '🥳'
];

const filteredFriends = computed(() => {
  if (!searchKeyword.value) return friends.value;
  return friends.value.filter(f => 
    f.name.toLowerCase().includes(searchKeyword.value.toLowerCase())
  );
});

// 格式化时间
const formatTime = (timeString) => {
  if (!timeString) return '';
  const date = new Date(timeString);
  return date.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' });
};

// 切换emoji选择器
const toggleEmojiPicker = () => {
  showEmojiPicker.value = !showEmojiPicker.value;
};

// 选择emoji
const selectEmoji = (emoji) => {
  inputMessage.value += emoji;
  showEmojiPicker.value = false;
};

// 触发图片上传
const triggerImageUpload = () => {
  if (imageInput.value) {
    imageInput.value.click();
  }
};

// 处理图片上传
const handleImageUpload = async (event) => {
  const file = event.target.files[0];
  if (!file) return;

  const formData = new FormData();
  formData.append('image', file);

  try {
    const response = await fetch('http://localhost:8080/api/v1/chat/upload-image', {
      method: 'POST',
      credentials: 'include',
      body: formData
    });

    if (response.ok) {
      const data = await response.json();
      if (data.code === 0) {
        // 发送图片消息
        sendMessage('image', data.data.image_url);
      }
    }
  } catch (error) {
    console.error('Upload image failed:', error);
  } finally {
    // 重置文件输入
    if (imageInput.value) {
      imageInput.value.value = '';
    }
  }
};

// 选择好友
const selectFriend = async (friend) => {
  selectedFriend.value = friend;
  friend.unread = 0;
  // 获取消息历史
  await getMessageHistory(friend.id);
};

// 打开好友个人主页
const openFriendProfile = async () => {
  if (selectedFriend.value) {
    console.log('打开好友个人主页:', selectedFriend.value.name);
    
    try {
      await getUserProfileInfo(selectedFriend.value.id);
      await getUserPosts(selectedFriend.value.id);
      await getUserCollectedPosts(selectedFriend.value.id);
      showUserProfile.value = true;
      console.log('User profile opened successfully');
    } catch (error) {
      console.error('Open user profile failed:', error);
    }
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
          user_id: data.data.id,
          username: data.data.username,
          avatar: data.data.avatar,
          followers: data.data.followers,
          following: data.data.following,
          likes: data.data.like_count,
          comments: data.data.comment_count,
          collections: data.data.collect_count
        };
        
        // 检查是否已关注
        await checkFollowStatus(userId);
        
        // 检查是否是自己
        const currentUserID = localStorage.getItem('user_id');
        showFollowButtonInProfile.value = currentUserID !== String(userId);
      }
    }
  } catch (error) {
    console.error('获取用户信息失败:', error);
  }
};

// 获取用户帖子
const getUserPosts = async (userId) => {
  try {
    const response = await fetch(`http://localhost:8080/api/v1/post/user/${userId}`, {
      credentials: 'include'
    });
    
    if (response.ok) {
      const data = await response.json();
      if (data.code === 0) {
        userPosts.value = data.data;
      }
    }
  } catch (error) {
    console.error('获取用户帖子失败:', error);
  }
};

// 获取用户收藏的帖子
const getUserCollectedPosts = async (userId) => {
  try {
    const response = await fetch(`http://localhost:8080/api/v1/post/collected/${userId}`, {
      credentials: 'include'
    });
    
    if (response.ok) {
      const data = await response.json();
      if (data.code === 0) {
        userCollectedPosts.value = data.data;
      }
    }
  } catch (error) {
    console.error('获取用户收藏失败:', error);
  }
};

// 检查关注状态
const checkFollowStatus = async (userId) => {
  try {
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
    console.error('检查关注状态失败:', error);
  }
};

// 关注/取消关注
const toggleFollow = async () => {
  if (!userProfileInfo.value.user_id) return;
  
  try {
    const url = isFollowing.value 
      ? `http://localhost:8080/api/v1/user/unfollow/${userProfileInfo.value.user_id}`
      : `http://localhost:8080/api/v1/user/follow/${userProfileInfo.value.user_id}`;
    
    const response = await fetch(url, {
      method: 'POST',
      credentials: 'include'
    });
    
    if (response.ok) {
      const data = await response.json();
      if (data.code === 0) {
        isFollowing.value = !isFollowing.value;
        // 更新关注数
        if (isFollowing.value) {
          userProfileInfo.value.followers++;
        } else {
          userProfileInfo.value.followers--;
        }
      }
    }
  } catch (error) {
    console.error('关注/取消关注失败:', error);
  }
};

// 发送消息
const sendMessage = async (messageType = 'text', content = inputMessage.value) => {
  if (!selectedFriend.value) return;
  if (messageType === 'text' && !content.trim()) return;

  try {
    // 确保所有参数类型正确
    const finalMessageType = typeof messageType === 'string' ? messageType : 'text';
    const receiverId = Number(selectedFriend.value.id);
    const messageContent = String(content);
    
    const messageData = {
      receiver_id: receiverId,
      content: messageContent,
      type: finalMessageType
    };
    
    console.log('Sending message:', messageData);
    
    // 尝试使用WebSocket发送消息
    const wsSent = sendWSMessage('chat', messageData);
    
    if (!wsSent) {
      // WebSocket未连接，使用HTTP请求作为备用
      console.log('WebSocket未连接，使用HTTP请求发送消息');
      const response = await fetch('http://localhost:8080/api/v1/chat/send', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        credentials: 'include',
        body: JSON.stringify(messageData)
      });

      if (response.ok) {
        const data = await response.json();
        console.log('Send message response:', data);
        if (data.code === 0) {
          // 添加到消息列表
          const newMessage = {
            id: data.data.message_id,
            content: messageContent,
            type: finalMessageType,
            isSelf: true,
            create_time: data.data.create_time
          };
          messages.value.push(newMessage);
          console.log('Message added to array:', newMessage);
          console.log('Messages array length:', messages.value.length);

          // 清空输入
          if (finalMessageType === 'text') {
            inputMessage.value = '';
          }

          // 滚动到底部
          nextTick(() => {
            if (messageContainer.value) {
              messageContainer.value.scrollTop = messageContainer.value.scrollHeight;
              console.log('Scrolled to bottom');
            }
          });
        }
      }
    } else {
      // WebSocket发送成功，立即添加消息到列表
      console.log('WebSocket消息发送成功');
      
      // 添加到消息列表
      const newMessage = {
        id: Date.now(), // 临时ID，服务器会返回真实ID
        content: messageContent,
        type: finalMessageType,
        isSelf: true,
        create_time: new Date().toISOString().replace('T', ' ').substring(0, 19)
      };
      messages.value.push(newMessage);
      console.log('Message added to array:', newMessage);
      console.log('Messages array length:', messages.value.length);
      
      // 清空输入
      if (finalMessageType === 'text') {
        inputMessage.value = '';
      }
      
      // 滚动到底部
      nextTick(() => {
        if (messageContainer.value) {
          messageContainer.value.scrollTop = messageContainer.value.scrollHeight;
          console.log('Scrolled to bottom');
        }
      });
    }
  } catch (error) {
    console.error('Send message failed:', error);
  }
};

// 获取消息历史
const getMessageHistory = async (friendID) => {
  try {
    const response = await fetch(`http://localhost:8080/api/v1/chat/history/${friendID}`, {
      credentials: 'include'
    });

    if (response.ok) {
      const data = await response.json();
      if (data.code === 0) {
        // 转换消息格式
        messages.value = data.data.map(msg => ({
          id: msg.id,
          content: msg.content,
          type: msg.type,
          isSelf: msg.sender_id === userID.value,
          create_time: msg.create_time
        })).reverse(); // 反转顺序，使最新消息在底部

        // 重置对应好友的未读消息数量
        friends.value = friends.value.map(friend => {
          if (friend.id === friendID) {
            return {
              ...friend,
              unread: 0
            };
          }
          return friend;
        });

        // 滚动到底部
        nextTick(() => {
          if (messageContainer.value) {
            messageContainer.value.scrollTop = messageContainer.value.scrollHeight;
          }
        });
      }
    }
  } catch (error) {
    console.error('Get message history failed:', error);
  }
};

// 获取好友列表
const getFriends = async () => {
  try {
    // 首先获取当前登录用户的ID
    const loginResponse = await fetch('http://localhost:8080/api/v1/user/auto-login', {
      credentials: 'include'
    });
    
    if (loginResponse.ok) {
      const loginData = await loginResponse.json();
      if (loginData.code === 0) {
        userID.value = loginData.data.user_id;
        
        // 建立WebSocket连接
        connectWebSocket();
        
        // 获取当前用户信息
        const userInfoResponse = await fetch('http://localhost:8080/api/v1/user/info', {
          credentials: 'include'
        });
        
        if (userInfoResponse.ok) {
          const userInfoData = await userInfoResponse.json();
          if (userInfoData.code === 0) {
            userAvatar.value = userInfoData.data.avatar ? `http://localhost:8080${userInfoData.data.avatar}` : defaultAvatar;
          }
        }
        
        // 然后获取好友列表
        const friendsResponse = await fetch(`http://localhost:8080/api/v1/user/friends/${userID.value}`, {
          credentials: 'include'
        });
        
        if (friendsResponse.ok) {
          const friendsData = await friendsResponse.json();
          if (friendsData.code === 0) {
            friends.value = friendsData.data.friends.map(friend => ({
              ...friend,
              lastMessage: '',
              time: '',
              unread: 0
            }));
            
            // 获取按好友分组的未读消息数量
            const unreadResponse = await fetch('http://localhost:8080/api/v1/chat/unread-count-by-friend', {
              credentials: 'include'
            });
            
            if (unreadResponse.ok) {
              const unreadData = await unreadResponse.json();
              if (unreadData.code === 0) {
                const unreadCounts = unreadData.data.unread_counts;
                // 更新好友列表中的未读消息数量
                friends.value = friends.value.map(friend => ({
                  ...friend,
                  unread: unreadCounts[friend.id] || 0
                }));
              }
            }
            
            // 默认选择第一个好友
            if (friends.value.length > 0) {
              selectedFriend.value = friends.value[0];
              // 获取消息历史
              await getMessageHistory(friends.value[0].id);
            }
          }
        }
      }
    }
  } catch (error) {
    console.error('Get friends failed:', error);
  }
};

// 监听选中好友变化
watch(selectedFriend, (newFriend) => {
  if (newFriend) {
    getMessageHistory(newFriend.id);
  }
});

// 建立WebSocket连接
const connectWebSocket = () => {
  if (!userID.value) return;
  
  // 关闭现有的连接
  if (ws.value) {
    ws.value.close();
  }
  
  // 建立新的WebSocket连接
  const wsUrl = `ws://localhost:8080/api/v1/chat/ws`;
  ws.value = new WebSocket(wsUrl);
  // WebSocket连接会自动携带cookie，因为域名相同（都是localhost）
  
  ws.value.onopen = () => {
    console.log('WebSocket连接成功');
    wsConnected.value = true;
  };
  
  ws.value.onmessage = (event) => {
    handleWSMessage(event.data);
  };
  
  ws.value.onerror = (error) => {
    console.error('WebSocket错误:', error);
    wsConnected.value = false;
  };
  
  ws.value.onclose = () => {
    console.log('WebSocket连接关闭');
    wsConnected.value = false;
    // 尝试重连
    setTimeout(connectWebSocket, 3000);
  };
};

// 处理WebSocket消息
const handleWSMessage = (message) => {
  try {
    const wsMsg = JSON.parse(message);
    console.log('收到WebSocket消息:', wsMsg);
    
    switch (wsMsg.type) {
      case 'chat':
        handleChatMessage(wsMsg.payload);
        break;
      case 'status':
        handleStatusMessage(wsMsg.payload);
        break;
      case 'error':
        handleErrorMessage(wsMsg.payload);
        break;
      default:
        console.log('未知消息类型:', wsMsg.type);
    }
  } catch (error) {
    console.error('解析WebSocket消息失败:', error);
  }
};

// 处理聊天消息
const handleChatMessage = (message) => {
  console.log('收到聊天消息:', message);
  const senderId = Number(message.sender_id);
  const receiverId = Number(message.receiver_id);
  const currentUserId = Number(userID.value);
  const isSelf = message.is_self;
  
  console.log('消息信息:', {
    senderId,
    receiverId,
    currentUserId,
    isSelf,
    selectedFriend: selectedFriend.value ? selectedFriend.value.id : null
  });
  
  // 更新好友列表中的未读消息数量（如果不是当前聊天窗口）
  if (message.sender_id && !isSelf) {
    const isCurrentChat = selectedFriend.value && senderId === Number(selectedFriend.value.id);
    if (!isCurrentChat) {
      friends.value = friends.value.map(friend => {
        if (Number(friend.id) === senderId) {
          return {
            ...friend,
            unread: (friend.unread || 0) + 1
          };
        }
        return friend;
      });
    }
  }
  
  // 检查消息是否属于当前聊天
  if (selectedFriend.value) {
    const selectedId = Number(selectedFriend.value.id);
    console.log('检查消息是否属于当前聊天:', {
      senderId,
      selectedId,
      receiverId,
      currentUserId,
      senderMatch: senderId === selectedId,
      receiverMatch: receiverId === selectedId,
      selfMatch: senderId === currentUserId && receiverId === selectedId
    });
    
    // 消息属于当前聊天的情况：
    // 1. 消息是从当前选中的好友发送的（senderId === selectedId）
    // 2. 消息是发送给当前选中的好友的（receiverId === selectedId）
    if (senderId === selectedId || receiverId === selectedId) {
      // 添加到消息列表
      const newMessage = {
        id: message.id,
        content: message.content,
        type: message.type,
        isSelf: isSelf,
        create_time: message.create_time
      };
      console.log('添加消息到列表:', newMessage);
      // 使用展开运算符创建新数组，触发Vue的响应式更新
      messages.value = [...messages.value, newMessage];
      console.log('消息列表长度:', messages.value.length);
      
      // 滚动到底部
      nextTick(() => {
        if (messageContainer.value) {
          console.log('滚动前的scrollTop:', messageContainer.value.scrollTop);
          console.log('滚动前的scrollHeight:', messageContainer.value.scrollHeight);
          messageContainer.value.scrollTop = messageContainer.value.scrollHeight;
          console.log('滚动后的scrollTop:', messageContainer.value.scrollTop);
          console.log('滚动到底部');
        } else {
          console.error('messageContainer is null');
        }
      });
    } else {
      console.log('消息不属于当前聊天，忽略');
    }
  } else {
    console.log('没有选中好友，忽略消息');
  }
};

// 处理状态消息
const handleStatusMessage = (status) => {
  console.log('状态消息:', status);
  // 可以在这里更新好友的在线状态
};

// 处理错误消息
const handleErrorMessage = (error) => {
  console.error('错误消息:', error);
  // 可以在这里显示错误提示
};

// 发送WebSocket消息
const sendWSMessage = (type, payload) => {
  if (!ws.value || ws.value.readyState !== WebSocket.OPEN) {
    console.error('WebSocket未连接');
    return false;
  }
  
  const message = {
    type: type,
    payload: payload
  };
  
  try {
    ws.value.send(JSON.stringify(message));
    return true;
  } catch (error) {
    console.error('发送WebSocket消息失败:', error);
    return false;
  }
};

onMounted(() => {
  getFriends().then(() => {
    // 等获取用户ID后再连接WebSocket
    if (userID.value) {
      connectWebSocket();
    }
  });
});

onUnmounted(() => {
  if (ws.value) {
    ws.value.close();
  }
});
</script>

<style scoped>
.chat-container {
  display: flex;
  height: 100%;
  background-color: white;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.09);
}

/* 左侧好友列表 */
.friend-list {
  width: 280px;
  border-right: 1px solid #e8e8e8;
  display: flex;
  flex-direction: column;
  background-color: #fafafa;
}

.search-header {
  padding: 12px;
  border-bottom: 1px solid #e8e8e8;
}

.search-box {
  display: flex;
  gap: 8px;
  align-items: center;
}

.search-box input {
  flex: 1;
  padding: 8px 12px;
  border: 1px solid #d9d9d9;
  border-radius: 4px;
  font-size: 14px;
  background-color: white;
}

.search-box input:focus {
  outline: none;
  border-color: #1890ff;
}

.add-btn {
  width: 32px;
  height: 32px;
  border: none;
  background-color: #e8e8e8;
  border-radius: 4px;
  font-size: 18px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}

.add-btn:hover {
  background-color: #d9d9d9;
}

.friend-items {
  flex: 1;
  overflow-y: auto;
}

.friend-item {
  display: flex;
  align-items: center;
  padding: 12px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.friend-item:hover {
  background-color: #f0f0f0;
}

.friend-item.active {
  background-color: #c8e6c9;
}

.avatar {
  width: 40px;
  height: 40px;
  border-radius: 4px;
  background-color: #ddd;
  margin-right: 12px;
  object-fit: cover;
}

.friend-info {
  flex: 1;
  min-width: 0;
}

.friend-name {
  font-size: 14px;
  font-weight: 500;
  color: #333;
  margin-bottom: 4px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.last-message {
  font-size: 12px;
  color: #999;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.friend-meta {
  text-align: right;
  min-width: 50px;
}

.message-time {
  font-size: 12px;
  color: #999;
  margin-bottom: 4px;
}

.unread-badge {
  display: inline-block;
  background-color: #ff4d4f;
  color: white;
  font-size: 11px;
  width: 18px;
  height: 18px;
  line-height: 18px;
  border-radius: 50%;
  text-align: center;
  margin-left: 8px;
}

/* 右侧聊天窗口 */
.chat-window {
  flex: 1;
  display: flex;
  flex-direction: column;
  background-color: #f5f5f5;
}

.chat-window.empty {
  align-items: center;
  justify-content: center;
}

.empty-tip {
  color: #999;
  font-size: 14px;
}

.chat-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 20px;
  background-color: #f5f5f5;
  border-bottom: 1px solid #e8e8e8;
}

.chat-title {
  font-size: 16px;
  font-weight: 500;
  color: #333;
}

.chat-actions {
  display: flex;
  gap: 12px;
}

.action-btn {
  background: none;
  border: none;
  font-size: 18px;
  cursor: pointer;
  padding: 4px;
}

.chat-messages {
  flex: 1;
  padding: 20px;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.message-item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
}

.message-item.self {
  flex-direction: row-reverse;
}

.msg-avatar {
  width: 36px;
  height: 36px;
  border-radius: 4px;
  background-color: #ddd;
  object-fit: cover;
}

.message-content {
  max-width: 60%;
}

.message-bubble {
  padding: 10px 14px;
  border-radius: 4px;
  font-size: 14px;
  line-height: 1.5;
  word-break: break-word;
}

.message-item:not(.self) .message-bubble {
  background-color: white;
  color: #333;
}

.message-item.self .message-bubble {
  background-color: #95ec69;
  color: #333;
}

.message-time {
  font-size: 11px;
  color: #999;
  margin-top: 4px;
  text-align: right;
}

.message-item.self .message-time {
  text-align: left;
}

.message-image {
  max-width: 200px;
  max-height: 200px;
  border-radius: 4px;
  object-fit: cover;
}

/* Emoji选择器 */
.emoji-picker {
  position: absolute;
  bottom: 140px;
  left: 20px;
  background-color: white;
  border: 1px solid #e8e8e8;
  border-radius: 8px;
  padding: 10px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  z-index: 100;
}

.emoji-list {
  display: grid;
  grid-template-columns: repeat(10, 1fr);
  gap: 10px;
}

.emoji-item {
  font-size: 20px;
  cursor: pointer;
  text-align: center;
  padding: 5px;
  border-radius: 4px;
  transition: background-color 0.2s;
}

.emoji-item:hover {
  background-color: #f0f0f0;
}

.send-btn:disabled {
  background-color: #e8e8e8;
  color: #999;
  cursor: not-allowed;
}

.send-btn:not(:disabled) {
  background-color: #1890ff;
  color: white;
}

.send-btn:not(:disabled):hover {
  background-color: #40a9ff;
}

/* 输入区域 */
.chat-input-area {
  background-color: #f5f5f5;
  border-top: 1px solid #e8e8e8;
  padding: 12px 20px;
}

.input-toolbar {
  display: flex;
  gap: 16px;
  margin-bottom: 8px;
}

.toolbar-btn {
  background: none;
  border: none;
  font-size: 20px;
  cursor: pointer;
  padding: 4px;
}

.message-input {
  width: 100%;
  min-height: 80px;
  padding: 10px;
  border: none;
  background-color: white;
  border-radius: 4px;
  font-size: 14px;
  resize: none;
  outline: none;
}

.input-actions {
  display: flex;
  justify-content: flex-end;
  margin-top: 8px;
}

.send-btn {
  padding: 8px 24px;
  background-color: #e8e8e8;
  border: none;
  border-radius: 4px;
  font-size: 14px;
  color: #999;
  cursor: pointer;
  transition: all 0.3s;
}

.send-btn:hover {
  background-color: #d9d9d9;
}

.send-btn:active {
  background-color: #c8c8c8;
}

/* 滚动条样式 */
.friend-items::-webkit-scrollbar,
.chat-messages::-webkit-scrollbar {
  width: 6px;
}

.friend-items::-webkit-scrollbar-track,
.chat-messages::-webkit-scrollbar-track {
  background: transparent;
}

.friend-items::-webkit-scrollbar-thumb,
.chat-messages::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 3px;
}

.friend-items::-webkit-scrollbar-thumb:hover,
.chat-messages::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8;
}
  /* 个人主页弹窗 */
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
    z-index: 1000;
  }
  
  .profile-modal-content {
    background-color: white;
    border-radius: 8px;
    width: 90%;
    max-width: 1000px;
    max-height: 90vh;
    overflow-y: auto;
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.15);
  }
  
  .profile-modal-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 20px;
    border-bottom: 1px solid #e8e8e8;
  }
  
  .profile-modal-header h3 {
    margin: 0;
    font-size: 18px;
    font-weight: 500;
    color: #333;
  }
  
  .close-btn {
    background: none;
    border: none;
    font-size: 24px;
    cursor: pointer;
    color: #999;
  }
  
  .close-btn:hover {
    color: #333;
  }
  
  .profile-modal-body {
    display: flex;
    padding: 30px;
  }
  
  .profile-left {
    flex: 1;
    margin-right: 40px;
  }
  
  .profile-right {
    flex: 2;
  }
  
  .user-info-section {
    margin-bottom: 30px;
  }
  
  .user-avatar {
    width: 120px;
    height: 120px;
    border-radius: 50%;
    overflow: hidden;
    margin-bottom: 20px;
  }
  
  .user-avatar img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
  
  .user-basic-info h4 {
    margin: 0 0 10px 0;
    font-size: 20px;
    font-weight: 500;
  }
  
  .user-stats {
    margin-bottom: 15px;
  }
  
  .stat {
    margin-right: 15px;
    font-size: 14px;
    color: #666;
  }
  
  .user-actions {
    margin-bottom: 20px;
  }
  
  .follow-btn {
    padding: 8px 20px;
    background-color: #1890ff;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 14px;
  }
  
  .follow-btn:hover {
    background-color: #40a9ff;
  }
  
  .follow-btn.following {
    background-color: #e8e8e8;
    color: #333;
  }
  
  .follow-btn.following:hover {
    background-color: #d9d9d9;
  }
  
  .user-achievements {
    background-color: #f9f9f9;
    padding: 20px;
    border-radius: 8px;
  }
  
  .user-achievements h4 {
    margin: 0 0 15px 0;
    font-size: 16px;
    font-weight: 500;
  }
  
  .achievement-item {
    display: flex;
    align-items: center;
    margin-bottom: 10px;
  }
  
  .achievement-icon {
    font-size: 16px;
    margin-right: 10px;
  }
  
  .achievement-text {
    font-size: 14px;
    color: #666;
  }
  
  .profile-tabs {
    display: flex;
    border-bottom: 1px solid #e8e8e8;
    margin-bottom: 20px;
  }
  
  .tab-btn {
    padding: 10px 20px;
    background: none;
    border: none;
    cursor: pointer;
    font-size: 14px;
    color: #666;
    border-bottom: 2px solid transparent;
  }
  
  .tab-btn:hover {
    color: #1890ff;
  }
  
  .tab-btn.active {
    color: #1890ff;
    border-bottom-color: #1890ff;
  }
  
  .profile-content {
    min-height: 300px;
  }
  
  .user-posts-section h4,
  .user-collected-section h4 {
    margin: 0 0 20px 0;
    font-size: 16px;
    font-weight: 500;
  }
  
  .empty-state {
    text-align: center;
    color: #999;
    padding: 40px 0;
  }
  
  .posts-list {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
    gap: 20px;
  }
  
  .post-item {
    padding: 20px;
    border: 1px solid #e8e8e8;
    border-radius: 8px;
    background-color: #f9f9f9;
  }
  
  .post-item h5 {
    margin: 0 0 10px 0;
    font-size: 14px;
    font-weight: 500;
  }
  
  .post-meta {
    margin-bottom: 10px;
    font-size: 12px;
    color: #999;
  }
  
  .post-time {
    margin-right: 15px;
  }
  
  .post-stats {
    margin-bottom: 15px;
    font-size: 12px;
    color: #666;
  }
  
  .stat-item {
    margin-right: 15px;
  }
  
  .view-btn {
    padding: 6px 12px;
    background-color: #f0f0f0;
    color: #333;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 12px;
  }
  
  .view-btn:hover {
    background-color: #e0e0e0;
  }
  
  /* 响应式设计 */
  @media (max-width: 768px) {
    .profile-modal-body {
      flex-direction: column;
    }
    
    .profile-left {
      margin-right: 0;
      margin-bottom: 30px;
    }
  }
</style>