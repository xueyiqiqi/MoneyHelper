import React, { useState, useEffect } from 'react'
import api from '../api/client'
import Sidebar from '../components/Sidebar'
import StatsCard from '../components/StatsCard'
import { PlusCircle, List, MessageSquare, CreditCard, TrendingUp, Wallet, ArrowUpRight, Users, Briefcase, ExternalLink } from 'lucide-react'
import { Link } from 'react-router-dom'

interface Bill {
  id: number
  amount: number
  category: string
  remarks?: string
  date: string
}

interface AnalysisReport {
  id: number
  content: string
  created_at: string
}

interface Space {
  id: number
  name: string
  description?: string
  total_expense?: number
  role?: string
}

interface DashboardProps {
  initialMode?: 'personal' | 'collaboration'
}

const Dashboard: React.FC<DashboardProps> = ({ initialMode = 'personal' }) => {
  const [viewMode, setViewMode] = useState<'personal' | 'collaboration'>(initialMode)
  const [bills, setBills] = useState<Bill[]>([])
  const [reports, setReports] = useState<AnalysisReport[]>([])
  const [spaces, setSpaces] = useState<Space[]>([])
  const [amount, setAmount] = useState('')
  const [category, setCategory] = useState('')
  const [remarks, setRemarks] = useState('')
  const [selectedSpaceId, setSelectedSpaceId] = useState<string>('') // 新增：协作模式下选中的空间ID
  const [isAnalyzing, setIsAnalyzing] = useState(false)

  useEffect(() => {
    setViewMode(initialMode)
  }, [initialMode])

  useEffect(() => {
    fetchBills()
    fetchReports()
    if (viewMode === 'personal') {
      fetchSpaces()
    }
  }, [viewMode])

  useEffect(() => {
    if (viewMode === 'collaboration' && selectedSpaceId) {
      fetchBills()
      fetchReports()
    }
  }, [selectedSpaceId])

  useEffect(() => {
    if (viewMode === 'collaboration') {
      fetchSpaces()
    }
  }, [])

  const fetchSpaces = async () => {
    try {
      const res = await api.get('/spaces')
      const enrichedSpaces = (res.data || []).map((s: any) => ({
        ...s,
        total_expense: Math.floor(Math.random() * 5000) + 1000,
        role: '创建者'
      }))
      setSpaces(enrichedSpaces)
      if (enrichedSpaces.length > 0 && !selectedSpaceId) {
        setSelectedSpaceId(enrichedSpaces[0].id.toString())
      }
    } catch (err) {
      console.error('Error fetching spaces:', err)
      setSpaces([])
    }
  }

  const fetchBills = async () => {
    try {
      const params: any = { is_personal: viewMode === 'personal' }
      if (viewMode === 'collaboration' && selectedSpaceId) {
        params.space_id = selectedSpaceId
      }
      const res = await api.get('/bills', { params })
      setBills(res.data || [])
    } catch (err) {
      console.error('Error fetching bills:', err)
    }
  }

  const fetchReports = async () => {
    try {
      const params: any = { is_personal: viewMode === 'personal' }
      if (viewMode === 'collaboration' && selectedSpaceId) {
        params.space_id = selectedSpaceId
      }
      const res = await api.get('/reports', { params })
      setReports(res.data || [])
    } catch (err) {
      console.error('Error fetching reports:', err)
    }
  }

  const handleAddBill = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      const isPersonal = viewMode === 'personal'
      await api.post('/bills', {
        amount: parseFloat(amount),
        category,
        remarks,
        is_personal: isPersonal,
        space_id: !isPersonal ? Number(selectedSpaceId) : undefined
      })
      setAmount('')
      setCategory('')
      setRemarks('')
      fetchBills()
    } catch (err) {
      console.error('Error adding bill:', err)
    }
  }

  const handleAddSpace = async () => {
    const name = prompt('请输入空间名称:');
    if (name) {
      try {
        await api.post('/spaces', { name });
        fetchSpaces();
      } catch (err) {
        console.error('Error creating space:', err);
      }
    }
  }

  const handleAnalyze = async () => {
    setIsAnalyzing(true)
    try {
      await api.post('/analyze', { 
        is_personal: viewMode === 'personal',
        space_id: viewMode === 'collaboration' ? Number(selectedSpaceId) : undefined 
      })
      fetchReports()
    } catch (err) {
      console.error('Error analyzing bills:', err)
    } finally {
      setIsAnalyzing(false)
    }
  }

  const handleLogout = () => {
    localStorage.removeItem('token');
    window.location.reload();
  }

  const totalPersonalExpense = bills.reduce((sum, bill) => sum + bill.amount, 0);
  const totalSharedExpense = spaces.reduce((sum, space) => sum + (space.total_expense || 0), 0);

  return (
    <div className="flex min-h-screen bg-slate-50">
      <Sidebar 
        spaces={spaces} 
        onAddSpace={handleAddSpace} 
        onLogout={handleLogout} 
      />

      <main className="flex-1 ml-64 p-10">
        <header className="flex flex-col mb-10">
          <div className="flex justify-between items-center mb-6">
            <div>
              <div className="flex items-center space-x-3">
                <h2 className="text-3xl font-bold text-slate-900">
                  {viewMode === 'personal' ? '个人主页' : '协作空间示例'}
                </h2>
                {viewMode === 'collaboration' && spaces.length > 0 && (
                  <select 
                    className="ml-4 px-3 py-1.5 text-sm font-semibold text-blue-600 bg-blue-50 border-none rounded-lg focus:ring-2 focus:ring-blue-500/20 cursor-pointer hover:bg-blue-100 transition-colors"
                    value={selectedSpaceId}
                    onChange={(e) => setSelectedSpaceId(e.target.value)}
                  >
                    {spaces.map(s => (
                      <option key={s.id} value={s.id}>{s.name}</option>
                    ))}
                  </select>
                )}
              </div>
              <p className="text-slate-500 mt-1">
                {viewMode === 'personal' ? '欢迎回来！这是您的财务状况概览。' : `正在查看：${spaces.find(s => s.id.toString() === selectedSpaceId)?.name || '选择一个空间'}`}
              </p>
            </div>
            <div className="flex space-x-4">
              <button 
                onClick={handleAnalyze}
                disabled={isAnalyzing}
                className="flex items-center px-5 py-2.5 text-blue-600 bg-blue-50 rounded-xl font-semibold hover:bg-blue-100 transition-colors disabled:opacity-50"
              >
                <MessageSquare className="w-4 h-4 mr-2" />
                {isAnalyzing ? '分析中...' : 'AI 财务建议'}
              </button>
            </div>
          </div>
        </header>

        {/* Stats Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-10">
          <StatsCard 
            title={viewMode === 'personal' ? "总支出" : "空间累计支出"} 
            value={`¥${(viewMode === 'personal' ? totalPersonalExpense : totalSharedExpense).toLocaleString()}`} 
            icon={Wallet} 
            color="blue" 
            trend={viewMode === 'personal' ? { value: '12%', isUp: false } : undefined}
          />
          <StatsCard 
            title={viewMode === 'personal' ? "本月预算" : "协作空间总数"} 
            value={viewMode === 'personal' ? "¥8,000" : `${spaces.length}`} 
            icon={viewMode === 'personal' ? CreditCard : Users} 
            color="purple" 
          />
          <StatsCard 
            title={viewMode === 'personal' ? "储蓄率" : "活跃空间"} 
            value={viewMode === 'personal' ? "24%" : `${spaces.filter(s => (s.total_expense || 0) > 0).length}`} 
            icon={TrendingUp} 
            color="green" 
            trend={viewMode === 'personal' ? { value: '5%', isUp: true } : undefined}
          />
          <StatsCard 
            title={viewMode === 'personal' ? "健康评分" : "空间 AI 评分"} 
            value="85/100" 
            icon={ArrowUpRight} 
            color="blue" 
          />
        </div>

        <div className="grid grid-cols-1 xl:grid-cols-3 gap-10">
          <div className="xl:col-span-2 space-y-10">
            {/* Recent Bills */}
            <section className="bg-white rounded-2xl shadow-sm border border-slate-100 overflow-hidden">
              <div className="p-6 border-b border-slate-100 flex justify-between items-center">
                <h3 className="text-lg font-bold text-slate-900 flex items-center">
                  <List className="w-5 h-5 mr-2 text-blue-600" />
                  {viewMode === 'personal' ? '最近交易记录' : '最近协作交易'}
                </h3>
                <button className="text-sm font-semibold text-blue-600 hover:text-blue-700">查看全部</button>
              </div>
              <div className="overflow-x-auto">
                <table className="w-full text-left">
                  <thead>
                    <tr className="bg-slate-50/50">
                      <th className="py-4 px-6 text-xs font-bold text-slate-500 uppercase tracking-wider">日期</th>
                      {viewMode === 'collaboration' && (
                        <th className="py-4 px-6 text-xs font-bold text-slate-500 uppercase tracking-wider">创建人</th>
                      )}
                      <th className="py-4 px-6 text-xs font-bold text-slate-500 uppercase tracking-wider">分类</th>
                      <th className="py-4 px-6 text-xs font-bold text-slate-500 uppercase tracking-wider">备注</th>
                      <th className="py-4 px-6 text-xs font-bold text-slate-500 uppercase tracking-wider text-right">金额</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-slate-100">
                    {bills.map(bill => (
                      <tr key={bill.id} className="hover:bg-slate-50/50 transition-colors">
                        <td className="py-4 px-6 text-sm text-slate-500">{new Date(bill.date).toLocaleDateString()}</td>
                        {viewMode === 'collaboration' && (
                          <td className="py-4 px-6">
                            <div className="flex items-center">
                              <div className="w-6 h-6 rounded-full bg-blue-100 text-blue-600 flex items-center justify-center text-[10px] font-bold mr-2">
                                {(bill.creator_name || 'U').charAt(0).toUpperCase()}
                              </div>
                              <span className="text-sm text-slate-700 font-medium">{bill.creator_name || '演示用户'}</span>
                            </div>
                          </td>
                        )}
                        <td className="py-4 px-6">
                          <span className="px-3 py-1 bg-slate-100 text-slate-700 rounded-full text-xs font-semibold">
                            {bill.category}
                          </span>
                        </td>
                        <td className="py-4 px-6 text-sm text-slate-600">{bill.remarks || '-'}</td>
                        <td className="py-4 px-6 text-sm font-bold text-right text-slate-900">¥{bill.amount.toFixed(2)}</td>
                      </tr>
                    ))}
                    {bills.length === 0 && (
                      <tr>
                        <td colSpan={viewMode === 'collaboration' ? 5 : 4} className="py-12 text-center">
                          <div className="flex flex-col items-center text-slate-400">
                            <Wallet className="w-12 h-12 mb-3 opacity-20" />
                            <p>暂无交易记录</p>
                          </div>
                        </td>
                      </tr>
                    )}
                  </tbody>
                </table>
              </div>
            </section>

            {/* AI Reports */}
            <section className="bg-white rounded-2xl shadow-sm border border-slate-100 p-6">
              <h3 className="text-lg font-bold text-slate-900 mb-6 flex items-center">
                <MessageSquare className="w-5 h-5 mr-2 text-blue-600" />
                最新 AI 分析报告
              </h3>
              <div className="space-y-6">
                {reports.map(report => (
                  <div key={report.id} className={`p-5 rounded-xl border ${viewMode === 'personal' ? 'bg-blue-50/30 border-blue-100' : 'bg-purple-50/30 border-purple-100'}`}>
                    <div className="flex justify-between items-start mb-3">
                      <span className={`px-2 py-1 text-white text-[10px] font-bold uppercase rounded ${viewMode === 'personal' ? 'bg-blue-600' : 'bg-purple-600'}`}>
                        {viewMode === 'personal' ? 'AI 报告' : '空间报告'}
                      </span>
                      <span className="text-xs text-slate-400">{new Date(report.created_at).toLocaleString()}</span>
                    </div>
                    <div className="text-sm text-slate-700 leading-relaxed whitespace-pre-wrap">
                      {report.content}
                    </div>
                  </div>
                ))}
                {reports.length === 0 && (
                  <div className="text-center py-10 bg-slate-50 rounded-xl border border-dashed border-slate-200">
                    <p className="text-slate-400 text-sm">点击“AI 财务建议”生成您的第一份报告</p>
                  </div>
                )}
              </div>
            </section>

            {/* 协作模式下的空间网格 - 放在底部作为选择器 */}
            {viewMode === 'collaboration' && (
              <section className="space-y-6">
                <h3 className="text-lg font-bold text-slate-900 flex items-center">
                  <Briefcase className="w-5 h-5 mr-2 text-blue-600" />
                  您的协作空间
                </h3>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                  {spaces.map(space => (
                    <Link 
                      key={space.id} 
                      to={`/space/${space.id}`}
                      className="bg-white rounded-2xl shadow-sm border border-slate-100 p-6 hover:shadow-md transition-all group"
                    >
                      <div className="flex justify-between items-start mb-4">
                        <div className="w-10 h-10 bg-blue-50 rounded-lg flex items-center justify-center text-blue-600 group-hover:bg-blue-600 group-hover:text-white transition-colors">
                          <Briefcase className="w-5 h-5" />
                        </div>
                        <span className="px-2 py-0.5 bg-slate-100 text-slate-600 rounded-md text-[10px] font-bold uppercase">
                          {space.role}
                        </span>
                      </div>
                      <h4 className="text-lg font-bold text-slate-900 mb-1">{space.name}</h4>
                      <div className="flex items-center justify-between mt-4">
                        <p className="text-xs text-slate-400">本月支出: <span className="text-blue-600 font-bold">¥{(space.total_expense || 0).toLocaleString()}</span></p>
                        <div className="text-blue-600 font-bold text-xs flex items-center">
                          进入 <ExternalLink className="w-3 h-3 ml-1" />
                        </div>
                      </div>
                    </Link>
                  ))}
                </div>
              </section>
            )}
          </div>

          {/* Quick Action - Add Bill */}
          <aside>
            <section className="bg-white rounded-2xl shadow-sm border border-slate-100 p-6 sticky top-10">
              <h3 className="text-lg font-bold text-slate-900 mb-6 flex items-center">
                <PlusCircle className="w-5 h-5 mr-2 text-blue-600" />
                {viewMode === 'personal' ? '快速记账' : `快速记账 (${spaces.find(s => s.id.toString() === selectedSpaceId)?.name || '空间'})`}
              </h3>
              <form onSubmit={handleAddBill} className="space-y-5">
                <div>
                  <label className="block text-sm font-semibold text-slate-700 mb-2">金额 (¥)</label>
                  <input 
                    type="number" 
                    step="0.01" 
                    required 
                    placeholder="0.00"
                    className="w-full px-4 py-3 border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 bg-slate-50/50"
                    value={amount}
                    onChange={(e) => setAmount(e.target.value)}
                  />
                </div>
                <div>
                  <label className="block text-sm font-semibold text-slate-700 mb-2">分类</label>
                  <select 
                    required 
                    className="w-full px-4 py-3 border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 bg-slate-50/50"
                    value={category}
                    onChange={(e) => setCategory(e.target.value)}
                  >
                    <option value="">选择分类</option>
                    <option value="餐饮">🍔 餐饮</option>
                    <option value="交通">🚗 交通</option>
                    <option value="购物">🛍️ 购物</option>
                    <option value="账单">⚡ 账单</option>
                    <option value="娱乐">🎬 娱乐</option>
                    <option value="办公">🏢 办公</option>
                    <option value="团队">👥 团队</option>
                    <option value="其他">✨ 其他</option>
                  </select>
                </div>
                <div>
                  <label className="block text-sm font-semibold text-slate-700 mb-2">备注</label>
                  <textarea 
                    placeholder="这笔钱花在哪了？"
                    className="w-full px-4 py-3 border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 bg-slate-50/50"
                    rows={3}
                    value={remarks}
                    onChange={(e) => setRemarks(e.target.value)}
                  />
                </div>
                <button 
                  type="submit" 
                  className="w-full py-4 text-white bg-blue-600 rounded-xl font-bold hover:bg-blue-700 transition-all shadow-lg shadow-blue-100"
                >
                  确认记账
                </button>
              </form>
            </section>
          </aside>
        </div>
      </main>
    </div>
  )
}

export default Dashboard
