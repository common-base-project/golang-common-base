<template>
  <el-dialog v-model="visible" title="云存储配置" :close-on-click-modal="false" :close-on-press-escape="false">
    <el-form :model="dataForm" :rules="rules" ref="dataFormRef" @keyup.enter="dataFormSubmitHandle()" label-width="140px">
      <el-form-item label="类型">
        <el-radio-group v-model="dataForm.type">
          <el-radio :label="1">七牛</el-radio>
          <el-radio :label="2">阿里云</el-radio>
          <el-radio :label="3">腾讯云</el-radio>
        </el-radio-group>
      </el-form-item>
      <template v-if="dataForm.type === 1">
        <el-form-item prop="qiniuDomain" label="域名">
          <el-input v-model="dataForm.qiniuDomain" placeholder="七牛绑定的域名"></el-input>
        </el-form-item>
        <el-form-item prop="qiniuPrefix" label="路径前缀">
          <el-input v-model="dataForm.qiniuPrefix" placeholder="不设置默认为空"></el-input>
        </el-form-item>
        <el-form-item prop="qiniuAccessKey" label="AccessKey">
          <el-input v-model="dataForm.qiniuAccessKey" placeholder="七牛AccessKey"></el-input>
        </el-form-item>
        <el-form-item prop="qiniuSecretKey" label="SecretKey">
          <el-input v-model="dataForm.qiniuSecretKey" placeholder="七牛SecretKey"></el-input>
        </el-form-item>
        <el-form-item prop="qiniuBucketName" label="空间名">
          <el-input v-model="dataForm.qiniuBucketName" placeholder="七牛存储空间名"></el-input>
        </el-form-item>
      </template>
      <template v-else-if="dataForm.type === 2">
        <el-form-item prop="aliyunDomain" label="域名">
          <el-input v-model="dataForm.aliyunDomain" placeholder="阿里云绑定的域名，如：http://cdn.renren.io"></el-input>
        </el-form-item>
        <el-form-item prop="aliyunPrefix" label="路径前缀">
          <el-input v-model="dataForm.aliyunPrefix" placeholder="不设置默认为空"></el-input>
        </el-form-item>
        <el-form-item prop="aliyunEndPoint" label="EndPoint">
          <el-input v-model="dataForm.aliyunEndPoint" placeholder="阿里云EndPoint"></el-input>
        </el-form-item>
        <el-form-item prop="aliyunAccessKeyId" label="AccessKeyId">
          <el-input v-model="dataForm.aliyunAccessKeyId" placeholder="阿里云AccessKeyId"></el-input>
        </el-form-item>
        <el-form-item prop="aliyunAccessKeySecret" label="AccessKeySecret">
          <el-input v-model="dataForm.aliyunAccessKeySecret" placeholder="阿里云AccessKeySecret"></el-input>
        </el-form-item>
        <el-form-item prop="aliyunBucketName" label="BucketName">
          <el-input v-model="dataForm.aliyunBucketName" placeholder="阿里云BucketName"></el-input>
        </el-form-item>
      </template>
      <template v-else-if="dataForm.type === 3">
        <el-form-item prop="qcloudDomain" label="域名">
          <el-input v-model="dataForm.qcloudDomain" placeholder="腾讯云绑定的域名"></el-input>
        </el-form-item>
        <el-form-item prop="qcloudPrefix" label="路径前缀">
          <el-input v-model="dataForm.qcloudPrefix" placeholder="不设置默认为空"></el-input>
        </el-form-item>
        <el-form-item prop="qcloudAppId" label="AppId">
          <el-input v-model="dataForm.qcloudAppId" placeholder="腾讯云AppId"></el-input>
        </el-form-item>
        <el-form-item prop="qcloudSecretId" label="SecretId">
          <el-input v-model="dataForm.qcloudSecretId" placeholder="腾讯云SecretId"></el-input>
        </el-form-item>
        <el-form-item prop="qcloudSecretKey" label="SecretKey">
          <el-input v-model="dataForm.qcloudSecretKey" placeholder="腾讯云SecretKey"></el-input>
        </el-form-item>
        <el-form-item prop="qcloudBucketName" label="BucketName">
          <el-input v-model="dataForm.qcloudBucketName" placeholder="腾讯云BucketName"></el-input>
        </el-form-item>
        <el-form-item prop="qcloudRegion" label="所属地区">
          <el-select v-model="dataForm.qcloudRegion" clearable placeholder="请选择" class="w-percent-100">
            <el-option value="ap-beijing-1" label="北京一区（华北）"></el-option>
            <el-option value="ap-beijing" label="北京"></el-option>
            <el-option value="ap-shanghai" label="上海（华东）"></el-option>
            <el-option value="ap-guangzhou" label="广州（华南）"></el-option>
            <el-option value="ap-chengdu" label="成都（西南）"></el-option>
            <el-option value="ap-chongqing" label="重庆"></el-option>
            <el-option value="ap-singapore" label="新加坡"></el-option>
            <el-option value="ap-hongkong" label="香港"></el-option>
            <el-option value="na-toronto" label="多伦多"></el-option>
            <el-option value="eu-frankfurt" label="法兰克福"></el-option>
          </el-select>
        </el-form-item>
      </template>
    </el-form>
    <template v-slot:footer>
      <el-button @click="visible = false">取消</el-button>
      <el-button type="primary" @click="dataFormSubmitHandle()">确定</el-button>
    </template>
  </el-dialog>
</template>

<script lang="ts" setup>
import { reactive, ref } from "vue";
import baseService from "@/service/baseService";
import { ElMessage } from "element-plus";
const emit = defineEmits(["refreshDataList"]);

const visible = ref(false);
const dataFormRef = ref();

const dataForm = reactive({
  type: 0,
  qiniuDomain: "",
  qiniuPrefix: "",
  qiniuAccessKey: "",
  qiniuSecretKey: "",
  qiniuBucketName: "",
  aliyunDomain: "",
  aliyunPrefix: "",
  aliyunEndPoint: "",
  aliyunAccessKeyId: "",
  aliyunAccessKeySecret: "",
  aliyunBucketName: "",
  qcloudDomain: "",
  qcloudPrefix: "",
  qcloudAppId: 0,
  qcloudSecretId: "",
  qcloudSecretKey: "",
  qcloudBucketName: "",
  qcloudRegion: ""
});

const rules = ref({
  qiniuDomain: [{ required: true, message: "必填项不能为空", trigger: "blur" }],
  qiniuAccessKey: [{ required: true, message: "必填项不能为空", trigger: "blur" }],
  qiniuSecretKey: [{ required: true, message: "必填项不能为空", trigger: "blur" }],
  qiniuBucketName: [{ required: true, message: "必填项不能为空", trigger: "blur" }],
  aliyunDomain: [{ required: true, message: "必填项不能为空", trigger: "blur" }],
  aliyunEndPoint: [{ required: true, message: "必填项不能为空", trigger: "blur" }],
  aliyunAccessKeyId: [{ required: true, message: "必填项不能为空", trigger: "blur" }],
  aliyunAccessKeySecret: [{ required: true, message: "必填项不能为空", trigger: "blur" }],
  aliyunBucketName: [{ required: true, message: "必填项不能为空", trigger: "blur" }],
  qcloudDomain: [{ required: true, message: "必填项不能为空", trigger: "blur" }],
  qcloudAppId: [{ required: true, message: "必填项不能为空", trigger: "blur" }],
  qcloudSecretId: [{ required: true, message: "必填项不能为空", trigger: "blur" }],
  qcloudSecretKey: [{ required: true, message: "必填项不能为空", trigger: "blur" }],
  qcloudBucketName: [{ required: true, message: "必填项不能为空", trigger: "blur" }],
  qcloudRegion: [{ required: true, message: "必填项不能为空", trigger: "blur" }]
});

const init = () => {
  visible.value = true;

  // 重置表单数据
  if (dataFormRef.value) {
    dataFormRef.value.resetFields();
  }

  getInfo();
};

// 获取信息
const getInfo = () => {
  baseService.get("/sys/oss/info").then((res) => {
    Object.assign(dataForm, res.data);
  });
};

// 表单提交
const dataFormSubmitHandle = () => {
  dataFormRef.value.validate((valid: boolean) => {
    if (!valid) {
      return false;
    }
    baseService.post("/sys/oss", dataForm).then((res) => {
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
