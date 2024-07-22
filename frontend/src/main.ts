import axios from "axios";
import isEmpty from "lodash-es/isEmpty";
import Long from "long";
import protobufjs from "protobufjs";
import { computed, createApp } from "vue";
import App from "./App.vue";
import "./assets/css/github-markdown-style.css";
import "./assets/css/inter.css";
import "./assets/css/tailwind.css";
import dataSourceType from "./directives/data-source-type";
import dayjs from "./plugins/dayjs";
import highlight from "./plugins/highlight";
import i18n from "./plugins/i18n";
import NaiveUI from "./plugins/naive-ui";
import { isSilent } from "./plugins/silent-request";
import { router } from "./router";
import { AUTH_SIGNIN_MODULE } from "./router/auth";
import type { PageMode } from "./store";
import {
  pinia,
  pushNotification,
  useActuatorV1Store,
  useAuthStore,
} from "./store";
import {
  environmentName,
  humanizeTs,
  humanizeDuration,
  humanizeDurationV1,
  humanizeDate,
  instanceName,
  isDev,
  isRelease,
  projectName,
  sizeToFit,
  urlfy,
} from "./utils";
import { useCustomTheme } from "./utils/customTheme";

protobufjs.util.Long = Long;
protobufjs.configure();

console.debug("dev:", isDev());
console.debug("release:", isRelease());

axios.defaults.timeout = 30000;
axios.interceptors.request.use((request) => {
  if (isDev() && request.url!.startsWith("/api")) {
    console.debug(
      request.method?.toUpperCase() + " " + request.url + " request",
      JSON.stringify(request, null, 2)
    );
  }

  return request;
});

axios.interceptors.response.use(
  (response) => {
    if (isDev() && response.config.url!.startsWith("/api")) {
      console.debug(
        response.config.method?.toUpperCase() +
          " " +
          response.config.url +
          " response",
        JSON.stringify(response.data, null, 2)
      );
    }
    return response;
  },
  async (error) => {
    if (error.response) {
      // When receiving 401 and is returned by our server, it means the current
      // login user's token becomes invalid. Thus we force a logout.
      // We could receive 401 when calling external service such as VCS provider,
      // in such case, we shouldn't logout.
      if (error.response.status == 401) {
        const origin = location.origin;
        const pathname = location.pathname;
        if (
          pathname !== "/auth/mfa" &&
          error.response.request.responseURL.startsWith(origin)
        ) {
          // If the request URL starts with the browser's location origin
          // e.g. http://localhost:3000/
          // we know this is a request to Bytebase API endpoint (not an external service).
          // Means that the auth session is error or expired.
          // So we need to "kick out" here.
          try {
            await useAuthStore().logout();
          } finally {
            router.push({ name: AUTH_SIGNIN_MODULE });
          }
        }
      }

      // in such case, we shouldn't logout.
      if (error.response.status == 403) {
        const origin = location.origin;
        if (error.response.request.responseURL.startsWith(origin)) {
          // If the request URL starts with the browser's location origin
          // e.g. http://localhost:3000/
          // we know this is a request to Bytebase API endpoint (not an external service).
          // Means that the API request is denied by authorization reasons.
          router.push({ name: "error.403" });
        }
      }

      if (error.response.data?.message && !isSilent()) {
        pushNotification({
          module: "bytebase",
          style: "CRITICAL",
          title: error.response.data.message,
          // If server enables --debug, then the response will include the detailed error.
          description: error.response.data.error
            ? error.response.data.error
            : undefined,
        });
      }
    } else if (error.code == "ECONNABORTED" && !isSilent()) {
      pushNotification({
        module: "bytebase",
        style: "CRITICAL",
        title: "Connecting server timeout. Make sure the server is running.",
      });
    }

    throw error;
  }
);
const app = createApp(App);
// Allow template to access various function
app.config.globalProperties.window = window;
app.config.globalProperties.console = console;
app.config.globalProperties.dayjs = dayjs;
app.config.globalProperties.humanizeTs = humanizeTs;
app.config.globalProperties.humanizeDuration = humanizeDuration;
app.config.globalProperties.humanizeDurationV1 = humanizeDurationV1;
app.config.globalProperties.humanizeDate = humanizeDate;
app.config.globalProperties.isDev = isDev();
app.config.globalProperties.isRelease = isRelease();
app.config.globalProperties.sizeToFit = sizeToFit;
app.config.globalProperties.urlfy = urlfy;
app.config.globalProperties.isEmpty = isEmpty;
app.config.globalProperties.environmentName = environmentName;
app.config.globalProperties.projectName = projectName;
app.config.globalProperties.instanceName = instanceName;

app
  // Need to use a directive on the element.
  // The normal hljs.initHighlightingOnLoad() won't work because router change would cause vue
  // to re-render the page and remove the event listener required for
  .directive("data-source-type", dataSourceType)
  .use(pinia);

const overrideAppProfile = () => {
  const query = new URLSearchParams(window.location.search);
  const actuatorStore = useActuatorV1Store();
  const mode = query.get("mode") as PageMode;
  if (mode === "STANDALONE") {
    // mode=STANDALONE is not easy to read, but for legacy support we keep it as
    // some customers are using it.
    actuatorStore.overrideAppProfile({
      "bb.feature.embedded-in-iframe": true,
      "bb.feature.hide-help": true,
      "bb.feature.hide-quick-start": true,
      "bb.feature.hide-release-remind": true,
      "bb.feature.disallow-share-worksheet": true,
      "bb.feature.disallow-navigate-to-console": true,
      "bb.feature.disallow-navigate-away-sql-editor": true,
    });
  }
  const customTheme = query.get("customTheme");
  if (customTheme === "lixiang") {
    actuatorStore.overrideAppProfile({
      "bb.feature.custom-query-datasource": true,
      "bb.feature.disallow-export-query-data": true,
      "bb.feature.custom-color-scheme": {
        "--color-accent": "#00665f",
        "--color-accent-hover": "#00554f",
        "--color-accent-disabled": "#b8c3c3",
      },
    });
    if (actuatorStore.appProfile["bb.feature.embedded-in-iframe"]) {
      actuatorStore.overrideAppProfile({
        "bb.feature.hide-issue-review-actions": true,
      });
    }
  }

  useCustomTheme(
    computed(
      () => actuatorStore.appProfile["bb.feature.custom-color-scheme"]
    )
  );
};

const overrideLang = () => {
  const query = new URLSearchParams(window.location.search);
  const lang = query.get("lang");
  if (lang) {
    i18n.global.locale.value = lang;
  }
};

const initSearchParams = () => {
  overrideAppProfile();
  overrideLang();
};

initSearchParams();

app.use(router).use(highlight).use(i18n).use(NaiveUI);
app.mount("#app");
