import { head } from "lodash-es";
import {
  composeInstanceResourceForDatabase,
  useDatabaseV1Store,
  useDBGroupStore,
  useDeploymentConfigV1Store,
  useEnvironmentV1Store,
  useInstanceV1Store,
} from "@/store";
import {
  type ComposedProject,
  UNKNOWN_ID,
  type ComposedIssue,
  unknownDatabase,
  unknownEnvironment,
} from "@/types";
import { Engine, State } from "@/types/proto/v1/common";
import type { Plan_Spec } from "@/types/proto/v1/plan_service";
import {
  extractDatabaseResourceName,
  extractDatabaseGroupName,
  extractDeploymentConfigName,
  getPipelineFromDeploymentScheduleV1,
} from "@/utils";

export const databaseForSpec = (issue: ComposedIssue, spec: Plan_Spec) => {
  const { createDatabaseConfig, changeDatabaseConfig, exportDataConfig } = spec;
  if (createDatabaseConfig !== undefined) {
    const instance = createDatabaseConfig.target;
    const databaseName = createDatabaseConfig.database;
    const instanceEntity = useInstanceV1Store().getInstanceByName(instance);
    return {
      ...unknownDatabase(),
      name: `${instance}/databases/${databaseName}`,
      uid: String(UNKNOWN_ID),
      databaseName,
      instance,
      instanceEntity,
      project: issue.project,
      projectEntity: issue.projectEntity,
      effectiveEnvironment: instanceEntity.environment,
      effectiveEnvironmentEntity: instanceEntity.environmentEntity,
    };
  } else if (
    changeDatabaseConfig !== undefined ||
    exportDataConfig !== undefined
  ) {
    // TODO(steven): handle db group and deployment config.
    const target = (changeDatabaseConfig?.target ??
      exportDataConfig?.target) as string;
    const db = useDatabaseV1Store().getDatabaseByName(target);
    if (db.uid === String(UNKNOWN_ID)) {
      // Database not found, it's probably NOT_FOUND (maybe dropped actually)
      // Mock a database using all known resources
      db.project = issue.project;
      db.projectEntity = issue.projectEntity;
      db.name = target;
      const { instance, databaseName } = extractDatabaseResourceName(db.name);
      db.databaseName = databaseName;
      db.instance = instance;
      const ir = composeInstanceResourceForDatabase(instance, db);
      db.instanceResource = ir;
      db.environment = ir.environment;
      db.effectiveEnvironment = ir.environment;
      db.effectiveEnvironmentEntity =
        useEnvironmentV1Store().getEnvironmentByName(ir.environment) ??
        unknownEnvironment();
      db.syncState = State.DELETED;
    }
    return db;
  }
  return unknownDatabase();
};

/**
 *
 * @returns empty string if no sheet found
 */
export const sheetNameForSpec = (spec: Plan_Spec): string => {
  return spec.changeDatabaseConfig?.sheet ?? spec.exportDataConfig?.sheet ?? "";
};

export const databaseEngineForSpec = async (
  project: ComposedProject,
  specOrTarget?: Plan_Spec | string
) => {
  if (!specOrTarget) return Engine.ENGINE_UNSPECIFIED;

  const getTarget = (specOrTarget: Plan_Spec | string) => {
    if (typeof specOrTarget === "string") {
      return specOrTarget;
    }
    const config =
      specOrTarget.changeDatabaseConfig || specOrTarget.exportDataConfig;
    if (!config) {
      return Engine.ENGINE_UNSPECIFIED;
    }
    return config.target;
  };
  const target = getTarget(specOrTarget);

  if (extractDatabaseResourceName(target).databaseName !== String(UNKNOWN_ID)) {
    const db = await useDatabaseV1Store().getOrFetchDatabaseByName(
      target,
      /* silent */ true
    );
    if (db && db.uid !== String(UNKNOWN_ID)) {
      return db.instanceResource.engine;
    }
  }
  if (extractDatabaseGroupName(target)) {
    const dbGroupStore = useDBGroupStore();
    const dbGroup = await dbGroupStore.getOrFetchDBGroupByName(target);
    // Might be flaky: use the first database in the db group
    const dbName = head(dbGroup.matchedDatabases)?.name;
    if (dbName) {
      const db = await useDatabaseV1Store().getOrFetchDatabaseByName(
        dbName,
        /* silent */ true
      );
      if (db && db.uid !== String(UNKNOWN_ID)) {
        return db.instanceResource.engine;
      }
    }
  }
  if (extractDeploymentConfigName(target)) {
    const deploymentConfig =
      await useDeploymentConfigV1Store().fetchDeploymentConfigByProjectName(
        project.name
      );
    if (deploymentConfig) {
      const databaseList = useDatabaseV1Store().databaseListByProject(
        project.name
      );
      const pipeline = getPipelineFromDeploymentScheduleV1(
        databaseList,
        deploymentConfig.schedule
      );
      const db = head(head(pipeline));
      if (db && db.uid !== String(UNKNOWN_ID)) {
        return db.instanceResource.engine;
      }
    }
  }
  return Engine.ENGINE_UNSPECIFIED;
};
