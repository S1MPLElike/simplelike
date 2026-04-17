<template>
  <div class="login-container">
    <div class="login-left">
      <div class="welcome-text">欢迎来到冲击波网站</div>
      <img src="../assets/login_1.jpg" alt="登录背景图" class="login-image" />
    </div>
    <div class="login-right">
      <div class="login-form">
        <h1>{{ isLogin ? '登录' : '注册' }}</h1>
        <form @submit.prevent="handleSubmit">
          <div class="form-group">
            <label for="username">用户名</label>
            <input 
              type="text" 
              id="username" 
              v-model="form.username" 
              required 
              placeholder="请输入用户名"
              :maxlength="20"
            />
            <div v-if="errors.username" class="error-message">{{ errors.username }}</div>
          </div>
          <div v-if="!isLogin" class="form-group">
            <label for="email">邮箱</label>
            <input 
              type="email" 
              id="email" 
              v-model="form.email" 
              required 
              placeholder="请输入邮箱"
            />
            <div v-if="errors.email" class="error-message">{{ errors.email }}</div>
          </div>
          <div v-if="!isLogin" class="form-group">
            <label for="phone">电话</label>
            <input 
              type="tel" 
              id="phone" 
              v-model="form.phone" 
              required 
              placeholder="请输入电话"
              maxlength="11"
            />
            <div v-if="errors.phone" class="error-message">{{ errors.phone }}</div>
          </div>
          <div class="form-group">
            <label for="password">密码</label>
            <input 
              type="password" 
              id="password" 
              v-model="form.password" 
              required 
              placeholder="请输入密码"
              :minlength="6"
            />
            <div v-if="errors.password" class="error-message">{{ errors.password }}</div>
          </div>
          <button type="submit" class="submit-btn">{{ isLogin ? '登录' : '注册' }}</button>
        </form>
        <div class="toggle-form">
          {{ isLogin ? '还没有账号？' : '已有账号？' }}
          <span @click="isLogin = !isLogin">{{ isLogin ? '立即注册' : '立即登录' }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';

const router = useRouter();

const isLogin = ref(true);
const form = ref({
  username: '',
  password: '',
  email: '',
  phone: ''
});

const errors = ref({
  username: '',
  password: '',
  email: '',
  phone: ''
});

// 表单验证
const validateForm = () => {
  let isValid = true;
  
  // 重置错误信息
  errors.value = {
    username: '',
    password: '',
    email: ''
  };
  
  // 验证用户名
  if (!form.value.username) {
    errors.value.username = '用户名不能为空';
    isValid = false;
  } else if (form.value.username.length < 3 || form.value.username.length > 20) {
    errors.value.username = '用户名长度必须在3-20之间';
    isValid = false;
  }
  
  // 验证密码
  if (!form.value.password) {
    errors.value.password = '密码不能为空';
    isValid = false;
  } else if (form.value.password.length < 6) {
    errors.value.password = '密码长度不能少于6位';
    isValid = false;
  }
  
  // 验证邮箱和电话（仅注册时）
  if (!isLogin.value) {
    if (!form.value.email) {
      errors.value.email = '邮箱不能为空';
      isValid = false;
    } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.value.email)) {
      errors.value.email = '请输入有效的邮箱地址';
      isValid = false;
    }
    
    if (!form.value.phone) {
      errors.value.phone = '电话不能为空';
      isValid = false;
    } else if (!/^1[3-9]\d{9}$/.test(form.value.phone)) {
      errors.value.phone = '请输入有效的11位手机号';
      isValid = false;
    }
  }
  
  return isValid;
};

// 检查token并自动登录
const checkToken = async () => {
  try {
    // 调用后端自动登录接口
    const response = await fetch('http://localhost:8080/api/v1/user/auto-login', {
      method: 'GET',
      credentials: 'include' // 包含Cookie
    });
    
    if (response.ok) {
      const data = await response.json();
      if (data.code === 0) {
        console.log('Auto login successful:', data.data);
        // 跳转到主页
        router.push('/home');
      }
    }
  } catch (error) {
    console.error('Auto login failed:', error);
  }
};

const handleSubmit = async () => {
  if (!validateForm()) {
    return;
  }
  
  try {
    const url = isLogin.value 
      ? 'http://localhost:8080/api/v1/user/login' 
      : 'http://localhost:8080/api/v1/user/register';
    
    const response = await fetch(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      credentials: 'include', // 包含Cookie
      body: JSON.stringify(form.value)
    });
    
    const data = await response.json();
    
    if (data.code === 0) {
      console.log(isLogin.value ? 'Login successful' : 'Register successful', data.data);
      // 跳转到主页
      router.push('/home');
    } else {
      console.error(isLogin.value ? 'Login failed' : 'Register failed', data.msg);
    }
  } catch (error) {
    console.error('Request failed:', error);
  }
};

onMounted(() => {
  checkToken();
});
</script>

<style scoped>
.login-container {
  display: flex;
  height: 100vh;
  width: 100vw;
  overflow: hidden;
  background-color: #e6f7ff;
}

.login-left {
  flex: 3;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.welcome-text {
  color: #1890ff;
  font-size: 28px;
  font-weight: 500;
  margin-bottom: 20px;
  text-align: center;
  width: 70%;
}

.login-image {
  width: 70%;
  height: 70%;
  object-fit: cover;
  border-radius: 8px;
}

.login-right {
  flex: 2;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0 2rem;
}

.login-form {
  width: 100%;
  max-width: 400px;
  background-color: white;
  padding: 40px;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.login-form h1 {
  color: #1890ff;
  margin-bottom: 30px;
  text-align: center;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  color: #333;
  font-weight: 500;
}

.form-group input {
  width: 100%;
  padding: 12px;
  border: 1px solid #d9d9d9;
  border-radius: 4px;
  font-size: 16px;
  transition: all 0.3s;
}

.form-group input:focus {
  outline: none;
  border-color: #1890ff;
  box-shadow: 0 0 0 2px rgba(24, 144, 255, 0.2);
}

.error-message {
  color: #ff4d4f;
  font-size: 12px;
  margin-top: 4px;
}

.submit-btn {
  width: 100%;
  padding: 12px;
  background-color: #1890ff;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 16px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s;
  margin-top: 10px;
}

.submit-btn:hover {
  background-color: #40a9ff;
}

.submit-btn:active {
  background-color: #096dd9;
}

.toggle-form {
  margin-top: 20px;
  text-align: center;
  color: #666;
  font-size: 14px;
}

.toggle-form span {
  color: #1890ff;
  cursor: pointer;
  margin-left: 5px;
}

.toggle-form span:hover {
  text-decoration: underline;
}
</style>