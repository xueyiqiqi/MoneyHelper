import React, { useState, useEffect } from 'react'
import { useParams, Link } from 'react-router-dom'
import api from '../api/client'
import Sidebar from '../components/Sidebar'
import StatsCard from '../components/StatsCard'
import { ArrowLeft, MessageSquare, PlusCircle, List, User, Shield, Users, Eye, Wallet, CreditCard, TrendingUp } from 'lucide-react'

interface Space {
  id: number
  name: string
  description?: string
}

interface Bill {
  id: number
  amount: number
  category: string
  remarks?: string
  date: string
  creator_name?: string // 新增：创建人姓名
}

const SpaceDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>()
  const [space, setSpace] = useState<Space | null>(null)
  const [bills, setBills] = useState<Bill[]>([])
  const [reports, setReports] = useState<any[]>([])
  const [spaces, setSpaces] = useState<any[]>([])
  const [amount, setAmount] = useState('')
  const [category, setCategory] = useState('')
  const [remarks, setRemarks] = useState('')
  const [isAnalyzing, setIsAnalyzing] = useState(false)
  const [role, setRole] = useState<string>('observer')

  useEffect(() => {
    fetchSpace()
    fetchBills()
    fetchReports()
    fetchSpaces()
  }, [id])

  const fetchSpaces = async () => {
    try {
      const res = await api.get('/spaces')
      setSpaces(res.data || [])
    } catch (err) {
      console.error('Error fetching spaces:', err)
    }
  }

  const fetchSpace = async () => {
    try {
      const res = await api.get('/spaces')
      const currentSpace = res.data.find((s: Space) => s.id === Number(id))
      if (currentSpace) {
        setSpace(currentSpace)
        setRole('creator') // Mocking role
      }
    } catch (err) {
      console.error('Error fetching space:', err)
    }
  }

  const fetchBills = async () => {
    try {
      const res = await api.get(`/bills?space_id=${id}&is_personal=false`)
      setBills(res.data || [])
    } catch (err) {
      console.error('Error fetching bills:', err)
    }
  }

  const fetchReports = async () => {
    try {
      const res = await api.get(`/reports?space_id=${id}&is_personal=false`)
      setReports(res.data || [])
    } catch (err) {
      console.error('Error fetching reports:', err)
    }
  }

  const handleAddBill = async (e: React.FormEvent) => {
    e.preventDefault()
    if (role === 'observer') return

    try {
      await api.post('/bills', {
        amount: parseFloat(amount),
        category,
        remarks,
        is_personal: false,
        space_id: Number(id)
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
      await api.post('/analyze', { space_id: Number(id), is_personal: false })
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

  const getRoleIcon = () => {
    switch (role) {
      case 'creator': return <Shield className="w-5 h-5 text-purple-600" title="创建者" />
      case 'admin': return <Users className="w-5 h-5 text-blue-600" title="管理员" />
      case 'member': return <User className="w-5 h-5 text-green-600" title="成员" />
      case 'observer': return <Eye className="w-5 h-5 text-gray-400" title="观察者 (只读)" />
      default: return null
    }
  }

  const totalExpense = bills.reduce((sum, bill) => sum + bill.amount, 0);

  return (
    <div className="flex min-h-screen bg-slate-50">
      <Sidebar 
        spaces={spaces} 
        onAddSpace={handleAddSpace} 
        onLogout={handleLogout} 
      />

      <main className="flex-1 ml-64 p-10">
        <header className="mb-10">
          <Link to="/" className="inline-flex items-center text-sm font-semibold text-blue-600 hover:text-blue-700 mb-6 group">
            <ArrowLeft className="w-4 h-4 mr-2 transition-transform group-hover:-translate-x-1" />
            返回个人看板
          </Link>
          <div className="flex justify-between items-start">
            <div>
              <div className="flex items-center space-x-3 mb-2">
                <h1 className="text-3xl font-bold text-slate-900">{space?.name}</h1>
                <div className="flex items-center px-3 py-1 bg-white rounded-full border border-slate-200 shadow-sm">
                  {getRoleIcon()}
                  <span className="ml-2 text-xs font-bold text-slate-500 uppercase tracking-wider">{role === 'creator' ? '创建者' : role === 'admin' ? '管理员' : role === 'member' ? '成员' : '观察者'}</span>
                </div>
              </div>
              <p className="text-slate-500">{space?.description || '您的团队协作财务空间。'}</p>
            </div>
            {role !== 'observer' && (
              <button 
                onClick={handleAnalyze}
                disabled={isAnalyzing}
                className="flex items-center px-5 py-2.5 text-blue-600 bg-blue-50 rounded-xl font-semibold hover:bg-blue-100 transition-colors disabled:opacity-50"
              >
                <MessageSquare className="w-4 h-4 mr-2" />
                {isAnalyzing ? '分析中...' : '空间财务建议'}
              </button>
            )}
          </div>
        </header>

        {/* Stats Grid */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-10">
          <StatsCard 
            title="空间总支出" 
            value={`¥${totalExpense.toLocaleString()}`} 
            icon={Wallet} 
            color="blue" 
          />
          <StatsCard 
            title="参与成员" 
            value="4 人" 
            icon={Users} 
            color="purple" 
          />
          <StatsCard 
            title="预算执行率" 
            value="68%" 
            icon={TrendingUp} 
            color="green" 
          />
        </div>

        <div className="grid grid-cols-1 xl:grid-cols-3 gap-10">
          <div className="xl:col-span-2 space-y-10">
            {/* Space Bills */}
            <section className="bg-white rounded-2xl shadow-sm border border-slate-100 overflow-hidden">
              <div className="p-6 border-b border-slate-100 flex justify-between items-center">
                <h3 className="text-lg font-bold text-slate-900 flex items-center">
                  <List className="w-5 h-5 mr-2 text-blue-600" />
                  空间交易明细
                </h3>
              </div>
              <div className="overflow-x-auto">
                <table className="w-full text-left">
                  <thead>
                    <tr className="bg-slate-50/50">
                      <th className="py-4 px-6 text-xs font-bold text-slate-500 uppercase tracking-wider">日期</th>
                      <th className="py-4 px-6 text-xs font-bold text-slate-500 uppercase tracking-wider">创建人</th>
                      <th className="py-4 px-6 text-xs font-bold text-slate-500 uppercase tracking-wider">分类</th>
                      <th className="py-4 px-6 text-xs font-bold text-slate-500 uppercase tracking-wider">备注</th>
                      <th className="py-4 px-6 text-xs font-bold text-slate-500 uppercase tracking-wider text-right">金额</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-slate-100">
                    {bills.map(bill => (
                      <tr key={bill.id} className="hover:bg-slate-50/50 transition-colors">
                        <td className="py-4 px-6 text-sm text-slate-500">{new Date(bill.date).toLocaleDateString()}</td>
                        <td className="py-4 px-6">
                          <div className="flex items-center">
                            <div className="w-6 h-6 rounded-full bg-blue-100 text-blue-600 flex items-center justify-center text-[10px] font-bold mr-2">
                              {(bill.creator_name || 'U').charAt(0).toUpperCase()}
                            </div>
                            <span className="text-sm text-slate-700 font-medium">{bill.creator_name || '演示用户'}</span>
                          </div>
                        </td>
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
                        <td colSpan={4} className="py-12 text-center text-slate-400">
                          该空间暂无交易记录。
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
                空间 AI 洞察
              </h3>
              <div className="space-y-6">
                {reports.map(report => (
                  <div key={report.id} className="p-5 rounded-xl bg-purple-50/30 border border-purple-100">
                    <div className="flex justify-between items-start mb-3">
                      <span className="px-2 py-1 bg-purple-600 text-white text-[10px] font-bold uppercase rounded">空间报告</span>
                      <span className="text-xs text-slate-400">{new Date(report.created_at).toLocaleString()}</span>
                    </div>
                    <div className="text-sm text-slate-700 leading-relaxed whitespace-pre-wrap">{report.content}</div>
                  </div>
                ))}
                {reports.length === 0 && (
                  <div className="text-center py-10 bg-slate-50 rounded-xl border border-dashed border-slate-200 text-slate-400 text-sm">
                    暂无空间分析报告。
                  </div>
                )}
              </div>
            </section>
          </div>

          {/* Quick Action - Add Bill (if not observer) */}
          {role !== 'observer' && (
            <aside>
              <section className="bg-white rounded-2xl shadow-sm border border-slate-100 p-6 sticky top-10">
                <h3 className="text-lg font-bold text-slate-900 mb-6 flex items-center">
                  <PlusCircle className="w-5 h-5 mr-2 text-blue-600" />
                  快速记账 (空间)
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
                      <option value="办公">🏢 办公</option>
                      <option value="团队">👥 团队</option>
                      <option value="购物">🛍️ 购物</option>
                      <option value="其他">✨ 其他</option>
                    </select>
                  </div>
                  <div>
                    <label className="block text-sm font-semibold text-slate-700 mb-2">备注</label>
                    <textarea 
                      placeholder="这笔钱在空间里花在哪了？"
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
                    确认空间记账
                  </button>
                </form>
              </section>
            </aside>
          )}
        </div>
      </main>
    </div>
  )
}

export default SpaceDetail
