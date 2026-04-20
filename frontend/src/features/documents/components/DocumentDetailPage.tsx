import { useParams } from 'react-router-dom'
import { Card, List, Spin, Tag } from '@arco-design/web-react'
import { IconLoading } from '@arco-design/web-react/icon'
import { useDocumentDetail } from '../hooks/useDocuments'
import { StatusTag } from '@/components/business/StatusTag'
import type { DocumentChunk } from '@/types/api'

export function DocumentDetailPage() {
  const { documentId } = useParams<{ documentId: string }>()
  const { data: document, isLoading } = useDocumentDetail(documentId || '')

  if (!documentId) {
    return (
      <div className="p-6">
        <Card>文档 ID 无效</Card>
      </div>
    )
  }

  if (isLoading) {
    return (
      <div className="p-6 flex justify-center items-center" style={{ minHeight: 400 }}>
        <Spin icon={<IconLoading spin />} tip="加载中..." />
      </div>
    )
  }

  if (!document) {
    return (
      <div className="p-6">
        <Card>文档加载失败</Card>
      </div>
    )
  }

  return (
    <div className="p-6">
      <Card>
        <div className="mb-4">
          <h2 className="text-xl font-semibold">{document.name}</h2>
        </div>

        {/* Document Info */}
        <div className="grid grid-cols-2 gap-4 mb-6">
          <div className="flex items-center gap-2">
            <span className="text-gray-600">文档类型:</span>
            <StatusTag status={document.type} category="document_type" />
          </div>
          <div className="flex items-center gap-2">
            <span className="text-gray-600">处理状态:</span>
            <StatusTag status={document.status} category="document_status" />
          </div>
          <div className="flex items-center gap-2">
            <span className="text-gray-600">分块数量:</span>
            <span>{document.chunkCount}</span>
          </div>
          <div className="flex items-center gap-2">
            <span className="text-gray-600">上传者:</span>
            <span>{document.uploadedByName || '-'}</span>
          </div>
          <div className="flex items-center gap-2">
            <span className="text-gray-600">创建时间:</span>
            <span>{new Date(document.createdAt).toLocaleString('zh-CN')}</span>
          </div>
          <div className="flex items-center gap-2">
            <span className="text-gray-600">更新时间:</span>
            <span>{new Date(document.updatedAt).toLocaleString('zh-CN')}</span>
          </div>
        </div>

        {/* Show processing indicator if status is processing */}
        {document.status === 'processing' && (
          <div className="mt-4 p-4 bg-blue-50 rounded flex items-center gap-2">
            <Spin icon={<IconLoading spin />} />
            <span>文档正在解析中，请稍候...</span>
          </div>
        )}
      </Card>

      {/* Document Chunks */}
      <Card className="mt-6" title="文档分块">
        {document.chunks && document.chunks.length > 0 ? (
          <List
            dataSource={document.chunks}
            render={(item, index) => (
              <List.Item key={item.id}>
                <div className="flex gap-4">
                  <Tag color="blue" className="shrink-0">
                    #{index + 1}
                  </Tag>
                  <div className="flex-1">
                    <p className="text-sm leading-relaxed whitespace-pre-wrap">
                      {item.content}
                    </p>
                  </div>
                </div>
              </List.Item>
            )}
          />
        ) : (
          <div className="text-center text-gray-400 py-8">
            {document.status === 'processing' ? '文档解析中，分块将在完成后显示' : '暂无分块数据'}
          </div>
        )}
      </Card>
    </div>
  )
}
