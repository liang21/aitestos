import { useState } from 'react'
import { Button, Card, Input, Select, Space } from '@arco-design/web-react'
import { IconPlus, IconSearch } from '@arco-design/web-react/icon'
import { useSearchParams, useNavigate } from 'react-router-dom'
import { SearchTable } from '@/components/business/SearchTable'
import { StatusTag } from '@/components/business/StatusTag'
import { useDocumentList } from '../hooks/useDocuments'
import { UploadDocumentModal } from './UploadDocumentModal'
import type { Document, DocumentType, DocumentStatus } from '@/types/api'
import type { Columns } from '@arco-design/web-react/es/table'

const { Option } = Select

// Document type options
const DOCUMENT_TYPE_OPTIONS = [
  { label: 'PRD', value: 'prd' },
  { label: 'API Spec', value: 'api_spec' },
  { label: 'Swagger', value: 'swagger' },
  { label: 'Figma', value: 'figma' },
  { label: 'Markdown', value: 'markdown' },
]

// Document status options
const DOCUMENT_STATUS_OPTIONS = [
  { label: '待处理', value: 'pending' },
  { label: '解析中', value: 'processing' },
  { label: '已完成', value: 'completed' },
  { label: '失败', value: 'failed' },
]

export function KnowledgeListPage() {
  const [searchParams] = useSearchParams()
  const navigate = useNavigate()
  const projectId = searchParams.get('projectId') || ''

  // Local state for filters
  const [typeFilter, setTypeFilter] = useState<string>('')
  const [statusFilter, setStatusFilter] = useState<string>('')
  const [page, setPage] = useState(1)
  const [uploadModalVisible, setUploadModalVisible] = useState(false)

  // Query document list
  const { data, isLoading } = useDocumentList({
    projectId,
    type: typeFilter || undefined,
    status: statusFilter || undefined,
    offset: (page - 1) * 10,
    limit: 10,
  })

  // Table columns
  const columns: Columns<Document> = [
    {
      title: '文档名称',
      dataIndex: 'name',
      width: 300,
    },
    {
      title: '类型',
      dataIndex: 'type',
      width: 120,
      render: (type: DocumentType) => (
        <StatusTag status={type} category="document_type" />
      ),
    },
    {
      title: '状态',
      dataIndex: 'status',
      width: 120,
      render: (status: DocumentStatus) => (
        <StatusTag status={status} category="document_status" />
      ),
    },
    {
      title: '分块数量',
      dataIndex: 'chunkCount',
      width: 100,
    },
    {
      title: '上传时间',
      dataIndex: 'createdAt',
      width: 180,
      render: (date: string) => new Date(date).toLocaleString('zh-CN'),
    },
    {
      title: '操作',
      dataIndex: 'id',
      width: 100,
      render: (id: string) => (
        <Button type="text" size="small" onClick={() => navigate(id)}>
          查看
        </Button>
      ),
    },
  ]

  // Handle filter change
  const handleFilterChange = () => {
    setPage(1)
  }

  return (
    <div className="p-6">
      <Card>
        <div className="mb-4 flex items-center justify-between">
          <Space>
            <Select
              placeholder="文档类型"
              style={{ width: 150 }}
              value={typeFilter || undefined}
              onChange={setTypeFilter}
              allowClear
              onClear={() => {
                setTypeFilter('')
                handleFilterChange()
              }}
            >
              {DOCUMENT_TYPE_OPTIONS.map((opt) => (
                <Option key={opt.value} value={opt.value}>
                  {opt.label}
                </Option>
              ))}
            </Select>
            <Select
              placeholder="文档状态"
              style={{ width: 150 }}
              value={statusFilter || undefined}
              onChange={setStatusFilter}
              allowClear
              onClear={() => {
                setStatusFilter('')
                handleFilterChange()
              }}
            >
              {DOCUMENT_STATUS_OPTIONS.map((opt) => (
                <Option key={opt.value} value={opt.value}>
                  {opt.label}
                </Option>
              ))}
            </Select>
          </Space>
          <Button
            type="primary"
            icon={<IconPlus />}
            onClick={() => setUploadModalVisible(true)}
          >
            上传文档
          </Button>
        </div>

        <SearchTable
          columns={columns}
          data={data?.data ?? []}
          total={data?.total ?? 0}
          loading={isLoading}
          current={page}
          onPageChange={setPage}
        />
      </Card>

      {/* Upload Document Modal */}
      <UploadDocumentModal
        visible={uploadModalVisible}
        projectId={projectId}
        onCancel={() => setUploadModalVisible(false)}
        onSuccess={() => {
          setUploadModalVisible(false)
        }}
      />
    </div>
  )
}
