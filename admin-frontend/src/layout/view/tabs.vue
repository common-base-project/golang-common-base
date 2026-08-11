<script lang="ts">
import SvgIcon from "@/components/base/svg-icon";
import { EMitt } from "@/constants/enum";
import { IObject } from "@/types/interface";
import emits from "@/utils/emits";
import { arrayToObject } from "@/utils/utils";
import { ElMessage } from "element-plus";
import { findIndex } from "lodash";
import { defineComponent, reactive, watch } from "vue";
import { RouteLocationMatched, useRouter } from "vue-router";
import { useAppStore } from "@/store";

/**
 * tab标签页
 */
export default defineComponent({
  name: "Tabs",
  components: { SvgIcon },
  props: {
    tabs: Array,
    activeTabName: String
  },
  setup(props) {
    const ops = [
      { label: "关闭当前标签页", value: 5, icon: "close" },
      { label: "关闭其他标签页", value: 1, icon: "close" },
      { label: "关闭全部标签页", value: 4, icon: "circle-close" }
    ];
    const router = useRouter();
    const store = useAppStore();
    const firstRoute = (router.options.routes[0] || {}) as RouteLocationMatched;
    const home: RouteLocationMatched = firstRoute.children && firstRoute.children.length > 0 ? (firstRoute.children[0] as RouteLocationMatched) : firstRoute;
    const defaultTab = { label: "", value: home.path };
    const state = reactive({
      activeTabName: props.activeTabName || defaultTab.value,
      tabs: (props.tabs && props.tabs.length ? props.tabs : [defaultTab]) as IObject[]
    });
    watch(
      () => state.tabs,
      (res) => {
        store.updateState({ tabs: res });
      },
      { deep: true }
    );
    emits.on(EMitt.OnPushMenuToTabs, (route) => {
      const path: string = route.value;
      if (path.includes("/error")) {
        return;
      }
      const tabKeys: IObject<number> = arrayToObject(state.tabs, "value", () => 1);
      if (!tabKeys[path]) {
        state.tabs.push(route);
      }
      if (state.activeTabName !== path) {
        state.activeTabName = path;
      }
    });
    emits.on(EMitt.OnCloseCurrTab, () => {
      onClose(5);
    });
    const onTabClick = (tab: any) => {
      tab.props.name && router.push(tab.props.name);
    };
    const onTabRemove = (targetName: string) => {
      const index = findIndex(state.tabs, (x) => x.value === targetName);
      if (state.tabs.length > 1) {
        updateClosedTabs([...store.state.closedTabs, targetName], false);
        if (state.activeTabName === targetName) {
          const toIndex = index === 0 ? index + 1 : index - 1;
          state.activeTabName = state.tabs[toIndex].value;
          router.push(state.activeTabName);
        }
        state.tabs.splice(index, 1);
      } else {
        ElMessage({ type: "error", message: "只剩下一个标签页，不支持关闭", offset: 0 });
      }
    };
    const updateClosedTabs = (closedTabs: any[], isTransform = true) => {
      if (isTransform) {
        closedTabs = closedTabs.map((x) => x.value);
      }
      store.updateState({ closedTabs });
    };
    const onClose = (value: number) => {
      let index = null;
      const rawTabs = state.tabs;
      switch (value) {
        case 1:
          //其他
          state.tabs = state.tabs.filter((x) => [home.path, state.activeTabName].includes(x.value));
          updateClosedTabs(rawTabs.filter((x) => ![home.path, state.activeTabName].includes(x.value)));
          break;
        case 2:
          //右侧
          index = findIndex(state.tabs, (x) => x.value === state.activeTabName);
          state.tabs.splice(index + 1, state.tabs.length - (index + 1));
          updateClosedTabs(rawTabs.slice(index + 1));
          break;
        case 3:
          //左侧
          index = findIndex(state.tabs, (x) => x.value === state.activeTabName);
          state.tabs.splice(1, index - 1);
          updateClosedTabs(rawTabs.slice(1, index - 1));
          break;
        case 4:
          //全部
          state.tabs = [defaultTab];
          state.activeTabName = defaultTab.value;
          updateClosedTabs(rawTabs);
          router.push(state.activeTabName);
          break;
        case 5:
          //当前
          if (state.activeTabName !== defaultTab.value) {
            updateClosedTabs([...store.state.closedTabs, state.activeTabName], false);
            index = findIndex(state.tabs, (x) => x.value === state.activeTabName);
            state.tabs.splice(index, 1);
            state.activeTabName = state.tabs[state.tabs.length - 1].value;
            router.push(state.activeTabName);
          }
          break;
        default:
          break;
      }
    };
    return { state, onTabClick, onTabRemove, home, onClose, ops };
  }
});
</script>
<template>
  <div class="rr-view-tab-wrap">
    <el-tabs class="rr-view-tab" v-model="state.activeTabName" @tab-click="onTabClick" @tab-remove="onTabRemove">
      <el-tab-pane :name="home.path" :closable="false">
        <template #label>
          <!-- 文字主页和图标主页tab -->
          <!-- {{ t("ui.router.pageHome") }} -->
          <svg-icon name="home"></svg-icon>
        </template>
      </el-tab-pane>
      <el-tab-pane v-for="x in state.tabs.slice(1)" :key="x.value" :label="x.label" :name="x.value" :closable="true"></el-tab-pane>
    </el-tabs>
    <el-dropdown trigger="click" placement="bottom-end" class="rr-view-tab-ops" @command="onClose">
      <template #dropdown>
        <el-dropdown-menu>
          <el-dropdown-item v-for="x in ops" :key="x.value" :icon="x.icon" :command="x.value">
            {{ x.label }}
          </el-dropdown-item>
        </el-dropdown-menu>
      </template>
      <span class="el-dropdown-link">
        <el-icon class="el-icon--right"><arrow-down /></el-icon>
      </span>
    </el-dropdown>
  </div>
</template>
