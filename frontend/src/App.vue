<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import DefaultLayout from '@/layouts/DefaultLayout.vue'
import AuthLayout from '@/layouts/AuthLayout.vue'

const route = useRoute()

const layoutName = computed(() => route.meta.layout ?? 'default')
</script>

<template>
  <DefaultLayout v-if="layoutName === 'default'" />
  <AuthLayout v-else-if="layoutName === 'auth'" />
  <RouterView v-else />
</template>
