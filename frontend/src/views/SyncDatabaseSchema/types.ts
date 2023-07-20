import { ChangeHistory } from "@/types/proto/v1/database_service";

export type SourceSchemaType = "DATABASE_SCHEMA" | "SCHEMA_DESIGN";

export interface ChangeHistorySourceSchema {
  projectId?: string;
  environmentId?: string;
  databaseId?: string;
  changeHistory?: ChangeHistory;
}
