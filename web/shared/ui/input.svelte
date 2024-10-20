<script lang="ts">
  import { CloseIcon } from './icons'

  export let value = ''
  export let disabled = false

  const clear = () => {
    value = ''
  }
</script>

<div class="input flex relative items-center">
  {#if $$slots.icon}
    <span class="absolute inset-y-0 left-0 flex items-center pl-2">
      <slot name="icon" />
    </span>
  {/if}

  <input class="flex-1 py-2 rounded-md pl-10" bind:value {disabled} />

  {#if value && !disabled}
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <!-- svelte-ignore a11y-no-static-element-interactions -->
    <div class="cursor-pointer px-2 z-10" on:click={clear} data-testid="clear">
      <CloseIcon />
    </div>
  {/if}
</div>
