import { render, screen } from '@testing-library/react'
import { StatusTag } from './StatusTag'

describe('StatusTag', () => {
  describe('case_status category', () => {
    it('should render correct color for unexecuted status', () => {
      const { container } = render(
        <StatusTag status="unexecuted" category="case_status" />
      )
      const tag = container.querySelector('[class*="arco-tag"]') as HTMLElement
      expect(tag).toHaveStyle({ color: '#86909C' })
      expect(tag).toHaveStyle({ backgroundColor: 'rgba(134,144,156,0.10)' })
      expect(tag).toHaveTextContent('未执行')
    })

    it('should render correct color for pass status', () => {
      const { container } = render(<StatusTag status="pass" category="case_status" />)
      const tag = container.querySelector('[class*="arco-tag"]') as HTMLElement
      expect(tag).toHaveStyle({ color: '#00B42A' })
      expect(tag).toHaveStyle({ backgroundColor: 'rgba(0,180,42,0.10)' })
      expect(tag).toHaveTextContent('通过')
    })

    it('should render correct color for block status', () => {
      const { container } = render(
        <StatusTag status="block" category="case_status" />
      )
      const tag = container.querySelector('[class*="arco-tag"]') as HTMLElement
      expect(tag).toHaveStyle({ color: '#FF7D00' })
      expect(tag).toHaveStyle({ backgroundColor: 'rgba(255,125,0,0.10)' })
      expect(tag).toHaveTextContent('阻塞')
    })

    it('should render correct color for fail status', () => {
      const { container } = render(<StatusTag status="fail" category="case_status" />)
      const tag = container.querySelector('[class*="arco-tag"]') as HTMLElement
      expect(tag).toHaveStyle({ color: '#F53F3F' })
      expect(tag).toHaveStyle({ backgroundColor: 'rgba(245,63,63,0.10)' })
      expect(tag).toHaveTextContent('失败')
    })
  })

  describe('priority category', () => {
    it('should render correct color for P0 priority', () => {
      const { container } = render(<StatusTag status="P0" category="priority" />)
      const tag = container.querySelector('[class*="arco-tag"]') as HTMLElement
      expect(tag).toHaveStyle({ color: '#F53F3F' })
      expect(tag).toHaveTextContent('P0 紧急')
    })

    it('should render correct color for P1 priority', () => {
      const { container } = render(<StatusTag status="P1" category="priority" />)
      const tag = container.querySelector('[class*="arco-tag"]') as HTMLElement
      expect(tag).toHaveStyle({ color: '#FF7D00' })
      expect(tag).toHaveTextContent('P1 高')
    })

    it('should render correct color for P2 priority', () => {
      const { container } = render(<StatusTag status="P2" category="priority" />)
      const tag = container.querySelector('[class*="arco-tag"]') as HTMLElement
      expect(tag).toHaveStyle({ color: '#7B61FF' })
      expect(tag).toHaveTextContent('P2 中')
    })

    it('should render correct color for P3 priority', () => {
      const { container } = render(<StatusTag status="P3" category="priority" />)
      const tag = container.querySelector('[class*="arco-tag"]') as HTMLElement
      expect(tag).toHaveStyle({ color: '#86909C' })
      expect(tag).toHaveTextContent('P3 低')
    })
  })

  describe('edge cases', () => {
    it('should return null for non-existent status', () => {
      const { container } = render(
        <StatusTag status="invalid" category="case_status" />
      )
      expect(container.firstChild).toBeNull()
    })

    it('should use custom label when provided', () => {
      const { container } = render(
        <StatusTag status="pass" category="case_status" label="成功" />
      )
      const tag = container.querySelector('[class*="arco-tag"]') as HTMLElement
      expect(tag).toHaveTextContent('成功')
    })
  })
})
