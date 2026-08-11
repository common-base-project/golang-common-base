<script lang="ts">
import { defineComponent, PropType } from "vue";
import classNames from "classnames";
import SvgIcon from "@/components/base/svg-icon";
import { RouteRecordRaw } from "vue-router";

export default defineComponent({
  name: "SidebarMenusItems",
  components: { SvgIcon },
  props: {
    menus: Array as PropType<RouteRecordRaw[]>,
    hiddenIndex: Number,
    className: String
  },
  setup(props) {
    const getStyle = (index: number): string => {
      const styles: Array<any> = [];
      const isHidden = props.hiddenIndex ? props.hiddenIndex > -1 && index > props.hiddenIndex : false;
      styles.push("display:" + (isHidden ? "none" : "block"));
      return styles.join(";");
    };
    return { props, classNames, getStyle };
  }
});
</script>
<template>
  <template v-for="(x, index) in props.menus || []" :key="x.path">
    <el-sub-menu v-if="x.children && x.children.length > 0" :index="x.path" :popper-class="props.className" :class="classNames({ isMore: x.meta?.isMore })" :style="getStyle(index)">
      <template #title>
        <el-icon v-if="x.meta?.icon !== false">
          <svg-icon :name="`${x.meta?.icon || 'icon-file-fill'}`"></svg-icon>
        </el-icon>
        <span>
          <a>{{ x.meta?.title }}</a>
        </span>
      </template>
      <sidebar-menus-items :menus="x.children"></sidebar-menus-items>
    </el-sub-menu>
    <el-menu-item v-else :index="x.meta?.isNewPage ? x.path : x.path" :class="classNames({ isLink: !!x.meta?.isNewPage, isMore: x.meta?.isMore })" :style="getStyle(index)">
      <template #title>
        <a v-if="x.meta?.isNewPage" :href="`${x.meta.url}`" target="_blank" rel="opener">
          {{ x.meta.title }}
        </a>
        <a v-else>{{ x.meta?.title }}</a>
      </template>
      <el-icon v-if="x.meta?.icon !== false">
        <svg-icon :name="`${x.meta?.icon || 'icon-file-fill'}`"></svg-icon>
      </el-icon>
    </el-menu-item>
  </template>
</template>
