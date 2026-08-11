<script lang="ts">
import { IObject } from "@/types/interface";
import { getValueByKeys } from "@/utils/utils";
import { defineComponent, ref, watch } from "vue";
import { RouteLocationMatched, useRouter } from "vue-router";

/**
 * 顶部面包屑
 */
export default defineComponent({
  name: "Breadcrumb",
  setup() {
    const router = useRouter();
    const breadcrumbs = ref<IObject[]>([]);
    const { currentRoute } = router;
    const firstRoute = (router.options.routes[0] || {}) as RouteLocationMatched;
    const home: RouteLocationMatched = firstRoute.children && firstRoute.children.length > 0 ? (firstRoute.children[0] as RouteLocationMatched) : firstRoute;
    watch(
      () => currentRoute.value,
      () => {
        breadcrumbs.value = currentRoute.value.path !== home.path ? getValueByKeys(currentRoute.value, "meta.matched", []) : [];
      }
    );

    return { breadcrumbs, currentRoute, home };
  }
});
</script>
<template>
  <el-breadcrumb separator="/" style="padding-top: 4px">
    <el-breadcrumb-item :to="{ path: home.path }"> 工作台 </el-breadcrumb-item>
    <el-breadcrumb-item v-for="x in breadcrumbs" :key="x.path">{{ currentRoute.query._mt || x.title || "" }} </el-breadcrumb-item>
  </el-breadcrumb>
</template>
