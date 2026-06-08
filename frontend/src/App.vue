<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import DefaultLayout from '@/layouts/DefaultLayout.vue'
import AuthLayout from '@/layouts/AuthLayout.vue'
import AdminLayout from '@/layouts/AdminLayout.vue'

const route = useRoute()

const layoutName = computed(() => route.meta.layout ?? 'default')
</script>

<template>
  <DefaultLayout v-if="layoutName === 'default'" />
  <AuthLayout v-else-if="layoutName === 'auth'" />
  <AdminLayout v-else-if="layoutName === 'admin'" />
  <RouterView v-else />
</template>
