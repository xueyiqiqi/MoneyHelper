import React from 'react';
import { Link, useLocation } from 'react-router-dom';
import { Layout, User, Briefcase, PlusCircle, LogOut, ChevronRight } from 'lucide-react';

interface Space {
  id: number;
  name: string;
}

interface SidebarProps {
  spaces: Space[];
  onAddSpace: () => void;
  onLogout: () => void;
}

const Sidebar: React.FC<SidebarProps> = ({ spaces, onAddSpace, onLogout }) => {
  const location = useLocation();

  const isActive = (path: string) => location.pathname === path;

  return (
    <aside className="w-64 bg-slate-900 text-slate-300 h-screen flex flex-col fixed left-0 top-0">
      <div className="p-6 border-b border-slate-800">
        <h1 className="text-xl font-bold text-white flex items-center">
          <div className="w-8 h-8 bg-blue-600 rounded-lg mr-3 flex items-center justify-center">
            <span className="text-white text-lg">¥</span>
          </div>
          生活财务助手
        </h1>
      </div>

      <nav className="flex-1 overflow-y-auto py-6">
        <div className="px-4 mb-8">
          <p className="text-xs font-semibold text-slate-500 uppercase tracking-wider mb-4 px-2">主导航</p>
          <Link
            to="/"
            className={`flex items-center px-4 py-3 rounded-lg transition-colors ${
              isActive('/') ? 'bg-blue-600 text-white' : 'hover:bg-slate-800 hover:text-white'
            }`}
          >
            <Layout className="w-5 h-5 mr-3" />
            个人看板
          </Link>
        </div>

        <div className="px-4">
          <div className="flex items-center justify-between mb-4 px-2">
            <p className="text-xs font-semibold text-slate-500 uppercase tracking-wider">协作空间</p>
            <button
              onClick={onAddSpace}
              className="p-1 hover:bg-slate-800 rounded-md transition-colors text-slate-400 hover:text-blue-400"
              title="创建空间"
            >
              <PlusCircle className="w-4 h-4" />
            </button>
          </div>
          
          <div className="space-y-1">
            <Link
              to="/collaboration"
              className={`flex items-center px-4 py-2.5 rounded-lg transition-colors group mb-2 ${
                isActive('/collaboration') ? 'bg-slate-800 text-white' : 'hover:bg-slate-800/50 hover:text-white'
              }`}
            >
              <Briefcase className="w-4 h-4 mr-3 text-blue-400" />
              <span className="font-semibold text-blue-400">协作空间示例</span>
            </Link>

            {spaces.map((space) => (
              <Link
                key={space.id}
                to={`/space/${space.id}`}
                className={`flex items-center justify-between px-4 py-2.5 rounded-lg transition-colors group ${
                  isActive(`/space/${space.id}`) ? 'bg-slate-800 text-white' : 'hover:bg-slate-800/50 hover:text-white'
                }`}
              >
                <div className="flex items-center overflow-hidden">
                  <Briefcase className="w-4 h-4 mr-3 flex-shrink-0 text-slate-500 group-hover:text-blue-400" />
                  <span className="truncate">{space.name}</span>
                </div>
                {isActive(`/space/${space.id}`) && <ChevronRight className="w-4 h-4 text-blue-500" />}
              </Link>
            ))}
            {spaces.length === 0 && (
              <p className="px-4 py-2 text-sm text-slate-600 italic">暂无空间</p>
            )}
          </div>
        </div>
      </nav>

      <div className="p-4 border-t border-slate-800">
        <div className="flex items-center px-4 py-3 mb-4 rounded-lg bg-slate-800/50">
          <div className="w-8 h-8 rounded-full bg-blue-500 flex items-center justify-center text-white font-semibold mr-3">
            U
          </div>
          <div className="flex-1 min-w-0">
            <p className="text-sm font-medium text-white truncate">演示用户</p>
            <p className="text-xs text-slate-500 truncate">demo@example.com</p>
          </div>
        </div>
        <button
          onClick={onLogout}
          className="flex items-center w-full px-4 py-2 text-sm text-slate-400 hover:text-red-400 transition-colors"
        >
          <LogOut className="w-4 h-4 mr-3" />
          退出登录
        </button>
      </div>
    </aside>
  );
};

export default Sidebar;
