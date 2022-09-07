import { markRaw, reactive } from "vue";
import { isEmpty } from "lodash-es";
import { useI18n } from "vue-i18n";

import {
  parseSQL,
  transformSQL,
  isSelectStatement,
  isMultipleStatements,
  isDDLStatement,
  isDMLStatement,
} from "../components/MonacoEditor/sqlParser";
import { pushNotification, useTabStore, useSQLEditorStore } from "@/store";
import { BBNotificationStyle } from "@/bbkit/types";
import { ExecuteConfig, ExecuteOption } from "@/types";

const useExecuteSQL = () => {
  const { t } = useI18n();
  const tabStore = useTabStore();
  const sqlEditorStore = useSQLEditorStore();

  const state = reactive({
    isLoadingData: false,
  });

  const notify = (
    type: BBNotificationStyle,
    title: string,
    description?: string
  ) => {
    pushNotification({
      module: "bytebase",
      style: type,
      title,
      description,
    });
  };

  const execute = async (
    query: string,
    config: ExecuteConfig,
    option?: Partial<ExecuteOption>
  ) => {
    if (state.isLoadingData) {
      notify("INFO", t("common.tips"), t("sql-editor.can-not-execute-query"));
    }

    const isDisconnected = tabStore.isDisconnected;
    if (isDisconnected) {
      notify("CRITICAL", t("sql-editor.select-connection"));
      return;
    }

    if (isEmpty(query)) {
      notify("CRITICAL", t("sql-editor.notify-empty-statement"));
      return;
    }

    const { data } = parseSQL(query);

    if (data === undefined) {
      notify("CRITICAL", t("sql-editor.notify-invalid-sql-statement"));
      return;
    }

    if (data !== null && !isSelectStatement(data)) {
      if (isMultipleStatements(data)) {
        notify(
          "INFO",
          t("common.tips"),
          t("sql-editor.notify-multiple-statements")
        );
        return;
      }
      // only DDL and DML statements are allowed
      if (isDDLStatement(data) || isDMLStatement(data)) {
        sqlEditorStore.setSQLEditorState({
          isShowExecutingHint: true,
        });
        return;
      }
    }

    if (isMultipleStatements(data)) {
      notify(
        "INFO",
        t("common.tips"),
        t("sql-editor.notify-multiple-statements")
      );
    }

    let selectStatement =
      data !== null ? transformSQL(data, config.databaseType) : query;
    if (option?.explain) {
      selectStatement = `EXPLAIN ${selectStatement}`;
    }

    try {
      state.isLoadingData = true;
      sqlEditorStore.setIsExecuting(true);
      const sqlResultSet = await sqlEditorStore.executeQuery({
        statement: selectStatement,
      });
      // TODO(steven): use BBModel instead of notify to show the advice from SQL review.
      let adviceStatus = "SUCCESS";
      let adviceNotifyMessage = "";
      for (const advice of sqlResultSet.adviceList) {
        if (advice.status === "ERROR") {
          adviceStatus = "ERROR";
        } else if (adviceStatus !== "ERROR") {
          adviceStatus = advice.status;
        }

        adviceNotifyMessage += `${advice.status}: ${advice.title}\n`;
        if (advice.content) {
          adviceNotifyMessage += `${advice.content}\n`;
        }
      }
      if (adviceStatus !== "SUCCESS") {
        const notifyStyle = adviceStatus === "ERROR" ? "CRITICAL" : "WARN";
        notify(
          notifyStyle,
          t("sql-editor.sql-review-result"),
          adviceNotifyMessage
        );
      }
      tabStore.updateCurrentTab({
        // use `markRaw` to prevent vue from monitoring the object change deeply
        queryResult: markRaw(sqlResultSet.data) as any,
        adviceList: sqlResultSet.adviceList,
        executeParams: {
          query,
          config,
          option,
        },
      });
      sqlEditorStore.fetchQueryHistoryList();
    } catch (error) {
      tabStore.updateCurrentTab({
        queryResult: undefined,
        adviceList: undefined,
        executeParams: {
          query,
          config,
          option,
        },
      });
      notify("CRITICAL", error as string);
    } finally {
      state.isLoadingData = false;
      sqlEditorStore.setIsExecuting(false);
    }
  };

  return {
    state,
    execute,
  };
};

export { useExecuteSQL };
