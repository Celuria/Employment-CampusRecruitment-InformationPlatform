<script setup lang="ts">
import { useRoute } from 'vue-router'
import PagePlaceholder from '@/components/common/PagePlaceholder.vue'

const route = useRoute()
</script>

<template>
  <PagePlaceholder :title="`宣讲会详情 #${route.params.id}`">
    详情页框架已就绪，待接入 getCareerTalkDetailApi 与「加入日历」操作
  </PagePlaceholder>
</template>
