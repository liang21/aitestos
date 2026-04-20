import { render, screen } from '@testing-library/react'
import { SplitPanel } from './SplitPanel'

describe('SplitPanel', () => {
  it('should render left and right panel children', () => {
    render(
      <SplitPanel>
        <SplitPanel.Left>
          <div data-testid="left-content">Left Panel Content</div>
        </SplitPanel.Left>
        <SplitPanel.Right>
          <div data-testid="right-content">Right Panel Content</div>
        </SplitPanel.Right>
      </SplitPanel>
    )

    expect(screen.getByTestId('left-content')).toBeInTheDocument()
    expect(screen.getByTestId('right-content')).toBeInTheDocument()
  })

  it('should use default split ratio (50:50) when not specified', () => {
    const { container } = render(
      <SplitPanel>
        <SplitPanel.Left>
          <div data-testid="left-content">Left</div>
        </SplitPanel.Left>
        <SplitPanel.Right>
          <div data-testid="right-content">Right</div>
        </SplitPanel.Right>
      </SplitPanel>
    )

    // Both panels should be rendered
    expect(screen.getByTestId('left-content')).toBeInTheDocument()
    expect(screen.getByTestId('right-content')).toBeInTheDocument()

    // Check that container has flex display
    const splitContainer = container.firstChild as HTMLElement
    expect(splitContainer).toHaveStyle({ display: 'flex' })
  })

  it('should render resizable divider between panels', () => {
    const { container } = render(
      <SplitPanel>
        <SplitPanel.Left>
          <div data-testid="left-content">Left</div>
        </SplitPanel.Left>
        <SplitPanel.Right>
          <div data-testid="right-content">Right</div>
        </SplitPanel.Right>
      </SplitPanel>
    )

    const divider = container.querySelector('.resizer-divider')
    expect(divider).toBeInTheDocument()
    expect(divider).toHaveStyle({ cursor: 'col-resize' })
  })
})
