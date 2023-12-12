<template>
  <NCollapse
    class="bb-select-sync-object-collapse"
    :default-value="['created', 'updated', 'dropped']"
  >
    <SourceSection
      title="新增的对象"
      status="created"
      :options="optionsByStatus.get('created') ?? []"
      :flatten-options="flattenOptionsByStatus('created')"
      :checked-keys="checkedKeysByStatus('created')"
      @update:checked-keys="updateCheckedKeysStatus('created', $event)"
    />
    <SourceSection
      title="修改的对象"
      status="updated"
      :options="optionsByStatus.get('updated') ?? []"
      :flatten-options="flattenOptionsByStatus('updated')"
      :checked-keys="checkedKeysByStatus('updated')"
      @update:checked-keys="updateCheckedKeysStatus('updated', $event)"
    />
    <SourceSection
      title="删除的对象"
      status="dropped"
      :options="optionsByStatus.get('dropped') ?? []"
      :flatten-options="flattenOptionsByStatus('dropped')"
      :checked-keys="checkedKeysByStatus('dropped')"
      @update:checked-keys="updateCheckedKeysStatus('dropped', $event)"
    />
  </NCollapse>
</template>

<script setup lang="ts">
import { NCollapse } from "naive-ui";
import { computed } from "vue";
import { groupBy } from "@/utils";
import { SyncSchemaTransferOption } from "./common";

const props = defineProps<{
  options: SyncSchemaTransferOption[];
  flattenOptions: SyncSchemaTransferOption[];
  checkedKeys: string[];
}>();
const emit = defineEmits<{
  (event: "update:checked-keys", keys: string[]): void;
}>();

const optionsByStatus = computed(() => {
  return groupBy(props.options, (opt) => opt.status);
});
const flattenOptionsByStatus = (status: SyncSchemaTransferOption["status"]) => {
  return props.flattenOptions.filter((opt) => opt.status === status);
};
const checkedKeysByStatus = (status: SyncSchemaTransferOption["status"]) => {
  return props.checkedKeys.filter((key) => key.startsWith(`${status}--`));
};
const updateCheckedKeysStatus = (
  status: SyncSchemaTransferOption["status"],
  keys: string[]
) => {
  const irrelevantKeys = props.checkedKeys.filter(
    (key) => !key.startsWith(`${status}--`)
  );
  emit("update:checked-keys", [...irrelevantKeys, ...keys]);
};
</script>

<style lang="postcss" scoped>
.bb-select-sync-object-collapse :deep(.n-collapse-item__header) {
  padding: 0.25rem 0 !important;
}
</style>
