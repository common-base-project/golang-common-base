<template>
  <el-form :model="dataForm" :rules="rules" ref="dataFormRef" @keyup.enter="dataFormSubmitHandle()" label-width="120px">
    <el-form-item label="账号">
      <span>{{ user.username }}</span>
    </el-form-item>
    <el-form-item prop="password" label="原密码">
      <el-input v-model="dataForm.password" type="password" placeholder="原密码"></el-input>
    </el-form-item>
    <el-form-item prop="newPassword" label="新密码">
      <el-input v-model="dataForm.newPassword" type="password" placeholder="新密码"></el-input>
    </el-form-item>
    <el-form-item prop="confirmPassword" label="确认密码">
      <el-input v-model="dataForm.confirmPassword" type="password" placeholder="确认密码"></el-input>
    </el-form-item>
    <el-form-item>
      <el-button type="primary" @click="dataFormSubmitHandle">确定</el-button>
    </el-form-item>
  </el-form>
</template>

<script lang="ts" setup>
import { computed, reactive, ref } from "vue";
import { IObject } from "@/types/interface";
import baseService from "@/service/baseService";
import { useAppStore } from "@/store";
import { useRouter } from "vue-router";
import { ElMessage } from "element-plus";
const router = useRouter();

const dataFormRef = ref();
const dataForm = reactive({
  password: "",
  newPassword: "",
  confirmPassword: ""
});

const store = useAppStore();
const user = computed(() => store.state.user);

const validateConfirmPassword = (rule: IObject, value: string, callback: (e?: Error) => any) => {
  if (dataForm.newPassword !== value) {
    return callback(new Error("确认密码与新密码输入不一致"));
  }
  callback();
};

const rules = ref({
  password: [{ required: true, message: "必填项不能为空", trigger: "blur" }],
  newPassword: [{ required: true, message: "必填项不能为空", trigger: "blur" }],
  confirmPassword: [
    { required: true, message: "必填项不能为空", trigger: "blur" },
    { validator: validateConfirmPassword, trigger: "blur" }
  ]
});

// 表单提交
const dataFormSubmitHandle = () => {
  dataFormRef.value.validate((valid: boolean) => {
    if (!valid) {
      return false;
    }
    baseService.put("/sys/user/password", dataForm).then((res) => {
      ElMessage.success({
        message: "成功",
        duration: 500,
        onClose: () => {
          router.replace("/login");
        }
      });
    });
  });
};
</script>
