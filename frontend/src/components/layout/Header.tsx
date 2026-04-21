import { Breadcrumb, Button, Dropdown, Avatar, Badge, Menu as ArcoMenu } from '@arco-design/web-react'
import { Bell, Menu as MenuIcon } from 'lucide-react'
import { useNavigate } from 'react-router-dom'
import { useAuthStore } from '@/features/auth/hooks/useAuthStore'
import { useLogout } from '@/features/auth/hooks/useAuth'
import { useAppStore } from '@/store/useAppStore'
import './header.css'

const BreadcrumbItem = Breadcrumb.Item
const MenuItem = ArcoMenu.Item

interface HeaderProps {
  title: string
  breadcrumbs?: Array<{ title: string; path: string }>
}

export function Header({ title, breadcrumbs = [] }: HeaderProps) {
  const { user } = useAuthStore()
  const handleLogout = useLogout()
  const { sidebarCollapsed, toggleSidebar } = useAppStore()
  const navigate = useNavigate()

  return (
    <header data-testid="app-header" className="app-header">
      <div className="header-left">
        {/* Collapse Button */}
        <Button
          type="text"
          icon={<MenuIcon size={18} />}
          onClick={toggleSidebar}
          className="collapse-btn"
          aria-label="折叠侧边栏"
        />

        {/* Breadcrumb */}
        {breadcrumbs.length > 0 && (
          <Breadcrumb className="breadcrumb">
            {breadcrumbs.map((crumb, index) => (
              <BreadcrumbItem key={index}>
                <button onClick={() => navigate(crumb.path)}>{crumb.title}</button>
              </BreadcrumbItem>
            ))}
            <BreadcrumbItem>{title}</BreadcrumbItem>
          </Breadcrumb>
        )}

        {/* Title when no breadcrumbs */}
        {breadcrumbs.length === 0 && (
          <h1 className="page-title">{title}</h1>
        )}
      </div>

      <div className="header-right">
        {/* Notifications */}
        <Badge count={1} dot>
          <Button
            type="text"
            icon={<Bell size={18} />}
            aria-label="通知"
          />
        </Badge>

        {/* User Dropdown */}
        <Dropdown
          trigger="click"
          position="br"
          droplist={
            <ArcoMenu onClickMenuItem={(key) => key === 'logout' && handleLogout()}>
              <MenuItem key="settings">设置</MenuItem>
              <MenuItem key="logout">退出登录</MenuItem>
            </ArcoMenu>
          }
        >
          <div className="user-info" role="button" tabIndex={0} aria-label="用户菜单">
            <Avatar size={32} style={{ cursor: 'pointer' }}>
              {user?.username?.charAt(0).toUpperCase()}
            </Avatar>
            <span className="username">{user?.username}</span>
          </div>
        </Dropdown>
      </div>
    </header>
  )
}
