/**
 * E2E Testing: data-testid 添加指南
 *
 * 本文件展示如何为组件添加 data-testid 以支持 E2E 测试
 */

/* ============================================
   1. 登录表单示例 (LoginPage.tsx)
   ============================================ */

// ❌ 之前：没有 data-testid
<form onSubmit={handleSubmit}>
  <Input placeholder="邮箱" />
  <InputPassword placeholder="密码" />
  <Button type="primary">登录</Button>
</form>

// ✅ 改进：添加 data-testid
<form onSubmit={handleSubmit} data-testid="login-form">
  <Input
    placeholder="邮箱"
    data-testid="login-email-input"
  />
  <InputPassword
    placeholder="密码"
    data-testid="login-password-input"
  />
  <Button
    type="primary"
    htmlType="submit"
    data-testid="login-submit-button"
  >
    登录
  </Button>
</form>

/* ============================================
   2. 项目列表表格示例 (ProjectListPage.tsx)
   ============================================ */

// ❌ 之前
<Table>
  <Column title="项目名称" dataIndex="name" />
  <Column title="前缀" dataIndex="prefix" />
</Table>

// ✅ 改进：添加 data-testid 到表格和单元格
<Table data-testid="project-list-table">
  <Column
    title="项目名称"
    dataIndex="name"
    render={(name) => (
      <span data-project-name={name}>{name}</span>
    )}
  />
  <Column title="前缀" dataIndex="prefix" />
</Table>

/* ============================================
   3. 按钮组示例
   ============================================ */

// ❌ 之前
<Button onClick={onCreate}>新建项目</Button>
<Input.Search placeholder="搜索项目" />

// ✅ 改进
<Button
  onClick={onCreate}
  data-testid="new-project-button"
>
  新建项目
</Button>
<Input.Search
  placeholder="搜索项目"
  data-testid="project-search-input"
/>

/* ============================================
   4. 模态框示例 (CreateProjectModal.tsx)
   ============================================ */

// ❌ 之前
<Modal visible={visible} title="新建项目">
  <Form>
    <Form.Item label="项目名称">
      <Input />
    </Form.Item>
  </Form>
</Modal>

// ✅ 改进
<Modal
  visible={visible}
  title="新建项目"
  data-testid="create-project-modal"
>
  <Form data-testid="create-project-form">
    <Form.Item label="项目名称">
      <Input data-testid="project-name-input" />
    </Form.Item>
  </Form>
</Modal>

/* ============================================
   5. 状态指示器示例
   ============================================ */

// ❌ 之前
<Tag color="green">通过</Tag>
<Tag color="red">失败</Tag>

// ✅ 改进
<Tag
  color="green"
  data-testid="status-pass"
  data-status="pass"
>
  通过
</Tag>
<Tag
  color="red"
  data-testid="status-fail"
  data-status="fail"
>
  失败
</Tag>

/* ============================================
   6. 导航菜单示例 (Sidebar.tsx)
   ============================================ */

// ❌ 之前
<Menu>
  <MenuItem key="/projects">项目列表</MenuItem>
  <MenuItem key="/testcases">测试用例</MenuItem>
</Menu>

// ✅ 改进
<Menu data-testid="app-sidebar-menu">
  <MenuItem
    key="/projects"
    data-testid="menu-projects"
  >
    项目列表
  </MenuItem>
  <MenuItem
    key="/testcases"
    data-testid="menu-testcases"
  >
    测试用例
  </MenuItem>
</Menu>

/* ============================================
   7. 命名规范总结
   ============================================ */

/**
 * data-testid 命名规范：
 *
 * 1. 使用 kebab-case (小写 + 连字符)
 *    ✅ login-submit-button
 *    ❌ loginSubmitButton
 *
 * 2. 描述性名称，基于功能而非实现
 *    ✅ email-input, project-table
 *    ❌ input1, arco-table-2
 *
 * 3. 组件级别命名
 *    - 按钮：xxx-button, xxx-submit-button, xxx-cancel-button
 *    - 输入：xxx-input, xxx-email-input
 *    - 模态框：xxx-modal
 *    - 表格：xxx-table
 *    - 菜单：xxx-menu, xxx-menu-item
 *    - 状态：xxx-status, data-status (用于断言)
 *
 * 4. 列表项使用 data 属性
 *    ✅ data-project-name="ECommerce"
 *    ✅ data-user-id="123"
 *
 * 5. 表单字段命名
 *    ✅ project-name-input
 *    ✅ user-email-input
 *    ❌ name-input (太通用)
 */

/* ============================================
   8. 迁移检查清单
   ============================================ */

/**
 * 组件 data-testid 添加检查清单：
 *
 * [ ] 表单：所有输入框添加 data-testid
 * [ ] 按钮：submit/cancel 按钮添加 data-testid
 * [ ] 模态框：modal 容器添加 data-testid
 * [ ] 表格：table 容器添加 data-testid
 * [ ] 列表：列表容器添加 data-testid
 * [ ] 导航：菜单项添加 data-testid
 * [ ] 状态：状态标签添加 data-status 属性
 * [ ] 错误：错误消息添加 data-testid
 * [ ] 加载：loading 状态添加 data-testid
 *
 * 优先级：
 * 🔴 高：认证、项目创建、AI 生成核心流程
 * 🟡 中：测试用例、测试计划相关
 * 🟢 低：统计图表、配置页面
 */
