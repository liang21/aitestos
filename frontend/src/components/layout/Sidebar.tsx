import { Menu, Dropdown, Avatar, Badge, Button } from '@arco-design/web-react'
import {
  PanelLeftClose,
  PanelLeftOpen,
  Inbox,
  FolderOpen,
  Settings,
  BarChart3,
  PlayCircle,
  FileText,
  BookOpen,
  ListTodo,
} from 'lucide-react'
import { NavLink, useLocation } from 'react-router-dom'
import { useMemo } from 'react'
import { useAppStore } from '@/store/useAppStore'
import { usePendingDraftCount } from '@/features/drafts/hooks/useDrafts'
import { useAuthStore } from '@/features/auth/hooks/useAuthStore'
import { useLogout } from '@/features/auth/hooks/useAuth'
import './sidebar.css'

const MenuItem = Menu.Item

export function Sidebar() {
  const { sidebarCollapsed, toggleSidebar } = useAppStore()
  const { data: pendingCount } = usePendingDraftCount()
  const { user } = useAuthStore()
  const handleLogout = useLogout()
  const location = useLocation()

  // Memoize menu items to prevent unnecessary re-renders
  const menuItems = useMemo(
    () => (
      <>
        {/* Projects */}
        <MenuItem key="/projects">
          <NavLink to="/projects" aria-label="项目列表">
            <FolderOpen size={18} aria-hidden="true" />
            {!sidebarCollapsed && <span>项目列表</span>}
          </NavLink>
        </MenuItem>

        {/* Drafts with Badge */}
        <MenuItem key="/drafts">
          <NavLink to="/drafts" aria-label="草稿箱">
            <Badge count={pendingCount ?? 0} offset={[8, 0]}>
              <Inbox size={18} aria-hidden="true" />
            </Badge>
            {!sidebarCollapsed && <span>草稿箱</span>}
          </NavLink>
        </MenuItem>

        {/* Test Cases */}
        <MenuItem key="/testcases">
          <NavLink to="/testcases" aria-label="测试用例">
            <FileText size={18} aria-hidden="true" />
            {!sidebarCollapsed && <span>测试用例</span>}
          </NavLink>
        </MenuItem>

        {/* Test Plans */}
        <MenuItem key="/plans">
          <NavLink to="/plans" aria-label="测试计划">
            <ListTodo size={18} aria-hidden="true" />
            {!sidebarCollapsed && <span>测试计划</span>}
          </NavLink>
        </MenuItem>

        {/* AI Generation */}
        <MenuItem key="/generation">
          <NavLink to="/generation" aria-label="AI 生成">
            <PlayCircle size={18} aria-hidden="true" />
            {!sidebarCollapsed && <span>AI 生成</span>}
          </NavLink>
        </MenuItem>

        {/* Documents */}
        <MenuItem key="/documents">
          <NavLink to="/documents" aria-label="知识库">
            <BookOpen size={18} aria-hidden="true" />
            {!sidebarCollapsed && <span>知识库</span>}
          </NavLink>
        </MenuItem>

        {/* Dashboard */}
        <MenuItem key="/dashboard">
          <NavLink to="/dashboard" aria-label="仪表盘">
            <BarChart3 size={18} aria-hidden="true" />
            {!sidebarCollapsed && <span>仪表盘</span>}
          </NavLink>
        </MenuItem>
      </>
    ),
    [sidebarCollapsed, pendingCount] // Only re-render when these change
  )

  return (
    <div
      data-testid="app-sidebar"
      className={`sidebar ${sidebarCollapsed ? 'collapsed' : ''}`}
    >
      {/* Collapse Button */}
      <div className="sidebar-header">
        <Button
          type="text"
          icon={
            sidebarCollapsed ? (
              <PanelLeftOpen size={16} />
            ) : (
              <PanelLeftClose size={16} />
            )
          }
          onClick={toggleSidebar}
          className="collapse-btn"
        />
      </div>

      {/* Main Menu */}
      <Menu
        selectedKeys={[location.pathname]}
        style={{ width: '100%' }}
        className="sidebar-menu"
      >
        {menuItems}
      </Menu>

      {/* User Section */}
      {!sidebarCollapsed && (
        <div className="sidebar-footer">
          <Dropdown
            trigger="click"
            position="top"
            droplist={
              <Menu onClickMenuItem={handleLogout}>
                <MenuItem key="settings">
                  <Settings size={14} />
                  <span>设置</span>
                </MenuItem>
                <MenuItem key="logout">
                  <span>退出登录</span>
                </MenuItem>
              </Menu>
            }
          >
            <div className="user-info">
              <Avatar size={32} style={{ marginRight: 8 }}>
                {user?.username?.charAt(0).toUpperCase()}
              </Avatar>
              <span className="username">{user?.username}</span>
            </div>
          </Dropdown>
        </div>
      )}

      {/* Collapsed User Avatar */}
      {sidebarCollapsed && (
        <div className="sidebar-footer collapsed">
          <Dropdown
            trigger="click"
            position="top"
            droplist={
              <Menu onClickMenuItem={handleLogout}>
                <MenuItem key="settings">设置</MenuItem>
                <MenuItem key="logout">退出登录</MenuItem>
              </Menu>
            }
          >
            <Avatar size={32}>{user?.username?.charAt(0).toUpperCase()}</Avatar>
          </Dropdown>
        </div>
      )}
    </div>
  )
}
