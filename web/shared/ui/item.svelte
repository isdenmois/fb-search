<script lang="ts">
  export let title: string
  export let suptitle = ''
  export let subtitle = ''
</script>

<li class="flex flex-row items-center cursor-pointer" role="listitem" aria-label={title}>
  {#if $$slots.icon}
    <div class="mr-1"><slot name="icon" /></div>
  {/if}

  <div class="flex-1">
    {#if subtitle}
      <div class="text-xs text-secondary">{suptitle}</div>
    {/if}

    <div class="text-primary">{title}</div>

    {#if subtitle}
      <div class="text-xs text-secondary">{subtitle}</div>
    {/if}
  </div>
</li>
