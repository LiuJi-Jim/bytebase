<template>
  <div class="flex flex-row items-center justify-start gap-x-1">
    <span class="w-4 h-4 flex items-center justify-center">
      <TableIcon v-if="option.type === 'table'" />
      <ColumnIcon v-if="option.type === 'column'" />
      <PrimaryKeyIcon v-if="option.type === 'primary-key'" class="w-3 h-3" />
      <IndexIcon v-if="option.type === 'index'" class="w-3.5 h-3.5" />
      <SettingsIcon v-if="option.type === 'config'" class="w-3.5 h-3.5" />
    </span>

    <span v-if="option.type === 'primary-key'" class="text-control-placeholder"
      >(Primary Key)</span
    >
    <span v-if="option.type === 'index'" class="text-control-placeholder"
      >(Index)</span
    >
    <span v-if="option.type === 'config'" class="text-control-placeholder"
      >(Table Configs)</span
    >
    <span
      :class="option.isLeaf && option.status === 'dropped' && 'line-through'"
      >{{ option.label }}</span
    >
  </div>
</template>

<script setup lang="ts">
import { SettingsIcon } from "lucide-vue-next";
import {
  TableIcon,
  ColumnIcon,
  IndexIcon,
  PrimaryKeyIcon,
} from "@/components/Icon";
import { SyncSchemaTransferOption } from "./common";

defineProps<{
  option: SyncSchemaTransferOption;
  checked?: boolean;
}>();
</script>
