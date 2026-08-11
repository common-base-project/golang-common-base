<script lang="ts">
import { EMitt } from "@/constants/enum";
import emits from "@/utils/emits";
import { defineComponent, reactive } from "vue";
import BaseSidebar from "./base-sidebar.vue";
import Logo from "../header/logo.vue";
import logoUrl from "@/assets/images/logo.png";
import app from "@/constants/app";

/**
 * 移动端侧边栏菜单
 */
export default defineComponent({
  name: "MobileSidebar",
  components: { BaseSidebar, Logo },
  setup() {
    const state = reactive({
      show: true
    });
    emits.on(EMitt.OnMobileOpenSidebar, () => {
      state.show = true;
    });
    const onSelect = () => {
      state.show = false;
    };
    return { state, onSelect, logoUrl, app };
  }
});
</script>

<template>
  <el-drawer v-model="state.show" :append-to-body="false" size="230" :withHeader="false" direction="ltr" class="rr-setting-wrap">
    <div class="rr-header-ctx-logo rr-header-ctx-logo-mobile">
      <logo :logoUrl="logoUrl" :logoName="app.appName"></logo>
    </div>
    <div class="rr-sidebar-mobile-inner" style="overflow: auto; height: calc(100vh - 50px); width: initial !important">
      <base-sidebar :router="true" mode="vertical" :isMobile="true" :onSelect="onSelect"></base-sidebar>
    </div>
  </el-drawer>
</template>
