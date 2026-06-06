<script setup lang="ts">
import PageHeader from '@/components/business/PageHeader.vue'
import PagePlaceholder from '@/components/common/PagePlaceholder.vue'
</script>

<template>
  <div>
    <PageHeader
      title="个性化推荐"
      description="根据您的学院、专业与偏好，为您匹配最合适的校招活动"
    />
    <PagePlaceholder title="推荐列表">
      待接入 getRecommendationsApi、匹配原因标签、「加入日历」「不感兴趣」交互
    </PagePlaceholder>
  </div>
</template>
