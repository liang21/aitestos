import { Layout } from '@arco-design/web-react'
import { Outlet } from 'react-router-dom'
import { useAppStore } from '@/store/useAppStore'
import { Header } from './Header'
import { Sidebar } from './Sidebar'
import './layout.css'

const { Sider, Content } = Layout

interface AppLayoutProps {
  title?: string
  breadcrumbs?: Array<{ title: string; path: string }>
  children?: React.ReactNode
}

export function AppLayout({
  title = 'Aitestos',
  breadcrumbs = [],
  children,
}: AppLayoutProps) {
  const { sidebarCollapsed } = useAppStore()

  return (
    <Layout className="app-layout">
      <Sider
        width={220}
        collapsedWidth={64}
        collapsed={sidebarCollapsed}
        className="app-sider"
      >
        <Sidebar />
      </Sider>
      <Layout>
        <Header title={title} breadcrumbs={breadcrumbs} />
        <Content className="app-content">
          {children ?? <Outlet />}
        </Content>
      </Layout>
    </Layout>
  )
}
