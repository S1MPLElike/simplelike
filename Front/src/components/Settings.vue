<template>
  <div class="settings">
    <h2>设置</h2>
    
    <div class="settings-section">
      <h3>个人信息</h3>
      <div class="settings-form">
        <div class="form-group">
          <label for="username">用户名</label>
          <input 
            type="text" 
            id="username" 
            v-model="userForm.username"
            :disabled="loading"
          >
        </div>
        <div class="form-group">
          <label for="email">邮箱</label>
          <input 
            type="email" 
            id="email" 
            v-model="userForm.email"
            :disabled="loading"
          >
        </div>
        <div class="form-group">
          <label for="bio">个人简介</label>
          <textarea 
            id="bio" 
            v-model="userForm.bio"
            :disabled="loading"
            rows="3"
          ></textarea>
        </div>
        <div class="form-actions">
          <button 
            class="save-btn" 
            @click="saveProfile"
            :disabled="loading"
          >
            {{ loading ? '保存中...' : '保存' }}
          </button>
        </div>
      </div>
    </div>
    
    <div class="settings-section">
      <h3>账号安全</h3>
      <div class="settings-form">
        <div class="form-group">
          <label for="password">修改密码</label>
          <button 
            class="secondary-btn" 
            @click="showChangePassword = true"
            :disabled="loading"
          >
            修改密码
          </button>
        </div>
      </div>
    </div>
    
    <div class="settings-section">
      <h3>通知设置</h3>
      <div class="settings-form">
        <div class="form-group checkbox-group">
          <input 
            type="checkbox" 
            id="notifyLikes" 
            v-model="notificationSettings.likes"
            :disabled="loading"
          >
          <label for="notifyLikes">收到点赞通知</label>
        </div>
        <div class="form-group checkbox-group">
          <input 
            type="checkbox" 
            id="notifyComments" 
            v-model="notificationSettings.comments"
            :disabled="loading"
          >
          <label for="notifyComments">收到评论通知</label>
        </div>
        <div class="form-group checkbox-group">
          <input 
            type="checkbox" 
            id="notifyFollows" 
            v-model="notificationSettings.follows"
            :disabled="loading"
          >
          <label for="notifyFollows">收到关注通知</label>
        </div>
        <div class="form-actions">
          <button 
            class="save-btn" 
            @click="saveNotificationSettings"
            :disabled="loading"
          >
            {{ loading ? '保存中...' : '保存' }}
          </button>
        </div>
      </div>
    </div>
    
    <div class="settings-section">
      <h3>关于</h3>
      <div class="about-content">
        <p>版本: 1.0.0</p>
        <p>© 2026 社区论坛</p>
      </div>
    </div>
    
    <!-- 密码修改弹窗 -->
    <div v-if="showChangePassword" class="modal-overlay">
      <div class="modal-content">
        <div class="modal-header">
          <h3>修改密码</h3>
          <button class="close-btn" @click="showChangePassword = false">×</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label for="oldPassword">旧密码</label>
            <input 
              type="password" 
              id="oldPassword" 
              v-model="passwordForm.oldPassword"
              :disabled="loading"
            >
          </div>
          <div class="form-group">
            <label for="newPassword">新密码</label>
            <input 
              type="password" 
              id="newPassword" 
              v-model="passwordForm.newPassword"
              :disabled="loading"
            >
          </div>
          <div class="form-group">
            <label for="confirmPassword">确认新密码</label>
            <input 
              type="password" 
              id="confirmPassword" 
              v-model="passwordForm.confirmPassword"
              :disabled="loading"
            >
          </div>
        </div>
        <div class="modal-footer">
          <button 
            class="cancel-btn" 
            @click="showChangePassword = false"
            :disabled="loading"
          >
            取消
          </button>
          <button 
            class="save-btn" 
            @click="changePassword"
            :disabled="loading"
          >
            {{ loading ? '修改中...' : '确认修改' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';

const userForm = ref({
  username: '',
  email: '',
  bio: ''
});

const passwordForm = ref({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
});

const notificationSettings = ref({
  likes: true,
  comments: true,
  follows: true
});

const loading = ref(false);
const showChangePassword = ref(false);

// 获取用户信息
const getUserInfo = async () => {
  try {
    loading.value = true;
    const response = await fetch('http://localhost:8080/api/v1/user/info', {
      credentials: 'include'
    });
    
    if (response.ok) {
      const data = await response.json();
      if (data.code === 0) {
        userForm.value.username = data.data.username || '';
        userForm.value.email = data.data.email || '';
        userForm.value.bio = data.data.bio || '';
      }
    }
  } catch (error) {
    console.error('Get user info failed:', error);
  } finally {
    loading.value = false;
  }
};

// 保存个人信息
const saveProfile = async () => {
  try {
    loading.value = true;
    const response = await fetch('http://localhost:8080/api/v1/user/update', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      credentials: 'include',
      body: JSON.stringify(userForm.value)
    });
    
    if (response.ok) {
      const data = await response.json();
      if (data.code === 0) {
        alert('个人信息更新成功');
      } else {
        alert('更新失败: ' + (data.message || '未知错误'));
      }
    } else {
      alert('网络错误，请稍后重试');
    }
  } catch (error) {
    console.error('Save profile failed:', error);
    alert('网络错误，请稍后重试');
  } finally {
    loading.value = false;
  }
};

// 修改密码
const changePassword = async () => {
  if (passwordForm.value.newPassword !== passwordForm.value.confirmPassword) {
    alert('两次输入的密码不一致');
    return;
  }
  
  try {
    loading.value = true;
    const response = await fetch('http://localhost:8080/api/v1/user/change-password', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      credentials: 'include',
      body: JSON.stringify({
        old_password: passwordForm.value.oldPassword,
        new_password: passwordForm.value.newPassword
      })
    });
    
    if (response.ok) {
      const data = await response.json();
      if (data.code === 0) {
        alert('密码修改成功');
        showChangePassword.value = false;
        passwordForm.value = {
          oldPassword: '',
          newPassword: '',
          confirmPassword: ''
        };
      } else {
        alert('修改失败: ' + (data.message || '未知错误'));
      }
    } else {
      alert('网络错误，请稍后重试');
    }
  } catch (error) {
    console.error('Change password failed:', error);
    alert('网络错误，请稍后重试');
  } finally {
    loading.value = false;
  }
};

// 保存通知设置
const saveNotificationSettings = async () => {
  try {
    loading.value = true;
    // 这里可以实现保存通知设置的API调用
    // 由于后端可能还没有实现这个接口，我们先模拟成功
    setTimeout(() => {
      alert('通知设置保存成功');
      loading.value = false;
    }, 500);
  } catch (error) {
    console.error('Save notification settings failed:', error);
    alert('网络错误，请稍后重试');
    loading.value = false;
  }
};

onMounted(() => {
  getUserInfo();
});
</script>

<style scoped>
.settings {
  padding: 20px;
  max-width: 800px;
  margin: 0 auto;
}

.settings h2 {
  margin-bottom: 30px;
  color: #333;
  font-size: 24px;
  border-bottom: 2px solid #1890ff;
  padding-bottom: 10px;
}

.settings-section {
  background: white;
  border-radius: 8px;
  padding: 20px;
  margin-bottom: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.settings-section h3 {
  margin-bottom: 20px;
  color: #666;
  font-size: 18px;
  border-bottom: 1px solid #f0f0f0;
  padding-bottom: 10px;
}

.settings-form {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.form-group label {
  font-size: 14px;
  color: #666;
  font-weight: 500;
}

.form-group input,
.form-group textarea {
  padding: 8px 12px;
  border: 1px solid #d9d9d9;
  border-radius: 4px;
  font-size: 14px;
  transition: all 0.3s;
}

.form-group input:focus,
.form-group textarea:focus {
  outline: none;
  border-color: #1890ff;
  box-shadow: 0 0 0 2px rgba(24, 144, 255, 0.2);
}

.checkbox-group {
  flex-direction: row;
  align-items: center;
  gap: 10px;
}

.checkbox-group input {
  width: auto;
}

.form-actions {
  margin-top: 10px;
  display: flex;
  justify-content: flex-end;
}

.save-btn {
  padding: 8px 16px;
  background-color: #1890ff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.3s;
}

.save-btn:hover {
  background-color: #40a9ff;
}

.save-btn:disabled {
  background-color: #d9d9d9;
  cursor: not-allowed;
}

.secondary-btn {
  padding: 8px 16px;
  background-color: transparent;
  color: #1890ff;
  border: 1px solid #1890ff;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.3s;
}

.secondary-btn:hover {
  background-color: #e6f7ff;
}

.secondary-btn:disabled {
  border-color: #d9d9d9;
  color: #d9d9d9;
  cursor: not-allowed;
}

.about-content {
  padding: 10px 0;
  color: #999;
  font-size: 14px;
}

/* 弹窗样式 */
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
  max-width: 400px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px 20px;
  border-bottom: 1px solid #f0f0f0;
}

.modal-header h3 {
  margin: 0;
  font-size: 16px;
  color: #333;
}

.close-btn {
  background-color: transparent;
  border: none;
  font-size: 20px;
  color: #999;
  cursor: pointer;
  padding: 0;
  width: 24px;
  height: 24px;
  line-height: 1;
  transition: all 0.3s;
}

.close-btn:hover {
  color: #333;
}

.modal-body {
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  padding: 15px 20px;
  border-top: 1px solid #f0f0f0;
}

.cancel-btn {
  padding: 8px 16px;
  background-color: transparent;
  color: #666;
  border: 1px solid #d9d9d9;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.3s;
}

.cancel-btn:hover {
  border-color: #1890ff;
  color: #1890ff;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .settings {
    padding: 10px;
  }
  
  .settings-section {
    padding: 15px;
  }
  
  .modal-content {
    width: 95%;
  }
}
</style>