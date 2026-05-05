// 无限循环诊断脚本 - 在浏览器控制台运行
// 修改 React 的 useState 来追踪哪个组件导致了问题

(() => {
  const originalUseState = React.useState;
  let renderCount = 0;
  const componentStack = [];

  React.useState = function(initialState) {
    // 获取调用堆栈
    const stack = new Error().stack;
    const componentMatch = stack.match(/at (\w+) \(/);
    const componentName = componentMatch ? componentMatch[1] : 'Unknown';

    if (componentStack.length === 0 || componentStack[componentStack.length - 1] !== componentName) {
      componentStack.push(componentName);
    }

    renderCount++;

    if (renderCount > 100) {
      console.error('=== Infinite Loop Detected ===');
      console.error('Component stack:', [...new Set(componentStack)]);
      console.error('Total renders:', renderCount);
      console.error('Current state:', initialState);
      console.trace('Full stack trace');
      throw new Error('Stopping to prevent browser crash');
    }

    return originalUseState(initialState);
  };

  console.log('Diagnostic started. Open the URL and wait for the error.');
  console.log('The script will stop after 100 renders and show which component caused it.');
})();
