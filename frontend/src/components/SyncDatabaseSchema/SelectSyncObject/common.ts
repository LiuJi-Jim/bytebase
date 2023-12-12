import { TransferOption } from "naive-ui";

export type SyncSchemaTransferOption = TransferOption & {
  value: string;
  status: "created" | "dropped" | "updated";
  type:
    | "schema"
    | "table"
    | "view"
    | "column"
    | "primary-key"
    | "index"
    | "config";
  isLeaf: boolean;
  children?: SyncSchemaTransferOption[];
};

export const flattenOption = (option: SyncSchemaTransferOption) => {
  const result = [option];
  if (option.children) {
    const descendants = option.children.flatMap(flattenOption);
    result.push(...descendants);
  }
  return result;
};

export const flattenOptions = (options: SyncSchemaTransferOption[]) => {
  return options.flatMap(flattenOption);
};
