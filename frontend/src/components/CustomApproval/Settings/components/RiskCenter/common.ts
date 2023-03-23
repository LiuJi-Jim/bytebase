import {
  Risk,
  Risk_Source,
  risk_SourceToJSON,
} from "@/types/proto/v1/risk_service";
import { useI18n } from "vue-i18n";

export const sourceText = (source: Risk_Source) => {
  const { t, te } = useI18n();
  const name = risk_SourceToJSON(source);
  const keypath = `custom-approval.security-rule.risk.namespace.${name.toLowerCase()}`;
  if (te(keypath)) {
    return t(keypath);
  }
  return name;
};

export const levelText = (level: number) => {
  const { t, te } = useI18n();
  const keypath = `custom-approval.security-rule.risk.risk-level.${level}`;
  if (te(keypath)) {
    return t(keypath);
  }
  return String(level);
};

export const orderByLevelDesc = (a: Risk, b: Risk): number => {
  if (a.level !== b.level) return -(a.level - b.level);
  if (a.name === b.name) return 0;
  return a.name < b.name ? -1 : 1;
};
