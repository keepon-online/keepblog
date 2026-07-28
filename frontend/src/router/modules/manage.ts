export default {
  path: "/manage",
  redirect: "/manage/post",
  meta: {
    icon: "majesticons:paper-fold-text",
    title: "内容管理",
    // showLink: false,
    rank: 9
  },
  children: [
    {
      path: "/manage/post",
      name: "内容管理",
      component: () => import("@/views/manage/post/index.vue"),
      meta: {
        title: "内容管理",
        icon: "majesticons:paper-fold-text"
      }
    },
    {
      path: "/manage/editor",
      name: "内容新增",
      component: () => import("@/views/manage/post/editor.vue"),
      meta: {
        title: "内容编辑",
        showLink: false
      }
    },
    {
      path: "/manage/preview/:id",
      name: "内容预览",
      component: () => import("@/views/manage/post/preview.vue"),
      meta: {
        title: "内容预览",
        showLink: false
      }
    },
    {
      path: "/manage/editor/:id",
      name: "内容编辑",
      component: () => import("@/views/manage/post/editor.vue"),
      meta: {
        title: "内容编辑",
        showLink: false
      }
    },
    {
      path: "/manage/category",
      name: "分类管理",
      component: () => import("@/views/manage/category/index.vue"),
      meta: {
        title: "分类管理",
        icon: "ant-design:ordered-list-outlined"
      }
    },
    {
      path: "/manage/tag",
      name: "标签管理",
      component: () => import("@/views/manage/tag/index.vue"),
      meta: {
        title: "标签管理",
        icon: "ant-design:tags-outlined"
      }
    },
    {
      path: "/manage/music",
      name: "音乐管理",
      component: () => import("@/views/manage/music/index.vue"),
      meta: {
        title: "音乐管理",
        icon: "mdi:music"
      }
    }
  ]
} as RouteConfigsTable;
