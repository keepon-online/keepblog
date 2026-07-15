export default {
  path: "/system",
  redirect: "/system/log",
  meta: {
    icon: "ant-design:setting-filled",
    title: "系统管理",
    rank: 9
  },
  children: [
    {
      path: "/log",
      name: "登录日志",
      component: () => import("@/views/system/loginlog/index.vue"),
      meta: {
        title: "登录日志",
        icon: "ant-design:eye-outlined"
      }
    },
    {
      path: "/access",
      name: "访问日志",
      component: () => import("@/views/system/accesslog/index.vue"),
      meta: {
        title: "访问日志",
        icon: "ant-design:control-filled"
      }
    },
    {
      path: "/monitor",
      name: "系统监控",
      component: () => import("@/views/system/monitor/index.vue"),
      meta: {
        title: "系统监控",
        icon: "ant-design:monitor-outlined"
      }
    },
    {
      path: "/notice",
      name: "通知管理",
      component: () => import("@/views/system/notice/index.vue"),
      meta: {
        title: "通知管理",
        icon: "ant-design:bell-outlined"
      }
    }
  ]
} as RouteConfigsTable;
