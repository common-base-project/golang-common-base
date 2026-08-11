<template>
  <el-card shadow="never" class="aui-card--fill">
    <div class="mod-home">
      <h3>项目介绍</h3>
      <ul>
    <li>当前项目为 Go 后端 + Vue3 管理前端的一体化权限管理平台。</li>
    <li>后端基于 Gin、Gorm、Casbin，支持用户、角色、菜单、部门、日志、任务调度与对象存储等管理能力。</li>
    <li>已兼容 renren 风格前端接口，支持分页、导出、上传等常见管理后台操作。</li>
    <li>支持 MySQL 与 Redis，支持 MinIO 对象存储，便于本地开发与生产部署。</li>
    <li>支持基于 Docker Compose 的快速部署，并提供 Casdoor/Caddy 相关部署方案。</li>
    <li>默认演示账号：admin / admin。</li>
      </ul>
      <h3>获取帮助</h3>
      <ul>
    <li>后端说明请参考：backend-admin/README.md</li>
    <li>部署说明请参考：deploy/README.md</li>
    <li>如需 Casdoor 生产部署，请参考：deploy/casdoor/README.md</li>
      </ul>
    </div>
  </el-card>
</template>

<style scoped>
.mod-home {
  line-height: 1.5;
}
</style>
<script setup lang="ts"></script>
