import { Database, DatabaseId, Table, TableId } from ".";

export enum UIEditorTabType {
  TabForDatabase = "database",
  TabForTable = "table",
}

// Tab context for editing database.
export interface DatabaseTabContext {
  id: string;
  type: UIEditorTabType.TabForDatabase;
  databaseId: DatabaseId;
}

// Tab context for editing table.
export interface TableTabContext {
  id: string;
  type: UIEditorTabType.TabForTable;
  databaseId: DatabaseId;
  tableId: TableId;
  // Save the editing table cache in tab.
  tableCache: Table;
}

export type TabContext = DatabaseTabContext | TableTabContext;

type TabId = string;

export interface UIEditorState {
  tabState: {
    tabMap: Map<TabId, TabContext>;
    currentTabId: TabId;
  };
  databaseList: Database[];
  tableList: Table[];
}

/**
 * Type definition for API message.
 */
export interface DatabaseEdit {
  databaseId: DatabaseId;

  CreateTableList: CreateTableContext[];
}

export interface CreateTableContext {
  name: string;
  type: string;
  engine: string;
  characterSet: string;
  collation: string;
  comment: string;

  addColumnList: AddColumnContext[];
}

export interface AddColumnContext {
  name: string;
  type: string;
  characterSet: string;
  collation: string;
  comment: string;
  nullable: boolean;
  default?: string;
}
