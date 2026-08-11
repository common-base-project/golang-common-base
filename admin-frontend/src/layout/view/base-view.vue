<script lang="ts">
import app from "@/constants/app";
import { EMitt, EThemeSetting } from "@/constants/enum";
import emits from "@/utils/emits";
import { getThemeConfigCacheByKey } from "@/utils/theme";
import { defineComponent, reactive, ref } from "vue";
import { useRoute } from "vue-router";
import { useAppStore } from "@/store";
import Tabs from "./tabs.vue";

/**
 * 业务内容视图框架
 */
export default defineComponent({
  name: "View",
  components: { Tabs },
  setup() {
    const store = useAppStore();
    const route = useRoute();
    const state = reactive({
      openTabsPage: getThemeConfigCacheByKey(EThemeSetting.OpenTabsPage)
    });
    const routerKeys = ref({} as any);
    emits.on(EMitt.OnSetThemeTabsPage, (vl) => {
      state.openTabsPage = vl;
    });
    emits.on(EMitt.OnReloadTabPage, () => {
      routerKeys.value[route.fullPath] = new Date().getTime();
    });
    return { state, store, enabledKeepAlive: app.enabledKeepAlive, routerKeys };
  }
});
</script>

<template>
  <tabs v-if="state.openTabsPage" :tabs="store.state.tabs" :activeTabName="store.state.activeTabName"></tabs>
  <div class="rr-view-ctx">
    <el-card shadow="never" class="rr-view-ctx-card">
      <router-view v-slot="{ Component }">
        <keep-alive v-if="enabledKeepAlive">
          <component :is="Component" :key="routerKeys[$route.fullPath] || $route.fullPath" />
        </keep-alive>
        <component :is="Component" v-if="!enabledKeepAlive" />
      </router-view>
    </el-card>
  </div>
</template>
