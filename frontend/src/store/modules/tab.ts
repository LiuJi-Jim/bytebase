import { defineStore } from "pinia";
import { TabInfo, AnyTabInfo, TabState, Connection } from "@/types";
import { getDefaultTab } from "@/utils";

export const useTabStore = defineStore("tab", {
  state: (): TabState => ({
    tabList: [],
    currentTabId: "",
  }),

  getters: {
    currentTab(state): TabInfo {
      const idx = state.tabList.findIndex(
        (tab: TabInfo) => tab.id === state.currentTabId
      );
      return (idx === -1 ? {} : state.tabList[idx]) as TabInfo;
    },
    hasTabs(state: TabState): boolean {
      return state.tabList.length > 0;
    },
  },

  actions: {
    setTabState(payload: Partial<TabState>) {
      Object.assign(this, payload);
    },
    addTab(payload?: AnyTabInfo) {
      const defaultTab = getDefaultTab();

      const newTab = {
        ...defaultTab,
        ...payload,
      };

      this.setTabState({
        currentTabId: newTab.id,
      });
      this.tabList.push(newTab);
    },
    removeTab(payload: TabInfo) {
      this.tabList.splice(this.tabList.indexOf(payload), 1);
    },
    updateCurrentTab(payload: AnyTabInfo) {
      Object.assign(this.currentTab, payload);
    },
    updateCurrentConnection(connection: Partial<Connection>) {
      Object.assign(this.currentTab.connection, connection);
    },
    setCurrentTabId(payload: string) {
      this.currentTabId = payload;
    },
  },
});
