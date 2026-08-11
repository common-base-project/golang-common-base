# Renren-Security 对标清单（Backend Admin 模板化）

> 目标：将当前 `backend-admin` 作为后续二次开发模板，逐项覆盖 `renren-security` 的后台管理能力。

参考项目：
- https://gitee.com/renrenio/renren-security

## 1. 当前项目验证结果

### 1.1 后端验证
- 编译验证：`go test ./...` 通过。
- 运行验证：服务初始化（MySQL、Redis、MinIO、Casbin）已通过。
- 启动失败原因：`9080` 端口被其他进程占用（本机 `java` 进程），并非后端初始化错误。

### 1.2 前端验证
- 依赖安装与构建：`npm install && npm run build` 通过。
- 现状：前端当前为基础管理模板（登录、仪表盘、Demo 页面、权限示例），非完整业务后台。

## 2. Renren-Security 后台功能基线（按控制器）

`renren-admin` 主要控制器：
- `LoginController`：登录/认证
- `SysUserController`：用户管理
- `SysRoleController`：角色管理
- `SysMenuController`：菜单管理
- `SysDeptController`：部门管理
- `SysDictTypeController`：字典类型
- `SysDictDataController`：字典数据
- `SysParamsController`：系统参数
- `ScheduleJobController`：定时任务
- `ScheduleJobLogController`：任务日志
- `SysLogOperationController`：操作日志
- `SysLogLoginController`：登录日志
- `SysLogErrorController`：异常日志
- `SysOssController`：对象存储管理
- `IndexController`：系统首页/聚合信息

此外还有：
- `renren-generator`：代码生成器
- 数据权限体系（部门/本人/子部门）
- 按钮级权限与动态菜单

## 3. 当前 backend-admin 覆盖度

已实现：
- 用户管理（基础 CRUD）
- 邮件模板/邮件发送相关模块
- Casdoor OIDC（身份认证接入）
- Casbin（授权框架接入）
- MySQL + SQLite、Redis、MinIO 连接能力

未覆盖（相对 renren-security）：
- 角色管理
- 菜单管理（动态路由源）
- 部门管理
- 字典管理（type/data）
- 系统参数管理
- 定时任务与任务日志
- 操作/登录/异常日志体系
- OSS 文件管理后台（列表、策略、上传审计）
- 数据权限规则体系（本部门/子部门/本人）
- 代码生成器

## 4. 建议模板化落地顺序（强制分阶段）

### Phase A：权限与组织架构核心（必须先做）
1. 组织模型：用户、角色、菜单、部门（4 大表 + 关联表）
2. Casbin 策略模型升级：支持 RBAC + 菜单/按钮资源
3. 后端 API：
   - `/sys/user`
   - `/sys/role`
   - `/sys/menu`
   - `/sys/dept`
4. 前端：动态菜单与按钮级权限（从后端菜单树驱动）

验收标准：
- 不同角色看到不同菜单与按钮。
- 无权限请求被统一拒绝（403）。

### Phase B：系统治理能力
1. 字典类型/字典数据
2. 系统参数
3. 操作日志、登录日志、异常日志

验收标准：
- 关键配置可后台维护。
- 审计日志可检索、分页、按时间过滤。

### Phase C：运维能力
1. 定时任务（新增/暂停/恢复/触发）
2. 任务执行日志
3. OSS 管理页（MinIO 对接 + 元数据管理）

验收标准：
- 任务调度全流程可视化。
- 文件上传下载链路可追踪。

### Phase D：效率能力
1. 代码生成器（可先做简化版）
2. 统一 CRUD 脚手架（Go 模板）

验收标准：
- 输入表结构后，自动生成 model/handler/router 前后端基础代码。

## 5. 推荐数据库模型最小集合（模板必备）

- `sys_user`
- `sys_role`
- `sys_menu`
- `sys_dept`
- `sys_user_role`
- `sys_role_menu`
- `sys_dict_type`
- `sys_dict_data`
- `sys_params`
- `sys_log_operation`
- `sys_log_login`
- `sys_log_error`
- `sys_schedule_job`
- `sys_schedule_job_log`
- `sys_oss`

## 6. 本项目作为“母版模板”的硬规范

- 模块命名统一用 `sys/*` 与 `biz/*` 分层。
- 路由遵循：`/api/v1/sys/...`（平台能力）与 `/api/v1/biz/...`（业务能力）。
- 所有管理模块必须带：分页、审计字段、软删除、权限点。
- 所有新增接口必须补齐：
  - OpenAPI 注释
  - Casbin 资源点
  - 前端权限按钮点
  - 最少一个集成验证用例

## 7. 立即下一步（执行级）

1. 先落地 Phase A 的数据库表与迁移。
2. 新建 `sys` 模块目录并接入 4 个核心控制器：user/role/menu/dept。
3. 前端把 demo 权限页替换为真实菜单权限页。
4. 完成首轮后再推进 Phase B/C。

---

如果按此清单执行，本仓库可以稳定演进为“可复用后台母版”，后续业务项目只需在 `biz` 层做增量开发。
