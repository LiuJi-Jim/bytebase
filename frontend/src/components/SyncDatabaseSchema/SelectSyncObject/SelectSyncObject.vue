<template>
  <NTransfer
    v-model:value="checkedKeys"
    :options="flattenSourceOptions"
    :render-source-list="renderSourceList"
    :render-target-list="renderTargetList"
    style="width: 761px; height: 441px"
  />
</template>

<script setup lang="ts">
import { NTransfer } from "naive-ui";
import { v1 as uuidv1 } from "uuid";
import { computed, h, ref } from "vue";
import Source from "./Source.vue";
import Target from "./Target.vue";
import { SyncSchemaTransferOption, flattenOptions } from "./common";

const generateKey = (
  status: SyncSchemaTransferOption["status"],
  type: SyncSchemaTransferOption["type"]
) => {
  return [status, type, uuidv1()].join("--");
};

const sourceOptions: SyncSchemaTransferOption[] = [
  {
    label: "employee",
    value: generateKey("created", "table"),
    status: "created",
    type: "table",
    isLeaf: false,
    children: [
      {
        label: "phone",
        value: generateKey("created", "column"),
        status: "created",
        type: "column",
        isLeaf: true,
      },
      {
        label: "idx_employee_phone",
        value: generateKey("created", "index"),
        status: "created",
        type: "index",
        isLeaf: true,
      },
      {
        label: "semantic_type_phone",
        value: generateKey("created", "config"),
        status: "created",
        type: "config",
        isLeaf: true,
      },
    ],
  },
  {
    label: "employee",
    value: generateKey("updated", "table"),
    status: "updated",
    type: "table",
    isLeaf: false,
    children: [
      {
        label: "(`emp_no`, `phone`)",
        value: generateKey("updated", "primary-key"),
        status: "updated",
        type: "primary-key",
        isLeaf: true,
      },
      {
        label: "hire_date",
        value: generateKey("updated", "column"),
        status: "updated",
        type: "column",
        isLeaf: true,
      },
    ],
  },
  {
    label: "title",
    value: generateKey("dropped", "table"),
    status: "dropped",
    type: "table",
    isLeaf: false,
    children: [
      {
        label: "level",
        value: generateKey("dropped", "column"),
        status: "dropped",
        type: "column",
        isLeaf: true,
      },
    ],
  },
];

const checkedKeys = ref<string[]>([]);

const flattenSourceOptions = computed((): SyncSchemaTransferOption[] => {
  return flattenOptions(sourceOptions);
});

const renderSourceList = (props: {
  onCheck: (values: string[]) => void;
  checkedOptions: SyncSchemaTransferOption[];
  pattern: string;
}) => {
  return h(Source, {
    options: sourceOptions,
    flattenOptions: flattenSourceOptions.value,
    checkedKeys: checkedKeys.value,
    "onUpdate:checked-keys": (keys: string[]) => props.onCheck(keys),
  });
};

const renderTargetList = (props: {
  onCheck: (values: string[]) => void;
  checkedOptions: SyncSchemaTransferOption[];
  pattern: string;
}) => {
  return h(Target, {
    checkedOptions: props.checkedOptions,
    checkedKeys: checkedKeys.value,
    "onUpdate:checked-keys": (keys: string[]) => props.onCheck(keys),
  });
};
</script>
