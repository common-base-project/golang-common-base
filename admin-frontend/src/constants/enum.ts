/**
 * 页面渲染布局
 */
export enum EPageLayoutEnum {
  "page",
  "fullscreen"
}

/**
 * 导航模式
 */
export enum ESidebarLayoutEnum {
  /**
   * 左侧导航
   */
  Left = "left",
  /**
   * 顶部导航
   */
  Top = "top",
  /**
   * 混合导航
   */
  Mix = "mix"
}

/**
 * 主题设置
 */
export enum EThemeSetting {
  /**
   * 侧边栏风格
   */
  Sidebar = "sidebar",
  /**
   * 顶部风格
   */
  TopHeader = "topHeader",
  /**
   * 主题色
   */
  ThemeColor = "themeColor",
  //---
  /**
   * 布局模式
   */
  NavLayout = "navLayout",
  /**
   * 内容是否铺满
   */
  ContentFull = "contentFull",
  //---
  /**
   * logo宽度自动
   */
  LogoAuto = "logoAuto",
  /**
   * 多彩图标
   */
  ColorIcon = "colorIcon",
  /**
   * 侧边栏排他展开
   */
  SidebarUniOpened = "sidebarUniOpened",
  /**
   * 开启tab标签页
   */
  OpenTabsPage = "openTabsPage",
  /**
   * tab标签风格
   */
  TabStyle = "tabStyle",
  //---
  /**
   * 侧边栏展开收起
   */
  SidebarCollapse = "sidebarCollapse"
}

/**
 * 系统框架事件枚举
 */
export enum EMitt {
  /**
   * 全局加载
   */
  OnLoading = "onLoading",
  /**
   * 切换左侧侧边栏
   */
  OnSwitchLeftSidebar = "onSwitchLeftSidebar",
  /**
   * 推送菜单到tab标签页
   */
  OnPushMenuToTabs = "onPushMenuToTabs",
  /**
   * 设置主题
   */
  OnSetTheme = "onSetTheme",
  /**
   * 设置侧边栏排他展开
   */
  OnSetThemeNotUniqueOpened = "onSetTheme_not_uniqueOpened",
  /**
   * 设置开启标签页
   */
  OnSetThemeTabsPage = "onSetTheme_tabsPage",
  /**
   * 设置导航模式
   */
  OnSetNavLayout = "onSetNavLayout",
  /**
   * 刷新tab标签页
   */
  OnReloadTabPage = "onReloadTabPage",

  //
  /**
   * 移动端打开侧边栏
   */
  OnMobileOpenSidebar = "onMobileOpenSidebar",

  //
  /**
   * 混合导航选中顶部主菜单
   */
  OnSelectHeaderNavMenusByMixNav = "onSelectHeaderNavMenusByMixNav",

  /**
   * 关闭当前tab页
   */
  OnCloseCurrTab = "onCloseCurrTab"
}

/**
 * 主题是key
 */
export enum EThemeColor {
  /**
   * 主题色
   */
  ThemeColor = "--color-primary"
}
