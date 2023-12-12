<template>
  <NTree
    key-field="value"
    :data="checkedLeafOptions"
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
</template>

<script setup lang="ts">
import { TreeOption } from "naive-ui";
import { computed, h } from "vue";
import Label from "./Label.vue";
import Suffix from "./Suffix.vue";
import { SyncSchemaTransferOption, flattenOptions } from "./common";

const props = defineProps<{
  checkedOptions: SyncSchemaTransferOption[];
  checkedKeys: string[];
}>();
defineEmits<{
  (event: "update:checked-keys", keys: string[]): void;
}>();

const checkedLeafOptions = computed(() => {
  return props.checkedOptions.filter((opt) => opt.isLeaf);
});

const renderLabel = (info: { option: TreeOption; checked: boolean }) => {
  const option = info.option as any as SyncSchemaTransferOption;
  return h(Label, { option, checked: info.checked });
};
const renderSuffix = (info: { option: TreeOption; checked: boolean }) => {
  const option = info.option as any as SyncSchemaTransferOption;
  return h(Suffix, { option, checked: info.checked });
};
</script>
