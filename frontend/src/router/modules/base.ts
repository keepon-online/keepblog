export default {
  path: "/base",
  redirect: "/base/link",
  meta: {
    icon: "ant-design:aim-outlined",
    title: "基础信息",
    rank: 9
  },
  children: [
    {
      path: "/base/site",
      name: "站点设置",
      component: () => import("@/views/base/setting/index.vue"),
      meta: {
        title: "站点设置",
        icon: "ant-design:ie-circle-filled"
      }
    },
    {
      path: "/base/link",
      name: "友情链接",
      component: () => import("@/views/base/link/index.vue"),
      meta: {
        title: "友情链接",
        icon: "ant-design:link-outlined"
      }
    },
    {
      path: "/base/about",
      name: "关于",
      component: () => import("@/views/base/about/index.vue"),
      meta: {
        title: "关于",
        icon: "ant-design:send-outlined"
      }
    }
  ]
} as RouteConfigsTable;
