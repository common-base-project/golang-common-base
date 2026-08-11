<template>
  <svg aria-hidden="true" :class="`iconfont ${className}`" :style="`width:${width};height:${height};color:${color};${style}`">
    <use :xlink:href="symbolId" />
  </svg>
</template>
<script lang="ts">
import { computed, defineComponent } from "vue";

/**
 * 自定义svg图标，可自行将svg图标下载后存放在/src/assets/icons/svg目录下
 * `使用方法：<svg-icon name="earth" color="red"></svg-icon>`
 */
export default defineComponent({
  name: "SvgIcon",
  props: {
    prefix: {
      type: String,
      default: "icon"
    },
    name: {
      type: String,
      required: true
    },
    color: {
      type: String,
      default: ""
    },
    width: String,
    height: String,
    className: { type: String, default: "" },
    style: { type: String, default: "" }
  },
  setup(props) {
    const symbolId = computed(() => `#${props.prefix}-${props.name.replace("icon-", "")}`);
    return { symbolId };
  }
});
</script>
