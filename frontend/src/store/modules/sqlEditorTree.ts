import { useLocalStorage } from "@vueuse/core";
import { cloneDeep, isFunction, orderBy } from "lodash-es";
import { defineStore, storeToRefs } from "pinia";
import { v4 as uuidv4 } from "uuid";
import { computed, reactive, ref, watch } from "vue";
import {
  ComposedDatabase,
  ComposedInstance,
  ComposedProject,
  SQLEditorTreeFactor as Factor,
  SQLEditorTreeNode as TreeNode,
  SQLEditorTreeNodeTarget as NodeTarget,
  SQLEditorTreeNodeType as NodeType,
  SQLEditorTreeState as TreeState,
  isValidSQLEditorTreeFactor as isValidFactor,
  extractSQLEditorLabelFactor as extractLabelFactor,
  unknownEnvironment,
  UNKNOWN_ID,
  RichSchemaMetadata,
  RichTableMetadata,
  StatefulSQLEditorTreeFactor as StatefulFactor,
  LeafTreeNodeTypes,
  Connection,
  DEFAULT_PROJECT_V1_NAME,
  RichViewMetadata,
  TextTarget,
  RichPartitionTableMetadata,
  RichExternalTableMetadata,
} from "@/types";
import { Environment } from "@/types/proto/v1/environment_service";
import { emptyConnection, getSemanticLabelValue, groupBy } from "@/utils";
import { useFilterStore } from "./filter";
import { useSQLEditorV2Store } from "./sqlEditorV2";
import {
  useDBSchemaV1Store,
  useDatabaseV1Store,
  useEnvironmentV1Store,
  useInstanceV1Store,
  useProjectV1Store,
} from "./v1";

export const ROOT_NODE_ID = "ROOT";

const defaultEnvironmentFactor: StatefulFactor = {
  factor: "environment",
  disabled: false,
};

export const useSQLEditorTreeStore = defineStore("SQL-Editor-Tree", () => {
  const { filter } = useFilterStore();

  const defaultFactorList = (): StatefulFactor[] => {
    return [defaultEnvironmentFactor];
  };

  const factorListInLocalStorage = useLocalStorage<StatefulFactor[]>(
    "bb.sql-editor.tree-factor-list",
    defaultFactorList(),
    {
      serializer: {
        read: (raw: string): StatefulFactor[] => {
          try {
            const array = JSON.parse(raw) as StatefulFactor[];
            if (!Array.isArray(array)) {
              throw new Error();
            }
            const factorList: StatefulFactor[] = [];
            array.forEach((sf) => {
              if (isValidFactor(sf.factor)) {
                factorList.push({
                  factor: sf.factor,
                  disabled: !!sf.disabled,
                });
              }
            });
            if (factorList.length === 0) {
              throw new Error();
            }
            return factorList;
          } catch {
            return defaultFactorList();
          }
        },
        write: (factorList) => JSON.stringify(factorList),
      },
    }
  );
  const nodeListMapById = reactive(new Map<string, TreeNode[]>());
  // states
  // just re-expose `databaseList` and `project` in sqlEditor store
  const { databaseList, project } = storeToRefs(useSQLEditorV2Store());
  const factorList = ref<StatefulFactor[]>(
    cloneDeep(factorListInLocalStorage.value)
  );
  const selectedProject = computed(() => {
    if (project.value) {
      return useProjectV1Store().getProjectByName(project.value);
    }
    return undefined;
  });

  const filteredDatabaseList = computed(() => {
    if (filter.database) {
      return databaseList.value.filter((database) => {
        return database.name === filter.database;
      });
    }
    if (filter.project || selectedProject.value) {
      const projectName = filter.project ?? selectedProject.value?.name;
      return databaseList.value.filter((database) => {
        return database.project === projectName;
      });
    }

    return databaseList.value;
  });
  const filteredFactorList = computed(() => {
    return factorList.value
      .filter((sf) => {
        if (!selectedProject.value) {
          return true;
        }
        return sf.factor !== "project";
      })
      .filter((sf) => !sf.disabled)
      .map((sf) => sf.factor);
  });
  const state = ref<TreeState>("UNSET");
  const tree = ref<TreeNode[]>([]);

  const collectNode = <T extends NodeType>(node: TreeNode<T>) => {
    const { type, target } = node.meta;
    const id = idForSQLEditorTreeNodeTarget(type, target);
    const nodeList = nodeListMapById.get(id) ?? [];
    nodeList.push(node);
    nodeListMapById.set(id, nodeList);
  };
  const nodesByTarget = <T extends NodeType>(
    type: T,
    target: NodeTarget<T>
  ) => {
    const id = idForSQLEditorTreeNodeTarget(type, target);
    return nodeListMapById.get(id) ?? [];
  };

  const buildTree = () => {
    nodeListMapById.clear();
    tree.value = buildTreeImpl(
      filteredDatabaseList.value,
      filteredFactorList.value
    );
  };
  const fetchConnectionByInstanceIdAndDatabaseId = async (
    instanceId: string,
    databaseId: string
  ): Promise<Connection> => {
    try {
      const [db, _] = await Promise.all([
        useDatabaseV1Store().getOrFetchDatabaseByUID(
          databaseId,
          true /* silent */
        ),
        useInstanceV1Store().getOrFetchInstanceByUID(instanceId),
      ]);
      await useDBSchemaV1Store().getOrFetchTableList(db.name);

      return {
        instanceId,
        databaseId,
      };
    } catch {
      // Fallback to disconnected if error occurs such as 404.
      return { instanceId: String(UNKNOWN_ID), databaseId: String(UNKNOWN_ID) };
    }
  };
  const fetchConnectionByInstanceId = async (
    instanceId: string
  ): Promise<Connection> => {
    try {
      await useInstanceV1Store().getOrFetchInstanceByUID(instanceId);

      return {
        instanceId,
        databaseId: String(UNKNOWN_ID),
      };
    } catch {
      // Fallback to disconnected if error occurs such as 404.
      return { instanceId: String(UNKNOWN_ID), databaseId: String(UNKNOWN_ID) };
    }
  };
  const cleanup = () => {
    tree.value = [];
    // expandedKeys.value = [];
    factorList.value = defaultFactorList();
    nodeListMapById.clear();
    state.value = "UNSET";
  };

  watch(
    factorList,
    (factorList) => {
      factorListInLocalStorage.value = factorList;
    },
    { immediate: true, deep: true }
  );
  watch(
    project,
    (project) => {
      if (project) {
        const position = factorList.value.findIndex(
          (sf) => sf.factor === "project"
        );
        if (position > 0) {
          factorList.value.splice(position, 1);
          if (factorList.value.length === 0) {
            factorList.value = defaultFactorList();
          }
        }
      }
    },
    {
      immediate: true,
    }
  );

  return {
    databaseList,
    factorList,
    filteredFactorList,
    selectedProject,
    state,
    tree,
    collectNode,
    nodesByTarget,
    buildTree,
    fetchConnectionByInstanceIdAndDatabaseId,
    fetchConnectionByInstanceId,
    cleanup,
  };
});

export const searchConnectionByName = (
  instanceId: string,
  databaseId: string,
  instanceName: string,
  databaseName: string
): Connection => {
  const connection = emptyConnection();
  const store = useSQLEditorTreeStore();

  if (instanceId !== String(UNKNOWN_ID)) {
    // If we found instanceId and/or databaseId, use the IDs first.
    connection.instanceId = instanceId;
    if (databaseId !== String(UNKNOWN_ID)) {
      connection.databaseId = databaseId;
    }

    return connection;
  }

  // Search the instance and database by name otherwise.
  // Remain this part for legacy sheet support.
  const rootNodes = store.tree;
  for (let i = 0; i < rootNodes.length; i++) {
    const maybeInstanceNode = rootNodes[i];
    if (maybeInstanceNode.meta.type !== "instance") {
      // Skip if we met dirty data.
      continue;
    }
    if (maybeInstanceNode.label === instanceName) {
      connection.instanceId = (
        maybeInstanceNode.meta.target as ComposedInstance
      ).uid;
      if (databaseName) {
        const { children = [] } = maybeInstanceNode;
        for (let j = 0; j < children.length; j++) {
          const maybeDatabaseNode = children[j];
          if (maybeDatabaseNode.meta.type !== "database") {
            // Skip if we met dirty data.
            continue;
          }
          if (maybeDatabaseNode.label === databaseName) {
            connection.databaseId = (
              maybeDatabaseNode.meta.target as ComposedDatabase
            ).uid;
            // Don't go further since we've found the databaseId
            break;
          }
        }
      }
      // Don't go further since we've found the instanceId
      break;
    }
  }

  return connection;
};

export const keyForSQLEditorTreeNodeTarget = <T extends NodeType>(
  type: T,
  target: NodeTarget<T>
): string => {
  let prefix: string = type;
  if (type === "label") {
    prefix = `label:${(target as NodeTarget<"label">).key}`;
  }
  return `${prefix}-${uuidv4()}`;
};

export const idForSQLEditorTreeNodeTarget = <T extends NodeType>(
  type: T,
  target: NodeTarget<T>
) => {
  if (
    type === "project" ||
    type === "instance" ||
    type === "environment" ||
    type === "database"
  ) {
    return (
      target as
        | ComposedProject
        | ComposedInstance
        | Environment
        | ComposedDatabase
    ).name;
  }
  if (type === "schema") {
    const { database, schema } = target as RichSchemaMetadata;
    return `${database.name}/schemas/${schema.name || "-"}`;
  }
  if (type === "table") {
    const { database, schema, table } = target as RichTableMetadata;
    return `${database.name}/schemas/${schema.name || "-"}/tables/${
      table.name
    }`;
  }
  if (type === "external-table") {
    const { database, schema, externalTable } =
      target as RichExternalTableMetadata;
    return `${database.name}/schemas/${schema.name || "-"}/externalTables/${
      externalTable.name
    }`;
  }
  if (type === "partition-table") {
    const { database, schema, table, partition } =
      target as RichPartitionTableMetadata;
    return `${database.name}/schemas/${schema.name || "-"}/tables/${
      table.name
    }/partitions/${partition.name}`;
  }
  if (type === "view") {
    const { database, schema, view } = target as RichViewMetadata;
    return `${database.name}/schemas/${schema.name || "-"}/views/${view.name}`;
  }
  if (type === "label") {
    const kv = target as NodeTarget<"label">;
    return `labels/${kv.key}:${kv.value}`;
  }
  if (type === "expandable-text") {
    const { text, type } = target as NodeTarget<"expandable-text">;
    return `texts-${type}/${typeof text === "function" ? text() : text}`;
  }
  if (type === "dummy") {
    const dummyType = (target as NodeTarget<"dummy">).type;
    return `dummy-${dummyType}`;
  }

  throw new Error(
    `should never reach this line, type=${type}, target=${target}, parent=${parent}`
  );
};

const buildSubTree = (
  databaseList: ComposedDatabase[],
  parent: TreeNode | undefined,
  factorList: Factor[],
  factorIndex: number
): TreeNode[] => {
  if (factorIndex === factorList.length) {
    // the last dimension is database nodes
    return databaseList.map((db) => {
      const node = mapTreeNodeByType("database", db, parent);
      return node;
    });
  }

  // group (project, instance, environment, label) nodes
  const nodes: TreeNode[] = [];
  const factor = factorList[factorIndex];

  const groups = groupBy(databaseList, (db) =>
    getSemanticFactorValue(db, factor)
  );
  for (const [value, childrenDatabaseList] of groups) {
    const groupNode = mapGroupNode(factor, value, parent);
    groupNode.children = buildSubTree(
      childrenDatabaseList,
      groupNode,
      factorList,
      factorIndex + 1
    );
    nodes.push(groupNode);
  }
  return sortNodesIfNeeded(nodes, factor);
};

const sortNodesIfNeeded = (nodes: TreeNode[], factor: Factor) => {
  if (factor === "environment") {
    // Sort by environment order DESC. Put production envs to the first.
    return orderBy(
      nodes as TreeNode<"environment">[],
      [(node) => node.meta.target.order],
      ["desc"]
    );
  }
  if (factor.startsWith("label:")) {
    // Sort in lexicographical order, and put <empty value> to the last
    return orderBy(
      nodes as TreeNode<"label">[],
      [
        (node) => (node.meta.target.value ? -1 : 1), // Empty value to the last,
        (node) => node.meta.target.value, // lexicographical order then
      ],
      ["asc", "asc"]
    );
  }
  if (factor === "project") {
    // Put unassigned project to the last
    return orderBy(
      nodes as TreeNode<"project">[],
      [(node) => (node.meta.target.name === DEFAULT_PROJECT_V1_NAME ? 1 : -1)],
      ["asc"]
    );
  }

  return nodes;
};

const buildTreeImpl = (
  databaseList: ComposedDatabase[],
  factorList: Factor[]
): TreeNode[] => {
  return buildSubTree(
    databaseList,
    undefined /* parent */,
    factorList,
    0 /* factorIndex */
  );
};

const mapGroupNode = (
  factor: Factor,
  value: string,
  parent: TreeNode | undefined
) => {
  if (factor === "environment") {
    const environment =
      useEnvironmentV1Store().getEnvironmentByName(value) ??
      unknownEnvironment();
    return mapTreeNodeByType("environment", environment, parent);
  }
  if (factor === "instance") {
    const instance = useInstanceV1Store().getInstanceByName(value);
    return mapTreeNodeByType("instance", instance, parent);
  }
  if (factor === "project") {
    const project = useProjectV1Store().getProjectByName(value);
    return mapTreeNodeByType("project", project, parent);
  }
  // factor is label
  const key = extractLabelFactor(factor);
  if (key) {
    return mapTreeNodeByType(
      "label",
      {
        key,
        value,
      },
      parent
    );
  }

  throw new Error(
    `mapGroupNode: should never reach this line. factor=${factor}, factorValue=${value}`
  );
};

export const mapTreeNodeByType = <T extends NodeType>(
  type: T,
  target: NodeTarget<T>,
  parent: TreeNode | undefined,
  overrides: Partial<TreeNode<T>> | undefined = undefined
): TreeNode<T> => {
  const key = keyForSQLEditorTreeNodeTarget(type, target);
  const node: TreeNode<T> = {
    key,
    meta: { type, target },
    parent,
    label: readableTargetByType(type, target),
    isLeaf: isLeafNodeType(type),
    ...overrides,
  };

  useSQLEditorTreeStore().collectNode(node);

  return node;
};

const readableTargetByType = <T extends NodeType>(
  type: T,
  target: NodeTarget<T>
): string => {
  if (type === "environment") {
    return (target as Environment).title;
  }
  if (type === "project") {
    return (target as ComposedProject).title;
  }
  if (type === "instance") {
    return (target as ComposedInstance).title;
  }
  if (type === "database") {
    return (target as ComposedDatabase).databaseName;
  }
  if (type === "schema") {
    return (target as RichSchemaMetadata).schema.name;
  }
  if (type === "table") {
    return (target as RichTableMetadata).table.name;
  }
  if (type === "external-table") {
    return (target as RichExternalTableMetadata).externalTable.name;
  }
  if (type === "partition-table") {
    return (target as RichPartitionTableMetadata).partition.name;
  }
  if (type === "view") {
    return (target as RichViewMetadata).view.name;
  }
  if (type === "expandable-text") {
    const { text, searchable } = target as TextTarget<true>;
    if (!searchable) return "";
    return isFunction(text) ? text() : text;
  }
  if (type === "dummy") {
    // Use empty strings for dummy nodes to make them unsearchable
    return "";
  }
  return (target as NodeTarget<"label">).value;
};

const isLeafNodeType = (type: NodeType) => {
  return LeafTreeNodeTypes.includes(type);
};

const getSemanticFactorValue = (db: ComposedDatabase, factor: Factor) => {
  switch (factor) {
    case "environment":
      return db.instanceEntity.environment;
    case "instance":
      return db.instance;
    case "project":
      return db.project;
  }

  const key = extractLabelFactor(factor);
  if (key) {
    return getSemanticLabelValue(db, key);
  }
  console.error("should never reach this line", db, factor);
  return "";
};
