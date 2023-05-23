<template>
  <router-link
    v-if="user"
    :to="`/u/${extractUserUID(user.name)}`"
    class="hover:underline"
  >
    {{ user.email }}
  </router-link>
  <template v-else>
    {{ auditLog.creator }}
  </template>
</template>

import { computed } from "vue";
<script setup lang="ts">
import { useUserStore } from "@/store";
import { AuditLog } from "@/types";
import { extractUserUID } from "@/utils";
import { computed } from "vue";

const props = defineProps<{
  auditLog: AuditLog;
}>();

const userStore = useUserStore();

const user = computed(() => {
  return userStore.getUserByEmail(props.auditLog.creator);
});
</script>
