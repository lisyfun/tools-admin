<template>
  <div class="login-container">
    <div class="login-box">
      <el-form
        ref="loginFormRef"
        :model="loginForm"
        :rules="loginRules"
        class="login-form"
        autocomplete="on"
        label-position="top"
      >
        <div class="title-container">
          <div class="header">
            <div class="left-title">登录</div>
            <div class="right-title">运维管理系统</div>
          </div>
        </div>

        <el-form-item prop="username">
          <el-input
            v-model="loginForm.username"
            placeholder="请输入用户名"
            type="text"
            tabindex="1"
            :prefix-icon="User"
            size="large"
          >
            <template #label>
              <span class="input-label" :class="{ 'label-active': loginForm.username }">用户名</span>
            </template>
          </el-input>
        </el-form-item>

        <el-form-item prop="password">
          <el-input
            v-model="loginForm.password"
            :type="passwordVisible ? 'text' : 'password'"
            placeholder="请输入密码"
            tabindex="2"
            :prefix-icon="Lock"
            size="large"
            @keyup.enter="handleLogin(loginFormRef)"
          >
            <template #label>
              <span class="input-label" :class="{ 'label-active': loginForm.password }">密码</span>
            </template>
            <template #suffix>
              <el-icon 
                class="show-pwd" 
                @click="passwordVisible = !passwordVisible"
              >
                <component :is="passwordVisible ? 'View' : 'Hide'" />
              </el-icon>
            </template>
          </el-input>
        </el-form-item>

        <el-form-item>
          <el-button
            :loading="loading"
            type="primary"
            size="large"
            style="width: 100%"
            @click="handleLogin(loginFormRef)"
            tabindex="3"
          >
            登录
          </el-button>
        </el-form-item>

        <div class="tips-container">
          <span class="reset-link" @click="showResetDialog = true">忘记密码？</span>
        </div>

      </el-form>
    </div>

    <!-- 重置密码对话框 -->
    <el-dialog
      v-model="showResetDialog"
      title="重置密码"
      width="400px"
      :close-on-click-modal="false"
      @close="handleDialogClose"
      class="reset-dialog"
      destroy-on-close
      center
      append-to-body
    >
      <div class="reset-container">
        <el-form
          ref="resetFormRef"
          :model="resetForm"
          :rules="resetRules"
          class="reset-form"
          size="large"
          label-position="top"
        >
          <el-form-item prop="username" label="用户名">
            <el-input
              ref="usernameRef"
              v-model="resetForm.username"
              placeholder="请输入用户名"
              :prefix-icon="User"
              clearable
              @keyup.enter="$refs.passwordRef?.$el.querySelector('input').focus()"
            />
          </el-form-item>

          <el-form-item prop="newPassword" label="新密码">
            <el-input
              ref="passwordRef"
              v-model="resetForm.newPassword"
              :type="passwordVisible ? 'text' : 'password'"
              placeholder="请输入新密码"
              :prefix-icon="Lock"
              @keyup.enter="$refs.confirmPasswordRef?.$el.querySelector('input').focus()"
            >
              <template #suffix>
                <el-icon 
                  class="show-pwd" 
                  @click="passwordVisible = !passwordVisible"
                >
                  <component :is="passwordVisible ? 'View' : 'Hide'" />
                </el-icon>
              </template>
            </el-input>
          </el-form-item>

          <el-form-item prop="confirmPassword" label="确认密码">
            <el-input
              ref="confirmPasswordRef"
              v-model="resetForm.confirmPassword"
              :type="passwordVisible ? 'text' : 'password'"
              placeholder="请再次输入新密码"
              :prefix-icon="Lock"
              @keyup.enter="handleResetPassword"
            >
              <template #suffix>
                <el-icon 
                  class="show-pwd" 
                  @click="passwordVisible = !passwordVisible"
                >
                  <component :is="passwordVisible ? 'View' : 'Hide'" />
                </el-icon>
              </template>
            </el-input>
          </el-form-item>

          <div class="reset-button">
            <el-button 
              :loading="resetLoading" 
              type="primary" 
              class="submit-btn"
              @click.prevent="handleResetPassword"
            >
              重置密码
            </el-button>
          </div>
        </el-form>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useMenuStore } from '@/stores/menu'
import { useUserStore } from '@/stores/user'
import type { FormInstance, FormRules } from 'element-plus'
import { User, Lock, View, Hide } from '@element-plus/icons-vue'

const router = useRouter()
const menuStore = useMenuStore()
const userStore = useUserStore()

const loginForm = ref({
  username: 'admin',
  password: '123456'
})

const loginRules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
} as FormRules

const loading = ref(false)
const loginFormRef = ref<FormInstance>()

const handleLogin = async (formEl: FormInstance | undefined) => {
  if (!formEl) return
  await formEl.validate(async (valid: boolean) => {
    if (valid) {
      loading.value = true
      try {
        // 登录
        await userStore.login(loginForm.value)
        // 获取菜单
        await menuStore.fetchMenuList()
        router.push({ path: '/' })
      } catch (error) {
        console.error('Login error:', error)
      } finally {
        loading.value = false
      }
    }
  })
}

const passwordVisible = ref(false)
const showResetDialog = ref(false)
const resetForm = ref({
  username: '',
  newPassword: '',
  confirmPassword: ''
})
const resetRules = {
  username: [
    { required: true, trigger: 'blur', message: '请输入用户名' },
    { min: 3, max: 20, message: '用户名长度应在3-20个字符之间', trigger: 'blur' }
  ],
  newPassword: [
    { required: true, trigger: 'blur', message: '请输入新密码' },
    { min: 6, max: 20, message: '密码长度应在6-20个字符之间', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, trigger: 'blur', message: '请再次输入新密码' },
    {
      validator: (rule: any, value: string, callback: any) => {
        if (value !== resetForm.value.newPassword) {
          callback(new Error('两次输入的密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ]
}
const resetLoading = ref(false)
const resetFormRef = ref<FormInstance>()

const handleDialogClose = () => {
  showResetDialog.value = false
  resetForm.value.username = ''
  resetForm.value.newPassword = ''
  resetForm.value.confirmPassword = ''
  if (resetFormRef.value) {
    resetFormRef.value.resetFields()
  }
}

const handleResetPassword = async () => {
  if (!resetFormRef.value) return

  try {
    resetLoading.value = true
    await resetFormRef.value.validate()

    const success = await userStore.resetPassword({
      username: resetForm.value.username,
      newPassword: resetForm.value.newPassword
    })

    if (success) {
      handleDialogClose()
    }
  } catch (error) {
    console.error('重置密码失败:', error)
  } finally {
    resetLoading.value = false
  }
}
</script>

<style lang="scss" scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  justify-content: center;
  background: #f7f8fa;
  padding-top: 271px;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;

  .login-box {
    width: 440px;
    height: 333px;
    background: #fff;
    box-shadow: 0 10px 50px rgba(0, 0, 0, 0.08);

    .login-form {
      padding: 35px;

      .title-container {
        margin-bottom: 35px;

        .header {
          display: flex;
          justify-content: space-between;
          align-items: center;

          .left-title {
            font-size: 24px;
            font-weight: 500;
            color: #1f2937;
          }

          .right-title {
            font-size: 14px;
            color: #6b7280;
          }
        }
      }

      .tips-container {
        margin-top: 16px;
        text-align: right;

        .reset-link {
          font-size: 14px;
          color: #6b7280;
          cursor: pointer;
          transition: all 0.3s ease;
          text-decoration: none;

          &:hover {
            color: #409eff;
          }
        }
      }
    }
  }
}

.reset-dialog {
  :deep(.el-dialog) {
    border-radius: 20px;
    background: #fff;
    
    .el-dialog__header {
      padding: 15px 20px;
      margin: 0;
      background-color: #2d3a4b;
      border-top-left-radius: 20px;
      border-top-right-radius: 20px;
      
      .el-dialog__title {
        color: #fff;
        font-size: 16px;
        font-weight: 500;
      }
      
      .el-dialog__headerbtn {
        top: 12px;
        right: 16px;
        z-index: 999;
        .el-dialog__close {
          color: #fff;
          font-weight: bold;
          font-size: 18px;
        }
      }
    }
    
    .el-dialog__body {
      padding: 0;
    }

    .reset-container {
      position: relative;
      width: 100%;
      padding: 35px 25px 15px;
      margin: 0 auto;
      overflow: hidden;
      background-color: #2d3a4b;

      .reset-form {
        width: 85%;
        margin: 0 auto;
        
        .el-form-item {
          margin-bottom: 20px;
          border: 1px solid rgba(255, 255, 255, 0.1);
          background: rgba(0, 0, 0, 0.1);
          border-radius: 5px;
          color: #454545;

          &:last-child {
            margin-bottom: 0;
          }
        }

        :deep(.el-input) {
          display: inline-block;
          height: 47px;
          width: 100%;

          .el-input__wrapper {
            padding: 0;
            background: transparent;
            box-shadow: none;
          }

          input {
            background: transparent;
            border: 0;
            -webkit-appearance: none;
            border-radius: 0;
            padding: 12px 5px 12px 15px;
            color: #eee;
            height: 47px;
            caret-color: #fff;
            outline: none;
            
            &:-webkit-autofill {
              box-shadow: 0 0 0 1000px #283443 inset !important;
              -webkit-text-fill-color: #fff !important;
            }

            &:focus {
              outline: none;
            }
          }
        }

        .reset-button {
          margin-top: 15px;
          text-align: center;

          .submit-btn {
            width: 100%;
            height: 47px;
            border-radius: 5px;
            font-size: 16px;
          }
        }
      }
    }
  }
}
</style>
