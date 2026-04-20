import { render, screen } from '@testing-library/react'
import { SearchTable } from './SearchTable'

interface MockData {
  id: string
  name: string
  status: string
}

const mockColumns = [
  {
    title: 'ID',
    dataIndex: 'id' as const,
    width: 100,
  },
  {
    title: '名称',
    dataIndex: 'name' as const,
  },
  {
    title: '状态',
    dataIndex: 'status' as const,
  },
]

const mockData: MockData[] = [
  { id: '1', name: '测试用例1', status: 'pending' },
  { id: '2', name: '测试用例2', status: 'completed' },
]

describe('SearchTable', () => {
  it('should render table columns and data rows', () => {
    const { container } = render(
      <SearchTable<MockData>
        columns={mockColumns}
        data={mockData}
        total={2}
        loading={false}
      />
    )

    // Check column headers
    expect(screen.getByText('ID')).toBeInTheDocument()
    expect(screen.getByText('名称')).toBeInTheDocument()
    expect(screen.getByText('状态')).toBeInTheDocument()

    // Check data rows
    expect(screen.getByText('测试用例1')).toBeInTheDocument()
    expect(screen.getByText('测试用例2')).toBeInTheDocument()
  })

  it('should display Spin component when loading is true', () => {
    const { container } = render(
      <SearchTable<MockData>
        columns={mockColumns}
        data={mockData}
        total={2}
        loading={true}
      />
    )

    const spinElement = container.querySelector('[class*="arco-spin"]')
    expect(spinElement).toBeInTheDocument()
  })

  it('should render pagination controls', () => {
    const onPageChange = vi.fn()

    const { container } = render(
      <SearchTable<MockData>
        columns={mockColumns}
        data={mockData}
        total={20}
        loading={false}
        onPageChange={onPageChange}
      />
    )

    const paginationElement = container.querySelector(
      '[class*="arco-pagination"]'
    )
    expect(paginationElement).toBeInTheDocument()
  })

  it('should pass columns and data correctly to display', () => {
    render(
      <SearchTable<MockData>
        columns={mockColumns}
        data={mockData}
        total={2}
        loading={false}
      />
    )

    // Verify first row data
    expect(screen.getByText('1')).toBeInTheDocument()
    expect(screen.getByText('pending')).toBeInTheDocument()

    // Verify second row data
    expect(screen.getByText('2')).toBeInTheDocument()
    expect(screen.getByText('completed')).toBeInTheDocument()
  })
})
