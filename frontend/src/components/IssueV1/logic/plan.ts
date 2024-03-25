import { Plan_Spec } from "@/types/proto/v1/rollout_service";

/**
 *
 * @returns empty string if no sheet found
 */
export const sheetNameForSpec = (spec: Plan_Spec): string => {
  return spec.changeDatabaseConfig?.sheet ?? "";
};

export const targetForSpec = (spec: Plan_Spec | undefined) => {
  if (!spec) return undefined;
  return (
    spec.changeDatabaseConfig?.target ??
    spec.createDatabaseConfig?.target ??
    undefined
  );
};
