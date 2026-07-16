export default {
  path: "/profile",
  redirect: "/profile/center",
  meta: {
    icon: "ant-design:setting-filled",
    title: "个人中心",
    rank: 9,
    showLink: false
  },
  children: [
    {
      path: "/profile/center",
      name: "个人中心",
      component: () => import("@/views/base/profile/index.vue"),
      meta: {
        title: "个人中心",
        icon: "ant-design:eye-outlined"
      }
    }
  ]
} as RouteConfigsTable;
