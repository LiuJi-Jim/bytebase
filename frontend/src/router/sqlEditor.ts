import type { RouteRecordRaw } from "vue-router";
import SQLEditorLayout from "@/layouts/SQLEditorLayout.vue";

export const SQL_EDITOR_HOME_MODULE = "sql-editor.home";
export const SQL_EDITOR_PROJECT_MODULE = "sql-editor.project";
export const SQL_EDITOR_INSTANCE_MODULE = "sql-editor.instance";
export const SQL_EDITOR_DATABASE_MODULE = "sql-editor.database";
export const SQL_EDITOR_WORKSHEET_MODULE = "sql-editor.worksheet";
export const SQL_EDITOR_DETAIL_MODULE_LEGACY = "sql-editor.legacy-detail";
export const SQL_EDITOR_SHARE_MODULE_LEGACY = "sql-editor.legacy-share";
export const SQL_EDITOR_SETTING_MODULE = "sql-editor.setting";
export const SQL_EDITOR_SETTING_GENERAL_MODULE = "sql-editor.setting.general";
export const SQL_EDITOR_SETTING_SUBSCRIPTION_MODULE =
  "sql-editor.setting.subscription";
export const SQL_EDITOR_SETTING_INSTANCE_MODULE = "sql-editor.setting.instance";
export const SQL_EDITOR_SETTING_ENVIRONMENT_MODULE =
  "sql-editor.setting.environment";
export const SQL_EDITOR_SETTING_PROJECT_MODULE = "sql-editor.setting.project";
export const SQL_EDITOR_SETTING_DATABASES_MODULE =
  "sql-editor.setting.databases";
export const SQL_EDITOR_SETTING_USERS_MODULE = "sql-editor.setting.users";
export const SQL_EDITOR_SETTING_MEMBERS_MODULE = "sql-editor.setting.members";

const sqlEditorRoutes: RouteRecordRaw[] = [
  {
    path: "/sql-editor",
    name: "sql-editor",
    component: SQLEditorLayout,
    children: [
      {
        path: "",
        name: SQL_EDITOR_HOME_MODULE,
        meta: { title: () => "Bytebase SQL Editor" },
        component: () => import("../views/sql-editor/SQLEditorHomePage.vue"),
      },
      {
        path: "projects/:project",
        name: SQL_EDITOR_PROJECT_MODULE,
        meta: { title: () => "Bytebase SQL Editor" },
        component: () => import("../views/sql-editor/SQLEditorHomePage.vue"),
      },
      {
        path: "projects/:project/instances/:instance/databases/:database",
        name: SQL_EDITOR_DATABASE_MODULE,
        meta: { title: () => "Bytebase SQL Editor" },
        component: () => import("../views/sql-editor/SQLEditorHomePage.vue"),
      },
      {
        path: "projects/:project/instances/:instance",
        name: SQL_EDITOR_INSTANCE_MODULE,
        meta: { title: () => "Bytebase SQL Editor" },
        component: () => import("../views/sql-editor/SQLEditorHomePage.vue"),
      },
      {
        path: "projects/:project/sheets/:sheet",
        name: SQL_EDITOR_WORKSHEET_MODULE,
        meta: { title: () => "Bytebase SQL Editor" },
        component: () => import("../views/sql-editor/SQLEditorHomePage.vue"),
      },
      {
        path: ":connectionSlug",
        name: SQL_EDITOR_DETAIL_MODULE_LEGACY,
        meta: { title: () => "Bytebase SQL Editor" },
        component: () => import("../views/sql-editor/SQLEditorHomePage.vue"),
      },
      {
        path: "sheet/:sheetSlug",
        name: SQL_EDITOR_SHARE_MODULE_LEGACY,
        meta: { title: () => "Bytebase SQL Editor" },
        component: () => import("../views/sql-editor/SQLEditorHomePage.vue"),
      },
      {
        path: "sheet/:sheetSlug",
        name: SQL_EDITOR_SHARE_MODULE_LEGACY,
        meta: { title: () => "Bytebase SQL Editor" },
        component: () => import("../views/sql-editor/SQLEditorHomePage.vue"),
      },
      {
        path: "setting",
        name: SQL_EDITOR_SETTING_MODULE,
        meta: { title: () => "Bytebase SQL Editor" },
        component: () => import("../views/sql-editor/SQLEditorSettingPage.vue"),
        children: [
          {
            path: "general",
            name: SQL_EDITOR_SETTING_GENERAL_MODULE,
            meta: {
              requiredWorkspacePermissionList: () => ["bb.settings.get"],
            },
            component: () => import("../views/sql-editor/Setting/General"),
          },
          {
            path: "subscription",
            name: SQL_EDITOR_SETTING_SUBSCRIPTION_MODULE,
            meta: {
              requiredWorkspacePermissionList: () => ["bb.settings.get"],
            },
            component: () => import("../views/sql-editor/Setting/Subscription"),
          },
          {
            path: "instance",
            name: SQL_EDITOR_SETTING_INSTANCE_MODULE,
            meta: {
              requiredWorkspacePermissionList: () => [
                "bb.instances.list",
                "bb.instances.create",
                "bb.instances.update",
              ],
            },
            component: () => import("../views/sql-editor/Setting/Instance"),
          },
          {
            path: "project",
            name: SQL_EDITOR_SETTING_PROJECT_MODULE,
            meta: {
              requiredWorkspacePermissionList: () => [
                "bb.projects.list",
                "bb.projects.create",
              ],
            },
            component: () => import("../views/sql-editor/Setting/Project"),
          },
          {
            path: "environment",
            name: SQL_EDITOR_SETTING_ENVIRONMENT_MODULE,
            meta: {
              requiredWorkspacePermissionList: () => [
                "bb.environments.list",
                "bb.environments.create",
                "bb.environments.update",
              ],
            },
            component: () => import("../views/sql-editor/Setting/Environment"),
          },
          {
            path: "database",
            name: SQL_EDITOR_SETTING_DATABASES_MODULE,
            meta: {
              requiredWorkspacePermissionList: () => ["bb.databases.list"],
            },
            component: () => import("../views/sql-editor/Setting/Database"),
          },
          {
            path: "users",
            name: SQL_EDITOR_SETTING_USERS_MODULE,
            component: () => import("../views/sql-editor/Setting/Users"),
          },
          {
            path: "members",
            name: SQL_EDITOR_SETTING_MEMBERS_MODULE,
            meta: {
              requiredWorkspacePermissionList: () => ["bb.policies.get"],
            },
            component: () => import("../views/sql-editor/Setting/Members"),
          },
        ],
      },
    ],
  },
];

export default sqlEditorRoutes;
