// Package ierrors defines HTTP error response structures
package ierrors

import "encoding/json"

// ErrorResponse is the unified API error response
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	TraceID string `json:"traceId,omitempty"`
}

// ToJSON converts ErrorResponse to JSON bytes
func (e *ErrorResponse) ToJSON() []byte {
	data, _ := json.Marshal(e)
	return data
}

// NewErrorResponse creates a new error response
func NewErrorResponse(code int, traceID string) *ErrorResponse {
	return &ErrorResponse{
		Code:    code,
		Message: CodeToMessage(code),
		TraceID: traceID,
	}
}

// CodeToMessage maps error codes to human-readable messages
func CodeToMessage(code int) string {
	messages := map[int]string{
		// Identity
		CodeUserNotFound:      "用户不存在",
		CodeEmailDuplicate:    "邮箱已被注册",
		CodeUsernameDuplicate: "用户名已被占用",
		CodePasswordMismatch:  "密码错误",
		CodeInvalidEmail:      "邮箱格式无效",
		CodePermissionDenied:  "权限不足",

		// Project
		CodeProjectNotFound:        "项目不存在",
		CodeProjectNameDuplicate:   "项目名称已存在",
		CodeProjectPrefixDuplicate: "项目前缀已存在",
		CodeInvalidProjectPrefix:   "项目前缀格式无效，需2-4位大写字母",
		CodeModuleNotFound:         "模块不存在",
		CodeModuleNameDuplicate:    "模块名称重复",
		CodeModuleAbbrevDuplicate:  "模块缩写重复",
		CodeInvalidModuleAbbrev:    "模块缩写格式无效，需2-4位大写字母",
		CodeConfigNotFound:         "配置项不存在",
		CodeConfigKeyDuplicate:     "配置键重复",

		// Knowledge
		CodeDocumentNotFound:        "文档不存在",
		CodeDocumentParseFailed:     "文档解析失败",
		CodeUnsupportedDocumentType: "不支持的文档类型",
		CodeEmptyChunks:             "文档分块为空",
		CodeEmbeddingFailed:         "向量化失败",
		CodeVectorSearchFailed:      "向量检索失败",
		CodeKnowledgeBaseEmpty:      "知识库为空，请先上传文档",
		CodeDocumentProcessing:      "文档处理中，请稍后",

		// TestCase
		CodeCaseNotFound:        "测试用例不存在",
		CodeCaseNumberDuplicate: "用例编号已存在",
		CodeInvalidCaseNumber:   "用例编号格式无效",
		CodeEmptySteps:          "用例步骤不能为空",
		CodeInvalidPriority:     "优先级无效",
		CodeInvalidCaseType:     "用例类型无效",

		// TestPlan
		CodePlanNotFound:       "测试计划不存在",
		CodePlanNameDuplicate:  "计划名称已存在",
		CodePlanArchived:       "计划已归档",
		CodeResultNotFound:     "执行结果不存在",
		CodeCaseNotInPlan:      "用例未关联到计划",
		CodeDuplicateExecution: "重复执行",

		// Generation
		CodeTaskNotFound:           "生成任务不存在",
		CodeTaskAlreadyProcessed:   "任务已处理",
		CodeDraftNotFound:          "草稿不存在",
		CodeDraftAlreadyConfirmed:  "草稿已确认",
		CodeDraftAlreadyRejected:   "草稿已拒绝",
		CodeInvalidDraftStatus:     "草稿状态无效",
		CodeLLMCallFailed:          "LLM 调用失败",
		CodeRAGNoResult:            "RAG 检索无结果",
		CodeConcurrentModification: "并发修改冲突",
		CodeLLMTimeout:             "LLM 调用超时，请重试",
		CodeGenerationQueueFull:    "生成任务队列已满，请稍后重试",

		// System
		CodeInternalError:   "系统内部错误",
		CodeDatabaseError:   "数据库错误",
		CodeValidationError: "参数校验失败",
		CodeUnauthorized:    "未授权",
		CodeRateLimited:     "请求过于频繁",
	}

	if msg, ok := messages[code]; ok {
		return msg
	}
	return "未知错误"
}
