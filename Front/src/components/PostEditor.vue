<template>
  <div class="post-editor-container">
    <div class="editor-header">
      <h2>发布帖子</h2>
      <div class="editor-actions">
        <button class="cancel-btn" @click="resetForm">重置</button>
        <button class="submit-btn" @click="submitPost">发布</button>
      </div>
    </div>
    
    <div class="editor-content">
      <!-- 标题输入 -->
      <input 
        type="text" 
        v-model="postForm.title" 
        class="title-input" 
        placeholder="请输入帖子标题（3-100字）"
        maxlength="100"
      />
      
      <!-- 工具栏 -->
      <div class="editor-toolbar">
        <div class="toolbar-group">
          <button 
            class="toolbar-btn" 
            :class="{ active: activeFormats.bold }"
            @click="format('bold')" 
            title="加粗"
          >
            <b>B</b>
          </button>
          <button 
            class="toolbar-btn" 
            :class="{ active: activeFormats.italic }"
            @click="format('italic')" 
            title="斜体"
          >
            <i>I</i>
          </button>
          <button 
            class="toolbar-btn" 
            :class="{ active: activeFormats.underline }"
            @click="format('underline')" 
            title="下划线"
          >
            <u>U</u>
          </button>
        </div>
        
        <div class="toolbar-group">
          <button class="toolbar-btn" @click="saveCursorPosition(); showEmojiPicker = !showEmojiPicker" title="表情">
            😊
          </button>
          <button class="toolbar-btn" @click="triggerImageUpload" title="上传图片">
            🖼️
          </button>
          <input 
            type="file" 
            ref="fileInput" 
            style="display: none" 
            accept="image/*"
            @change="handleImageUpload"
          />
        </div>
      </div>
      
      <!-- 表情选择器 -->
      <div v-if="showEmojiPicker" class="emoji-picker" style="z-index: 1000;">
        <div class="emoji-grid">
          <span 
            v-for="emoji in emojis" 
            :key="emoji"
            class="emoji-item"
            @click="insertEmoji(emoji)"
            style="cursor: pointer;"
          >
            {{ emoji }}
          </span>
        </div>
      </div>
      
      <!-- 内容编辑器 -->
      <div 
        ref="editor" 
        class="content-editor"
        contenteditable="true"
        @input="updateContent"
        @focus="editorFocused = true; updateActiveFormats()"
        @blur="editorFocused = false"
        @mouseup="updateActiveFormats()"
        @keyup="updateActiveFormats()"
        placeholder="开始编写帖子内容..."
      ></div>
      

    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';

const postForm = ref({
  title: '',
  content: ''
});

const editor = ref(null);
const fileInput = ref(null);
const showEmojiPicker = ref(false);
const editorFocused = ref(false);
const activeFormats = ref({
  bold: false,
  italic: false,
  underline: false
});

// 保存光标位置
let savedRange = null;

// 常用emoji表情
const emojis = [
  '😊', '😂', '😍', '🤔', '😮', '😢', '😡', '👍',
  '👎', '👌', '✌️', '❤️', '🎉', '🔥', '🌟', '💯'
];

// 保存光标位置
const saveCursorPosition = () => {
  if (window.getSelection) {
    const sel = window.getSelection();
    if (sel.getRangeAt && sel.rangeCount) {
      savedRange = sel.getRangeAt(0);
    }
  } else if (document.selection && document.selection.createRange) {
    savedRange = document.selection.createRange();
  }
};

// 恢复光标位置
const restoreCursorPosition = () => {
  if (savedRange) {
    if (window.getSelection) {
      const sel = window.getSelection();
      sel.removeAllRanges();
      sel.addRange(savedRange);
    } else if (document.selection && savedRange.select) {
      savedRange.select();
    }
  }
};

// 格式化文本
const format = (command) => {
  document.execCommand(command, false, null);
  updateActiveFormats();
  editor.value.focus();
};

// 更新激活的格式
const updateActiveFormats = () => {
  activeFormats.value.bold = document.queryCommandState('bold');
  activeFormats.value.italic = document.queryCommandState('italic');
  activeFormats.value.underline = document.queryCommandState('underline');
};



// 插入表情
const insertEmoji = (emoji) => {
  try {
    // 确保编辑器获得焦点
    editor.value.focus();
    
    // 恢复光标位置
    restoreCursorPosition();
    
    // 使用更简单的方法插入表情
    document.execCommand('insertText', false, emoji);
    
    // 关闭表情选择器
    showEmojiPicker.value = false;
    
    // 更新内容
    updateContent();
  } catch (error) {
    console.error('插入表情失败:', error);
    //  fallback方法
    if (editor.value) {
      const currentContent = editor.value.innerHTML;
      editor.value.innerHTML = currentContent + emoji;
      showEmojiPicker.value = false;
      updateContent();
    }
  }
};

// 触发图片上传
const triggerImageUpload = () => {
  fileInput.value.click();
};

// 处理图片上传
const handleImageUpload = (event) => {
  const file = event.target.files[0];
  if (file) {
    // 这里可以实现真实的图片上传逻辑
    // 现在使用模拟的base64图片
    const reader = new FileReader();
    reader.onload = (e) => {
      const img = `<img src="${e.target.result}" style="max-width: 100%; height: auto;" />`;
      document.execCommand('insertHTML', false, img);
    };
    reader.readAsDataURL(file);
  }
};

// 更新内容
const updateContent = () => {
  if (editor.value) {
    postForm.value.content = editor.value.innerHTML;
  }
};

// 提交帖子
const submitPost = async () => {
  try {
    // 确保内容已更新
    updateContent();
    
    if (postForm.value.title.length < 3) {
      alert('标题至少3个字符');
      return;
    }
    
    // 计算纯文本长度，不包含HTML标签
    const plainText = postForm.value.content.replace(/<[^>]*>/g, '');
    if (plainText.length < 10) {
      alert('内容至少10个字符');
      return;
    }
    
    // 限制内容长度，防止超出数据库限制
    if (postForm.value.content.length > 1048576) { // 1MB 左右
      alert('内容过长，请精简后再发布');
      return;
    }
    
    const response = await fetch('http://localhost:8080/api/v1/post/create', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      credentials: 'include',
      body: JSON.stringify(postForm.value)
    });
    
    if (response.ok) {
      const data = await response.json();
      if (data.code === 0) {
        alert('发布成功！');
        resetForm();
        // 刷新帖子列表
        window.location.reload();
      } else {
        alert('发布失败：' + data.msg);
      }
    }
  } catch (error) {
    console.error('Create post failed:', error);
    alert('发布失败，请重试');
  }
};

// 重置表单
const resetForm = () => {
  postForm.value = {
    title: '',
    content: ''
  };
  if (editor.value) {
    editor.value.innerHTML = '';
  }
};

onMounted(() => {
  // 初始化编辑器
  if (editor.value) {
    editor.value.innerHTML = '';
  }
});
</script>

<style scoped>
.post-editor-container {
  background-color: white;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.09);
  min-height: 600px;
  display: flex;
  flex-direction: column;
}

.editor-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 12px;
  border-bottom: 1px solid #e8e8e8;
}

.editor-header h2 {
  font-size: 18px;
  font-weight: 500;
  color: #333;
  margin: 0;
}

.editor-actions {
  display: flex;
  gap: 10px;
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

.editor-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.title-input {
  padding: 12px;
  border: 1px solid #d9d9d9;
  border-radius: 4px;
  font-size: 16px;
  font-weight: 500;
  outline: none;
  transition: all 0.3s;
}

.title-input:focus {
  border-color: #1890ff;
  box-shadow: 0 0 0 2px rgba(24, 144, 255, 0.2);
}

.editor-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 0;
  border-bottom: 1px solid #e8e8e8;
  border-top: 1px solid #e8e8e8;
}

.toolbar-group {
  display: flex;
  align-items: center;
  gap: 8px;
}

.toolbar-btn {
  padding: 6px 10px;
  background-color: white;
  border: 1px solid #d9d9d9;
  border-radius: 4px;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.3s;
  display: flex;
  align-items: center;
  justify-content: center;
  min-width: 36px;
  height: 36px;
}

.toolbar-btn:hover {
  background-color: #f5f5f5;
  border-color: #1890ff;
}

.toolbar-btn.active {
  background-color: #e6f7ff;
  border-color: #1890ff;
  color: #1890ff;
}

.toolbar-divider {
  width: 1px;
  height: 20px;
  background-color: #e8e8e8;
  margin: 0 8px;
}

.emoji-picker {
  position: absolute;
  background-color: white;
  border: 1px solid #d9d9d9;
  border-radius: 4px;
  padding: 12px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  z-index: 100;
  margin-top: 8px;
}

.emoji-grid {
  display: grid;
  grid-template-columns: repeat(8, 1fr);
  gap: 8px;
}

.emoji-item {
  font-size: 20px;
  cursor: pointer;
  padding: 4px;
  text-align: center;
  border-radius: 4px;
  transition: all 0.2s;
}

.emoji-item:hover {
  background-color: #f0f8ff;
}

.content-editor {
  min-height: 300px;
  max-height: 600px;
  padding: 16px;
  border: 1px solid #d9d9d9;
  border-radius: 4px;
  outline: none;
  overflow-y: auto;
  font-size: 14px;
  line-height: 1.6;
  white-space: pre-wrap;
}

.content-editor:focus {
  border-color: #1890ff;
  box-shadow: 0 0 0 2px rgba(24, 144, 255, 0.2);
}

.content-editor:empty:before {
  content: attr(placeholder);
  color: #999;
  pointer-events: none;
}



/* 编辑器内容样式 */
.content-editor h1 {
  font-size: 24px;
  font-weight: 600;
  margin: 20px 0 10px;
}

.content-editor h2 {
  font-size: 20px;
  font-weight: 600;
  margin: 16px 0 8px;
}

.content-editor h3 {
  font-size: 16px;
  font-weight: 600;
  margin: 12px 0 6px;
}

.content-editor ul {
  margin: 8px 0;
  padding-left: 24px;
}

.content-editor ol {
  margin: 8px 0;
  padding-left: 24px;
}

.content-editor li {
  margin: 4px 0;
}

.content-editor pre {
  background-color: #f5f5f5;
  padding: 12px;
  border-radius: 4px;
  overflow-x: auto;
  margin: 8px 0;
}

.content-editor code {
  font-family: 'Courier New', Courier, monospace;
  font-size: 13px;
}

.content-editor img {
  max-width: 100%;
  height: auto;
  margin: 8px 0;
  border-radius: 4px;
}

/* 滚动条样式 */
.content-editor::-webkit-scrollbar {
  width: 8px;
}

.content-editor::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 4px;
}

.content-editor::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 4px;
}

.content-editor::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8;
}
</style>