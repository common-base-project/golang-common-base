<template>
  <div class="mod-sys__log-error">
    <el-form :inline="true" :model="state.dataForm" @keyup.enter="state.getDataList()">
      <el-form-item>
        <el-button type="info" @click="state.exportHandle()">导出</el-button>
      </el-form-item>
    </el-form>
    <el-table v-loading="state.dataListLoading" :data="state.dataList" border @sort-change="state.dataListSortChangeHandle" style="width: 100%">
      <el-table-column prop="requestUri" label="请求URI" header-align="center" align="center" :show-overflow-tooltip="true"></el-table-column>
      <el-table-column prop="requestMethod" label="请求方式" header-align="center" align="center" :show-overflow-tooltip="true"></el-table-column>
      <el-table-column prop="requestParams" label="请求参数" header-align="center" align="center" width="150" :show-overflow-tooltip="true"></el-table-column>
      <el-table-column prop="ip" label="操作IP" header-align="center" align="center"></el-table-column>
      <el-table-column prop="userAgent" label="用户代理" header-align="center" align="center" width="150" :show-overflow-tooltip="true"></el-table-column>
      <el-table-column prop="createDate" label="创建时间" sortable="custom" header-align="center" align="center" width="180"></el-table-column>
      <el-table-column label="操作" fixed="right" header-align="center" align="center" width="150">
        <template v-slot="scope">
          <el-button type="primary" link @click="infoHandle(scope.row.errorInfo)">异常信息</el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination :current-page="state.page" :page-sizes="[10, 20, 50, 100]" :page-size="state.limit" :total="state.total" layout="total, sizes, prev, pager, next, jumper" @size-change="state.pageSizeChangeHandle" @current-change="state.pageCurrentChangeHandle"> </el-pagination>

    <el-dialog v-model="infoDialogVisible" title="异常信息" center>
    <span>{{ infoMessage }}</span>
      <template #footer>
        <div class="dialog-footer">
          <el-button type="primary" @click="infoDialogVisible = false">确定</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script lang="ts" setup>
import useView from "@/hooks/useView";
import { reactive, toRefs, ref } from "vue";

const infoDialogVisible = ref(false)
const infoMessage = ref('')
const view = reactive({
  getDataListURL: "/sys/log/error/page",
  getDataListIsPage: true,
  exportURL: "/sys/log/error/export"
});

const state = reactive({ ...useView(view), ...toRefs(view) });

const infoHandle = (info: string) => {
  infoMessage.value = info
  infoDialogVisible.value = true
};
</script>

<style lang="less" scoped>
.mod-sys__log-error {
  &-view-info {
    width: 80%;
  }
}
</style>
