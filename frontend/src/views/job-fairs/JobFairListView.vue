<script setup lang="ts">
import PageHeader from '@/components/business/PageHeader.vue'
import PagePlaceholder from '@/components/common/PagePlaceholder.vue'
</script>

<template>
  <div>
    <PageHeader
      title="双选会"
      description="大型校园招聘会信息，把握与企业双向选择的机会"
    />
    <PagePlaceholder title="双选会列表">
      待接入 getJobFairListApi、日期区间筛选与列表卡片组件
    </PagePlaceholder>
  </div>
</template>
