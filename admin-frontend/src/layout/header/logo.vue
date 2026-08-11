<script lang="ts">
import app from "@/constants/app";
import { defineComponent } from "vue";
interface ILogo {
  logoUrl?: string;
  logoName?: string;
}

/**
 * 顶部logo
 */
export default defineComponent({
  name: "Logo",
  props: {
    logoUrl: String,
    logoName: {
      type: String,
      default: "logo"
    }
  },
  setup(props: ILogo) {
    return { props, app };
  }
});
</script>

<template>
  <span :class="`rr-header-ctx-logo-img-wrap ${'enabled-logo-' + app.enabledLogo}`">
    <!-- 支持显示图片logo或者产品名称缩写，二选一模式，通过注释开启功能，app.enabledLogo控制正常模式下图片logo是否显示，如果有图片logo，收起状态会强制显示图片logo -->
    <!-- <img :src="props.logoUrl" class="rr-header-ctx-logo-img" :alt="props.logoName" /> -->
    <span>人人</span>
    <span class="rr-header-ctx-logo-line"></span>
  </span>
  <span class="rr-header-ctx-logo-text">{{ props.logoName }}</span>
</template>
