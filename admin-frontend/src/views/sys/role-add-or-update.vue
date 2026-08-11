<template>
  <el-dialog v-model="visible" :title="!dataForm.id ? '新增' : '修改'" :close-on-click-modal="false" :close-on-press-escape="false">
    <el-form :model="dataForm" :rules="rules" ref="dataFormRef" @keyup.enter="dataFormSubmitHandle()" label-width="120px">
      <el-form-item prop="name" label="名称">
        <el-input v-model="dataForm.name" placeholder="名称"></el-input>
      </el-form-item>
      <el-form-item prop="remark" label="备注">
        <el-input v-model="dataForm.remark" placeholder="备注"></el-input>
      </el-form-item>
      <el-row>
        <el-col :span="12">
          <el-form-item size="small" label="菜单授权">
            <el-tree :data="menuList" :props="{ label: 'name', children: 'children' }" node-key="id" ref="menuListTree" accordion show-checkbox> </el-tree>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item size="small" label="数据授权">
            <el-tree :data="deptList" :props="{ label: 'name', children: 'children' }" node-key="id" ref="deptListTree" accordion show-checkbox> </el-tree>
          </el-form-item>
        </el-col>
      </el-row>
    </el-form>
    <template v-slot:footer>
      <el-button @click="visible = false">取消</el-button>
      <el-button type="primary" @click="dataFormSubmitHandle()">确定</el-button>
    </template>
  </el-dialog>
</template>

<script lang="ts" setup>
import { nextTick, reactive, ref } from "vue";
import baseService from "@/service/baseService";
import { IObject } from "@/types/interface";
import { ElMessage } from "element-plus";
const emit = defineEmits(["refreshDataList"]);

const visible = ref(false);
const menuList = ref([]);
const deptList = ref([]);
const dataFormRef = ref();
const menuListTree = ref();
const deptListTree = ref();

const dataForm = reactive({
  id: "",
  name: "",
  menuIdList: [] as IObject[],
  deptIdList: [],
  remark: ""
});

const rules = ref({
  name: [{ required: true, message: "必填项不能为空", trigger: "blur" }]
});

const init = (id?: number) => {
  visible.value = true;
  dataForm.id = "";

  // 重置表单数据
  if (dataFormRef.value) {
    dataFormRef.value.resetFields();
  }

  nextTick(() => {
    if (menuListTree.value) {
      menuListTree.value.setCheckedKeys([]);
    }

    if (deptListTree.value) {
      deptListTree.value.setCheckedKeys([]);
    }

    Promise.all([getMenuList(), getDeptList()]).then(() => {
      if (id) {
        getInfo(id);
      }
    });
  });
};

// 获取菜单列表
const getMenuList = () => {
  return baseService.get("/sys/menu/select").then((res) => {
    menuList.value = res.data;
  });
};

// 获取部门列表
const getDeptList = () => {
  return baseService.get("/sys/dept/list").then((res) => {
    deptList.value = res.data;
  });
};

// 获取信息
const getInfo = (id: number) => {
  baseService.get(`/sys/role/${id}`).then((res) => {
    Object.assign(dataForm, res.data);

    dataForm.menuIdList.forEach((item: IObject) => menuListTree.value.setChecked(item, true));
    deptListTree.value.setCheckedKeys(dataForm.deptIdList);
  });
};

// 表单提交
const dataFormSubmitHandle = () => {
  dataFormRef.value.validate((valid: boolean) => {
    if (!valid) {
      return false;
    }
    dataForm.menuIdList = [...menuListTree.value.getHalfCheckedKeys(), ...menuListTree.value.getCheckedKeys()];
    dataForm.deptIdList = deptListTree.value.getCheckedKeys();
    (!dataForm.id ? baseService.post : baseService.put)("/sys/role", dataForm).then((res) => {
      ElMessage.success({
        message: "成功",
        duration: 500,
        onClose: () => {
          visible.value = false;
          emit("refreshDataList");
        }
      });
    });
  });
};

defineExpose({
  init
});
</script>
