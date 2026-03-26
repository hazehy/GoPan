import { createRouter, createWebHistory } from 'vue-router';
import { useAuthStore } from '@/stores/auth';

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      redirect: '/disk',
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/LoginView.vue'),
    },
    {
      path: '/register',
      name: 'register',
      component: () => import('@/views/RegisterView.vue'),
    },
    {
      path: '/disk',
      name: 'disk',
      component: () => import('@/views/DiskView.vue'),
      meta: {
        requiresAuth: true,
      },
    },
    {
      path: '/admin',
      name: 'admin',
      component: () => import('@/views/AdminView.vue'),
      meta: {
        requiresAuth: true,
        requiresAdmin: true,
      },
    },
    {
      path: '/share/:identity',
      name: 'share',
      component: () => import('@/views/ShareView.vue'),
    },
  ],
});

router.beforeEach((to) => {
  const authStore = useAuthStore();

  if (to.meta.requiresAuth && !authStore.isLoggedIn) {
    return {
      path: '/login',
      query: {
        redirect: to.fullPath,
      },
    };
  }

  if (to.meta.requiresAdmin && !authStore.isAdmin) {
    return '/disk';
  }

  if (!to.meta.requiresAdmin && to.path === '/disk' && authStore.isAdmin) {
    return '/admin';
  }

  if ((to.path === '/login' || to.path === '/register') && authStore.isLoggedIn) {
    return authStore.isAdmin ? '/admin' : '/disk';
  }

  return true;
});

export default router;
