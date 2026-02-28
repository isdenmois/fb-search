<script setup lang="ts">
import { CloseIcon } from './icons'

defineProps<{
  modelValue: string
  disabled?: boolean
}>()

const emit = defineEmits(['update:modelValue'])

const clear = () => {
  emit('update:modelValue', '')
}
</script>

<template>
  <div class="input flex relative items-center gap-2">
    <slot />

    <input class="flex-1 py-2 rounded-md" :disabled="disabled" :value="modelValue"
      @input="$emit('update:modelValue', ($event.target as HTMLInputElement).value)" />

    <div v-if="modelValue && !disabled" data-testid="clear" class="cursor-pointer px-2 z-10" @click="clear">
      <CloseIcon />
    </div>
  </div>
</template>
