/**
 * UI messages for generation module
 * Centralized text constants for consistency and easier localization
 */

export const GENERATION_MESSAGES = {
  // Page titles
  pageTitle: {
    taskList: 'AI 生成任务',
    newTask: '新建 AI 生成任务',
    taskDetail: '任务详情',
  },

  // Descriptions
  description: {
    newTask: '描述测试需求，AI 将自动生成测试用例草稿',
    promptPlaceholder: '请描述测试需求，例如：测试用户注册功能，包括邮箱验证和密码强度校验',
    caseCountHelp: '生成 1-20 个测试用例',
  },

  // Buttons
  button: {
    createTask: '新建任务',
    submit: '立即生成',
    cancel: '取消',
    backToList: '返回列表',
    viewDetail: '查看详情',
    confirmInsufficient: '继续生成',
  },

  // Knowledge readiness
  knowledge: {
    sufficient: '🟢 就绪',
    insufficient: '🟡 内容有限',
    empty: '🔴 请先上传需求文档',
    documentCount: (count: number) => `📄 ${count} 份文档`,
    warning: '知识库内容较少，生成质量可能较低。建议上传更多需求文档后再试。',
    error: '暂无需求文档，请先上传 PRD、API 文档或 Figma 设计稿。',
    submitDisabled: '请先上传需求文档后再创建生成任务',
    degradationWarning: '知识库内容不足，生成质量可能较低。建议上传更多需求文档后再试。',
    degradationNote: '继续生成时，系统将使用较低的置信度设置，可能导致生成的测试用例质量下降。',
  },

  // Task status
  task: {
    status: {
      pending: '待处理',
      processing: '生成中',
      completed: '已完成',
      failed: '失败',
    },
    processing: {
      pending: '任务排队中...',
      processing: 'AI 正在生成用例，请稍候...',
      autoRefresh: (seconds: number) => `系统每 ${seconds} 秒自动刷新状态`,
    },
    failed: '生成失败',
  },

  // Form fields
  form: {
    targetModule: '目标模块',
    prompt: '需求描述',
    caseCount: '用例数量',
    advancedOptions: '高级选项（场景类型、优先级、用例类型）',
    sceneType: '场景类型',
    priority: '优先级',
    caseType: '用例类型',
  },

  // Scene types
  sceneType: {
    positive: '正向测试',
    negative: '负向测试',
    boundary: '边界测试',
  },

  // Priorities
  priority: {
    P0: 'P0 紧急',
    P1: 'P1 高',
    P2: 'P2 中',
    P3: 'P3 低',
  },

  // Case types
  caseType: {
    functionality: '功能测试',
    performance: '性能测试',
    api: 'API 测试',
    ui: 'UI 测试',
    security: '安全测试',
  },

  // Table columns
  table: {
    prompt: '需求描述',
    status: '状态',
    draftCount: '草稿数',
    createdAt: '创建时间',
    title: '标题',
    type: '类型',
    priority: '优先级',
    confidence: '置信度',
    action: '操作',
  },

  // Confidence levels
  confidence: {
    high: '高',
    medium: '中',
    low: '低',
  },

  // Empty states
  empty: {
    taskNotFound: '任务不存在或已被删除',
    noDrafts: '未生成任何草稿',
    draftsLoadFailed: '获取草稿列表失败，请刷新重试',
    taskListLoadFailed: '获取任务列表失败，请稍后重试',
  },

  // Success messages
  success: {
    taskCreated: '生成任务创建成功',
  },

  // Error messages
  error: {
    taskCreateFailed: '创建任务失败，请重试',
  },

  // Modal titles
  modal: {
    insufficientKnowledge: '知识库内容不足',
  },

  // Filters
  filter: {
    statusPlaceholder: '筛选状态',
  },

  // Drafts
  drafts: {
    title: '生成的草稿',
    count: (count: number) => `${count} 条`,
  },

  // Labels
  label: {
    taskId: '任务ID',
    status: '状态',
    promptDescription: '需求描述',
    createdAt: '创建时间',
    updatedAt: '更新时间',
    draftCount: '生成草稿数',
    errorMessage: '错误信息',
  },
} as const
