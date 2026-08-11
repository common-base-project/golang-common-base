<template>
  <div class="ren-region">
    <el-input v-model="showName" :placeholder="placeholder" @click="treeDialog">
      <template v-slot:append>
        <el-button icon="search" @click="treeDialog"></el-button>
      </template>
    </el-input>
    <el-dialog v-model="visibleTree" width="360px" :modal="false" :title="placeholder" :close-on-click-modal="false" :close-on-press-escape="false">
      <el-form size="small" :inline="true">
        <el-form-item label="关键字">
          <el-input v-model="filterText" :style="{ width: '150px' }"></el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="default">查询</el-button>
        </el-form-item>
      </el-form>
      <el-tree class="filter-tree" :data="dataList" :default-expanded-keys="expandedKeys" :props="{ label: 'name', children: 'children' }" :expand-on-click-node="false" :filter-node-method="filterNode" :highlight-current="true" node-key="id" ref="treeRef"> </el-tree>
      <template v-slot:footer>
        <el-button type="default" @click="cancelHandle()">取消</el-button>
        <el-button type="info" @click="clearHandle()">清除</el-button>
        <el-button type="primary" @click="commitHandle()">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script lang="ts" setup>
import { nextTick, ref, watch } from "vue";
import { treeDataTranslate } from "@/utils/utils";
import { IObject } from "@/types/interface";
import baseService from "@/service/baseService";
import { ElMessage } from "element-plus";

const filterText = ref("");
const visibleTree = ref(false);
const dataList = ref<any[]>([]);
const showName = ref("");
const expandedKeys = ref<any[]>([]);
const treeRef = ref();

const props = defineProps({
  modelValue: [Number, String],
  parentName: String,
  placeholder: String
});

watch(
  () => filterText.value,
  (val) => {
    treeRef.value.filter(val);
  }
);

watch(
  () => props.parentName,
  (val) => {
    showName.value = val as string;
  }
);

const treeDialog = () => {
  expandedKeys.value = [];
  if (treeRef.value) {
    treeRef.value.setCurrentKey(null);
  }
  visibleTree.value = true;
  getDataList(props.modelValue);
};

const filterNode = (value: string, data: IObject) => {
  if (!value) return true;
  return data.name.indexOf(value) !== -1;
};

const getDataList = (id: any) => {
  return baseService.get("/sys/region/tree").then((res) => {
    dataList.value = treeDataTranslate(res.data);
    nextTick(() => {
      treeRef.value.setCurrentKey(id);
      expandedKeys.value = [id];
    });
  });
};

const cancelHandle = () => {
  visibleTree.value = false;
  dataList.value = [];
  filterText.value = "";
};

const emit = defineEmits(["update:modelValue", "update:parentName"]);

const clearHandle = () => {
  emit("update:modelValue", "0");
  emit("update:parentName", "");
  showName.value = "";
  visibleTree.value = false;
  dataList.value = [];
  filterText.value = "";
};

const commitHandle = () => {
  const node = treeRef.value.getCurrentNode();
  if (!node) {
    ElMessage.error("请选择");
    return;
  }
  emit("update:modelValue", node.id);
  emit("update:parentName", node.name);

  showName.value = node.name;
  visibleTree.value = false;
  dataList.value = [];
  filterText.value = "";
};
</script>

<style lang="less" scoped>
.ren-region {
  .filter-tree {
    max-height: 230px;
    overflow: auto;
  }
  .el-dialog__body {
    padding: 0 0 0 20px;
  }
  .el-dialog__footer {
    padding: 10px 20px 8px 20px;
  }
}
</style>
