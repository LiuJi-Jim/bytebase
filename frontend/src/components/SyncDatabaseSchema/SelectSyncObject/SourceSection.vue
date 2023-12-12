<template>
  <NCollapseItem :title="title" :name="status">
    <NTree
      key-field="value"
      :data="options"
      :checkable="true"
      :cascade="true"
      :block-line="true"
      :selectable="false"
      :default-expand-all="true"
      :checked-keys="checkedKeys"
      :render-label="renderLabel"
      :render-suffix="renderSuffix"
      @update-checked-keys="$emit('update:checked-keys', $event)"
    />
  </NCollapseItem>
</template>

<script setup lang="ts">
import { NCollapseItem, NTree, TreeOption } from "naive-ui";
import { h } from "vue";
import Label from "./Label.vue";
import Suffix from "./Suffix.vue";
import { SyncSchemaTransferOption } from "./common";

defineProps<{
  title: string;
  status: SyncSchemaTransferOption["status"];
  options: SyncSchemaTransferOption[];
  flattenOptions: SyncSchemaTransferOption[];
  checkedKeys: string[];
}>();
defineEmits<{
  (event: "update:checked-keys", keys: string[]): void;
}>();

const renderLabel = (info: { option: TreeOption; checked: boolean }) => {
  const option = info.option as any as SyncSchemaTransferOption;
  return h(Label, { option, checked: info.checked });
};
const renderSuffix = (info: { option: TreeOption; checked: boolean }) => {
  const option = info.option as any as SyncSchemaTransferOption;
  return h(Suffix, { option, checked: info.checked });
};
</script>
