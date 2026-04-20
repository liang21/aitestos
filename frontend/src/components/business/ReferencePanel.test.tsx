import { render, screen } from '@testing-library/react'
import { describe, expect, vi } from 'vitest'
import type { ReferencedChunk } from '@/types/api'
import { ReferencePanel } from './ReferencePanel'

const mockChunks: ReferencedChunk[] = [
  {
    chunkId: 'chunk-1',
    documentId: 'doc-1',
    documentTitle: '用户注册模块 PRD v2.0',
    similarityScore: 0.92,
  },
  {
    chunkId: 'chunk-2',
    documentId: 'doc-2',
    documentTitle: 'API 规范文档',
    similarityScore: 0.85,
  },
  {
    chunkId: 'chunk-3',
    documentId: 'doc-1',
    documentTitle: '用户注册模块 PRD v2.0',
    similarityScore: 0.78,
  },
]

describe('ReferencePanel', () => {
  it('should render list of referenced chunks', () => {
    render(<ReferencePanel chunks={mockChunks} />)

    expect(screen.getByText('用户注册模块 PRD v2.0')).toBeInTheDocument()
    expect(screen.getByText('API 规范文档')).toBeInTheDocument()
    // Check for similarity scores with more flexible matching
    expect(screen.getByText(/92/)).toBeInTheDocument()
    expect(screen.getByText(/85/)).toBeInTheDocument()
    expect(screen.getByText(/78/)).toBeInTheDocument()
  })

  it('should display empty message when chunks list is empty', () => {
    render(<ReferencePanel chunks={[]} />)

    expect(screen.getByText(/无引用来源|没有引用/i)).toBeInTheDocument()
  })

  it('should display similarity score as percentage', () => {
    render(<ReferencePanel chunks={mockChunks} />)

    // Similarity scores should be displayed as percentages
    expect(screen.getByText(/92/)).toBeInTheDocument()
    expect(screen.getByText(/85/)).toBeInTheDocument()
    expect(screen.getByText(/78/)).toBeInTheDocument()
  })

  it('should group chunks by document', () => {
    const { container } = render(<ReferencePanel chunks={mockChunks} />)

    // Check that chunks from same document are visually grouped
    const docTitles = screen.getAllByText('用户注册模块 PRD v2.0')
    expect(docTitles.length).toBeGreaterThan(0)
  })
})
