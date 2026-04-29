# OpenAPI 规范更新补丁

## 需要添加的端点定义

### 1. GET /generation/tasks - 生成任务列表

在 `/generation/tasks` 路径下添加 GET 方法（在 POST 之前）：

```yaml
  /generation/tasks:
    get:
      tags: [Generation]
      summary: 获取生成任务列表
      operationId: listGenerationTasks
      security:
        - bearerAuth: []
      parameters:
        - name: project_id
          in: query
          required: true
          schema:
            type: string
            format: uuid
          description: 项目 ID
        - name: module_id
          in: query
          schema:
            type: string
            format: uuid
          description: 模块 ID（可选）
        - name: status
          in: query
          schema:
            type: string
            enum: [pending, processing, completed, failed]
          description: 任务状态
        - name: offset
          in: query
          schema:
            type: integer
            default: 0
          description: 偏移量
        - name: limit
          in: query
          schema:
            type: integer
            default: 10
          description: 每页数量
      responses:
        '200':
          description: 成功
          content:
            application/json:
              schema:
                type: object
                properties:
                  data:
                    type: array
                    items:
                      $ref: '#/components/schemas/GenerationTask'
                  total:
                    type: integer
                  offset:
                    type: integer
                  limit:
                    type: integer
```

### 2. GET /drafts - 全局草稿列表

在 `/generation/drafts` 路径下添加 GET 方法：

```yaml
  /generation/drafts:
    get:
      tags: [Generation]
      summary: 获取全局草稿列表
      operationId: listAllDrafts
      security:
        - bearerAuth: []
      parameters:
        - name: project_id
          in: query
          schema:
            type: string
            format: uuid
          description: 项目 ID
        - name: task_id
          in: query
          schema:
            type: string
            format: uuid
          description: 任务 ID
        - name: status
          in: query
          schema:
            type: string
            enum: [pending, confirmed, rejected]
          description: 草稿状态
        - name: offset
          in: query
          schema:
            type: integer
            default: 0
          description: 偏移量
        - name: limit
          in: query
          schema:
            type: integer
            default: 10
          description: 每页数量
      responses:
        '200':
          description: 成功
          content:
            application/json:
              schema:
                type: object
                properties:
                  data:
                    type: array
                    items:
                      $ref: '#/components/schemas/CaseDraft'
                  total:
                    type: integer
                  offset:
                    type: integer
                  limit:
                    type: integer
```

### 3. PUT /modules/{id} - 模块编辑

在 `/modules/{id}` 路径下添加 PUT 方法（在 DELETE 之前）：

```yaml
  /modules/{id}:
    put:
      tags: [Modules]
      summary: 更新模块信息
      operationId: updateModule
      security:
        - bearerAuth: []
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
            format: uuid
          description: 模块 ID
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/UpdateModuleRequest'
      responses:
        '200':
          description: 更新成功
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Module'
        '404':
          description: 模块不存在
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'
        '409':
          description: 名称或缩写冲突
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'
```

### 4. PATCH /plans/{id}/status - 计划状态变更

在 `/plans/{id}` 路径下添加 PATCH 方法：

```yaml
  /plans/{id}/status:
    patch:
      tags: [TestPlans]
      summary: 更新测试计划状态
      operationId: updatePlanStatus
      security:
        - bearerAuth: []
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
            format: uuid
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required:
                - status
              properties:
                status:
                  type: string
                  enum: [draft, active, completed, archived]
      responses:
        '200':
          description: 更新成功
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/TestPlan'
        '400':
          description: 状态转换不合法
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'
        '404':
          description: 计划不存在
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'
```

## 需要添加的 Schema

在 `components/schemas` 部分添加（约第 1520 行后）：

```yaml
    UpdateModuleRequest:
      type: object
      properties:
        name:
          type: string
          minLength: 2
          maxLength: 255
        abbreviation:
          type: string
          minLength: 2
          maxLength: 4
          pattern: '^[A-Z]+$'
        description:
          type: string
```

## 实现状态

✅ 端点 1: GET /generation/tasks - 已实现
✅ 端点 2: GET /drafts - 已实现
✅ 端点 3: PUT /modules/{id} - 已实现
✅ 端点 4: PATCH /plans/{id}/status - 已实现

所有端点的代码已完成并通过编译验证。请将上述 YAML 内容添加到 specs/openapi.yaml 文件中的相应位置。
