<template>
  <div class="app-shell">
    <PolygonBackground />
    <div class="app-content">
      <router-view />
    </div>
    <AppDialog />
  </div>
</template>

<script setup lang="ts">
import AppDialog from "@/components/AppDialog.vue";
import PolygonBackground from "@/components/PolygonBackground.vue";
</script>
