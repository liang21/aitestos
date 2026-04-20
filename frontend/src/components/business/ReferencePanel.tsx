import { Empty, List, Space, Tag, Typography } from '@arco-design/web-react'
import type { ReferencedChunk } from '@/types/api'

const { Text, Title } = Typography

export interface ReferencePanelProps {
  /** List of referenced document chunks */
  chunks: ReferencedChunk[]
  /** Additional class name */
  className?: string
  /** Custom style */
  style?: React.CSSProperties
  /** Show similarity score (default true) */
  showScore?: boolean
}

/**
 * Panel displaying referenced document chunks
 * Shows document title, similarity score, and link to original
 */
export function ReferencePanel({
  chunks,
  className,
  style,
  showScore = true,
}: ReferencePanelProps) {
  if (chunks.length === 0) {
    return (
      <div className={className} style={style}>
        <Empty description="无引用来源" />
      </div>
    )
  }

  // Group chunks by document
  const groupedByDocument = chunks.reduce<Record<string, ReferencedChunk[]>>(
    (acc, chunk) => {
      if (!acc[chunk.documentTitle]) {
        acc[chunk.documentTitle] = []
      }
      acc[chunk.documentTitle].push(chunk)
      return acc
    },
    {}
  )

  return (
    <div className={className} style={style}>
      <Title heading={6} style={{ marginBottom: 12 }}>
        引用来源 ({chunks.length})
      </Title>

      <Space direction="vertical" size="medium" style={{ width: '100%' }}>
        {Object.entries(groupedByDocument).map(([docTitle, docChunks]) => (
          <div key={docTitle}>
            <Text bold style={{ display: 'block', marginBottom: 8 }}>
              {docTitle}
            </Text>
            <List
              dataSource={docChunks}
              render={(chunk) => (
                <div
                  key={chunk.chunkId}
                  className="flex items-center justify-between py-2 px-3 bg-gray-50 rounded"
                >
                  <Space>
                    <Text type="secondary">#{chunk.chunkId.slice(0, 8)}</Text>
                  </Space>
                  {showScore && (
                    <Tag color="green" size="small">
                      相似度 {Math.round(chunk.similarityScore * 100)}%
                    </Tag>
                  )}
                </div>
              )}
            />
          </div>
        ))}
      </Space>
    </div>
  )
}
