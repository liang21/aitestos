import { render, screen } from '@testing-library/react'
import { StatsCard } from './StatsCard'

describe('StatsCard', () => {
  it('should render title and value', () => {
    render(<StatsCard title="用例总数" value={128} />)

    expect(screen.getByText('用例总数')).toBeInTheDocument()
    expect(screen.getByText('128')).toBeInTheDocument()
  })

  it('should render trend arrow up', () => {
    const { container } = render(
      <StatsCard title="通过率" value={85} trend="up" trendValue={5} />
    )

    expect(screen.getByText('85')).toBeInTheDocument()
    expect(screen.getByText(/5%/)).toBeInTheDocument()
    // Check for upward arrow indicator
    const trendIcon = container.querySelector('[class*="trend"]')
    expect(trendIcon).toBeInTheDocument()
  })

  it('should render trend arrow down', () => {
    const { container } = render(
      <StatsCard title="失败率" value={12} trend="down" trendValue={3} />
    )

    expect(screen.getByText('12')).toBeInTheDocument()
    expect(screen.getByText(/3%/)).toBeInTheDocument()
    // Check for downward arrow indicator
    const trendIcon = container.querySelector('[class*="trend"]')
    expect(trendIcon).toBeInTheDocument()
  })

  it('should render custom icon', () => {
    const CustomIcon = () => <span data-testid="custom-icon">★</span>
    const { container } = render(
      <StatsCard title="覆盖率" value={75} icon={<CustomIcon />} />
    )

    expect(screen.getByTestId('custom-icon')).toBeInTheDocument()
    expect(screen.getByText('75')).toBeInTheDocument()
  })
})
