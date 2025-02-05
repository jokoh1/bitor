<script lang="ts">
  export let title: string;
</script>

<div class="mt-10 flex items-center justify-center">
  <div class="my-4 flex w-screen max-w-lg flex-col gap-16">
    <div class="card bg-base-300 shadow-xl">
      <div class="card-body">
        <h2 class="card-title">{title}</h2>
        <slot />
      </div>
    </div>
  </div>
</div> 