import type { RemovableRef } from "@vueuse/core";
import { useLocalStorage } from "@vueuse/core";
import axios from "axios";
import { defineStore, storeToRefs } from "pinia";
import { computed, watchEffect } from "vue";
import { actuatorServiceClient } from "@/grpcweb";
import { useSilentRequest } from "@/plugins/silent-request";
import {
  defaultCustomFeatureMatrix,
  type CustomFeature,
  type CustomFeatureMatrix,
  type Release,
  type ReleaseInfo,
} from "@/types";
import type {
  ActuatorInfo,
  ResourcePackage,
  DebugLog,
} from "@/types/proto/v1/actuator_service";
import { semverCompare } from "@/utils";

const EXTERNAL_URL_PLACEHOLDER =
  "https://www.bytebase.com/docs/get-started/install/external-url";
const GITHUB_API_LIST_BYTEBASE_RELEASE =
  "https://api.github.com/repos/bytebase/bytebase/releases";

export type PageMode =
  // General mode. Console is full-featured and SQL Editor is bundled in the layout.
  | "BUNDLED"
  // Vender customized mode. Hide certain parts (e.g., headers, sidebars) and
  // some features are disabled or hidden.
  | "STANDALONE";

interface ActuatorState {
  pageMode: PageMode;
  serverInfo?: ActuatorInfo;
  resourcePackage?: ResourcePackage;
  releaseInfo: RemovableRef<ReleaseInfo>;
  debugLogList: DebugLog[];
  customFeatureMatrix: CustomFeatureMatrix;
}

export const useActuatorV1Store = defineStore("actuator_v1", {
  state: (): ActuatorState => ({
    pageMode: "BUNDLED",
    serverInfo: undefined,
    resourcePackage: undefined,
    releaseInfo: useLocalStorage("bytebase_release", {
      ignoreRemindModalTillNextRelease: false,
      nextCheckTs: 0,
    }),
    debugLogList: [],
    customFeatureMatrix: defaultCustomFeatureMatrix(),
  }),
  getters: {
    info: (state) => {
      return state.serverInfo;
    },
    brandingLogo: (state) => {
      if (!state.resourcePackage?.logo) {
        return "";
      }
      return new TextDecoder().decode(state.resourcePackage?.logo);
    },
    version: (state) => {
      return state.serverInfo?.version || "";
    },
    gitCommit: (state) => {
      return state.serverInfo?.gitCommit || "";
    },
    isDemo: (state) => {
      return state.serverInfo?.demoName;
    },
    isReadonly: (state) => {
      return state.serverInfo?.readonly || false;
    },
    isDebug: (state) => {
      return state.serverInfo?.debug || false;
    },
    isSaaSMode: (state) => {
      return state.serverInfo?.saas || false;
    },
    needAdminSetup: (state) => {
      return state.serverInfo?.needAdminSetup || false;
    },
    needConfigureExternalUrl: (state) => {
      if (!state.serverInfo) return false;
      const url = state.serverInfo?.externalUrl ?? "";
      return url === "" || url === EXTERNAL_URL_PLACEHOLDER;
    },
    disallowSignup: (state) => {
      return state.serverInfo?.disallowSignup || false;
    },
    hasNewRelease: (state) => {
      return (
        (state.serverInfo?.version === "development" &&
          !!state.releaseInfo.latest?.tag_name) ||
        semverCompare(
          state.releaseInfo.latest?.tag_name ?? "",
          state.serverInfo?.version ?? ""
        )
      );
    },
  },
  actions: {
    setLogo(logo: string) {
      if (this.resourcePackage) {
        this.resourcePackage.logo = new TextEncoder().encode(logo);
      }
    },
    setServerInfo(serverInfo: ActuatorInfo) {
      this.serverInfo = serverInfo;
    },
    async fetchServerInfo() {
      const [serverInfo, resourcePackage] = await Promise.all([
        actuatorServiceClient.getActuatorInfo({}),
        actuatorServiceClient.getResourcePackage({}),
      ]);
      this.setServerInfo(serverInfo);
      this.resourcePackage = resourcePackage;
      return serverInfo;
    },
    async patchDebug({ debug }: { debug: boolean }) {
      const serverInfo = await actuatorServiceClient.updateActuatorInfo({
        actuator: {
          debug,
        },
        updateMask: ["debug"],
      });
      this.setServerInfo(serverInfo);
    },
    async fetchDebugLogList() {
      const { logs } = await actuatorServiceClient.listDebugLog({});
      this.debugLogList = logs;
      return logs;
    },
    async tryToRemindRelease(): Promise<boolean> {
      if (this.serverInfo?.saas ?? false) {
        return false;
      }
      if (!this.releaseInfo.latest) {
        const release = await this.fetchLatestRelease();
        this.releaseInfo.latest = release;
      }
      if (!this.releaseInfo.latest) {
        return false;
      }

      // It's time to fetch the release
      if (new Date().getTime() >= this.releaseInfo.nextCheckTs) {
        const release = await this.fetchLatestRelease();
        if (!release) {
          return false;
        }

        // check till 24 hours later
        this.releaseInfo.nextCheckTs =
          new Date().getTime() + 24 * 60 * 60 * 1000;

        if (semverCompare(release.tag_name, this.releaseInfo.latest.tag_name)) {
          this.releaseInfo.ignoreRemindModalTillNextRelease = false;
        }

        this.releaseInfo.latest = release;
      }

      if (this.releaseInfo.ignoreRemindModalTillNextRelease) {
        return false;
      }

      return this.hasNewRelease;
    },
    async fetchLatestRelease(): Promise<Release | undefined> {
      try {
        const { data: releaseList } = await useSilentRequest(() =>
          axios.get<Release[]>(`${GITHUB_API_LIST_BYTEBASE_RELEASE}?per_page=1`)
        );
        return releaseList[0];
      } catch {
        // It's okay to ignore the failure and just return undefined.
        return;
      }
    },
    overrideCustomFeatureMatrix(overrides: Partial<CustomFeatureMatrix>) {
      Object.assign(this.customFeatureMatrix, overrides);
    },
  },
});

export const useDebugLogList = () => {
  const store = useActuatorV1Store();
  watchEffect(() => store.fetchDebugLogList());

  return storeToRefs(store).debugLogList;
};

export const usePageMode = () => {
  const actuatorStore = useActuatorV1Store();
  const { pageMode } = storeToRefs(actuatorStore);
  return pageMode;
};

export const useCustomFeature = <T extends CustomFeature>(feature: T) => {
  return computed(() => useActuatorV1Store().customFeatureMatrix[feature]);
};
