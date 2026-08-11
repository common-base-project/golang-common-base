<template>
  <div class="mod-sys__log-login">
    <el-form :inline="true" :model="state.dataForm" @keyup.enter="state.getDataList()">
      <el-form-item>
        <el-input v-model="state.dataForm.creatorName" placeholder="用户名" clearable></el-input>
      </el-form-item>
      <el-form-item>
        <el-select v-model="state.dataForm.status" placeholder="状态" clearable>
          <el-option label="失败" :value="0"></el-option>
          <el-option label="成功" :value="1"></el-option>
          <el-option label="账号已锁定" :value="2"></el-option>
        </el-select>
      </el-form-item>
      <el-form-item>
        <el-button @click="state.getDataList()">查询</el-button>
      </el-form-item>
      <el-form-item>
        <el-button type="info" @click="state.exportHandle()">导出</el-button>
      </el-form-item>
    </el-form>
    <el-table v-loading="state.dataListLoading" :data="state.dataList" border @sort-change="state.dataListSortChangeHandle" style="width: 100%">
      <el-table-column prop="creatorName" label="用户名" header-align="center" align="center"></el-table-column>
      <el-table-column prop="operation" label="操作类型" header-align="center" align="center">
        <template v-slot="scope">
          {{ scope.row.operation === 0 ? "登录" : "退出" }}
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" sortable="custom" header-align="center" align="center">
        <template v-slot="scope">
          <el-tag v-if="scope.row.status === 0" size="small" type="danger">失败</el-tag>
          <el-tag v-else-if="scope.row.status === 1" size="small" type="success">成功</el-tag>
          <el-tag v-else size="small" type="warning">账号已锁定</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="ip" label="操作IP" header-align="center" align="center"></el-table-column>
      <el-table-column prop="userAgent" label="用户代理" header-align="center" align="center" width="150" :show-overflow-tooltip="true"></el-table-column>
      <el-table-column prop="createDate" label="创建时间" sortable="custom" header-align="center" align="center" width="180"></el-table-column>
    </el-table>
    <el-pagination :current-page="state.page" :page-sizes="[10, 20, 50, 100]" :page-size="state.limit" :total="state.total" layout="total, sizes, prev, pager, next, jumper" @size-change="state.pageSizeChangeHandle" @current-change="state.pageCurrentChangeHandle"> </el-pagination>
  </div>
</template>

<script lang="ts" setup>
import useView from "@/hooks/useView";
import { reactive, toRefs } from "vue";

const view = reactive({
  getDataListURL: "/sys/log/login/page",
  getDataListIsPage: true,
  exportURL: "/sys/log/login/export",
  dataForm: {
    creatorName: "",
    status: ""
  }
});

const state = reactive({ ...useView(view), ...toRefs(view) });
</script>
