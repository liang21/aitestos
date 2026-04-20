import { Result, Button } from '@arco-design/web-react'

export function NotFoundPage() {
  return (
    <div className="flex h-screen items-center justify-center">
      <Result
        status="404"
        subTitle="Sorry, the page you visited does not exist."
      >
        <Button type="primary" onClick={() => (window.location.href = '/')}>
          Back Home
        </Button>
      </Result>
    </div>
  )
}
