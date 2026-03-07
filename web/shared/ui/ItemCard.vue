<script setup lang="ts">
import { useSlots } from 'vue'

defineProps<{
  title: string
  suptitle?: string
  subtitle?: string
}>()

const slots = useSlots()
</script>

<template>
  <div class="flex flex-row items-center cursor-pointer">
    <div class="mr-1" v-if="slots.icon">
      <slot name="icon" />
    </div>

    <div class="flex-1">
      <div v-if="suptitle" class="text-xs text-secondary">{{ suptitle }}</div>

      <div class="text-primary">{{ title }}</div>

      <div v-if="subtitle" class="text-xs text-secondary">{{ subtitle }}</div>
    </div>
  </div>
</template>
