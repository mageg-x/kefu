import { createRouter, createWebHistory } from "vue-router";
import Login from "@/views/Login.vue";
import Home from "@/views/Home.vue";
import Dashboard from "@/views/Dashboard.vue";
import Visitors from "@/views/Visitors.vue";
import Staff from "@/views/Staff.vue";
import Settings from "@/views/Settings.vue";
import Inbox from "@/views/Inbox.vue";
import Snippets from "@/views/Snippets.vue";
import Profile from "@/views/Profile.vue";
import { useStore } from "@/script/store";

const routes = [
  {
    path: "/",
    redirect: "/home",
  },
  {
    path: "/login",
    name: "Login",
    component: Login,
    meta: {
      requiresAuth: false,
    },
  },
  {
    path: "/home",
    component: Home,
    meta: {
      requiresAuth: true,
    },
    children: [
      {
        path: "dashboard",
        name: "Dashboard",
        component: Dashboard,
      },
      {
        path: "visitors",
        name: "Visitors",
        component: Visitors,
      },
      {
        path: "staff",
        name: "Staff",
        component: Staff,
      },
      {
        path: "settings",
        name: "Settings",
        component: Settings,
      },
      {
        path: "inbox",
        name: "Inbox",
        component: Inbox,
      },
      {
        path: "snippets",
        name: "Snippets",
        component: Snippets,
      },
      {
        path: "profile",
        name: "Profile",
        component: Profile,
      },
    ],
  },
  {
    path: "/:pathMatch(.*)*",
    redirect: "/",
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

const getDefaultPath = (store) => {
  const role = store.user?.role;
  return role === "admin" ? "/home/dashboard" : "/home/inbox";
};

router.beforeEach((to, from, next) => {
  const store = useStore();
  const requiresAuth = to.matched.some((record) => record.meta.requiresAuth !== false);

  if (requiresAuth && !store.token) {
    next("/login");
  } else if (to.path === "/login" && store.token) {
    next(getDefaultPath(store));
  } else if ((to.path === "/home" || to.path === "/home/") && store.token) {
    next(getDefaultPath(store));
  } else {
    next();
  }
});

export default router;
