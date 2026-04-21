import { Result, Button } from '@arco-design/web-react'
import { useNavigate } from 'react-router-dom'

export function NotFoundPage() {
  const navigate = useNavigate()

  return (
    <div className="flex h-screen items-center justify-center">
      <Result status="404" subTitle="抱歉，您访问的页面不存在。">
        <Button type="primary" onClick={() => navigate('/projects')}>
          返回首页
        </Button>
      </Result>
    </div>
  )
}
