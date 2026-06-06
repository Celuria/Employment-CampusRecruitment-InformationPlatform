<script setup lang="ts">
import { useRoute } from 'vue-router'
import PagePlaceholder from '@/components/common/PagePlaceholder.vue'

const route = useRoute()
</script>

<template>
  <PagePlaceholder :title="`双选会详情 #${route.params.id}`">
    待接入 getJobFairDetailApi 与报名外链跳转
  </PagePlaceholder>
</template>
