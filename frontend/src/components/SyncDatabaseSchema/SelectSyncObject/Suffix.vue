<template>
  <div class="flex flex-row items-center justify-start gap-x-1">
    <span v-if="option.isLeaf" class="w-4 h-4 flex items-center justify-center">
      <CreatedIcon
        v-if="option.status === 'created'"
        class="w-4 h-4 text-green-700"
      />
      <UpdatedIcon
        v-if="option.status === 'updated'"
        class="w-4 h-4 text-yellow-700"
      />
      <DroppedIcon
        v-if="option.status === 'dropped'"
        class="w-4 h-4 text-red-700"
      />
    </span>
  </div>
</template>

<script setup lang="ts">
import {
  FileXIcon as DroppedIcon,
  FilePlusIcon as CreatedIcon,
  FileCodeIcon as UpdatedIcon,
} from "lucide-vue-next";
import { SyncSchemaTransferOption } from "./common";

defineProps<{
  option: SyncSchemaTransferOption;
  checked?: boolean;
}>();
</script>
