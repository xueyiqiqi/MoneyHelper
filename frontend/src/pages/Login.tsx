import React, { useState } from 'react'
import api from '../api/client'

interface LoginProps {
  onLogin: () => void
}

const Login: React.FC<LoginProps> = ({ onLogin }) => {
  const [isRegister, setIsRegister] = useState(false)
  const [username, setUsername] = useState('demo_user')
  const [password, setPassword] = useState('password123')
  const [isLoading, setIsLoading] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setIsLoading(true)
    
    // Demo Mode: Simulate a delay and then login directly
    setTimeout(() => {
      localStorage.setItem('token', 'demo_token_12345')
      setIsLoading(false)
      onLogin()
    }, 800)
  }

  return (
    <div className="flex items-center justify-center min-h-screen bg-slate-50">
      <div className="w-full max-w-md p-10 bg-white rounded-2xl shadow-xl border border-slate-100">
        <div className="flex justify-center mb-8">
          <div className="w-16 h-16 bg-blue-600 rounded-2xl flex items-center justify-center shadow-lg shadow-blue-200">
            <span className="text-white text-3xl font-bold">¥</span>
          </div>
        </div>
        
        <h2 className="mb-2 text-3xl font-bold text-center text-slate-900">
          {isRegister ? '创建账号' : '欢迎回来'}
        </h2>
        <p className="text-slate-500 text-center mb-10">
          {isRegister ? '今天就开始管理您的财务' : '登录以访问您的控制面板'}
        </p>

        <form onSubmit={handleSubmit} className="space-y-6">
          <div>
            <label className="block mb-2 text-sm font-semibold text-slate-700">用户名</label>
            <input
              type="text"
              className="w-full px-4 py-3 border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 transition-all bg-slate-50/50"
              placeholder="请输入用户名"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              required
            />
          </div>
          <div>
            <label className="block mb-2 text-sm font-semibold text-slate-700">密码</label>
            <input
              type="password"
              className="w-full px-4 py-3 border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 transition-all bg-slate-50/50"
              placeholder="••••••••"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
            />
          </div>
          
          <button
            type="submit"
            disabled={isLoading}
            className="w-full py-3.5 text-white bg-blue-600 rounded-xl font-bold hover:bg-blue-700 focus:outline-none focus:ring-4 focus:ring-blue-500/20 transition-all shadow-lg shadow-blue-100 disabled:opacity-70 disabled:cursor-not-allowed"
          >
            {isLoading ? (
              <span className="flex items-center justify-center">
                <svg className="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                  <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                  <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                处理中...
              </span>
            ) : (isRegister ? '注册' : '登录')}
          </button>
        </form>

        <div className="mt-8 pt-8 border-t border-slate-100 text-center">
          <p className="text-sm text-slate-600">
            {isRegister ? '已经有账号了？' : "还没有账号？"}{' '}
            <button
              className="font-bold text-blue-600 hover:text-blue-700 transition-colors"
              onClick={() => setIsRegister(!isRegister)}
            >
              {isRegister ? '点击登录' : '立即创建账号'}
            </button>
          </p>
        </div>
      </div>
    </div>
  )
}

export default Login
