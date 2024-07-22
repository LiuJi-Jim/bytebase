import { useActuatorV1Store, useAppFeature } from "./store";
import { useCustomTheme } from "./utils/customTheme";

export const overrideAppProfile = () => {
  const query = new URLSearchParams(window.location.search);
  const actuatorStore = useActuatorV1Store();
  const mode = query.get("mode");
  if (mode === "STANDALONE") {
    // The webapp is embedded within iframe
    actuatorStore.appProfile.embedded = true;

    // mode=STANDALONE is not easy to read, but for legacy support we keep it as
    // some customers are using it.
    actuatorStore.overrideAppFeatures({
      "bb.feature.embedded-in-iframe": true,
      "bb.feature.disable-kbar": true,
      "bb.feature.disable-schema-editor": true,
      "bb.feature.databases.operations": new Set(["CHANGE-DATA", "EDIT-SCHEMA"]),
      "bb.feature.hide-banner": true,
      "bb.feature.hide-help": true,
      "bb.feature.hide-quick-start": true,
      "bb.feature.hide-release-remind": true,
      "bb.feature.disallow-navigate-to-console": true,
      "bb.feature.console.hide-sidebar": true,
      "bb.feature.console.hide-header": true,
      "bb.feature.console.hide-quick-action": true,
      "bb.feature.databases.hide-unassigned": true,
      "bb.feature.databases.hide-inalterable": true,
      "bb.feature.sql-editor.disallow-share-worksheet": true,
    });
  }

  const customTheme = query.get("customTheme");
  if (customTheme === "lixiang") {
    actuatorStore.overrideAppFeatures({
      "bb.feature.custom-color-scheme": {
        "--color-accent": "#00665f",
        "--color-accent-hover": "#00554f",
        "--color-accent-disabled": "#b8c3c3",
      },
      "bb.feature.sql-editor.custom-query-datasource": true,
      "bb.feature.sql-editor.disallow-export-query-data": true,
    });
    if (actuatorStore.appProfile.embedded) {
      actuatorStore.overrideAppFeatures({
        "bb.feature.hide-issue-review-actions": true,
      });
    }
  }

  useCustomTheme(useAppFeature("bb.feature.custom-color-scheme"));
};
