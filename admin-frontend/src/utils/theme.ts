import { CacheTheme } from "@/constants/cacheKey";
import { themeSetting } from "@/constants/config";
import { EMitt, EThemeSetting } from "@/constants/enum";
import { IFunction, IObject } from "@/types/interface";
import { getCache, setCache } from "./cache";
import emits from "./emits";
import chalkCss from "element-plus/theme-chalk/index.css?inline";

/**
 * 取主题设置缓存
 * @returns
 */
export const getThemeConfigCache = (): IObject => {
  const cache = getCache(CacheTheme, {}, {});
  return { ...themeSetting, ...cache };
};

/**
 * 取主题设置缓存
 * @param key
 * @param config
 * @returns
 */
export const getThemeConfigCacheByKey = (
  key: EThemeSetting,
  config?: IObject
): string | boolean | number => {
  config = config || getCache(CacheTheme, {}, {});
  return config ? config[key] ?? themeSetting[key] : themeSetting[key];
};

/**
 * 生成主题设置样式名称
 * @param config
 * @returns
 */
export const getThemeConfigToClass = (config: IObject = {}): IObject<string> => {
  const cl: IObject<string> = {};
  Object.keys(config).forEach((x) => {
    cl[x] = `ui-${x}-${config[x]}`;
  });
  return cl;
};

/**
 * 主题设置到缓存
 * @param key
 * @param value
 */
export const setThemeConfigToCache = (
  key: EThemeSetting,
  value: string | boolean | number
): void => {
  const theme = getCache(CacheTheme, {}, {});
  setCache(CacheTheme, { ...theme, [key]: value });
};

/**
 * 设置主题色
 * @param key
 * @param value
 */
export const setThemeColor = (key: string, value: string): void => {
  const elm = window.document.querySelector("body");
  if (elm) {
    elm.style.setProperty(key, value);
    elm.style.setProperty(key + "-light", value + "14");
  }
};

/**
 * 生成主题色
 * @param theme
 * @returns
 */
export const getThemeCluster = (theme: string): string[] => {
  const tintColor = (color: string, tint: number) => {
    let red: any = parseInt(color.slice(0, 2), 16);
    let green: any = parseInt(color.slice(2, 4), 16);
    let blue: any = parseInt(color.slice(4, 6), 16);

    if (tint === 0) {
      // when primary color is in its rgb space
      return [red, green, blue].join(",");
    } else {
      red += Math.round(tint * (255 - red));
      green += Math.round(tint * (255 - green));
      blue += Math.round(tint * (255 - blue));

      red = red.toString(16);
      green = green.toString(16);
      blue = blue.toString(16);

      return `#${red}${green}${blue}`;
    }
  };

  const shadeColor = (color: string, shade: number): string => {
    let red: any = parseInt(color.slice(0, 2), 16);
    let green: any = parseInt(color.slice(2, 4), 16);
    let blue: any = parseInt(color.slice(4, 6), 16);

    red = Math.round((1 - shade) * red);
    green = Math.round((1 - shade) * green);
    blue = Math.round((1 - shade) * blue);

    red = red.toString(16);
    green = green.toString(16);
    blue = blue.toString(16);

    return `#${red}${green}${blue}`;
  };

  const clusters = [theme];
  for (let i = 0; i <= 9; i++) {
    clusters.push(tintColor(theme, Number((i / 10).toFixed(2))));
  }
  clusters.push(shadeColor(theme, 0.1));
  return clusters;
};

/**
 * 获取主题css
 * @param url
 * @param callback
 * @param variable
 */
export const getCSSString = (callback: IFunction, variable?: string): void => {
  if (variable) {
    (window as any)[variable] = chalkCss;
  }
  callback(chalkCss);
};

export const updateStyle = (style: string, oldCluster: string[], newCluster: string[]): string => {
  let newStyle = style;
  oldCluster.forEach((color, index) => {
    newStyle = newStyle.replace(new RegExp(color, "ig"), newCluster[index]);
  });
  return newStyle;
};

/**
 * 更新主题色
 * @param themeColor
 * @param val
 * @returns
 */
export const updateTheme = (val: string): void => {
  emits.emit(EMitt.OnLoading, true);
  const head = document.getElementsByTagName("head")[0];
  const themeCluster = getThemeCluster(val.replace("#", ""));
  const getHandler = (variable: string, id: string) => {
    return () => {
      const originalCluster = getThemeCluster("#409eff".replace("#", ""));
      const newStyle = updateStyle((window as any)[variable], originalCluster, themeCluster);
      let styleTag = document.getElementById(id);
      if (!styleTag) {
        styleTag = document.createElement("style");
        styleTag.setAttribute("id", id);
        head.appendChild(styleTag);
      }
      styleTag.innerText = newStyle;
      emits.emit(EMitt.OnLoading, false);
    };
  };
  const chalkHandler = getHandler("__chalk", "chalk-style");
  if (!(window as any)["__chalk"]) {
    getCSSString(chalkHandler, "__chalk");
  } else {
    chalkHandler();
  }
};
